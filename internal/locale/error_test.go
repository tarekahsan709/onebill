// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package locale // import "onebill.app/internal/locale"

import (
	"testing"
)

func TestNewLocalizedError(t *testing.T) {
	translationKey := "error.validation"
	args := []any{"field1", "invalid"}

	localizedErr := NewLocalizedError(translationKey, args...)

	if localizedErr.translationKey != translationKey {
		t.Errorf("Expected translation key to be %q, got %q", translationKey, localizedErr.translationKey)
	}

	if len(localizedErr.translationArgs) != 2 {
		t.Errorf("Expected 2 translation args, got %d", len(localizedErr.translationArgs))
	}

	if localizedErr.translationArgs[0] != "field1" || localizedErr.translationArgs[1] != "invalid" {
		t.Errorf("Expected translation args [field1, invalid], got %v", localizedErr.translationArgs)
	}
}

func TestLocalizedError_String(t *testing.T) {
	defaultCatalog = catalog{
		"en_US": translationDict{
			singulars: map[string]string{
				"error.validation": "Validation failed for %s: %s",
			},
		},
	}

	localizedErr := NewLocalizedError("error.validation", "username", "too short")

	result := localizedErr.String()
	expected := "Validation failed for username: too short"
	if result != expected {
		t.Errorf("Expected String() result %q, got %q", expected, result)
	}
}

func TestLocalizedError_StringWithMissingTranslation(t *testing.T) {
	defaultCatalog = catalog{
		"en_US": translationDict{},
	}

	localizedErr := NewLocalizedError("error.missing", "arg1")

	result := localizedErr.String()
	expected := "error.missing"
	if result != expected {
		t.Errorf("Expected String() result %q, got %q", expected, result)
	}
}

func TestLocalizedError_Error(t *testing.T) {
	defaultCatalog = catalog{
		"en_US": translationDict{
			singulars: map[string]string{
				"error.database": "Database connection failed: %s",
			},
		},
	}

	localizedErr := NewLocalizedError("error.database", "timeout")

	result := localizedErr.Error()
	if result == nil {
		t.Error("Expected Error() to return a non-nil error")
	}

	expected := "Database connection failed: timeout"
	if result.Error() != expected {
		t.Errorf("Expected Error() message %q, got %q", expected, result.Error())
	}
}

func TestLocalizedError_Translate(t *testing.T) {
	defaultCatalog = catalog{
		"en_US": translationDict{
			singulars: map[string]string{
				"error.permission": "Permission denied for %s",
			},
		},
		"es_ES": translationDict{
			singulars: map[string]string{
				"error.permission": "Permiso denegado para %s",
			},
		},
	}

	localizedErr := NewLocalizedError("error.permission", "admin panel")

	// Test English translation
	result := localizedErr.Translate("en_US")
	expected := "Permission denied for admin panel"
	if result != expected {
		t.Errorf("Expected English translation %q, got %q", expected, result)
	}

	// Test Spanish translation
	result = localizedErr.Translate("es_ES")
	expected = "Permiso denegado para admin panel"
	if result != expected {
		t.Errorf("Expected Spanish translation %q, got %q", expected, result)
	}

	// Test with missing language
	result = localizedErr.Translate("invalid_lang")
	expected = "error.permission"
	if result != expected {
		t.Errorf("Expected fallback translation %q, got %q", expected, result)
	}
}

func TestLocalizedError_TranslateWithNoArgs(t *testing.T) {
	defaultCatalog = catalog{
		"en_US": translationDict{
			singulars: map[string]string{
				"error.generic": "An error occurred",
			},
		},
		"de_DE": translationDict{
			singulars: map[string]string{
				"error.generic": "Ein Fehler ist aufgetreten",
			},
		},
	}

	localizedErr := NewLocalizedError("error.generic")

	// Test English
	result := localizedErr.Translate("en_US")
	expected := "An error occurred"
	if result != expected {
		t.Errorf("Expected English translation %q, got %q", expected, result)
	}

	// Test German
	result = localizedErr.Translate("de_DE")
	expected = "Ein Fehler ist aufgetreten"
	if result != expected {
		t.Errorf("Expected German translation %q, got %q", expected, result)
	}
}

func TestLocalizedError_TranslateWithComplexArgs(t *testing.T) {
	defaultCatalog = catalog{
		"en_US": translationDict{
			singulars: map[string]string{
				"error.complex": "Error %d: %s occurred at %s with severity %s",
			},
		},
	}

	localizedErr := NewLocalizedError("error.complex", 500, "Internal Server Error", "2024-01-01", "high")

	result := localizedErr.Translate("en_US")
	expected := "Error 500: Internal Server Error occurred at 2024-01-01 with severity high"
	if result != expected {
		t.Errorf("Expected complex translation %q, got %q", expected, result)
	}
}

func TestLocalizedError_EmptyKey(t *testing.T) {
	localizedErr := NewLocalizedError("")

	result := localizedErr.String()
	expected := ""
	if result != expected {
		t.Errorf("Expected empty string for empty key, got %q", result)
	}

	result = localizedErr.Translate("en_US")
	if result != expected {
		t.Errorf("Expected empty string for empty key translation, got %q", result)
	}
}
