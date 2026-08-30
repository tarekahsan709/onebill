// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package model // import "onebill.app/internal/model"

import (
	"time"

	"onebill.app/internal/timezone"
)

// User represents a user in the system.
type User struct {
	ID                int64      `json:"id"`
	Username          string     `json:"username"`
	Password          string     `json:"-"`
	Theme             string     `json:"theme"`
	Language          string     `json:"language"`
	Timezone          string     `json:"timezone"`
	EntryDirection    string     `json:"entry_sorting_direction"`
	EntryOrder        string     `json:"entry_sorting_order"`
	Stylesheet        string     `json:"stylesheet"`
	CustomJS          string     `json:"custom_js"`
	ExternalFontHosts string     `json:"external_font_hosts"`
	GoogleID          string     `json:"google_id"`
	OpenIDConnectID   string     `json:"openid_connect_id"`
	EntriesPerPage    int        `json:"entries_per_page"`
	LastLoginAt       *time.Time `json:"last_login_at"`
	DefaultHomePage   string     `json:"default_home_page"`
	IsAdmin           bool       `json:"is_admin"`
}

// UserCreationRequest represents the request to create a user.
type UserCreationRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	IsAdmin         bool   `json:"is_admin"`
	GoogleID        string `json:"google_id"`
	OpenIDConnectID string `json:"openid_connect_id"`
}

// UserModificationRequest represents the request to update a user.
type UserModificationRequest struct {
	Username          *string `json:"username"`
	Password          *string `json:"password"`
	Theme             *string `json:"theme"`
	Language          *string `json:"language"`
	Timezone          *string `json:"timezone"`
	EntryDirection    *string `json:"entry_sorting_direction"`
	EntryOrder        *string `json:"entry_sorting_order"`
	Stylesheet        *string `json:"stylesheet"`
	CustomJS          *string `json:"custom_js"`
	ExternalFontHosts *string `json:"external_font_hosts"`
	EntriesPerPage    *int    `json:"entries_per_page"`
	IsAdmin           *bool   `json:"is_admin"`
	DefaultHomePage   *string `json:"default_home_page"`
}

// Patch updates the User object with the modification request.
func (u *UserModificationRequest) Patch(user *User) {
	if u.Username != nil {
		user.Username = *u.Username
	}

	if u.Password != nil {
		user.Password = *u.Password
	}

	if u.IsAdmin != nil {
		user.IsAdmin = *u.IsAdmin
	}

	if u.Theme != nil {
		user.Theme = *u.Theme
	}

	if u.Language != nil {
		user.Language = *u.Language
	}

	if u.Timezone != nil {
		user.Timezone = *u.Timezone
	}

	if u.EntryDirection != nil {
		user.EntryDirection = *u.EntryDirection
	}

	if u.EntryOrder != nil {
		user.EntryOrder = *u.EntryOrder
	}

	if u.Stylesheet != nil {
		user.Stylesheet = *u.Stylesheet
	}

	if u.CustomJS != nil {
		user.CustomJS = *u.CustomJS
	}

	if u.ExternalFontHosts != nil {
		user.ExternalFontHosts = *u.ExternalFontHosts
	}

	if u.EntriesPerPage != nil {
		user.EntriesPerPage = *u.EntriesPerPage
	}

	if u.DefaultHomePage != nil {
		user.DefaultHomePage = *u.DefaultHomePage
	}

}

// UseTimezone converts last login date to the given timezone.
func (u *User) UseTimezone(tz string) {
	if u.LastLoginAt != nil {
		*u.LastLoginAt = timezone.Convert(tz, *u.LastLoginAt)
	}
}

// Users represents a list of users.
type Users []*User

// UseTimezone converts last login timestamp of all users to the given timezone.
func (u Users) UseTimezone(tz string) {
	for _, user := range u {
		user.UseTimezone(tz)
	}
}
