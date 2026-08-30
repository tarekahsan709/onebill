// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package ui // import "onebill.app/internal/ui"

import (
	"net/http"

	"onebill.app/internal/http/request"
	"onebill.app/internal/http/response"
	"onebill.app/internal/locale"
	"onebill.app/internal/model"
)

func (h *handler) showWebManifest(w http.ResponseWriter, r *http.Request) {
	type webManifestIcon struct {
		Source  string `json:"src"`
		Sizes   string `json:"sizes,omitempty"`
		Type    string `json:"type,omitempty"`
		Purpose string `json:"purpose,omitempty"`
	}

	type webManifestShortcut struct {
		Name  string            `json:"name"`
		URL   string            `json:"url"`
		Icons []webManifestIcon `json:"icons,omitempty"`
	}

	type webManifest struct {
		Name            string                `json:"name"`
		Description     string                `json:"description"`
		ShortName       string                `json:"short_name"`
		StartURL        string                `json:"start_url"`
		Icons           []webManifestIcon     `json:"icons"`
		Display         string                `json:"display"`
		BackgroundColor string                `json:"background_color"`
		Shortcuts       []webManifestShortcut `json:"shortcuts"`
	}

	labelUnreadMenu := "Unread"
	labelHistoryMenu := "History"
	labelFeedsMenu := "Feeds"
	labelSettingsMenu := "Settings"
	if request.IsAuthenticated(r) {
		user, err := h.store.UserByID(request.UserID(r))
		if err != nil {
			response.JSONServerError(w, r, err)
			return
		}
		printer := locale.NewPrinter(user.Language)
		labelUnreadMenu = printer.Print("menu.unread")
		labelHistoryMenu = printer.Print("menu.history")
		labelFeedsMenu = printer.Print("menu.feeds")
		labelSettingsMenu = printer.Print("menu.settings")
	}
	themeColor := model.ThemeColor(request.WebSession(r).Theme(), "light")
	manifest := &webManifest{
		Name:            "OneBill",
		ShortName:       "OneBill",
		Description:     "Cloud and AI spend dashboard",
		Display:         "standalone",
		StartURL:        h.routePath("/"),
		BackgroundColor: themeColor,
		Icons: []webManifestIcon{
			{Source: h.iconPath("icon-120.png"), Sizes: "120x120", Type: "image/png", Purpose: "any"},
			{Source: h.iconPath("icon-192.png"), Sizes: "192x192", Type: "image/png", Purpose: "any"},
			{Source: h.iconPath("icon-512.png"), Sizes: "512x512", Type: "image/png", Purpose: "any"},
			{Source: h.iconPath("maskable-icon-120.png"), Sizes: "120x120", Type: "image/png", Purpose: "maskable"},
			{Source: h.iconPath("maskable-icon-192.png"), Sizes: "192x192", Type: "image/png", Purpose: "maskable"},
			{Source: h.iconPath("maskable-icon-512.png"), Sizes: "512x512", Type: "image/png", Purpose: "maskable"},
		},
		Shortcuts: []webManifestShortcut{
			{Name: labelUnreadMenu, URL: h.routePath("/unread"), Icons: []webManifestIcon{{Source: h.iconPath("unread-icon.png"), Sizes: "240x240", Type: "image/png"}}},
			{Name: labelHistoryMenu, URL: h.routePath("/history"), Icons: []webManifestIcon{{Source: h.iconPath("history-icon.png"), Sizes: "240x240", Type: "image/png"}}},
			{Name: labelFeedsMenu, URL: h.routePath("/feeds"), Icons: []webManifestIcon{{Source: h.iconPath("feeds-icon.png"), Sizes: "240x240", Type: "image/png"}}},
			{Name: labelSettingsMenu, URL: h.routePath("/settings"), Icons: []webManifestIcon{{Source: h.iconPath("settings-icon.png"), Sizes: "240x240", Type: "image/png"}}},
		},
	}

	response.JSON(w, r, manifest)
}
