// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package model // import "onebill.app/internal/model"

// HomePages returns the list of available home pages.
func HomePages() map[string]string {
	return map[string]string{
		"unread":  "menu.unread",
		"history": "menu.history",
		"feeds":   "menu.feeds",
	}
}
