// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package api // import "onebill.app/internal/api"

import (
	"onebill.app/internal/model"
)

type entriesResponse struct {
	Total   int           `json:"total"`
	Entries model.Entries `json:"entries"`
}

type entryIDResponse struct {
	ID int64 `json:"id"`
}

type entryIDsResponse struct {
	Total    int     `json:"total"`
	EntryIDs []int64 `json:"entry_ids"`
}

type entryImportRequest struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	PublishedAt int64  `json:"published_at"`
	Status      string `json:"status"`
	ExternalID  string `json:"external_id"`
}

type feedCreationResponse struct {
	FeedID int64 `json:"feed_id"`
}

type versionResponse struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
	Compiler  string `json:"compiler"`
	Arch      string `json:"arch"`
	OS        string `json:"os"`
}
