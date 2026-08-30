// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package storage // import "onebill.app/internal/storage"

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"

	"onebill.app/internal/config"
	"onebill.app/internal/model"
)

type byStateAndName struct{ f model.Feeds }

func (l byStateAndName) Len() int      { return len(l.f) }
func (l byStateAndName) Swap(i, j int) { l.f[i], l.f[j] = l.f[j], l.f[i] }
func (l byStateAndName) Less(i, j int) bool {
	// disabled test first, since we don't care about errors if disabled
	if l.f[i].Disabled != l.f[j].Disabled {
		return l.f[j].Disabled
	}
	if l.f[i].ParsingErrorCount != l.f[j].ParsingErrorCount {
		return l.f[i].ParsingErrorCount > l.f[j].ParsingErrorCount
	}
	if l.f[i].UnreadCount != l.f[j].UnreadCount {
		return l.f[i].UnreadCount > l.f[j].UnreadCount
	}
	return l.f[i].Title < l.f[j].Title
}

// FeedExists checks if the given feed exists. The returned error is
// non-nil only for a genuine query/connection failure; a feed that doesn't
// exist is reported as (false, nil), matching the underlying sql.ErrNoRows
// case.
func (s *Storage) FeedExists(userID, feedID int64) (bool, error) {
	var result bool
	query := `SELECT true FROM feeds WHERE user_id=$1 AND id=$2 LIMIT 1`
	if err := s.db.QueryRow(query, userID, feedID).Scan(&result); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, fmt.Errorf(`store: unable to check if feed exists: %w`, err)
	}
	return result, nil
}

// FeedURLExists returns true if the given feed URL already exists for the user.
func (s *Storage) FeedURLExists(userID int64, feedURL string) bool {
	var result bool
	query := `SELECT true FROM feeds WHERE user_id=$1 AND feed_url=$2 LIMIT 1`
	s.db.QueryRow(query, userID, feedURL).Scan(&result)
	return result
}

// AnotherFeedURLExists returns true if another feed with the same URL exists for the user.
func (s *Storage) AnotherFeedURLExists(userID, feedID int64, feedURL string) bool {
	var result bool
	query := `SELECT true FROM feeds WHERE id <> $1 AND user_id=$2 AND feed_url=$3 LIMIT 1`
	s.db.QueryRow(query, feedID, userID, feedURL).Scan(&result)
	return result
}

// CountAllFeeds returns the number of feeds keyed by enabled, disabled, and total.
func (s *Storage) CountAllFeeds() (map[string]int64, error) {
	rows, err := s.db.Query(`SELECT disabled, count(*) FROM feeds GROUP BY disabled`)
	if err != nil {
		return nil, fmt.Errorf("storage: unable to count feeds: %w", err)
	}
	defer rows.Close()

	results := map[string]int64{
		"enabled":  0,
		"disabled": 0,
		"total":    0,
	}

	for rows.Next() {
		var disabled bool
		var count int64

		if err := rows.Scan(&disabled, &count); err != nil {
			continue
		}

		if disabled {
			results["disabled"] = count
		} else {
			results["enabled"] = count
		}
	}

	results["total"] = results["disabled"] + results["enabled"]
	return results, nil
}

// CountAllFeedsWithErrors returns the number of feeds with sync errors.
func (s *Storage) CountAllFeedsWithErrors() (int, error) {
	pollingParsingErrorLimit := config.Opts.PollingParsingErrorLimit()
	if pollingParsingErrorLimit <= 0 {
		pollingParsingErrorLimit = 1
	}
	query := `SELECT count(*) FROM feeds WHERE parsing_error_count >= $1`
	var result int
	err := s.db.QueryRow(query, pollingParsingErrorLimit).Scan(&result)
	if err != nil {
		return 0, fmt.Errorf("storage: unable to count feeds with errors: %w", err)
	}

	return result, nil
}

// Feeds returns all feeds that belong to the given user.
func (s *Storage) Feeds(userID int64) (model.Feeds, error) {
	return s.NewFeedQueryBuilder(userID).
		WithSorting(model.DefaultFeedSorting, model.DefaultFeedSortingDirection).
		GetFeeds()
}

// FeedsWithCounters returns all feeds of the given user with read and unread entry counters.
func (s *Storage) FeedsWithCounters(userID int64) (model.Feeds, error) {
	result, err := s.NewFeedQueryBuilder(userID).
		WithCounters().
		WithSorting(model.DefaultFeedSorting, model.DefaultFeedSortingDirection).
		GetFeeds()
	if err != nil {
		return result, err
	}

	sort.Sort(byStateAndName{result})
	return result, nil
}

// FetchCounters returns the per-feed read and unread entry counts for the given user.
func (s *Storage) FetchCounters(userID int64) (model.FeedCounters, error) {
	reads, unreads, err := s.NewFeedQueryBuilder(userID).
		WithCounters().
		fetchFeedCounter()

	return model.FeedCounters{ReadCounters: reads, UnreadCounters: unreads}, err
}

// FeedByID returns the feed with the given ID.
func (s *Storage) FeedByID(userID, feedID int64) (*model.Feed, error) {
	feed, err := s.NewFeedQueryBuilder(userID).
		WithFeedID(feedID).
		GetFeed()

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf(`store: unable to fetch feed #%d: %v`, feedID, err)
	}

	return feed, nil
}

// CreateFeed creates a new feed.
func (s *Storage) CreateFeed(feed *model.Feed) error {
	query := `
		INSERT INTO feeds
			(feed_url, site_url, title, user_id, disabled)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING
			id
	`
	err := s.db.QueryRow(
		query,
		feed.FeedURL,
		feed.SiteURL,
		feed.Title,
		feed.UserID,
		feed.Disabled,
	).Scan(&feed.ID)
	if err != nil {
		return fmt.Errorf(`store: unable to create feed %q: %v`, feed.FeedURL, err)
	}

	return nil
}

// UpdateFeed updates an existing feed.
func (s *Storage) UpdateFeed(feed *model.Feed) (err error) {
	query := `
		UPDATE
			feeds
		SET
			feed_url=$1,
			site_url=$2,
			title=$3,
			checked_at=$4,
			next_check_at=$5,
			parsing_error_msg=$6,
			parsing_error_count=$7,
			disabled=$8
		WHERE
			id=$9 AND user_id=$10
	`
	_, err = s.db.Exec(query,
		feed.FeedURL,
		feed.SiteURL,
		feed.Title,
		feed.CheckedAt,
		feed.NextCheckAt,
		feed.ParsingErrorMsg,
		feed.ParsingErrorCount,
		feed.Disabled,
		feed.ID,
		feed.UserID,
	)
	if err != nil {
		return fmt.Errorf(`store: unable to update feed #%d (%s): %v`, feed.ID, feed.FeedURL, err)
	}

	return nil
}

// RemoveFeed removes the given feed along with its entries.
func (s *Storage) RemoveFeed(userID, feedID int64) error {
	if _, err := s.db.Exec(`DELETE FROM feeds WHERE id=$1 AND user_id=$2`, feedID, userID); err != nil {
		return fmt.Errorf(`store: unable to delete feed #%d: %v`, feedID, err)
	}
	return nil
}

// ResetFeedErrors clears the sync error fields for all feeds.
func (s *Storage) ResetFeedErrors() error {
	_, err := s.db.Exec(`UPDATE feeds SET parsing_error_count=0, parsing_error_msg=''`)
	return err
}

// ResetNextCheckAt schedules all feeds to be checked immediately.
func (s *Storage) ResetNextCheckAt() error {
	_, err := s.db.Exec(`UPDATE feeds SET next_check_at=now()`)
	return err
}
