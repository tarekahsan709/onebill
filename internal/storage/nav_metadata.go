// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package storage // import "onebill.app/internal/storage"

import (
	"log/slog"

	"onebill.app/internal/config"
)

type NavMetadata struct {
	CountUnread     int
	CountErrorFeeds int
}

// GetNavMetadata returns the navigation metadata for the given user in a
// single SQL query.
func (s *Storage) GetNavMetadata(userID int64) (NavMetadata, error) {
	query := `
		SELECT
			(SELECT count(*)
			   FROM entries e
			  WHERE e.user_id = $1
			    AND e.status = 'unread'
			) AS count_unread,
	`
	if config.Opts.PollingParsingErrorLimit() == 0 {
		// zero means unlimited amount of accepted errors
		query += `(SELECT $2) AS count_error_feeds`
	} else {
		query += `(SELECT count(*)
			   FROM feeds
			  WHERE user_id = $1
			    AND parsing_error_count >= $2
			) AS count_error_feeds
			 `
	}

	var countUnread, countErrorFeeds int

	err := s.db.QueryRow(query, userID, config.Opts.PollingParsingErrorLimit()).Scan(
		&countUnread,
		&countErrorFeeds,
	)
	if err != nil {
		slog.Error("Unable to fetch navigation metadata",
			slog.Int64("user_id", userID),
			slog.Any("error", err),
		)
		return NavMetadata{}, err
	}

	return NavMetadata{
		CountUnread:     countUnread,
		CountErrorFeeds: countErrorFeeds,
	}, nil
}
