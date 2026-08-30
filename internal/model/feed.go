// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package model // import "onebill.app/internal/model"

import (
	"fmt"
	"time"
)

// Default settings for the feed query builder.
const (
	DefaultFeedSorting          = "parsing_error_count"
	DefaultFeedSortingDirection = "desc"
)

// Feed represents a feed in the application.
type Feed struct {
	ID                int64     `json:"id"`
	UserID            int64     `json:"user_id"`
	FeedURL           string    `json:"feed_url"`
	SiteURL           string    `json:"site_url"`
	Title             string    `json:"title"`
	CheckedAt         time.Time `json:"checked_at"`
	NextCheckAt       time.Time `json:"next_check_at"`
	ParsingErrorMsg   string    `json:"parsing_error_message"`
	ParsingErrorCount int       `json:"parsing_error_count"`
	Disabled          bool      `json:"disabled"`

	// Non-persisted attributes
	Entries Entries `json:"entries,omitempty"`

	// Internal attributes (not exposed in the API and not persisted in the database)
	UnreadCount            int `json:"-"`
	ReadCount              int `json:"-"`
	NumberOfVisibleEntries int `json:"-"`
}

type FeedCounters struct {
	ReadCounters   map[int64]int `json:"reads"`
	UnreadCounters map[int64]int `json:"unreads"`
}

func (f *Feed) String() string {
	return fmt.Sprintf("ID=%d, UserID=%d, FeedURL=%s, SiteURL=%s, Title=%s",
		f.ID,
		f.UserID,
		f.FeedURL,
		f.SiteURL,
		f.Title,
	)
}

// WithTranslatedErrorMessage adds a new error message and increment the error counter.
func (f *Feed) WithTranslatedErrorMessage(message string) {
	f.ParsingErrorCount++
	f.ParsingErrorMsg = message
}

// ResetErrorCounter removes all previous errors.
func (f *Feed) ResetErrorCounter() {
	f.ParsingErrorCount = 0
	f.ParsingErrorMsg = ""
}

// FeedCreationRequest represents the request to create a feed.
type FeedCreationRequest struct {
	FeedURL  string `json:"feed_url"`
	Title    string `json:"title"`
	Disabled bool   `json:"disabled"`
}

// FeedModificationRequest represents the request to update a feed.
type FeedModificationRequest struct {
	FeedURL  *string `json:"feed_url"`
	SiteURL  *string `json:"site_url"`
	Title    *string `json:"title"`
	Disabled *bool   `json:"disabled"`
}

// Patch updates a feed with modified values.
func (f *FeedModificationRequest) Patch(feed *Feed) {
	if f.FeedURL != nil && *f.FeedURL != "" {
		feed.FeedURL = *f.FeedURL
	}

	if f.SiteURL != nil && *f.SiteURL != "" {
		feed.SiteURL = *f.SiteURL
	}

	if f.Title != nil && *f.Title != "" {
		feed.Title = *f.Title
	}

	if f.Disabled != nil {
		feed.Disabled = *f.Disabled
	}
}

// Feeds is a list of feed
type Feeds []*Feed
