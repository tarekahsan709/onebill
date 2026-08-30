// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package validator // import "onebill.app/internal/validator"

import (
	"strings"
	"unicode"

	"onebill.app/internal/locale"
	"onebill.app/internal/model"
	"onebill.app/internal/storage"
	"onebill.app/internal/timezone"
)

// ValidateUserCreationWithPassword validates user creation with a password.
func ValidateUserCreationWithPassword(store *storage.Storage, request *model.UserCreationRequest) *locale.LocalizedError {
	if request.Username == "" {
		return locale.NewLocalizedError("error.user_mandatory_fields")
	}

	if store.UserExists(request.Username) {
		return locale.NewLocalizedError("error.user_already_exists")
	}

	if err := validateUsername(request.Username); err != nil {
		return err
	}

	if err := validatePassword(request.Password); err != nil {
		return err
	}

	return nil
}

// ValidateUserModification validates user modifications.
func ValidateUserModification(store *storage.Storage, userID int64, changes *model.UserModificationRequest) *locale.LocalizedError {
	if changes.Username != nil {
		if *changes.Username == "" {
			return locale.NewLocalizedError("error.user_mandatory_fields")
		} else if store.AnotherUserExists(userID, *changes.Username) {
			return locale.NewLocalizedError("error.user_already_exists")
		}
	}

	if changes.Password != nil {
		if err := validatePassword(*changes.Password); err != nil {
			return err
		}
	}

	if changes.Theme != nil {
		if err := validateTheme(*changes.Theme); err != nil {
			return err
		}
	}

	if changes.Language != nil {
		if err := validateLanguage(*changes.Language); err != nil {
			return err
		}
	}

	if changes.Timezone != nil {
		if err := validateTimezone(*changes.Timezone); err != nil {
			return err
		}
	}

	if changes.EntryDirection != nil {
		if err := ValidateDirection(*changes.EntryDirection); err != nil {
			return locale.NewLocalizedError("error.invalid_entry_direction")
		}
	}

	if changes.EntryOrder != nil {
		if err := validateEntrySortingOrder(*changes.EntryOrder); err != nil {
			return err
		}
	}

	if changes.EntriesPerPage != nil {
		if err := validateEntriesPerPage(*changes.EntriesPerPage); err != nil {
			return err
		}
	}

	if changes.DefaultHomePage != nil {
		if err := validateDefaultHomePage(*changes.DefaultHomePage); err != nil {
			return err
		}
	}

	if changes.ExternalFontHosts != nil {
		if !IsValidDomainList(*changes.ExternalFontHosts) {
			return locale.NewLocalizedError("error.settings_invalid_domain_list")
		}
	}

	return nil
}

func validatePassword(password string) *locale.LocalizedError {
	if len(password) < 6 {
		return locale.NewLocalizedError("error.password_min_length")
	}
	return nil
}

// validateUsername return an error if the `username` argument contains
// a character that isn't alphanumerical nor `_` and `-`.
//
// Note: this validation should not be applied to previously created usernames,
// and cannot be applied to Google/OIDC accounts creation because the email
// address is used for the username field.
func validateUsername(username string) *locale.LocalizedError {
	if strings.ContainsFunc(username, func(r rune) bool {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return false
		}
		if r == '_' || r == '-' || r == '@' || r == '.' {
			return false
		}
		return true
	}) {
		return locale.NewLocalizedError("error.invalid_username")
	}
	return nil
}

func validateTheme(theme string) *locale.LocalizedError {
	themes := model.Themes()
	if _, found := themes[theme]; !found {
		return locale.NewLocalizedError("error.invalid_theme")
	}
	return nil
}

func validateLanguage(language string) *locale.LocalizedError {
	languages := locale.AvailableLanguages
	if _, found := languages[language]; !found {
		return locale.NewLocalizedError("error.invalid_language")
	}
	return nil
}

func validateTimezone(timezoneValue string) *locale.LocalizedError {
	if !timezone.IsValid(timezoneValue) {
		return locale.NewLocalizedError("error.invalid_timezone")
	}
	return nil
}

func validateEntriesPerPage(entriesPerPage int) *locale.LocalizedError {
	if entriesPerPage < 1 || entriesPerPage > model.MaxEntryLimit {
		return locale.NewLocalizedError("error.entries_per_page_invalid")
	}
	return nil
}

// validateEntrySortingOrder must accept only the values of the
// entry_sorting_order enum type in the database.
func validateEntrySortingOrder(order string) *locale.LocalizedError {
	if order != "published_at" && order != "created_at" {
		return locale.NewLocalizedError("error.invalid_entry_order")
	}
	return nil
}

func validateDefaultHomePage(defaultHomePage string) *locale.LocalizedError {
	defaultHomePages := model.HomePages()
	if _, found := defaultHomePages[defaultHomePage]; !found {
		return locale.NewLocalizedError("error.invalid_default_home_page")
	}
	return nil
}
