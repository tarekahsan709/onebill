// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package form // import "onebill.app/internal/ui/form"

import (
	"net/http"
	"strings"

	"onebill.app/internal/locale"
)

// authForm represents the authentication form.
type authForm struct {
	Username string
	Password string
}

// Validate makes sure the form values are valid.
func (a authForm) Validate() *locale.LocalizedError {
	if a.Username == "" || a.Password == "" {
		return locale.NewLocalizedError("error.fields_mandatory")
	}

	return nil
}

// NewAuthForm returns a new AuthForm.
func NewAuthForm(r *http.Request) *authForm {
	return &authForm{
		Username: strings.TrimSpace(r.FormValue("username")),
		Password: strings.TrimSpace(r.FormValue("password")),
	}
}
