// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package ui // import "onebill.app/internal/ui"

import (
	"fmt"

	"onebill.app/internal/storage"
	"onebill.app/internal/template"
	"onebill.app/internal/ui/static"
)

type handler struct {
	basePath string
	store    *storage.Storage
	tpl      *template.Engine
}

func (h *handler) routePath(format string, args ...any) string {
	if len(args) > 0 {
		return h.basePath + fmt.Sprintf(format, args...)
	}
	return h.basePath + format
}

func (h *handler) iconPath(filename string) string {
	return h.basePath + fmt.Sprintf("/icon/%s/%s", static.BinaryBundles[filename].Checksum, filename)
}
