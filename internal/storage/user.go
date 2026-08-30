// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package storage // import "onebill.app/internal/storage"

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"onebill.app/internal/crypto"
	"onebill.app/internal/model"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var dummyBcryptHash = []byte("$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy")

// CountUsers returns the total number of users.
func (s *Storage) CountUsers() (int, error) {
	var result int
	err := s.db.QueryRow(`SELECT count(*) FROM users`).Scan(&result)
	if err != nil {
		return 0, fmt.Errorf("storage: unable to count users: %w", err)
	}

	return result, nil
}

// SetLastLogin sets the user's last login timestamp to the current time.
func (s *Storage) SetLastLogin(userID int64) error {
	query := `UPDATE users SET last_login_at=now() WHERE id=$1`
	_, err := s.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf(`store: unable to update last login date: %v`, err)
	}

	return nil
}

// UserExists returns true if a user with the given username exists.
func (s *Storage) UserExists(username string) bool {
	var result bool
	s.db.QueryRow(`SELECT true FROM users WHERE username=LOWER($1) LIMIT 1`, username).Scan(&result)
	return result
}

// AnotherUserExists returns true if a user other than userID has the given username.
func (s *Storage) AnotherUserExists(userID int64, username string) bool {
	var result bool
	s.db.QueryRow(`SELECT true FROM users WHERE id != $1 AND username=LOWER($2) LIMIT 1`, userID, username).Scan(&result)
	return result
}

// CreateUser creates a new user.
func (s *Storage) CreateUser(userCreationRequest *model.UserCreationRequest) (*model.User, error) {
	var hashedPassword string
	if userCreationRequest.Password != "" {
		var err error
		hashedPassword, err = crypto.HashPassword(userCreationRequest.Password)
		if err != nil {
			return nil, err
		}
	}

	query := `
		INSERT INTO users
			(username, password, is_admin, google_id, openid_connect_id)
		VALUES
			(LOWER($1), $2, $3, $4, $5)
		RETURNING
			id,
			username,
			is_admin,
			language,
			theme,
			timezone,
			entry_direction,
			entries_per_page,
			stylesheet,
			custom_js,
			external_font_hosts,
			google_id,
			openid_connect_id,
			entry_order,
			default_home_page
	`

	var user model.User
	err := s.db.QueryRow(
		query,
		userCreationRequest.Username,
		hashedPassword,
		userCreationRequest.IsAdmin,
		userCreationRequest.GoogleID,
		userCreationRequest.OpenIDConnectID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.IsAdmin,
		&user.Language,
		&user.Theme,
		&user.Timezone,
		&user.EntryDirection,
		&user.EntriesPerPage,
		&user.Stylesheet,
		&user.CustomJS,
		&user.ExternalFontHosts,
		&user.GoogleID,
		&user.OpenIDConnectID,
		&user.EntryOrder,
		&user.DefaultHomePage,
	)
	if err != nil {
		return nil, fmt.Errorf(`store: unable to create user: %v`, err)
	}

	return &user, nil
}

// UpdateUser updates a user.
func (s *Storage) UpdateUser(user *model.User) error {
	user.ExternalFontHosts = strings.TrimSpace(user.ExternalFontHosts)

	if user.Password != "" {
		hashedPassword, err := crypto.HashPassword(user.Password)
		if err != nil {
			return err
		}

		query := `
			UPDATE users SET
				username=LOWER($1),
				password=$2,
				is_admin=$3,
				theme=$4,
				language=$5,
				timezone=$6,
				entry_direction=$7,
				entries_per_page=$8,
				stylesheet=$9,
				custom_js=$10,
				external_font_hosts=$11,
				google_id=$12,
				openid_connect_id=$13,
				entry_order=$14,
				default_home_page=$15
			WHERE
				id=$16
		`

		_, err = s.db.Exec(
			query,
			user.Username,
			hashedPassword,
			user.IsAdmin,
			user.Theme,
			user.Language,
			user.Timezone,
			user.EntryDirection,
			user.EntriesPerPage,
			user.Stylesheet,
			user.CustomJS,
			user.ExternalFontHosts,
			user.GoogleID,
			user.OpenIDConnectID,
			user.EntryOrder,
			user.DefaultHomePage,
			user.ID,
		)
		if err != nil {
			return fmt.Errorf(`store: unable to update user: %v`, err)
		}
	} else {
		query := `
			UPDATE users SET
				username=LOWER($1),
				is_admin=$2,
				theme=$3,
				language=$4,
				timezone=$5,
				entry_direction=$6,
				entries_per_page=$7,
				stylesheet=$8,
				custom_js=$9,
				external_font_hosts=$10,
				google_id=$11,
				openid_connect_id=$12,
				entry_order=$13,
				default_home_page=$14
			WHERE
				id=$15
		`

		_, err := s.db.Exec(
			query,
			user.Username,
			user.IsAdmin,
			user.Theme,
			user.Language,
			user.Timezone,
			user.EntryDirection,
			user.EntriesPerPage,
			user.Stylesheet,
			user.CustomJS,
			user.ExternalFontHosts,
			user.GoogleID,
			user.OpenIDConnectID,
			user.EntryOrder,
			user.DefaultHomePage,
			user.ID,
		)

		if err != nil {
			return fmt.Errorf(`store: unable to update user: %v`, err)
		}
	}

	return nil
}

// UserByID finds a user by the ID.
func (s *Storage) UserByID(userID int64) (*model.User, error) {
	query := `
		SELECT
			id,
			username,
			is_admin,
			theme,
			language,
			timezone,
			entry_direction,
			entries_per_page,
			last_login_at,
			stylesheet,
			custom_js,
			external_font_hosts,
			google_id,
			openid_connect_id,
			entry_order,
			default_home_page
		FROM
			users
		WHERE
			id = $1
	`
	return s.fetchUser(query, userID)
}

// UserByUsername finds a user by the username.
func (s *Storage) UserByUsername(username string) (*model.User, error) {
	query := `
		SELECT
			id,
			username,
			is_admin,
			theme,
			language,
			timezone,
			entry_direction,
			entries_per_page,
			last_login_at,
			stylesheet,
			custom_js,
			external_font_hosts,
			google_id,
			openid_connect_id,
			entry_order,
			default_home_page
		FROM
			users
		WHERE
			username=LOWER($1)
	`
	return s.fetchUser(query, username)
}

// UserByField returns the user matching the given column name and value.
func (s *Storage) UserByField(field, value string) (*model.User, error) {
	query := `
		SELECT
			id,
			username,
			is_admin,
			theme,
			language,
			timezone,
			entry_direction,
			entries_per_page,
			last_login_at,
			stylesheet,
			custom_js,
			external_font_hosts,
			google_id,
			openid_connect_id,
			entry_order,
			default_home_page
		FROM
			users
		WHERE
			%s=$1
	`
	return s.fetchUser(fmt.Sprintf(query, pq.QuoteIdentifier(field)), value)
}

// AnotherUserWithFieldExists returns true if a user other than userID has the given value in the given column.
func (s *Storage) AnotherUserWithFieldExists(userID int64, field, value string) bool {
	var result bool
	query := `SELECT true FROM users WHERE id <> $1 AND ` + pq.QuoteIdentifier(field) + `=$2 LIMIT 1`
	s.db.QueryRow(query, userID, value).Scan(&result)
	return result
}

// UserByAPIKey returns the user associated with the given API key.
func (s *Storage) UserByAPIKey(token string) (*model.User, error) {
	query := `
		SELECT
			u.id,
			u.username,
			u.is_admin,
			u.theme,
			u.language,
			u.timezone,
			u.entry_direction,
			u.entries_per_page,
			u.last_login_at,
			u.stylesheet,
			u.custom_js,
			u.external_font_hosts,
			u.google_id,
			u.openid_connect_id,
			u.entry_order,
			u.default_home_page
		FROM
			users u
		INNER JOIN
			api_keys ON api_keys.user_id=u.id
		WHERE
			api_keys.token = $1
	`
	return s.fetchUser(query, token)
}

func (s *Storage) fetchUser(query string, args ...any) (*model.User, error) {
	var user model.User
	err := s.db.QueryRow(query, args...).Scan(
		&user.ID,
		&user.Username,
		&user.IsAdmin,
		&user.Theme,
		&user.Language,
		&user.Timezone,
		&user.EntryDirection,
		&user.EntriesPerPage,
		&user.LastLoginAt,
		&user.Stylesheet,
		&user.CustomJS,
		&user.ExternalFontHosts,
		&user.GoogleID,
		&user.OpenIDConnectID,
		&user.EntryOrder,
		&user.DefaultHomePage,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf(`store: unable to fetch user: %v`, err)
	}

	return &user, nil
}

// RemoveUser deletes a user and all related data.
func (s *Storage) RemoveUser(userID int64) error {
	if _, err := s.db.Exec(`DELETE FROM users WHERE id=$1`, userID); err != nil {
		return fmt.Errorf(`store: unable to remove user #%d: %v`, userID, err)
	}
	return nil
}

// Users returns all users.
func (s *Storage) Users() (model.Users, error) {
	query := `
		SELECT
			id,
			username,
			is_admin,
			theme,
			language,
			timezone,
			entry_direction,
			entries_per_page,
			last_login_at,
			stylesheet,
			custom_js,
			external_font_hosts,
			google_id,
			openid_connect_id,
			entry_order,
			default_home_page
		FROM
			users
		ORDER BY username ASC
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf(`store: unable to fetch users: %v`, err)
	}
	defer rows.Close()

	var users model.Users
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.IsAdmin,
			&user.Theme,
			&user.Language,
			&user.Timezone,
			&user.EntryDirection,
			&user.EntriesPerPage,
			&user.LastLoginAt,
			&user.Stylesheet,
			&user.CustomJS,
			&user.ExternalFontHosts,
			&user.GoogleID,
			&user.OpenIDConnectID,
			&user.EntryOrder,
			&user.DefaultHomePage,
		)

		if err != nil {
			return nil, fmt.Errorf(`store: unable to fetch users row: %v`, err)
		}

		users = append(users, &user)
	}

	return users, nil
}

// CheckPassword returns nil if the given password matches the user's stored hash.
func (s *Storage) CheckPassword(username, password string) error {
	var hash string
	username = strings.ToLower(username)

	err := s.db.QueryRow("SELECT password FROM users WHERE username=$1", username).Scan(&hash)
	if errors.Is(err, sql.ErrNoRows) {
		// Perform a dummy bcrypt comparison against the hashed `password` string to avoid leaking whether the user exists via response timing.
		_ = bcrypt.CompareHashAndPassword(dummyBcryptHash, []byte(password))
		return fmt.Errorf(`store: unable to find this user: %s`, username)
	} else if err != nil {
		return fmt.Errorf(`store: unable to fetch user: %v`, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return fmt.Errorf(`store: invalid password for "%s" (%v)`, username, err)
	}

	return nil
}

// HasPassword returns true if the given user exists and has a non-empty password.
func (s *Storage) HasPassword(userID int64) (bool, error) {
	var result bool
	query := `SELECT true FROM users WHERE id=$1 AND password <> '' LIMIT 1`

	err := s.db.QueryRow(query, userID).Scan(&result)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf(`store: unable to execute query: %v`, err)
	}

	if result {
		return true, nil
	}
	return false, nil
}
