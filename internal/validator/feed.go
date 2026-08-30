// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package validator // import "onebill.app/internal/validator"

import (
	"onebill.app/internal/locale"
	"onebill.app/internal/model"
	"onebill.app/internal/storage"
	"onebill.app/internal/urllib"
)

// ValidateFeedCreation validates feed creation.
func ValidateFeedCreation(store *storage.Storage, userID int64, request *model.FeedCreationRequest) *locale.LocalizedError {
	if request.FeedURL == "" {
		return locale.NewLocalizedError("error.feed_mandatory_fields")
	}

	if !urllib.IsAbsoluteURL(request.FeedURL) {
		return locale.NewLocalizedError("error.invalid_feed_url")
	}

	if store.FeedURLExists(userID, request.FeedURL) {
		return locale.NewLocalizedError("error.feed_already_exists")
	}

	return nil
}

// ValidateFeedModification validates feed modification.
func ValidateFeedModification(store *storage.Storage, userID, feedID int64, request *model.FeedModificationRequest) *locale.LocalizedError {
	if request.FeedURL != nil {
		if *request.FeedURL == "" {
			return locale.NewLocalizedError("error.feed_url_not_empty")
		}

		if !urllib.IsAbsoluteURL(*request.FeedURL) {
			return locale.NewLocalizedError("error.invalid_feed_url")
		}

		if store.AnotherFeedURLExists(userID, feedID, *request.FeedURL) {
			return locale.NewLocalizedError("error.feed_already_exists")
		}
	}

	if request.SiteURL != nil {
		if *request.SiteURL == "" {
			return locale.NewLocalizedError("error.site_url_not_empty")
		}

		if !urllib.IsAbsoluteURL(*request.SiteURL) {
			return locale.NewLocalizedError("error.invalid_site_url")
		}
	}

	if request.Title != nil {
		if *request.Title == "" {
			return locale.NewLocalizedError("error.feed_title_not_empty")
		}
	}

	return nil
}
