// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package api // import "onebill.app/internal/api"

import (
	json_parser "encoding/json"
	"errors"
	"net/http"
	"time"

	"onebill.app/internal/http/request"
	"onebill.app/internal/http/response"
	"onebill.app/internal/model"
	"onebill.app/internal/validator"
)

func (h *handler) createFeedHandler(w http.ResponseWriter, r *http.Request) {
	userID := request.UserID(r)

	var feedCreationRequest model.FeedCreationRequest
	if err := json_parser.NewDecoder(r.Body).Decode(&feedCreationRequest); err != nil {
		response.JSONBadRequest(w, r, err)
		return
	}

	if validationErr := validator.ValidateFeedCreation(h.store, userID, &feedCreationRequest); validationErr != nil {
		response.JSONBadRequest(w, r, validationErr.Error())
		return
	}

	feed := &model.Feed{
		UserID:   userID,
		FeedURL:  feedCreationRequest.FeedURL,
		SiteURL:  feedCreationRequest.FeedURL,
		Title:    feedCreationRequest.Title,
		Disabled: feedCreationRequest.Disabled,
	}
	if feed.Title == "" {
		feed.Title = feed.FeedURL
	}

	if err := h.store.CreateFeed(feed); err != nil {
		response.JSONServerError(w, r, err)
		return
	}

	response.JSONCreated(w, r, &feedCreationResponse{FeedID: feed.ID})
}

func (h *handler) updateFeedHandler(w http.ResponseWriter, r *http.Request) {
	feedID := request.RouteInt64Param(r, "feedID")
	if feedID == 0 {
		response.JSONBadRequest(w, r, errors.New("invalid feed ID"))
		return
	}

	var feedModificationRequest model.FeedModificationRequest
	if err := json_parser.NewDecoder(r.Body).Decode(&feedModificationRequest); err != nil {
		response.JSONBadRequest(w, r, err)
		return
	}

	userID := request.UserID(r)
	originalFeed, err := h.store.FeedByID(userID, feedID)
	if err != nil {
		response.JSONServerError(w, r, err)
		return
	}

	if originalFeed == nil {
		response.JSONNotFound(w, r)
		return
	}

	if validationErr := validator.ValidateFeedModification(h.store, userID, originalFeed.ID, &feedModificationRequest); validationErr != nil {
		response.JSONBadRequest(w, r, validationErr.Error())
		return
	}

	feedModificationRequest.Patch(originalFeed)
	originalFeed.ResetErrorCounter()
	if err := h.store.UpdateFeed(originalFeed); err != nil {
		response.JSONServerError(w, r, err)
		return
	}

	originalFeed, err = h.store.FeedByID(userID, feedID)
	if err != nil {
		response.JSONServerError(w, r, err)
		return
	}

	response.JSONCreated(w, r, originalFeed)
}

func (h *handler) markFeedAsReadHandler(w http.ResponseWriter, r *http.Request) {
	userID := request.UserID(r)

	feedID := request.RouteInt64Param(r, "feedID")
	if feedID == 0 {
		response.JSONBadRequest(w, r, errors.New("invalid feed ID"))
		return
	}

	exists, err := h.store.FeedExists(userID, feedID)
	if err != nil {
		response.JSONServerError(w, r, err)
		return
	}
	if !exists {
		response.JSONNotFound(w, r)
		return
	}

	if err := h.store.MarkFeedAsRead(userID, feedID, time.Now()); err != nil {
		response.JSONServerError(w, r, err)
		return
	}

	response.NoContent(w, r)
}

func (h *handler) getFeedsHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := h.store.Feeds(request.UserID(r))
	if err != nil {
		response.JSONServerError(w, r, err)
		return
	}

	response.JSON(w, r, feeds)
}

func (h *handler) fetchCountersHandler(w http.ResponseWriter, r *http.Request) {
	counters, err := h.store.FetchCounters(request.UserID(r))
	if err != nil {
		response.JSONServerError(w, r, err)
		return
	}

	response.JSON(w, r, counters)
}

func (h *handler) getFeedHandler(w http.ResponseWriter, r *http.Request) {
	feedID := request.RouteInt64Param(r, "feedID")
	if feedID == 0 {
		response.JSONBadRequest(w, r, errors.New("invalid feed ID"))
		return
	}

	feed, err := h.store.FeedByID(request.UserID(r), feedID)
	if err != nil {
		response.JSONServerError(w, r, err)
		return
	}

	if feed == nil {
		response.JSONNotFound(w, r)
		return
	}

	response.JSON(w, r, feed)
}

func (h *handler) removeFeedHandler(w http.ResponseWriter, r *http.Request) {
	feedID := request.RouteInt64Param(r, "feedID")
	if feedID == 0 {
		response.JSONBadRequest(w, r, errors.New("invalid feed ID"))
		return
	}

	userID := request.UserID(r)
	exists, err := h.store.FeedExists(userID, feedID)
	if err != nil {
		response.JSONServerError(w, r, err)
		return
	}
	if !exists {
		response.JSONNotFound(w, r)
		return
	}

	if err := h.store.RemoveFeed(userID, feedID); err != nil {
		response.JSONServerError(w, r, err)
		return
	}

	response.NoContent(w, r)
}
