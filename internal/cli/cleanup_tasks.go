// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package cli // import "onebill.app/internal/cli"

import (
	"log/slog"
	"time"

	"onebill.app/internal/config"
	"onebill.app/internal/metric"
	"onebill.app/internal/model"
	"onebill.app/internal/storage"
)

func runCleanupTasks(store *storage.Storage) {
	if nbWebSessions, err := store.CleanOldWebSessions(config.Opts.CleanupRemoveSessionsInterval()); err != nil {
		slog.Error("Unable to clean old web sessions", slog.Any("error", err))
	} else {
		slog.Info("Sessions cleanup completed",
			slog.Int64("web_sessions_removed", nbWebSessions),
		)
	}

	startTime := time.Now()
	if rowsAffected, err := store.ArchiveEntries(model.EntryStatusRead, config.Opts.CleanupArchiveReadInterval(), config.Opts.CleanupArchiveBatchSize()); err != nil {
		slog.Error("Unable to archive read entries", slog.Any("error", err))
	} else {
		slog.Info("Archiving read entries completed",
			slog.Int64("read_entries_archived", rowsAffected),
		)

		if config.Opts.HasMetricsCollector() {
			metric.ArchiveEntriesDuration.WithLabelValues(model.EntryStatusRead).Observe(time.Since(startTime).Seconds())
		}
	}

	startTime = time.Now()
	if rowsAffected, err := store.ArchiveEntries(model.EntryStatusUnread, config.Opts.CleanupArchiveUnreadInterval(), config.Opts.CleanupArchiveBatchSize()); err != nil {
		slog.Error("Unable to archive unread entries", slog.Any("error", err))
	} else {
		slog.Info("Archiving unread entries completed",
			slog.Int64("unread_entries_archived", rowsAffected),
		)

		if config.Opts.HasMetricsCollector() {
			metric.ArchiveEntriesDuration.WithLabelValues(model.EntryStatusUnread).Observe(time.Since(startTime).Seconds())
		}
	}

}

func cleanupScheduler(store *storage.Storage, frequency time.Duration) {
	for range time.Tick(frequency) {
		runCleanupTasks(store)
	}
}
