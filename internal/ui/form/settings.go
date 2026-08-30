// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package form // import "onebill.app/internal/ui/form"

import (
	"net/http"
	"strconv"

	"onebill.app/internal/config"
	"onebill.app/internal/locale"
	"onebill.app/internal/model"
	"onebill.app/internal/validator"
)

// SettingsForm represents the settings form.
type SettingsForm struct {
	Username          string
	Password          string
	Confirmation      string
	Theme             string
	Language          string
	Timezone          string
	EntryDirection    string
	EntryOrder        string
	EntriesPerPage    int
	CustomCSS         string
	CustomJS          string
	ExternalFontHosts string
	DefaultHomePage   string
}

// Merge updates the fields of the given user.
func (s *SettingsForm) Merge(user *model.User) *model.User {
	if !config.Opts.DisableLocalAuth() {
		user.Username = s.Username
	}
	user.Theme = s.Theme
	user.Language = s.Language
	user.Timezone = s.Timezone
	user.EntryDirection = s.EntryDirection
	user.EntryOrder = s.EntryOrder
	user.EntriesPerPage = s.EntriesPerPage
	user.Stylesheet = s.CustomCSS
	user.CustomJS = s.CustomJS
	user.ExternalFontHosts = s.ExternalFontHosts
	user.DefaultHomePage = s.DefaultHomePage

	if s.Password != "" {
		user.Password = s.Password
	}

	return user
}

// Validate makes sure the form values are valid.
func (s *SettingsForm) Validate() *locale.LocalizedError {
	if (s.Username == "" && !config.Opts.DisableLocalAuth()) || s.Theme == "" || s.Language == "" || s.Timezone == "" || s.EntryDirection == "" || s.DefaultHomePage == "" {
		return locale.NewLocalizedError("error.settings_mandatory_fields")
	}

	if s.Confirmation == "" {
		// Firefox insists on auto-completing the password field.
		// If the confirmation field is blank, the user probably
		// didn't intend to change their password.
		s.Password = ""
	} else if s.Password != "" {
		if s.Password != s.Confirmation {
			return locale.NewLocalizedError("error.different_passwords")
		}
	}

	if s.ExternalFontHosts != "" {
		if !validator.IsValidDomainList(s.ExternalFontHosts) {
			return locale.NewLocalizedError("error.settings_invalid_domain_list")
		}
	}

	return nil
}

// NewSettingsForm returns a new SettingsForm.
func NewSettingsForm(r *http.Request) *SettingsForm {
	entriesPerPage, err := strconv.Atoi(r.FormValue("entries_per_page"))
	if err != nil {
		entriesPerPage = 0
	}
	return &SettingsForm{
		Username:          r.FormValue("username"),
		Password:          r.FormValue("password"),
		Confirmation:      r.FormValue("confirmation"),
		Theme:             r.FormValue("theme"),
		Language:          r.FormValue("language"),
		Timezone:          r.FormValue("timezone"),
		EntryDirection:    r.FormValue("entry_direction"),
		EntryOrder:        r.FormValue("entry_order"),
		EntriesPerPage:    entriesPerPage,
		CustomCSS:         r.FormValue("custom_css"),
		CustomJS:          r.FormValue("custom_js"),
		ExternalFontHosts: r.FormValue("external_font_hosts"),
		DefaultHomePage:   r.FormValue("default_home_page"),
	}
}
