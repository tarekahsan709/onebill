// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package form // import "onebill.app/internal/ui/form"

import (
	"net/http"

	"onebill.app/internal/model"
)

// FeedForm represents a feed form in the UI
type FeedForm struct {
	FeedURL  string
	SiteURL  string
	Title    string
	Disabled bool
}

// Merge updates the fields of the given feed.
func (f FeedForm) Merge(feed *model.Feed) *model.Feed {
	feed.Title = f.Title
	feed.SiteURL = f.SiteURL
	feed.FeedURL = f.FeedURL
	feed.Disabled = f.Disabled
	feed.ParsingErrorCount = 0
	feed.ParsingErrorMsg = ""
	return feed
}

// NewFeedForm parses the HTTP request and returns a FeedForm
func NewFeedForm(r *http.Request) *FeedForm {
	return &FeedForm{
		FeedURL:  r.FormValue("feed_url"),
		SiteURL:  r.FormValue("site_url"),
		Title:    r.FormValue("title"),
		Disabled: r.FormValue("disabled") == "1",
	}
}
