// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package ui // import "onebill.app/internal/ui"

import (
	"net/http"

	"onebill.app/internal/config"
	"onebill.app/internal/storage"
	"onebill.app/internal/template"
)

// Serve returns an http.Handler that serves the user interface.
// The returned handler expects the base path to be stripped from the request URL.
func Serve(store *storage.Storage) http.Handler {
	basePath := config.Opts.BasePath()
	webSessionMiddleware := newWebSessionMiddleware(basePath, store)
	csrfMiddleware := newCSRFMiddleware(basePath)
	authProxyMiddleware := newAuthProxyMiddleware(basePath, store)

	templateEngine := template.NewEngine(basePath)
	templateEngine.ParseTemplates()

	handler := &handler{basePath, store, templateEngine}

	mux := http.NewServeMux()

	// Static assets.
	mux.HandleFunc("GET /stylesheets/{checksum}/{filename}", handler.showStylesheet)
	mux.HandleFunc("GET /js/{checksum}/{filename}", handler.showJavascript)
	mux.HandleFunc("GET /icon/{checksum}/{filename}", handler.showAppIcon)
	mux.HandleFunc("GET /favicon.ico", handler.showFavicon)
	mux.HandleFunc("GET /manifest.json", handler.showWebManifest)

	// Unread page.
	mux.HandleFunc("POST /mark-all-as-read", handler.markAllAsRead)
	mux.HandleFunc("GET /unread", handler.showUnreadPage)

	// History pages.
	mux.HandleFunc("GET /history", handler.showHistoryPage)
	mux.HandleFunc("POST /history/flush", handler.flushHistory)

	// Feed listing pages.
	mux.HandleFunc("GET /feeds", handler.showFeedsPage)

	// Individual feed pages.
	mux.HandleFunc("GET /feed/{feedID}/edit", handler.showEditFeedPage)
	mux.HandleFunc("POST /feed/{feedID}/remove", handler.removeFeed)
	mux.HandleFunc("POST /feed/{feedID}/update", handler.updateFeed)
	mux.HandleFunc("GET /feed/{feedID}/entries", handler.showFeedEntriesPage)
	mux.HandleFunc("GET /feed/{feedID}/entries/all", handler.showFeedEntriesAllPage)
	mux.HandleFunc("POST /feed/{feedID}/mark-all-as-read", handler.markFeedAsRead)

	// Entry pages.
	mux.HandleFunc("POST /entry/status", handler.updateEntriesStatus)

	// User pages.
	mux.HandleFunc("GET /users", handler.showUsersPage)
	mux.HandleFunc("GET /user/create", handler.showCreateUserPage)
	mux.HandleFunc("POST /user/save", handler.saveUser)
	mux.HandleFunc("GET /users/{userID}/edit", handler.showEditUserPage)
	mux.HandleFunc("POST /users/{userID}/update", handler.updateUser)
	mux.HandleFunc("POST /users/{userID}/remove", handler.removeUser)

	// Settings pages.
	mux.HandleFunc("GET /settings", handler.showSettingsPage)
	mux.HandleFunc("POST /settings", handler.updateSettings)
	mux.HandleFunc("GET /about", handler.showAboutPage)

	// Session pages.
	mux.HandleFunc("GET /sessions", handler.showSessionsPage)
	mux.HandleFunc("POST /sessions/{sessionID}/remove", handler.removeSession)

	// API Keys pages.
	if config.Opts.HasAPI() {
		mux.HandleFunc("GET /keys", handler.showAPIKeysPage)
		mux.HandleFunc("POST /keys/{keyID}/delete", handler.deleteAPIKey)
		mux.HandleFunc("GET /keys/create", handler.showCreateAPIKeyPage)
		mux.HandleFunc("POST /keys/save", handler.saveAPIKey)
	}

	// OAuth2 flow.
	if config.Opts.OAuth2Provider() != "" {
		mux.HandleFunc("POST /oauth2/{provider}/unlink", handler.oauth2Unlink)
		mux.HandleFunc("GET /oauth2/{provider}/redirect", handler.oauth2Redirect)
		mux.HandleFunc("GET /oauth2/{provider}/callback", handler.oauth2Callback)
	}

	// Offline page.
	mux.HandleFunc("GET /offline", handler.showOfflinePage)

	// Authentication pages.
	mux.HandleFunc("POST /login", handler.checkLogin)
	mux.HandleFunc("POST /logout", handler.logout)
	mux.Handle("GET /{$}", authProxyMiddleware.handle(http.HandlerFunc(handler.showLoginPage)))

	// WebAuthn flow.
	if config.Opts.WebAuthn() {
		mux.HandleFunc("GET /webauthn/register/begin", handler.beginRegistration)
		mux.HandleFunc("POST /webauthn/register/finish", handler.finishRegistration)
		mux.HandleFunc("GET /webauthn/login/begin", handler.beginLogin)
		mux.HandleFunc("POST /webauthn/login/finish", handler.finishLogin)
		mux.HandleFunc("POST /webauthn/deleteall", handler.deleteAllCredentials)
		mux.HandleFunc("POST /webauthn/{credentialHandle}/delete", handler.deleteCredential)
		mux.HandleFunc("GET /webauthn/{credentialHandle}/rename", handler.renameCredential)
		mux.HandleFunc("POST /webauthn/{credentialHandle}/save", handler.saveCredential)
	}

	// robots.txt
	mux.HandleFunc("GET /robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("User-agent: *\nDisallow: /"))
	})

	// Apply middleware chain: web session -> CSRF validation -> handlers.
	return webSessionMiddleware.handle(csrfMiddleware.handle(mux))
}
