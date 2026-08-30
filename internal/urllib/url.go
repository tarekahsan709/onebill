// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package urllib // import "onebill.app/internal/urllib"

import (
	"net/netip"
	"net/url"
	"strings"
)

var rfc6598SharedAddressSpacePrefix = netip.MustParsePrefix("100.64.0.0/10")

// IsRelativePath reports whether the link is a relative path (no scheme, host, or scheme-relative // form).
func IsRelativePath(link string) bool {
	if link == "" {
		return false
	}

	// Reject surrounding whitespace: browsers strip leading and trailing C0
	// control characters and spaces before parsing a target, while Go's url.Parse
	// treats spaces as path characters. Without this check, a target like
	// " //evil.com" is accepted here and becomes //evil.com in the browser.
	if link != strings.TrimSpace(link) {
		return false
	}

	// Reject backslashes: Go's url.Parse treats them as ordinary path
	// characters, but browsers normalize them to forward slashes, so a target
	// like "/\evil.com" would parse as relative here yet redirect to
	// //evil.com in the browser (open redirect).
	if strings.Contains(link, "\\") {
		return false
	}

	if parsedURL, err := url.Parse(link); err == nil {
		// Only allow relative paths (not scheme-relative URLs like //example.org)
		// and ensure the URL doesn't have a host component
		if !parsedURL.IsAbs() && parsedURL.Host == "" && parsedURL.Scheme == "" {
			return true
		}
	}
	return false
}

// hasHTTPPrefix reports whether the URL string begins with an HTTP or HTTPS scheme.
func hasHTTPPrefix(inputURL string) bool {
	return strings.HasPrefix(inputURL, "https://") || strings.HasPrefix(inputURL, "http://")
}

// IsAbsoluteURL reports whether the link is absolute and starts with an HTTP or HTTPS scheme.
func IsAbsoluteURL(inputURL string) bool {
	if !hasHTTPPrefix(inputURL) {
		return false
	}
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return false
	}
	return parsedURL.IsAbs()
}

// Domain returns the host component of the given URL.
func Domain(websiteURL string) string {
	parsedURL, err := url.Parse(websiteURL)
	if err != nil {
		return websiteURL
	}

	return parsedURL.Host
}
