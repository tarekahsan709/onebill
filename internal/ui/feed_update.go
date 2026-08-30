// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package ui // import "onebill.app/internal/ui"

import (
	"net/http"

	"onebill.app/internal/http/request"
	"onebill.app/internal/http/response"
	"onebill.app/internal/model"
	"onebill.app/internal/ui/form"
	"onebill.app/internal/ui/view"
	"onebill.app/internal/validator"
)

func (h *handler) updateFeed(w http.ResponseWriter, r *http.Request) {
	loggedUser, err := h.store.UserByID(request.UserID(r))
	if err != nil {
		response.HTMLServerError(w, r, err)
		return
	}

	feedID := request.RouteInt64Param(r, "feedID")
	feed, err := h.store.FeedByID(loggedUser.ID, feedID)
	if err != nil {
		response.HTMLServerError(w, r, err)
		return
	}

	if feed == nil {
		response.HTMLNotFound(w, r)
		return
	}

	feedForm := form.NewFeedForm(r)

	view := view.New(h.tpl, r)
	view.Set("form", feedForm)
	view.Set("feed", feed)
	view.Set("menu", "feeds")
	view.Set("user", loggedUser)
	navMetadata, _ := h.store.GetNavMetadata(loggedUser.ID)
	view.Set("countUnread", navMetadata.CountUnread)
	view.Set("countErrorFeeds", navMetadata.CountErrorFeeds)

	feedModificationRequest := &model.FeedModificationRequest{
		FeedURL: model.OptionalString(feedForm.FeedURL),
		SiteURL: model.OptionalString(feedForm.SiteURL),
		Title:   model.OptionalString(feedForm.Title),
	}

	if validationErr := validator.ValidateFeedModification(h.store, loggedUser.ID, feed.ID, feedModificationRequest); validationErr != nil {
		view.Set("errorMessage", validationErr.Translate(loggedUser.Language))
		response.HTML(w, r, view.Render("edit_feed"))
		return
	}

	err = h.store.UpdateFeed(feedForm.Merge(feed))
	if err != nil {
		response.HTMLServerError(w, r, err)
		return
	}

	response.HTMLRedirect(w, r, h.routePath("/feed/%d/entries", feed.ID))
}
