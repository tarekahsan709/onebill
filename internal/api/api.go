// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package api // import "onebill.app/internal/api"

import (
	"net/http"

	"onebill.app/internal/storage"
)

type handler struct {
	store *storage.Storage
}

// NewHandler returns an http.Handler that handles API v1 calls.
// The returned handler expects the base path to be stripped from the request URL.
func NewHandler(store *storage.Storage) http.Handler {
	handler := &handler{store: store}
	middleware := newMiddleware(store)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/users", handler.createUserHandler)
	mux.HandleFunc("GET /v1/users", handler.usersHandler)
	mux.HandleFunc("GET /v1/users/{identifier}", handler.dispatchUserLookupHandler)
	mux.HandleFunc("PUT /v1/users/{userID}", handler.updateUserHandler)
	mux.HandleFunc("DELETE /v1/users/{userID}", handler.removeUserHandler)
	mux.HandleFunc("PUT /v1/users/{userID}/mark-all-as-read", handler.markUserAsReadHandler)
	mux.HandleFunc("GET /v1/me", handler.currentUserHandler)
	mux.HandleFunc("POST /v1/feeds", handler.createFeedHandler)
	mux.HandleFunc("GET /v1/feeds", handler.getFeedsHandler)
	mux.HandleFunc("GET /v1/feeds/counters", handler.fetchCountersHandler)
	mux.HandleFunc("GET /v1/feeds/{feedID}", handler.getFeedHandler)
	mux.HandleFunc("PUT /v1/feeds/{feedID}", handler.updateFeedHandler)
	mux.HandleFunc("DELETE /v1/feeds/{feedID}", handler.removeFeedHandler)
	mux.HandleFunc("PUT /v1/feeds/{feedID}/mark-all-as-read", handler.markFeedAsReadHandler)
	mux.HandleFunc("GET /v1/feeds/{feedID}/entries", handler.getFeedEntriesHandler)
	mux.HandleFunc("POST /v1/feeds/{feedID}/entries/import", handler.importFeedEntryHandler)
	mux.HandleFunc("GET /v1/feeds/{feedID}/entries/{entryID}", handler.getFeedEntryHandler)
	mux.HandleFunc("GET /v1/entries/ids", handler.getEntryIDsHandler)
	mux.HandleFunc("GET /v1/entries", handler.getEntriesHandler)
	mux.HandleFunc("PUT /v1/entries", handler.setEntryStatusHandler)
	mux.HandleFunc("GET /v1/entries/{entryID}", handler.getEntryHandler)
	mux.HandleFunc("PUT /v1/flush-history", handler.flushHistoryHandler)
	mux.HandleFunc("DELETE /v1/flush-history", handler.flushHistoryHandler)
	mux.HandleFunc("GET /v1/version", handler.versionHandler)
	mux.HandleFunc("POST /v1/api-keys", handler.createAPIKeyHandler)
	mux.HandleFunc("GET /v1/api-keys", handler.getAPIKeysHandler)
	mux.HandleFunc("DELETE /v1/api-keys/{apiKeyID}", handler.deleteAPIKeyHandler)

	return middleware.withCORSHeaders(middleware.validateAPIKeyAuth(middleware.validateBasicAuth(mux)))
}
