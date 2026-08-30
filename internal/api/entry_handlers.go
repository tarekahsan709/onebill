// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package api // import "onebill.app/internal/api"

import (
	json_parser "encoding/json"
	"errors"
	"net/http"
	"time"

	"onebill.app/internal/crypto"
	"onebill.app/internal/http/request"
	"onebill.app/internal/http/response"
	"onebill.app/internal/model"
	"onebill.app/internal/storage"
	"onebill.app/internal/validator"
)

func (h *handler) getEntryFromBuilder(w http.ResponseWriter, r *http.Request, b *storage.EntryQueryBuilder) {
	entry, err := b.GetEntry()
	if err != nil {
		response.JSONServerError(w, r, err)
		return
	}

	if entry == nil {
		response.JSONNotFound(w, r)
		return
	}

	response.JSON(w, r, entry)
}

func (h *handler) getFeedEntryHandler(w http.ResponseWriter, r *http.Request) {
	feedID := request.RouteInt64Param(r, "feedID")
	if feedID == 0 {
		response.JSONBadRequest(w, r, errors.New("invalid feed ID"))
		return
	}

	entryID := request.RouteInt64Param(r, "entryID")
	if entryID == 0 {
		response.JSONBadRequest(w, r, errors.New("invalid entry ID"))
		return
	}

	builder := h.store.NewEntryQueryBuilder(request.UserID(r)).
		WithFeedID(feedID).
		WithEntryIDs(entryID)

	h.getEntryFromBuilder(w, r, builder)
}

func (h *handler) getEntryHandler(w http.ResponseWriter, r *http.Request) {
	entryID := request.RouteInt64Param(r, "entryID")
	if entryID == 0 {
		response.JSONBadRequest(w, r, errors.New("invalid entry ID"))
		return
	}

	builder := h.store.NewEntryQueryBuilder(request.UserID(r)).
		WithEntryIDs(entryID)

	h.getEntryFromBuilder(w, r, builder)
}

func (h *handler) getFeedEntriesHandler(w http.ResponseWriter, r *http.Request) {
	feedID := request.RouteInt64Param(r, "feedID")
	if feedID == 0 {
		response.JSONBadRequest(w, r, errors.New("invalid feed ID"))
		return
	}

	h.findEntries(w, r, feedID)
}

func (h *handler) getEntriesHandler(w http.ResponseWriter, r *http.Request) {
	h.findEntries(w, r, 0)
}

func (h *handler) findEntries(w http.ResponseWriter, r *http.Request, feedID int64) {
	statuses := request.QueryStringParamList(r, "status")
	for _, status := range statuses {
		if err := validator.ValidateEntryStatus(status); err != nil {
			response.JSONBadRequest(w, r, err)
			return
		}
	}

	order := request.QueryStringParam(r, "order", model.DefaultSortingOrder)
	if err := validator.ValidateEntryOrder(order); err != nil {
		response.JSONBadRequest(w, r, err)
		return
	}

	direction := request.QueryStringParam(r, "direction", model.DefaultSortingDirection)
	if err := validator.ValidateDirection(direction); err != nil {
		response.JSONBadRequest(w, r, err)
		return
	}

	limit := request.QueryIntParam(r, "limit", 100)
	offset := request.QueryIntParam(r, "offset", 0)
	if err := validator.ValidateRange(offset, limit); err != nil {
		response.JSONBadRequest(w, r, err)
		return
	}

	userID := request.UserID(r)

	feedID = request.QueryInt64Param(r, "feed_id", feedID)
	if feedID > 0 {
		exists, err := h.store.FeedExists(userID, feedID)
		if err != nil {
			response.JSONServerError(w, r, err)
			return
		}
		if !exists {
			response.JSONBadRequest(w, r, errors.New("invalid feed ID"))
			return
		}
	}

	builder := h.store.NewEntryQueryBuilder(userID).
		WithFeedID(feedID).
		WithStatuses(statuses...).
		WithSorting(order, direction).
		WithOffset(offset).
		WithLimit(limit)

	builder = configureFilters(builder, r)

	entries, count, err := builder.GetEntriesWithCount()
	if err != nil {
		response.JSONServerError(w, r, err)
		return
	}

	response.JSON(w, r, &entriesResponse{Total: count, Entries: entries})
}

func (h *handler) setEntryStatusHandler(w http.ResponseWriter, r *http.Request) {
	var entriesStatusUpdateRequest model.EntriesStatusUpdateRequest
	if err := json_parser.NewDecoder(r.Body).Decode(&entriesStatusUpdateRequest); err != nil {
		response.JSONBadRequest(w, r, err)
		return
	}

	if err := validator.ValidateEntriesStatusUpdateRequest(&entriesStatusUpdateRequest); err != nil {
		response.JSONBadRequest(w, r, err)
		return
	}

	if _, err := h.store.SetEntriesStatus(request.UserID(r), entriesStatusUpdateRequest.EntryIDs, entriesStatusUpdateRequest.Status); err != nil {
		response.JSONServerError(w, r, err)
		return
	}

	response.NoContent(w, r)
}

func (h *handler) importFeedEntryHandler(w http.ResponseWriter, r *http.Request) {
	userID := request.UserID(r)

	feedID := request.RouteInt64Param(r, "feedID")
	if feedID <= 0 {
		response.JSONBadRequest(w, r, errors.New("invalid feed ID"))
		return
	}

	exists, err := h.store.FeedExists(userID, feedID)
	if err != nil {
		response.JSONServerError(w, r, err)
		return
	}
	if !exists {
		response.JSONBadRequest(w, r, errors.New("feed does not exist"))
		return
	}

	var importRequest entryImportRequest
	if err := json_parser.NewDecoder(r.Body).Decode(&importRequest); err != nil {
		response.JSONBadRequest(w, r, err)
		return
	}

	if importRequest.URL == "" {
		response.JSONBadRequest(w, r, errors.New("url is required"))
		return
	}

	if importRequest.Status == "" {
		importRequest.Status = model.EntryStatusRead
	}

	if err := validator.ValidateEntryStatus(importRequest.Status); err != nil {
		response.JSONBadRequest(w, r, err)
		return
	}

	entry := model.NewEntry()
	entry.URL = importRequest.URL
	entry.Author = importRequest.Author

	if importRequest.PublishedAt > 0 {
		entry.Date = time.Unix(importRequest.PublishedAt, 0).UTC()
	} else {
		entry.Date = time.Now().UTC()
	}

	if importRequest.Title == "" {
		entry.Title = entry.URL
	} else {
		entry.Title = importRequest.Title
	}

	hashInput := importRequest.ExternalID
	if hashInput == "" {
		hashInput = importRequest.URL
	}
	entry.Hash = crypto.HashFromBytes([]byte(hashInput))

	created, err := h.store.InsertEntryForFeed(userID, feedID, entry)
	if errors.Is(err, storage.ErrEntryTombstoned) {
		response.JSONBadRequest(w, r, err)
		return
	}
	if err != nil {
		response.JSONServerError(w, r, err)
		return
	}

	if _, err := h.store.SetEntriesStatus(userID, []int64{entry.ID}, importRequest.Status); err != nil {
		response.JSONServerError(w, r, err)
		return
	}
	entry.Status = importRequest.Status

	if created {
		response.JSONCreated(w, r, entryIDResponse{ID: entry.ID})
	} else {
		response.JSON(w, r, entryIDResponse{ID: entry.ID})
	}
}

func (h *handler) getEntryIDsHandler(w http.ResponseWriter, r *http.Request) {
	if request.HasQueryParam(r, "status") {
		statusValue := request.QueryStringParam(r, "status", "")
		if statusValue != model.EntryStatusRead && statusValue != model.EntryStatusUnread {
			response.JSONBadRequest(w, r, errors.New(`invalid status parameter, must be "read" or "unread"`))
			return
		}
	}

	limit, offset := parseEntryIDsParams(r)
	builder := h.store.NewEntryQueryBuilder(request.UserID(r)).
		WithSorting("id", "DESC").
		WithLimitAndMaximum(limit, model.MaxEntryIDsLimit).
		WithOffset(offset)

	if request.HasQueryParam(r, "status") {
		builder.WithStatuses(request.QueryStringParam(r, "status", ""))
	}

	entryIDs, total, err := builder.GetEntryIDsWithCount()
	if err != nil {
		response.JSONServerError(w, r, err)
		return
	}

	if entryIDs == nil {
		entryIDs = []int64{}
	}

	response.JSON(w, r, entryIDsResponse{Total: total, EntryIDs: entryIDs})
}

func (h *handler) flushHistoryHandler(w http.ResponseWriter, r *http.Request) {
	loggedUserID := request.UserID(r)
	go h.store.FlushHistory(loggedUserID)
	response.JSONAccepted(w, r)
}

func configureFilters(builder *storage.EntryQueryBuilder, r *http.Request) *storage.EntryQueryBuilder {
	if beforeEntryID := request.QueryInt64Param(r, "before_entry_id", 0); beforeEntryID > 0 {
		builder = builder.BeforeEntryID(beforeEntryID)
	}

	if afterEntryID := request.QueryInt64Param(r, "after_entry_id", 0); afterEntryID > 0 {
		builder = builder.AfterEntryID(afterEntryID)
	}

	if beforePublishedTimestamp := request.QueryInt64Param(r, "before", 0); beforePublishedTimestamp > 0 {
		builder = builder.BeforePublishedDate(time.Unix(beforePublishedTimestamp, 0))
	}

	if afterPublishedTimestamp := request.QueryInt64Param(r, "after", 0); afterPublishedTimestamp > 0 {
		builder = builder.AfterPublishedDate(time.Unix(afterPublishedTimestamp, 0))
	}

	if beforePublishedTimestamp := request.QueryInt64Param(r, "published_before", 0); beforePublishedTimestamp > 0 {
		builder = builder.BeforePublishedDate(time.Unix(beforePublishedTimestamp, 0))
	}

	if afterPublishedTimestamp := request.QueryInt64Param(r, "published_after", 0); afterPublishedTimestamp > 0 {
		builder = builder.AfterPublishedDate(time.Unix(afterPublishedTimestamp, 0))
	}

	if beforeChangedTimestamp := request.QueryInt64Param(r, "changed_before", 0); beforeChangedTimestamp > 0 {
		builder = builder.BeforeChangedDate(time.Unix(beforeChangedTimestamp, 0))
	}

	if afterChangedTimestamp := request.QueryInt64Param(r, "changed_after", 0); afterChangedTimestamp > 0 {
		builder = builder.AfterChangedDate(time.Unix(afterChangedTimestamp, 0))
	}

	return builder
}

func parseEntryIDsParams(r *http.Request) (limit, offset int) {
	limit = request.QueryIntParam(r, "limit", model.MaxEntryIDsLimit)
	if limit <= 0 || limit > model.MaxEntryIDsLimit {
		limit = model.MaxEntryIDsLimit
	}
	offset = request.QueryIntParam(r, "offset", 0)
	return limit, offset
}
