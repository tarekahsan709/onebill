// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package ui // import "onebill.app/internal/ui"

import (
	"context"
	"net/http"
	"time"

	"onebill.app/internal/config"
	"onebill.app/internal/http/request"
	"onebill.app/internal/model"
	"onebill.app/internal/oauth2"
	"onebill.app/internal/storage"
)

const sessionCookieName = "OneBillSessionID"

// authenticateWebSession binds the current browser session to the given user,
// rotates its identifier and secret, and refreshes the client cookie.
func authenticateWebSession(w http.ResponseWriter, r *http.Request, store *storage.Storage, user *model.User) error {
	session := request.WebSession(r)
	session.SetUser(user)

	oldID, secret := session.Rotate()
	if err := store.RotateWebSession(oldID, session); err != nil {
		return err
	}

	setSessionCookie(w, session, secret)
	return nil
}

// setSessionCookie writes the session cookie to the response with the
// security attributes used by onebill (HttpOnly, SameSite=Lax, Secure
// when HTTPS).
func setSessionCookie(w http.ResponseWriter, session *model.WebSession, secret string) {
	path := config.Opts.BasePath()
	if path == "" {
		path = "/"
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    session.ID + "." + secret,
		Path:     path,
		Secure:   config.Opts.HTTPS(),
		HttpOnly: true,
		Expires:  time.Now().Add(config.Opts.CleanupRemoveSessionsInterval()),
		SameSite: http.SameSiteLaxMode,
	})
}

func getOAuth2Manager(ctx context.Context) *oauth2.Manager {
	return oauth2.NewManager(
		ctx,
		config.Opts.OAuth2Provider(),
		config.Opts.OAuth2ClientID(),
		config.Opts.OAuth2ClientSecret(),
		config.Opts.OAuth2RedirectURL(),
		config.Opts.OAuth2OIDCDiscoveryEndpoint(),
	)
}
