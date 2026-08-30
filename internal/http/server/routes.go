// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package server // import "onebill.app/internal/http/server"

import (
	"net/http"

	"onebill.app/internal/api"
	"onebill.app/internal/config"
	"onebill.app/internal/storage"
	"onebill.app/internal/ui"
)

func newRouter(store *storage.Storage) http.Handler {
	readinessProbe := newReadinessProbe(store)

	// Application routes served under the base path.
	appMux := http.NewServeMux()

	appMux.HandleFunc("GET /healthcheck", readinessProbe)

	// REST API routing.
	if config.Opts.HasAPI() {
		appMux.Handle("/v1/", api.NewHandler(store))
	}

	// Metrics endpoint.
	if config.Opts.HasMetricsCollector() {
		appMux.Handle("GET /metrics", metricsHandler())
	}

	// UI routing (catch-all).
	appMux.Handle("/", ui.Serve(store))

	// Apply shared middleware.
	var appHandler http.Handler = appMux
	if config.Opts.HasMaintenanceMode() {
		appHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(config.Opts.MaintenanceMessage()))
		})
	}
	appHandler = middleware(appHandler)

	// Root router: health probes at root, app routes under base path.
	rootMux := http.NewServeMux()

	// These routes do not take the base path into consideration and are always available at the root of the server.
	rootMux.HandleFunc("/liveness", livenessProbe)
	rootMux.HandleFunc("/healthz", livenessProbe)
	rootMux.HandleFunc("/readiness", readinessProbe)
	rootMux.HandleFunc("/readyz", readinessProbe)

	basePath := config.Opts.BasePath()
	if basePath != "" {
		rootMux.Handle(basePath+"/", http.StripPrefix(basePath, appHandler))
	} else {
		rootMux.Handle("/", appHandler)
	}

	return rootMux
}
