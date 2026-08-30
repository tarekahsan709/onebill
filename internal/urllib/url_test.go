// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package urllib // import "onebill.app/internal/urllib"

import (
	"testing"
)

func TestIsRelativePath(t *testing.T) {
	scenarios := map[string]bool{
		// Valid relative paths
		"path/to/file.ext":    true,
		"./path/to/file.ext":  true,
		"../path/to/file.ext": true,
		"file.ext":            true,
		"./file.ext":          true,
		"../file.ext":         true,
		"/absolute/path":      true,
		"path?query=value":    true,
		"path#fragment":       true,
		"path?query#fragment": true,
		"/bookmarklet?uri=https://example.org/rss": true,

		// Not relative paths
		"https://example.org/file.ext": false,
		"http://example.org/file.ext":  false,
		"//example.org/file.ext":       false,
		"//example.org":                false,
		" //example.org":               false,
		"  //example.org/path":         false,
		"\t//example.org":              false,
		"\n//example.org":              false,
		" /path":                       false,
		"/path ":                       false,
		`/\example.org`:                false,
		` /\example.org`:               false,
		`\example.org`:                 false,
		`path\to\file.ext`:             false,
		"ftp://example.org/file.ext":   false,
		"mailto:user@example.org":      false,
		"magnet:?xt=urn:btih:example":  false,
		"":                             false,
		"magnet:?xt.1=urn:sha1:YNCKHTQCWBTRNJIV4WNAE52SJUQCZO5C": false,
	}

	for input, expected := range scenarios {
		actual := IsRelativePath(input)
		if actual != expected {
			t.Errorf(`Unexpected result for IsRelativePath, got %v instead of %v for %q`, actual, expected, input)
		}
	}
}

func TestIsAbsoluteURL(t *testing.T) {
	scenarios := map[string]bool{
		"https://example.org/file.pdf":                   true,
		"https://example.org/file.pdf?download=1#page=2": true,
		"mailto:user@example.org":                        false,
		"data:text/plain,hello":                          false,
		"magnet:?xt.1=urn:sha1:YNCKHTQCWBTRNJIV4WNAE52SJUQCZO5C&xt.2=urn:sha1:TXGCZQTH26NL6OUQAJJPFALHG2LTGBC7": false,
		"invalid url":                false,
		"/relative/path":             false,
		"//example.org/path":         false,
		" https://example.org/path":  false,
		"\thttps://example.org/path": false,
	}

	for input, expected := range scenarios {
		actual := IsAbsoluteURL(input)
		if actual != expected {
			t.Errorf(`Unexpected result, got %v instead of %v for %q`, actual, expected, input)
		}
	}
}

func TestDomain(t *testing.T) {
	scenarios := map[string]string{
		"https://static.example.org/": "static.example.org",
		"https://example|org/":        "https://example|org/",
	}

	for input, expected := range scenarios {
		actual := Domain(input)
		if actual != expected {
			t.Errorf(`Unexpected result, got %q instead of %q`, actual, expected)
		}
	}
}
