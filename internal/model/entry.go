// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package model // import "onebill.app/internal/model"

import (
	"time"
)

// Entry statuses and default sorting order.
const (
	EntryStatusUnread       = "unread"
	EntryStatusRead         = "read"
	DefaultSortingOrder     = "published_at"
	DefaultSortingDirection = "asc"
)

// MaxEntryLimit is the maximum allowed value for the "limit" query parameter
// and for the user "entries_per_page" preference.
const MaxEntryLimit = 1000

// MaxEntryIDsLimit is the maximum allowed value for the "limit" query parameter
// for the entry ID list endpoints.
const MaxEntryIDsLimit = 10000

// Entry represents a feed item in the system.
type Entry struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	FeedID    int64     `json:"feed_id"`
	Status    string    `json:"status"`
	Hash      string    `json:"hash"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Date      time.Time `json:"published_at"`
	CreatedAt time.Time `json:"created_at"`
	ChangedAt time.Time `json:"changed_at"`
	Author    string    `json:"author"`
	Feed      *Feed     `json:"feed,omitempty"`
}

func NewEntry() *Entry {
	return &Entry{Feed: &Feed{}}
}

// Entries represents a list of entries.
type Entries []*Entry

// EntriesStatusUpdateRequest represents a request to change entries status.
type EntriesStatusUpdateRequest struct {
	EntryIDs []int64 `json:"entry_ids"`
	Status   string  `json:"status"`
}
