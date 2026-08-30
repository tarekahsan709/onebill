// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package storage // import "onebill.app/internal/storage"

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"

	"onebill.app/internal/model"
	"onebill.app/internal/timezone"
)

// EntryQueryBuilder builds a SQL query to fetch entries.
type EntryQueryBuilder struct {
	store           *Storage
	args            []any
	conditions      []string
	sortExpressions []string
	limit           int
	offset          int
}

// BeforeChangedDate adds a condition < changed_at
func (e *EntryQueryBuilder) BeforeChangedDate(date time.Time) *EntryQueryBuilder {
	e.conditions = append(e.conditions, "e.changed_at < $"+strconv.Itoa(len(e.args)+1))
	e.args = append(e.args, date)
	return e
}

// AfterChangedDate adds a condition > changed_at
func (e *EntryQueryBuilder) AfterChangedDate(date time.Time) *EntryQueryBuilder {
	e.conditions = append(e.conditions, "e.changed_at > $"+strconv.Itoa(len(e.args)+1))
	e.args = append(e.args, date)
	return e
}

// BeforePublishedDate adds a condition < published_at
func (e *EntryQueryBuilder) BeforePublishedDate(date time.Time) *EntryQueryBuilder {
	e.conditions = append(e.conditions, "e.published_at < $"+strconv.Itoa(len(e.args)+1))
	e.args = append(e.args, date)
	return e
}

// AfterPublishedDate adds a condition > published_at
func (e *EntryQueryBuilder) AfterPublishedDate(date time.Time) *EntryQueryBuilder {
	e.conditions = append(e.conditions, "e.published_at > $"+strconv.Itoa(len(e.args)+1))
	e.args = append(e.args, date)
	return e
}

// BeforeEntryID adds a condition < entryID.
func (e *EntryQueryBuilder) BeforeEntryID(entryID int64) *EntryQueryBuilder {
	if entryID != 0 {
		e.conditions = append(e.conditions, "e.id < $"+strconv.Itoa(len(e.args)+1))
		e.args = append(e.args, entryID)
	}
	return e
}

// AfterEntryID adds a condition > entryID.
func (e *EntryQueryBuilder) AfterEntryID(entryID int64) *EntryQueryBuilder {
	if entryID != 0 {
		e.conditions = append(e.conditions, "e.id > $"+strconv.Itoa(len(e.args)+1))
		e.args = append(e.args, entryID)
	}
	return e
}

// WithEntryIDs filter by entry IDs.
func (e *EntryQueryBuilder) WithEntryIDs(entryIDs ...int64) *EntryQueryBuilder {
	if len(entryIDs) == 1 {
		e.conditions = append(e.conditions, fmt.Sprintf("e.id = $%d", len(e.args)+1))
		e.args = append(e.args, entryIDs[0])
	} else if len(entryIDs) > 1 {
		e.conditions = append(e.conditions, fmt.Sprintf("e.id = ANY($%d)", len(e.args)+1))
		e.args = append(e.args, pq.Int64Array(entryIDs))
	}
	return e
}

// WithFeedID filter by feed ID.
func (e *EntryQueryBuilder) WithFeedID(feedID int64) *EntryQueryBuilder {
	if feedID > 0 {
		e.conditions = append(e.conditions, "e.feed_id = $"+strconv.Itoa(len(e.args)+1))
		e.args = append(e.args, feedID)
	}
	return e
}

// WithStatuses filter by a list of entry statuses.
func (e *EntryQueryBuilder) WithStatuses(statuses ...string) *EntryQueryBuilder {
	if len(statuses) == 1 {
		e.conditions = append(e.conditions, fmt.Sprintf("e.status = $%d", len(e.args)+1))
		e.args = append(e.args, statuses[0])
	} else if len(statuses) > 1 {
		e.conditions = append(e.conditions, fmt.Sprintf("e.status = ANY($%d)", len(e.args)+1))
		e.args = append(e.args, pq.StringArray(statuses))
	}
	return e
}

// WithoutStatus set the entry status that should not be returned.
func (e *EntryQueryBuilder) WithoutStatus(status string) *EntryQueryBuilder {
	if status != "" {
		e.conditions = append(e.conditions, "e.status <> $"+strconv.Itoa(len(e.args)+1))
		e.args = append(e.args, status)
	}
	return e
}

// WithSorting add a sort expression.
func (e *EntryQueryBuilder) WithSorting(column, direction string) *EntryQueryBuilder {
	switch {
	case strings.EqualFold(direction, "ASC"):
		e.sortExpressions = append(e.sortExpressions, pq.QuoteIdentifier(column)+" ASC")
	case strings.EqualFold(direction, "DESC"):
		e.sortExpressions = append(e.sortExpressions, pq.QuoteIdentifier(column)+" DESC")
	}

	return e
}

// WithLimit sets the limit. A non-positive limit is clamped to
// model.MaxEntryLimit so callers cannot request an unbounded result set.
func (e *EntryQueryBuilder) WithLimit(limit int) *EntryQueryBuilder {
	return e.WithLimitAndMaximum(limit, model.MaxEntryLimit)
}

// WithLimitAndMaximum sets the limit, capped at the given maximum.
// A non-positive limit is clamped to the maximum.
func (e *EntryQueryBuilder) WithLimitAndMaximum(limit, maximum int) *EntryQueryBuilder {
	if limit <= 0 || limit > maximum {
		limit = maximum
	}
	e.limit = limit
	return e
}

// WithOffset set the offset.
func (e *EntryQueryBuilder) WithOffset(offset int) *EntryQueryBuilder {
	if offset > 0 {
		e.offset = offset
	}
	return e
}

// CountEntries count the number of entries that match the condition.
func (e *EntryQueryBuilder) CountEntries() (count int, err error) {
	query := `
		SELECT count(*)
		FROM entries e
			JOIN feeds f ON f.id = e.feed_id
		WHERE ` + e.buildCondition()

	err = e.store.db.QueryRow(query, e.args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("store: unable to count entries: %v", err)
	}

	return count, nil
}

// GetEntry returns a single entry that match the condition.
func (e *EntryQueryBuilder) GetEntry() (*model.Entry, error) {
	e.limit = 1
	entries, err := e.GetEntries()
	if err != nil {
		return nil, err
	}

	if len(entries) != 1 {
		return nil, nil
	}

	return entries[0], nil
}

// GetEntries returns a list of entries that match the condition.
func (e *EntryQueryBuilder) GetEntries() (model.Entries, error) {
	entries, _, err := e.fetchEntries(false)
	return entries, err
}

// GetEntriesWithCount returns a list of entries and the total count of matching
// rows, ignoring limit and offset. It uses a window function for non-empty pages
// and falls back to a separate count when the requested offset returns no rows.
func (e *EntryQueryBuilder) GetEntriesWithCount() (model.Entries, int, error) {
	entries, total, err := e.fetchEntries(true)
	if err != nil {
		return nil, 0, err
	}

	if len(entries) == 0 && e.offset > 0 {
		total, err = e.CountEntries()
		if err != nil {
			return nil, 0, err
		}
	}

	return entries, total, nil
}

// fetchEntries is the shared implementation for GetEntries and GetEntriesWithCount.
// When withCount is true, count(*) OVER() is included in the SELECT and the total
// count of matching rows is returned; otherwise the returned count is 0.
func (e *EntryQueryBuilder) fetchEntries(withCount bool) (model.Entries, int, error) {
	countColumn := ""
	if withCount {
		countColumn = "count(*) OVER(),"
	}

	query := `
		SELECT
			` + countColumn + `
			e.id,
			e.user_id,
			e.feed_id,
			e.hash,
			e.published_at at time zone u.timezone,
			e.title,
			e.url,
			e.author,
			e.status,
			e.created_at,
			e.changed_at,
			f.title as feed_title,
			f.feed_url,
			f.site_url,
			f.checked_at,
			u.timezone
		FROM
			entries e
		INNER JOIN
			feeds f ON f.id=e.feed_id
		INNER JOIN
			users u ON u.id=e.user_id
		WHERE ` + e.buildCondition() + " " + e.buildSorting()

	rows, err := e.store.db.Query(query, e.args...)
	if err != nil {
		return nil, 0, fmt.Errorf("store: unable to get entries: %v", err)
	}
	defer rows.Close()

	size := max(e.limit, 0)
	entries := make(model.Entries, 0, size)
	var totalCount int

	for rows.Next() {
		var tz string

		entry := model.NewEntry()

		dest := []any{
			&entry.ID,
			&entry.UserID,
			&entry.FeedID,
			&entry.Hash,
			&entry.Date,
			&entry.Title,
			&entry.URL,
			&entry.Author,
			&entry.Status,
			&entry.CreatedAt,
			&entry.ChangedAt,
			&entry.Feed.Title,
			&entry.Feed.FeedURL,
			&entry.Feed.SiteURL,
			&entry.Feed.CheckedAt,
			&tz,
		}

		if withCount {
			dest = append([]any{&totalCount}, dest...)
		}

		err := rows.Scan(dest...)
		if err != nil {
			return nil, 0, fmt.Errorf("store: unable to fetch entry row: %v", err)
		}

		// Make sure that timestamp fields contain timezone information (API)
		entry.Date = timezone.Convert(tz, entry.Date)
		entry.CreatedAt = timezone.Convert(tz, entry.CreatedAt)
		entry.ChangedAt = timezone.Convert(tz, entry.ChangedAt)
		entry.Feed.CheckedAt = timezone.Convert(tz, entry.Feed.CheckedAt)

		entry.Feed.ID = entry.FeedID
		entry.Feed.UserID = entry.UserID

		entries = append(entries, entry)
	}

	return entries, totalCount, nil
}

// GetEntryIDs returns a list of entry IDs that match the condition.
func (e *EntryQueryBuilder) GetEntryIDs() ([]int64, error) {
	query := `
		SELECT
			e.id
		FROM
			entries e
		LEFT JOIN
			feeds f
		ON
			f.id=e.feed_id
		WHERE ` + e.buildCondition() + " " + e.buildSorting()

	rows, err := e.store.db.Query(query, e.args...)
	if err != nil {
		return nil, fmt.Errorf("store: unable to get entries: %v", err)
	}
	defer rows.Close()

	var entryIDs []int64
	for rows.Next() {
		var entryID int64
		if err := rows.Scan(&entryID); err != nil {
			return nil, fmt.Errorf("store: unable to fetch entry row: %v", err)
		}
		entryIDs = append(entryIDs, entryID)
	}

	return entryIDs, nil
}

// GetEntryIDsWithCount returns a list of entry IDs and the total count of
// matching rows (ignoring limit/offset). It uses two queries: one to count
// all matching rows and one to fetch the paginated IDs.
func (e *EntryQueryBuilder) GetEntryIDsWithCount() ([]int64, int, error) {
	total, err := e.CountEntries()
	if err != nil {
		return nil, 0, err
	}

	entryIDs, err := e.GetEntryIDs()
	if err != nil {
		return nil, 0, err
	}

	return entryIDs, total, nil
}

func (e *EntryQueryBuilder) buildCondition() string {
	return strings.Join(e.conditions, " AND ")
}

func (e *EntryQueryBuilder) buildSorting() string {
	var parts string

	if len(e.sortExpressions) > 0 {
		parts += " ORDER BY " + strings.Join(e.sortExpressions, ", ")
	}

	if e.limit > 0 {
		parts += " LIMIT " + strconv.Itoa(e.limit)
	}

	if e.offset > 0 {
		parts += " OFFSET " + strconv.Itoa(e.offset)
	}

	return parts
}

// NewEntryQueryBuilder returns a new EntryQueryBuilder.
func (s *Storage) NewEntryQueryBuilder(userID int64) *EntryQueryBuilder {
	return &EntryQueryBuilder{
		store:      s,
		args:       []any{userID},
		conditions: []string{"e.user_id = $1"},
	}
}
