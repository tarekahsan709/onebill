// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package storage // import "onebill.app/internal/storage"

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/lib/pq"
	"onebill.app/internal/model"
	"onebill.app/internal/timezone"
)

// feedQueryBuilder builds a SQL query to fetch feeds.
type feedQueryBuilder struct {
	db                *sql.DB
	args              []any
	conditions        []string
	sortExpressions   []string
	limit             int
	offset            int
	withCounters      bool
	counterArgs       []any
	counterConditions []string
}

// NewFeedQueryBuilder returns a new FeedQueryBuilder.
func (s *Storage) NewFeedQueryBuilder(userID int64) *feedQueryBuilder {
	return &feedQueryBuilder{
		db:                s.db,
		args:              []any{userID},
		conditions:        []string{"f.user_id = $1"},
		counterArgs:       []any{userID, model.EntryStatusRead, model.EntryStatusUnread},
		counterConditions: []string{"e.user_id = $1", "e.status IN ($2, $3)"},
	}
}

// WithFeedID filter by feed ID.
func (f *feedQueryBuilder) WithFeedID(feedID int64) *feedQueryBuilder {
	if feedID > 0 {
		f.conditions = append(f.conditions, "f.id = $"+strconv.Itoa(len(f.args)+1))
		f.args = append(f.args, feedID)
	}
	return f
}

// WithCounters let the builder return feeds with counters of statuses of entries.
func (f *feedQueryBuilder) WithCounters() *feedQueryBuilder {
	f.withCounters = true
	return f
}

// WithSorting add a sort expression.
func (f *feedQueryBuilder) WithSorting(column, direction string) *feedQueryBuilder {
	switch {
	case strings.EqualFold(direction, "ASC"):
		f.sortExpressions = append(f.sortExpressions, pq.QuoteIdentifier(column)+" ASC")
	case strings.EqualFold(direction, "DESC"):
		f.sortExpressions = append(f.sortExpressions, pq.QuoteIdentifier(column)+" DESC")
	}

	return f
}

// WithLimit set the limit.
func (f *feedQueryBuilder) WithLimit(limit int) *feedQueryBuilder {
	f.limit = limit
	return f
}

// WithOffset set the offset.
func (f *feedQueryBuilder) WithOffset(offset int) *feedQueryBuilder {
	f.offset = offset
	return f
}

func (f *feedQueryBuilder) buildCondition() string {
	return strings.Join(f.conditions, " AND ")
}

func (f *feedQueryBuilder) buildCounterCondition() string {
	return strings.Join(f.counterConditions, " AND ")
}

func (f *feedQueryBuilder) buildSorting() string {
	var parts string

	if len(f.sortExpressions) > 0 {
		parts += " ORDER BY " + strings.Join(f.sortExpressions, ", ")
	}

	if len(parts) > 0 {
		parts += ", lower(f.title) ASC"
	}

	if f.limit > 0 {
		parts += " LIMIT " + strconv.Itoa(f.limit)
	}

	if f.offset > 0 {
		parts += " OFFSET " + strconv.Itoa(f.offset)
	}

	return parts
}

// GetFeed returns a single feed that match the condition.
func (f *feedQueryBuilder) GetFeed() (*model.Feed, error) {
	f.limit = 1
	feeds, err := f.GetFeeds()
	if err != nil {
		return nil, err
	}

	if len(feeds) != 1 {
		return nil, nil
	}

	return feeds[0], nil
}

// GetFeeds returns a list of feeds that match the condition.
func (f *feedQueryBuilder) GetFeeds() (model.Feeds, error) {
	query := `
		SELECT
			f.id,
			f.feed_url,
			f.site_url,
			f.title,
			f.user_id,
			f.checked_at at time zone u.timezone,
			f.next_check_at at time zone u.timezone,
			f.parsing_error_count,
			f.parsing_error_msg,
			f.disabled,
			u.timezone
		FROM
			feeds f
		LEFT JOIN
			users u ON u.id=f.user_id
		WHERE %s
		%s
	`

	query = fmt.Sprintf(query, f.buildCondition(), f.buildSorting())

	readCounters, unreadCounters, err := f.fetchFeedCounter()
	if err != nil {
		return nil, err
	}

	rows, err := f.db.Query(query, f.args...)
	if err != nil {
		return nil, fmt.Errorf(`store: unable to fetch feeds: %w`, err)
	}
	defer rows.Close()

	feeds := make(model.Feeds, 0)
	for rows.Next() {
		var feed model.Feed
		var tz string

		err := rows.Scan(
			&feed.ID,
			&feed.FeedURL,
			&feed.SiteURL,
			&feed.Title,
			&feed.UserID,
			&feed.CheckedAt,
			&feed.NextCheckAt,
			&feed.ParsingErrorCount,
			&feed.ParsingErrorMsg,
			&feed.Disabled,
			&tz,
		)
		if err != nil {
			return nil, fmt.Errorf(`store: unable to fetch feeds row: %w`, err)
		}

		if readCounters != nil {
			if count, found := readCounters[feed.ID]; found {
				feed.ReadCount = count
			}
		}
		if unreadCounters != nil {
			if count, found := unreadCounters[feed.ID]; found {
				feed.UnreadCount = count
			}
		}

		feed.NumberOfVisibleEntries = feed.ReadCount + feed.UnreadCount
		feed.CheckedAt = timezone.Convert(tz, feed.CheckedAt)
		feed.NextCheckAt = timezone.Convert(tz, feed.NextCheckAt)
		feeds = append(feeds, &feed)
	}

	return feeds, nil
}

func (f *feedQueryBuilder) fetchFeedCounter() (unreadCounters map[int64]int, readCounters map[int64]int, err error) {
	if !f.withCounters {
		return nil, nil, nil
	}
	query := `
		SELECT
			e.feed_id,
			e.status,
			count(*)
		FROM
			entries e
		WHERE
			%s
		GROUP BY
			e.feed_id, e.status
	`
	query = fmt.Sprintf(query, f.buildCounterCondition())

	rows, err := f.db.Query(query, f.counterArgs...)
	if err != nil {
		return nil, nil, fmt.Errorf(`store: unable to fetch feed counts: %w`, err)
	}
	defer rows.Close()

	readCounters = make(map[int64]int)
	unreadCounters = make(map[int64]int)
	for rows.Next() {
		var feedID int64
		var status string
		var count int
		if err := rows.Scan(&feedID, &status, &count); err != nil {
			return nil, nil, fmt.Errorf(`store: unable to fetch feed counter row: %w`, err)
		}

		switch status {
		case model.EntryStatusRead:
			readCounters[feedID] = count
		case model.EntryStatusUnread:
			unreadCounters[feedID] = count
		}
	}

	return readCounters, unreadCounters, nil
}
