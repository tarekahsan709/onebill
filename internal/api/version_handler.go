// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package api // import "onebill.app/internal/api"

import (
	"net/http"
	"runtime"

	"onebill.app/internal/http/response"
	"onebill.app/internal/version"
)

func (h *handler) versionHandler(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, r, &versionResponse{
		Version:   version.Version,
		Commit:    version.Commit,
		BuildDate: version.BuildDate,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Arch:      runtime.GOARCH,
		OS:        runtime.GOOS,
	})
}
