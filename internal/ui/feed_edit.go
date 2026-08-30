// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package ui // import "onebill.app/internal/ui"

import (
	"net/http"

	"onebill.app/internal/http/request"
	"onebill.app/internal/http/response"
	"onebill.app/internal/ui/form"
	"onebill.app/internal/ui/view"
)

func (h *handler) showEditFeedPage(w http.ResponseWriter, r *http.Request) {
	user, err := h.store.UserByID(request.UserID(r))
	if err != nil {
		response.HTMLServerError(w, r, err)
		return
	}

	feedID := request.RouteInt64Param(r, "feedID")
	feed, err := h.store.FeedByID(user.ID, feedID)
	if err != nil {
		response.HTMLServerError(w, r, err)
		return
	}

	if feed == nil {
		response.HTMLNotFound(w, r)
		return
	}

	feedForm := form.FeedForm{
		SiteURL:  feed.SiteURL,
		FeedURL:  feed.FeedURL,
		Title:    feed.Title,
		Disabled: feed.Disabled,
	}

	view := view.New(h.tpl, r)
	view.Set("form", feedForm)
	view.Set("feed", feed)
	view.Set("menu", "feeds")
	view.Set("user", user)
	navMetadata, _ := h.store.GetNavMetadata(user.ID)
	view.Set("countUnread", navMetadata.CountUnread)
	view.Set("countErrorFeeds", navMetadata.CountErrorFeeds)

	response.HTML(w, r, view.Render("edit_feed"))
}
