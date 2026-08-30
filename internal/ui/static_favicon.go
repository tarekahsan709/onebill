// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package ui // import "onebill.app/internal/ui"

import (
	"net/http"
	"time"

	"onebill.app/internal/http/response"

	"onebill.app/internal/ui/static"
)

func (h *handler) showFavicon(w http.ResponseWriter, r *http.Request) {
	value, ok := static.BinaryBundles["favicon.ico"]
	if !ok {
		response.HTMLNotFound(w, r)
		return
	}

	response.NewBuilder(w, r).WithCaching(value.Checksum, 48*time.Hour, func(b *response.Builder) {
		b.WithHeader("Content-Type", "image/x-icon")
		b.WithoutCompression()
		b.WithBodyAsBytes(value.Data)
		b.Write()
	})
}
