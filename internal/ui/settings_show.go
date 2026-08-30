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
	"onebill.app/internal/timezone"
	"onebill.app/internal/ui/form"
	"onebill.app/internal/ui/view"
)

func (h *handler) showSettingsPage(w http.ResponseWriter, r *http.Request) {
	user, err := h.store.UserByID(request.UserID(r))
	if err != nil {
		response.HTMLServerError(w, r, err)
		return
	}

	settingsForm := form.SettingsForm{
		Username:          user.Username,
		Theme:             user.Theme,
		Language:          user.Language,
		Timezone:          user.Timezone,
		EntryDirection:    user.EntryDirection,
		EntryOrder:        user.EntryOrder,
		EntriesPerPage:    user.EntriesPerPage,
		CustomCSS:         user.Stylesheet,
		CustomJS:          user.CustomJS,
		ExternalFontHosts: user.ExternalFontHosts,
		DefaultHomePage:   user.DefaultHomePage,
	}

	creds, err := h.store.WebAuthnCredentialsByUserID(user.ID)
	if err != nil {
		response.HTMLServerError(w, r, err)
		return
	}

	view := view.New(h.tpl, r)
	view.Set("form", settingsForm)
	view.Set("themes", model.Themes())
	view.Set("languages", locale.AvailableLanguages)
	view.Set("timezones", timezone.AvailableTimezones())
	view.Set("menu", "settings")
	view.Set("user", user)
	navMetadata, _ := h.store.GetNavMetadata(user.ID)
	view.Set("countUnread", navMetadata.CountUnread)
	view.Set("countErrorFeeds", navMetadata.CountErrorFeeds)
	view.Set("default_home_pages", model.HomePages())
	view.Set("maxEntriesPerPage", model.MaxEntryLimit)
	view.Set("countWebAuthnCerts", h.store.CountWebAuthnCredentialsByUserID(user.ID))
	view.Set("webAuthnCerts", creds)

	response.HTML(w, r, view.Render("settings"))
}
