// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package config // import "onebill.app/internal/config"

import (
	"slices"
	"testing"
)

func TestBaseURLOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.BaseURL() != "http://localhost" {
		t.Fatalf("Expected BASE_URL to be 'http://localhost' by default")
	}

	if configParser.options.RootURL() != "http://localhost" {
		t.Fatalf("Expected ROOT_URL to be 'http://localhost' by default")
	}

	if configParser.options.BasePath() != "" {
		t.Fatalf("Expected BASE_PATH to be empty by default")
	}

	if err := configParser.parseLines([]string{"BASE_URL=https://example.com/app"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.BaseURL() != "https://example.com/app" {
		t.Fatalf("Expected BASE_URL to be 'https://example.com/app', got '%s'", configParser.options.BaseURL())
	}

	if configParser.options.RootURL() != "https://example.com" {
		t.Fatalf("Expected ROOT_URL to be 'https://example.com', got '%s'", configParser.options.RootURL())
	}

	if configParser.options.BasePath() != "/app" {
		t.Fatalf("Expected BASE_PATH to be '/app', got '%s'", configParser.options.BasePath())
	}

	if err := configParser.parseLines([]string{"BASE_URL=https://example.com/app/"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.BaseURL() != "https://example.com/app" {
		t.Fatalf("Expected BASE_URL to be 'https://example.com/app', got '%s'", configParser.options.BaseURL())
	}

	if configParser.options.RootURL() != "https://example.com" {
		t.Fatalf("Expected ROOT_URL to be 'https://example.com', got '%s'", configParser.options.RootURL())
	}

	if configParser.options.BasePath() != "/app" {
		t.Fatalf("Expected BASE_PATH to be '/app', got '%s'", configParser.options.BasePath())
	}

	if err := configParser.parseLines([]string{"BASE_URL=example.com/app/"}); err == nil {
		t.Fatal("Expected an error due to missing scheme in BASE_URL")
	}
}

func TestWatchdogOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if !configParser.options.Watchdog() {
		t.Fatal("Expected WATCHDOG to be enabled by default")
	}

	if !configParser.options.HasSchedulerService() {
		t.Fatal("Expected HAS_SCHEDULER_SERVICE to be enabled by default")
	}

	if err := configParser.parseLines([]string{"WATCHDOG=1"}); err != nil {
		t.Fatal("Unexpected error:", err)
	}

	if !configParser.options.Watchdog() {
		t.Fatal("Expected WATCHDOG to be enabled")
	}

	if !configParser.options.HasSchedulerService() {
		t.Fatal("Expected HAS_SCHEDULER_SERVICE to be enabled")
	}

	if err := configParser.parseLines([]string{"WATCHDOG=0"}); err != nil {
		t.Fatal("Unexpected error:", err)
	}

	if configParser.options.Watchdog() {
		t.Fatal("Expected WATCHDOG to be disabled")
	}

	if configParser.options.HasWatchdog() {
		t.Fatal("Expected HAS_WATCHDOG to be disabled")
	}
}

func TestWebAuthnOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.WebAuthn() {
		t.Fatalf("Expected WEBAUTHN to be disabled by default")
	}

	if err := configParser.parseLines([]string{"WEBAUTHN=1"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !configParser.options.WebAuthn() {
		t.Fatalf("Expected WEBAUTHN to be enabled")
	}
}

func TestAdminPasswordOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.AdminPassword() != "" {
		t.Fatalf("Expected ADMIN_PASSWORD to be empty by default")
	}

	if err := configParser.parseLines([]string{"ADMIN_PASSWORD=secret123"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.AdminPassword() != "secret123" {
		t.Fatalf("Expected ADMIN_PASSWORD to be 'secret123'")
	}
}

func TestAdminUsernameOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.AdminUsername() != "" {
		t.Fatalf("Expected ADMIN_USERNAME to be empty by default")
	}

	if err := configParser.parseLines([]string{"ADMIN_USERNAME=admin"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.AdminUsername() != "admin" {
		t.Fatalf("Expected ADMIN_USERNAME to be 'admin'")
	}
}

func TestAuthProxyHeaderOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.AuthProxyHeader() != "" {
		t.Fatalf("Expected AUTH_PROXY_HEADER to be empty by default")
	}

	if err := configParser.parseLines([]string{"AUTH_PROXY_HEADER=X-Forwarded-User"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.AuthProxyHeader() != "X-Forwarded-User" {
		t.Fatalf("Expected AUTH_PROXY_HEADER to be 'X-Forwarded-User'")
	}
}

func TestAuthProxyUserCreationOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.AuthProxyUserCreation() {
		t.Fatal("Expected AUTH_PROXY_USER_CREATION to be disabled by default")
	}

	if configParser.options.IsAuthProxyUserCreationAllowed() {
		t.Fatal("Expected HAS_AUTH_PROXY_USER_CREATION to be disabled by default")
	}

	if err := configParser.parseLines([]string{"AUTH_PROXY_USER_CREATION=1"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !configParser.options.AuthProxyUserCreation() {
		t.Fatal("Expected AUTH_PROXY_USER_CREATION to be enabled")
	}

	if !configParser.options.IsAuthProxyUserCreationAllowed() {
		t.Fatal("Expected HAS_AUTH_PROXY_USER_CREATION to be enabled")
	}

	if err := configParser.parseLines([]string{"AUTH_PROXY_USER_CREATION=0"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.AuthProxyUserCreation() {
		t.Fatal("Expected AUTH_PROXY_USER_CREATION to be disabled")
	}

	if configParser.options.IsAuthProxyUserCreationAllowed() {
		t.Fatal("Expected HAS_AUTH_PROXY_USER_CREATION to be disabled")
	}
}

func TestCertDomainOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.CertDomain() != "" {
		t.Fatalf("Expected CERT_DOMAIN to be empty by default")
	}

	if err := configParser.parseLines([]string{"CERT_DOMAIN=example.com"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.CertDomain() != "example.com" {
		t.Fatalf("Expected CERT_DOMAIN to be 'example.com'")
	}
}

func TestCertFileOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.CertFile() != "" {
		t.Fatalf("Expected CERT_FILE to be empty by default")
	}

	if err := configParser.parseLines([]string{"CERT_FILE=/path/to/cert.pem"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.CertFile() != "/path/to/cert.pem" {
		t.Fatalf("Expected CERT_FILE to be '/path/to/cert.pem'")
	}
}

func TestCleanupArchiveBatchSizeOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.CleanupArchiveBatchSize() != 10000 {
		t.Fatalf("Expected CLEANUP_ARCHIVE_BATCH_SIZE to be 10000 by default")
	}

	if err := configParser.parseLines([]string{"CLEANUP_ARCHIVE_BATCH_SIZE=5000"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.CleanupArchiveBatchSize() != 5000 {
		t.Fatalf("Expected CLEANUP_ARCHIVE_BATCH_SIZE to be 5000")
	}
}

func TestCreateAdminOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.CreateAdmin() {
		t.Fatalf("Expected CREATE_ADMIN to be disabled by default")
	}

	if err := configParser.parseLines([]string{"CREATE_ADMIN=1"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !configParser.options.CreateAdmin() {
		t.Fatalf("Expected CREATE_ADMIN to be enabled")
	}

	if err := configParser.parseLines([]string{"CREATE_ADMIN=0"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.CreateAdmin() {
		t.Fatalf("Expected CREATE_ADMIN to be disabled")
	}
}

func TestDatabaseMaxConnsOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.DatabaseMaxConns() != 20 {
		t.Fatalf("Expected DATABASE_MAX_CONNS to be 20 by default")
	}

	if err := configParser.parseLines([]string{"DATABASE_MAX_CONNS=10"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.DatabaseMaxConns() != 10 {
		t.Fatalf("Expected DATABASE_MAX_CONNS to be 10")
	}
}

func TestDatabaseMinConnsOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.DatabaseMinConns() != 1 {
		t.Fatalf("Expected DATABASE_MIN_CONNS to be 1 by default")
	}

	if err := configParser.parseLines([]string{"DATABASE_MIN_CONNS=2"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.DatabaseMinConns() != 2 {
		t.Fatalf("Expected DATABASE_MIN_CONNS to be 2")
	}
}

func TestDatabaseURLOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.DatabaseURL() != "user=postgres password=postgres dbname=onebill sslmode=disable" {
		t.Fatal("Expected DATABASE_URL to have default value")
	}

	if !configParser.options.IsDefaultDatabaseURL() {
		t.Fatal("Expected DATABASE_URL to be the default value")
	}

	if err := configParser.parseLines([]string{"DATABASE_URL=postgres://user:pass@localhost/db"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.DatabaseURL() != "postgres://user:pass@localhost/db" {
		t.Fatal("Expected DATABASE_URL to be 'postgres://user:pass@localhost/db'")
	}

	if configParser.options.IsDefaultDatabaseURL() {
		t.Fatal("Expected DATABASE_URL to not be the default value")
	}
}

func TestDisableHSTSOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.DisableHSTS() {
		t.Fatal("Expected DISABLE_HSTS to be disabled by default")
	}

	if !configParser.options.HasHSTS() {
		t.Fatal("Expected HAS_HSTS to be enabled by default")
	}

	if err := configParser.parseLines([]string{"DISABLE_HSTS=1"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !configParser.options.DisableHSTS() {
		t.Fatal("Expected DISABLE_HSTS to be enabled")
	}

	if configParser.options.HasHSTS() {
		t.Fatal("Expected HAS_HSTS to be disabled")
	}

	if err := configParser.parseLines([]string{"DISABLE_HSTS=0"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.DisableHSTS() {
		t.Fatal("Expected DISABLE_HSTS to be disabled")
	}

	if !configParser.options.HasHSTS() {
		t.Fatal("Expected HAS_HSTS to be enabled")
	}
}

func TestDisableHTTPServiceOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.DisableHTTPService() {
		t.Fatal("Expected DISABLE_HTTP_SERVICE to be disabled by default")
	}

	if !configParser.options.HasHTTPService() {
		t.Fatal("Expected HAS_HTTP_SERVICE to be enabled by default")
	}

	if err := configParser.parseLines([]string{"DISABLE_HTTP_SERVICE=1"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !configParser.options.DisableHTTPService() {
		t.Fatal("Expected DISABLE_HTTP_SERVICE to be enabled")
	}

	if configParser.options.HasHTTPService() {
		t.Fatal("Expected HAS_HTTP_SERVICE to be disabled")
	}

	if err := configParser.parseLines([]string{"DISABLE_HTTP_SERVICE=0"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.DisableHTTPService() {
		t.Fatal("Expected DISABLE_HTTP_SERVICE to be disabled")
	}

	if !configParser.options.HasHTTPService() {
		t.Fatal("Expected HAS_HTTP_SERVICE to be disabled")
	}
}

func TestDisableLocalAuthOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.DisableLocalAuth() {
		t.Fatalf("Expected DISABLE_LOCAL_AUTH to be disabled by default")
	}

	if err := configParser.parseLines([]string{"DISABLE_LOCAL_AUTH=1"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !configParser.options.DisableLocalAuth() {
		t.Fatalf("Expected DISABLE_LOCAL_AUTH to be enabled")
	}

	if err := configParser.parseLines([]string{"DISABLE_LOCAL_AUTH=0"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.DisableLocalAuth() {
		t.Fatalf("Expected DISABLE_LOCAL_AUTH to be disabled")
	}
}

func TestDisableSchedulerServiceOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.DisableSchedulerService() {
		t.Fatal("Expected DISABLE_SCHEDULER_SERVICE to be disabled by default")
	}

	if !configParser.options.HasSchedulerService() {
		t.Fatal("Expected HAS_SCHEDULER_SERVICE to be enabled by default")
	}

	if err := configParser.parseLines([]string{"DISABLE_SCHEDULER_SERVICE=1"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !configParser.options.DisableSchedulerService() {
		t.Fatal("Expected DISABLE_SCHEDULER_SERVICE to be enabled")
	}

	if configParser.options.HasSchedulerService() {
		t.Fatal("Expected HAS_SCHEDULER_SERVICE to be disabled")
	}

	if err := configParser.parseLines([]string{"DISABLE_SCHEDULER_SERVICE=0"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.DisableSchedulerService() {
		t.Fatal("Expected DISABLE_SCHEDULER_SERVICE to be disabled")
	}

	if !configParser.options.HasSchedulerService() {
		t.Fatal("Expected HAS_SCHEDULER_SERVICE to be enabled")
	}
}

func TestHTTPSOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.HTTPS() {
		t.Fatalf("Expected HTTPS to be disabled by default")
	}

	if err := configParser.parseLines([]string{"HTTPS=1"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !configParser.options.HTTPS() {
		t.Fatalf("Expected HTTPS to be enabled")
	}

	if err := configParser.parseLines([]string{"HTTPS=0"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.HTTPS() {
		t.Fatalf("Expected HTTPS to be disabled")
	}
}

func TestCertKeyFileOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.CertKeyFile() != "" {
		t.Fatalf("Expected KEY_FILE to be empty by default")
	}

	if err := configParser.parseLines([]string{"KEY_FILE=/path/to/key.pem"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.CertKeyFile() != "/path/to/key.pem" {
		t.Fatalf("Expected KEY_FILE to be '/path/to/key.pem'")
	}
}

func TestLogDateTimeOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.LogDateTime() {
		t.Fatalf("Expected LOG_DATE_TIME to be disabled by default")
	}

	if err := configParser.parseLines([]string{"LOG_DATE_TIME=1"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !configParser.options.LogDateTime() {
		t.Fatalf("Expected LOG_DATE_TIME to be enabled")
	}

	if err := configParser.parseLines([]string{"LOG_DATE_TIME=0"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.LogDateTime() {
		t.Fatalf("Expected LOG_DATE_TIME to be disabled")
	}
}

func TestLogFileOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.LogFile() != "stderr" {
		t.Fatalf("Expected LOG_FILE to be 'stderr' by default")
	}

	if err := configParser.parseLines([]string{"LOG_FILE=/var/log/onebill.log"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.LogFile() != "/var/log/onebill.log" {
		t.Fatalf("Expected LOG_FILE to be '/var/log/onebill.log'")
	}
}

func TestLogFormatOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.LogFormat() != "text" {
		t.Fatalf("Expected LOG_FORMAT to be 'text' by default")
	}

	if err := configParser.parseLines([]string{"LOG_FORMAT=json"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.LogFormat() != "json" {
		t.Fatalf("Expected LOG_FORMAT to be 'json'")
	}
}

func TestLogLevelOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.LogLevel() != "info" {
		t.Fatalf("Expected LOG_LEVEL to be 'info' by default")
	}

	if err := configParser.parseLines([]string{"LOG_LEVEL=debug"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.LogLevel() != "debug" {
		t.Fatalf("Expected LOG_LEVEL to be 'debug'")
	}
}

func TestMaintenanceMessageOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.MaintenanceMessage() != "OneBill is currently under maintenance" {
		t.Fatalf("Expected MAINTENANCE_MESSAGE to have default value")
	}

	if err := configParser.parseLines([]string{"MAINTENANCE_MESSAGE=System upgrade in progress"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.MaintenanceMessage() != "System upgrade in progress" {
		t.Fatalf("Expected MAINTENANCE_MESSAGE to be 'System upgrade in progress'")
	}
}

func TestMaintenanceModeOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.MaintenanceMode() {
		t.Fatal("Expected MAINTENANCE_MODE to be disabled by default")
	}

	if configParser.options.HasMaintenanceMode() {
		t.Fatal("Expected HAS_MAINTENANCE_MODE to be disabled by default")
	}

	if err := configParser.parseLines([]string{"MAINTENANCE_MODE=1"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !configParser.options.MaintenanceMode() {
		t.Fatal("Expected MAINTENANCE_MODE to be enabled")
	}

	if !configParser.options.HasMaintenanceMode() {
		t.Fatal("Expected HAS_MAINTENANCE_MODE to be enabled")
	}

	if err := configParser.parseLines([]string{"MAINTENANCE_MODE=0"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.MaintenanceMode() {
		t.Fatal("Expected MAINTENANCE_MODE to be disabled")
	}

	if configParser.options.HasMaintenanceMode() {
		t.Fatal("Expected HAS_MAINTENANCE_MODE to be disabled")
	}
}

func TestMetricsCollectorOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.MetricsCollector() {
		t.Fatal("Expected METRICS_COLLECTOR to be disabled by default")
	}

	if configParser.options.HasMetricsCollector() {
		t.Fatal("Expected HAS_METRICS_COLLECTOR to be disabled by default")
	}

	if err := configParser.parseLines([]string{"METRICS_COLLECTOR=1"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !configParser.options.MetricsCollector() {
		t.Fatal("Expected METRICS_COLLECTOR to be enabled")
	}

	if !configParser.options.HasMetricsCollector() {
		t.Fatal("Expected HAS_METRICS_COLLECTOR to be enabled")
	}

	if err := configParser.parseLines([]string{"METRICS_COLLECTOR=0"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.MetricsCollector() {
		t.Fatal("Expected METRICS_COLLECTOR to be disabled")
	}

	if configParser.options.HasMetricsCollector() {
		t.Fatal("Expected HAS_METRICS_COLLECTOR to be disabled")
	}
}

func TestMetricsPasswordOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.MetricsPassword() != "" {
		t.Fatalf("Expected METRICS_PASSWORD to be empty by default")
	}

	if err := configParser.parseLines([]string{"METRICS_PASSWORD=secret123"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.MetricsPassword() != "secret123" {
		t.Fatalf("Expected METRICS_PASSWORD to be 'secret123'")
	}
}

func TestMetricsUsernameOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.MetricsUsername() != "" {
		t.Fatalf("Expected METRICS_USERNAME to be empty by default")
	}

	if err := configParser.parseLines([]string{"METRICS_USERNAME=metrics_user"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.MetricsUsername() != "metrics_user" {
		t.Fatalf("Expected METRICS_USERNAME to be 'metrics_user'")
	}
}

func TestOAuth2ClientIDOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.OAuth2ClientID() != "" {
		t.Fatalf("Expected OAUTH2_CLIENT_ID to be empty by default")
	}

	if err := configParser.parseLines([]string{"OAUTH2_CLIENT_ID=client123"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.OAuth2ClientID() != "client123" {
		t.Fatalf("Expected OAUTH2_CLIENT_ID to be 'client123'")
	}
}

func TestOAuth2ClientSecretOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.OAuth2ClientSecret() != "" {
		t.Fatalf("Expected OAUTH2_CLIENT_SECRET to be empty by default")
	}

	if err := configParser.parseLines([]string{"OAUTH2_CLIENT_SECRET=secret456"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.OAuth2ClientSecret() != "secret456" {
		t.Fatalf("Expected OAUTH2_CLIENT_SECRET to be 'secret456'")
	}
}

func TestOAuth2OIDCDiscoveryEndpointOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.OAuth2OIDCDiscoveryEndpoint() != "" {
		t.Fatalf("Expected OAUTH2_OIDC_DISCOVERY_ENDPOINT to be empty by default")
	}

	if err := configParser.parseLines([]string{"OAUTH2_OIDC_DISCOVERY_ENDPOINT=https://example.com/.well-known/openid_configuration"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.OAuth2OIDCDiscoveryEndpoint() != "https://example.com/.well-known/openid_configuration" {
		t.Fatalf("Expected OAUTH2_OIDC_DISCOVERY_ENDPOINT to be 'https://example.com/.well-known/openid_configuration'")
	}
}

func TestOAuth2OIDCProviderNameOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.OAuth2OIDCProviderName() != "OpenID Connect" {
		t.Fatalf("Expected OAUTH2_OIDC_PROVIDER_NAME to be 'OpenID Connect' by default")
	}

	if err := configParser.parseLines([]string{"OAUTH2_OIDC_PROVIDER_NAME=My Provider"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.OAuth2OIDCProviderName() != "My Provider" {
		t.Fatalf("Expected OAUTH2_OIDC_PROVIDER_NAME to be 'My Provider'")
	}
}

func TestOAuth2ProviderOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.OAuth2Provider() != "" {
		t.Fatal("Expected OAUTH2_PROVIDER to be empty by default")
	}

	if err := configParser.parseLines([]string{"OAUTH2_PROVIDER=google"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.OAuth2Provider() != "google" {
		t.Fatal("Expected OAUTH2_PROVIDER to be 'google'")
	}

	if err := configParser.parseLines([]string{"OAUTH2_PROVIDER=oidc"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.OAuth2Provider() != "oidc" {
		t.Fatal("Expected OAUTH2_PROVIDER to be 'oidc'")
	}

	if err := configParser.parseLines([]string{"OAUTH2_PROVIDER=invalid"}); err == nil {
		t.Fatal("Expected error for invalid OAUTH2_PROVIDER value")
	}
}

func TestOAuth2RedirectURLOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.OAuth2RedirectURL() != "" {
		t.Fatalf("Expected OAUTH2_REDIRECT_URL to be empty by default")
	}

	if err := configParser.parseLines([]string{"OAUTH2_REDIRECT_URL=https://example.com/callback"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.OAuth2RedirectURL() != "https://example.com/callback" {
		t.Fatalf("Expected OAUTH2_REDIRECT_URL to be 'https://example.com/callback'")
	}
}

func TestOAuth2UserCreationOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.OAuth2UserCreation() {
		t.Fatal("Expected OAUTH2_USER_CREATION to be disabled by default")
	}

	if configParser.options.IsOAuth2UserCreationAllowed() {
		t.Fatal("Expected OAUTH2_USER_CREATION to be disabled by default")
	}

	if err := configParser.parseLines([]string{"OAUTH2_USER_CREATION=1"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !configParser.options.OAuth2UserCreation() {
		t.Fatal("Expected OAUTH2_USER_CREATION to be enabled")
	}

	if !configParser.options.IsOAuth2UserCreationAllowed() {
		t.Fatal("Expected OAUTH2_USER_CREATION to be enabled")
	}

	if err := configParser.parseLines([]string{"OAUTH2_USER_CREATION=0"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.OAuth2UserCreation() {
		t.Fatal("Expected OAUTH2_USER_CREATION to be disabled")
	}

	if configParser.options.IsOAuth2UserCreationAllowed() {
		t.Fatal("Expected OAUTH2_USER_CREATION to be disabled")
	}
}

func TestPollingParsingErrorLimitOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.PollingParsingErrorLimit() != 3 {
		t.Fatalf("Expected POLLING_PARSING_ERROR_LIMIT to be 3 by default")
	}

	if err := configParser.parseLines([]string{"POLLING_PARSING_ERROR_LIMIT=5"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.PollingParsingErrorLimit() != 5 {
		t.Fatalf("Expected POLLING_PARSING_ERROR_LIMIT to be 5")
	}
}

func TestPortOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.Port() != "" {
		t.Fatalf("Expected PORT to be empty by default")
	}

	if err := configParser.parseLines([]string{"PORT=1234"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.Port() != "1234" {
		t.Fatalf("Expected PORT to be '1234'")
	}

	addresses := configParser.options.ListenAddr()
	if len(addresses) != 1 || addresses[0] != ":1234" {
		t.Fatalf("Expected LISTEN_ADDR to be ':1234'")
	}
}

func TestRunMigrationsOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.RunMigrations() {
		t.Fatalf("Expected RUN_MIGRATIONS to be disabled by default")
	}

	if err := configParser.parseLines([]string{"RUN_MIGRATIONS=1"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !configParser.options.RunMigrations() {
		t.Fatalf("Expected RUN_MIGRATIONS to be enabled")
	}

	if err := configParser.parseLines([]string{"RUN_MIGRATIONS=0"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.RunMigrations() {
		t.Fatalf("Expected RUN_MIGRATIONS to be disabled")
	}
}

func TestCleanupArchiveReadIntervalOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.CleanupArchiveReadInterval().Hours() != 24*60 {
		t.Fatalf("Expected CLEANUP_ARCHIVE_READ_DAYS to be 60 days by default")
	}

	if err := configParser.parseLines([]string{"CLEANUP_ARCHIVE_READ_DAYS=30"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.CleanupArchiveReadInterval().Hours() != 24*30 {
		t.Fatalf("Expected CLEANUP_ARCHIVE_READ_DAYS to be 30 days")
	}
}

func TestCleanupArchiveUnreadIntervalOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.CleanupArchiveUnreadInterval().Hours() != 24*180 {
		t.Fatalf("Expected CLEANUP_ARCHIVE_UNREAD_DAYS to be 180 days by default")
	}

	if err := configParser.parseLines([]string{"CLEANUP_ARCHIVE_UNREAD_DAYS=90"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.CleanupArchiveUnreadInterval().Hours() != 24*90 {
		t.Fatalf("Expected CLEANUP_ARCHIVE_UNREAD_DAYS to be 90 days")
	}
}

func TestCleanupFrequencyOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.CleanupFrequency().Hours() != 24 {
		t.Fatalf("Expected CLEANUP_FREQUENCY_HOURS to be 24 hours by default")
	}

	if err := configParser.parseLines([]string{"CLEANUP_FREQUENCY_HOURS=12"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.CleanupFrequency().Hours() != 12 {
		t.Fatalf("Expected CLEANUP_FREQUENCY_HOURS to be 12 hours")
	}
}

func TestCleanupRemoveSessionsIntervalOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.CleanupRemoveSessionsInterval().Hours() != 24*30 {
		t.Fatalf("Expected CLEANUP_REMOVE_SESSIONS_DAYS to be 30 days by default")
	}

	if err := configParser.parseLines([]string{"CLEANUP_REMOVE_SESSIONS_DAYS=14"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.CleanupRemoveSessionsInterval().Hours() != 24*14 {
		t.Fatalf("Expected CLEANUP_REMOVE_SESSIONS_DAYS to be 14 days")
	}
}

func TestDatabaseConnectionLifetimeOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.DatabaseConnectionLifetime().Minutes() != 5 {
		t.Fatalf("Expected DATABASE_CONNECTION_LIFETIME to be 5 minutes by default")
	}

	if err := configParser.parseLines([]string{"DATABASE_CONNECTION_LIFETIME=10"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.DatabaseConnectionLifetime().Minutes() != 10 {
		t.Fatalf("Expected DATABASE_CONNECTION_LIFETIME to be 10 minutes")
	}
}

func TestIntegrationAllowPrivateNetworksOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.IntegrationAllowPrivateNetworks() {
		t.Fatalf("Expected INTEGRATION_ALLOW_PRIVATE_NETWORKS to be disabled by default")
	}

	if err := configParser.parseLines([]string{"INTEGRATION_ALLOW_PRIVATE_NETWORKS=1"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !configParser.options.IntegrationAllowPrivateNetworks() {
		t.Fatalf("Expected INTEGRATION_ALLOW_PRIVATE_NETWORKS to be enabled")
	}

	if err := configParser.parseLines([]string{"INTEGRATION_ALLOW_PRIVATE_NETWORKS=0"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.IntegrationAllowPrivateNetworks() {
		t.Fatalf("Expected INTEGRATION_ALLOW_PRIVATE_NETWORKS to be disabled")
	}
}

func TestHTTPServerTimeoutOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.HTTPServerTimeout().Seconds() != 300 {
		t.Fatal("Expected HTTP_SERVER_TIMEOUT to be 300 seconds by default")
	}

	if err := configParser.parseLines([]string{"HTTP_SERVER_TIMEOUT=60"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.HTTPServerTimeout().Seconds() != 60 {
		t.Fatal("Expected HTTP_SERVER_TIMEOUT to be 60 seconds")
	}
}

func TestListenAddrOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	addrs := configParser.options.ListenAddr()
	if len(addrs) != 1 || addrs[0] != "127.0.0.1:8080" {
		t.Fatalf("Expected LISTEN_ADDR to be '127.0.0.1:8080' by default")
	}

	if err := configParser.parseLines([]string{"LISTEN_ADDR=0.0.0.0:8080,127.0.0.1:8081,/unix.socket"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	addrs = configParser.options.ListenAddr()
	if len(addrs) != 3 || addrs[0] != "0.0.0.0:8080" || addrs[1] != "127.0.0.1:8081" || addrs[2] != "/unix.socket" {
		t.Fatalf("Expected LISTEN_ADDR to contain two addresses")
	}
}

func TestMetricsAllowedNetworksOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	networks := configParser.options.MetricsAllowedNetworks()
	if len(networks) != 1 || networks[0] != "127.0.0.1/8" {
		t.Fatalf("Expected METRICS_ALLOWED_NETWORKS to have default values")
	}

	if err := configParser.parseLines([]string{"METRICS_ALLOWED_NETWORKS=10.0.0.0/8,192.168.0.0/16"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	networks = configParser.options.MetricsAllowedNetworks()
	if len(networks) != 2 || networks[0] != "10.0.0.0/8" || networks[1] != "192.168.0.0/16" {
		t.Fatalf("Expected METRICS_ALLOWED_NETWORKS to contain specified networks")
	}
}

func TestMetricsRefreshIntervalOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	if configParser.options.MetricsRefreshInterval().Seconds() != 60 {
		t.Fatalf("Expected METRICS_REFRESH_INTERVAL to be 60 seconds by default")
	}

	if err := configParser.parseLines([]string{"METRICS_REFRESH_INTERVAL=120"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if configParser.options.MetricsRefreshInterval().Seconds() != 120 {
		t.Fatalf("Expected METRICS_REFRESH_INTERVAL to be 120 seconds")
	}
}

func TestTrustedReverseProxyNetworksOptionParsing(t *testing.T) {
	configParser := NewConfigParser()

	// Test default value
	defaultNetworks := configParser.options.TrustedReverseProxyNetworks()
	if len(defaultNetworks) != 0 {
		t.Fatalf("Expected 0 allowed networks by default, got %d", len(defaultNetworks))
	}

	// Test valid value
	if err := configParser.parseLines([]string{"TRUSTED_REVERSE_PROXY_NETWORKS=10.0.0.0/8,192.168.1.0/24"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	allowedNetworks := configParser.options.TrustedReverseProxyNetworks()
	if len(allowedNetworks) != 2 {
		t.Fatalf("Expected 2 allowed networks, got %d", len(allowedNetworks))
	}
	if !slices.Contains(allowedNetworks, "10.0.0.0/8") {
		t.Errorf("Expected 10.0.0.0/8 in allowed networks")
	}
	if !slices.Contains(allowedNetworks, "192.168.1.0/24") {
		t.Errorf("Expected 192.168.1.0/24 in allowed networks")
	}

	// Test value with spaces and trailing comma
	if err := configParser.parseLines([]string{"TRUSTED_REVERSE_PROXY_NETWORKS=192.168.0.0/16, 10.0.0.0/8,"}); err != nil {
		t.Fatalf("Unexpected error for value with spaces: %v", err)
	}

	allowedNetworks = configParser.options.TrustedReverseProxyNetworks()
	if len(allowedNetworks) != 2 {
		t.Fatalf("Expected 2 allowed networks, got %d", len(allowedNetworks))
	}
	if !slices.Contains(allowedNetworks, "192.168.0.0/16") {
		t.Errorf("Expected 192.168.0.0/16 in allowed networks")
	}
	if !slices.Contains(allowedNetworks, "10.0.0.0/8") {
		t.Errorf("Expected 10.0.0.0/8 in allowed networks")
	}

	// Test invalid value
	if err := configParser.parseLines([]string{"TRUSTED_REVERSE_PROXY_NETWORKS=127.0.0.1"}); err == nil {
		t.Fatal("Expected error when parsing invalid CIDR notation IP 127.0.0.1, got nil")
	}

	if err := configParser.parseLines([]string{"TRUSTED_REVERSE_PROXY_NETWORKS=,"}); err == nil {
		t.Fatal("Expected error for value without any network")
	}
}

func TestSetLogLevelFunction(t *testing.T) {
	configParser := NewConfigParser()

	// Test default log level
	if configParser.options.LogLevel() != "info" {
		t.Fatalf("Expected LOG_LEVEL to be 'info' by default, got '%s'", configParser.options.LogLevel())
	}

	// Test setting log level to debug
	configParser.options.SetLogLevel("debug")
	if configParser.options.LogLevel() != "debug" {
		t.Fatalf("Expected LOG_LEVEL to be 'debug' after SetLogLevel('debug'), got '%s'", configParser.options.LogLevel())
	}
	if configParser.options.options["LOG_LEVEL"].rawValue != "debug" {
		t.Fatalf("Expected LOG_LEVEL RawValue to be 'debug', got '%s'", configParser.options.options["LOG_LEVEL"].rawValue)
	}

	// Test setting log level to warning
	configParser.options.SetLogLevel("warning")
	if configParser.options.LogLevel() != "warning" {
		t.Fatalf("Expected LOG_LEVEL to be 'warning' after SetLogLevel('warning'), got '%s'", configParser.options.LogLevel())
	}
	if configParser.options.options["LOG_LEVEL"].rawValue != "warning" {
		t.Fatalf("Expected LOG_LEVEL RawValue to be 'warning', got '%s'", configParser.options.options["LOG_LEVEL"].rawValue)
	}
}

func TestSetHTTPSValueFunction(t *testing.T) {
	configParser := NewConfigParser()

	// Test setting HTTPS to true
	configParser.options.SetHTTPSValue(true)
	if !configParser.options.HTTPS() {
		t.Fatalf("Expected HTTPS to be true after SetHTTPSValue(true)")
	}

	// Test setting HTTPS to false
	configParser.options.SetHTTPSValue(false)
	if configParser.options.HTTPS() {
		t.Fatalf("Expected HTTPS to be false after SetHTTPSValue(false)")
	}

	// Test setting HTTPS to true again
	configParser.options.SetHTTPSValue(true)
	if !configParser.options.HTTPS() {
		t.Fatalf("Expected HTTPS to be true after second SetHTTPSValue(true)")
	}
}

func TestConfigMap(t *testing.T) {
	configMap := NewConfigOptions().ConfigMap(false)

	if len(configMap) == 0 {
		t.Fatal("Expected ConfigMap to contain configuration options")
	}

	// The first option should be "ADMIN_PASSWORD"
	if configMap[0].Key != "ADMIN_PASSWORD" {
		t.Fatalf("Expected first config option to be 'ADMIN_PASSWORD', got '%s'", configMap[0].Key)
	}
}

func TestConfigMapWithRedactedSecrets(t *testing.T) {
	configParser := NewConfigParser()

	if err := configParser.parseLines([]string{"ADMIN_PASSWORD=secret123"}); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	configMap := configParser.options.ConfigMap(true)

	if len(configMap) == 0 {
		t.Fatal("Expected ConfigMap to contain configuration options")
	}

	// The first option should be "ADMIN_PASSWORD"
	if configMap[0].Key != "ADMIN_PASSWORD" {
		t.Fatalf("Expected first config option to be 'ADMIN_PASSWORD', got '%s'", configMap[0].Key)
	}

	// The value should be redacted
	if configMap[0].Value != "<redacted>" {
		t.Fatalf("Expected ADMIN_PASSWORD value to be redacted, got '%s'", configMap[0].Value)
	}
}

func TestValidateOIDCProviderRequiresDiscoveryEndpoint(t *testing.T) {
	configParser := NewConfigParser()
	if err := configParser.parseLines([]string{"OAUTH2_PROVIDER=oidc"}); err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	err := configParser.options.Validate()
	if err == nil {
		t.Fatal("Expected error when OIDC provider is set without discovery endpoint")
	}
	if err.Error() != "OAUTH2_OIDC_DISCOVERY_ENDPOINT must be configured when using the OIDC provider" {
		t.Fatalf("Unexpected error message: %v", err)
	}
}

func TestValidateOIDCProviderWithDiscoveryEndpoint(t *testing.T) {
	configParser := NewConfigParser()
	if err := configParser.parseLines([]string{
		"OAUTH2_PROVIDER=oidc",
		"OAUTH2_OIDC_DISCOVERY_ENDPOINT=https://example.com/.well-known/openid-configuration",
	}); err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if err := configParser.options.Validate(); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestValidateDisableLocalAuthWithoutAlternative(t *testing.T) {
	configParser := NewConfigParser()
	if err := configParser.parseLines([]string{"DISABLE_LOCAL_AUTH=1"}); err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if err := configParser.options.Validate(); err == nil {
		t.Fatal("Expected error when local auth is disabled without alternative")
	}
}

func TestValidateDisableLocalAuthWithOAuth2ButNoUserCreation(t *testing.T) {
	configParser := NewConfigParser()
	if err := configParser.parseLines([]string{
		"DISABLE_LOCAL_AUTH=1",
		"OAUTH2_PROVIDER=oidc",
		"OAUTH2_OIDC_DISCOVERY_ENDPOINT=https://example.com/.well-known/openid-configuration",
	}); err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if err := configParser.options.Validate(); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestValidateDisableLocalAuthWithOAuth2AndUserCreation(t *testing.T) {
	configParser := NewConfigParser()
	if err := configParser.parseLines([]string{
		"DISABLE_LOCAL_AUTH=1",
		"OAUTH2_PROVIDER=oidc",
		"OAUTH2_OIDC_DISCOVERY_ENDPOINT=https://example.com/.well-known/openid-configuration",
		"OAUTH2_USER_CREATION=1",
	}); err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if err := configParser.options.Validate(); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestValidateDisableLocalAuthWithAuthProxyButNoUserCreation(t *testing.T) {
	configParser := NewConfigParser()
	if err := configParser.parseLines([]string{
		"DISABLE_LOCAL_AUTH=1",
		"AUTH_PROXY_HEADER=X-Forwarded-User",
		"AUTH_PROXY_USER_CREATION=0",
		"TRUSTED_REVERSE_PROXY_NETWORKS=10.0.0.0/8",
	}); err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if err := configParser.options.Validate(); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestValidateDisableLocalAuthWithAuthProxyAndUserCreation(t *testing.T) {
	configParser := NewConfigParser()
	if err := configParser.parseLines([]string{
		"DISABLE_LOCAL_AUTH=1",
		"AUTH_PROXY_HEADER=X-Forwarded-User",
		"AUTH_PROXY_USER_CREATION=1",
		"TRUSTED_REVERSE_PROXY_NETWORKS=10.0.0.0/8",
	}); err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if err := configParser.options.Validate(); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestValidateAuthProxyRequiresTrustedNetworks(t *testing.T) {
	configParser := NewConfigParser()
	if err := configParser.parseLines([]string{"AUTH_PROXY_HEADER=X-Forwarded-User"}); err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	err := configParser.options.Validate()
	if err == nil {
		t.Fatal("Expected error when auth proxy header is set without trusted networks")
	}
	if err.Error() != "TRUSTED_REVERSE_PROXY_NETWORKS must be configured when AUTH_PROXY_HEADER is used" {
		t.Fatalf("Unexpected error message: %v", err)
	}
}

func TestValidateAuthProxyWithTrustedNetworks(t *testing.T) {
	configParser := NewConfigParser()
	if err := configParser.parseLines([]string{
		"AUTH_PROXY_HEADER=X-Forwarded-User",
		"TRUSTED_REVERSE_PROXY_NETWORKS=10.0.0.0/8",
	}); err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if err := configParser.options.Validate(); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestValidateCertFileMissingKeyFile(t *testing.T) {
	configParser := NewConfigParser()
	if err := configParser.parseLines([]string{"CERT_FILE=/path/to/cert.pem"}); err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if err := configParser.options.Validate(); err == nil {
		t.Fatal("Expected error when CERT_FILE is set without KEY_FILE")
	}
}

func TestValidateKeyFileMissingCertFile(t *testing.T) {
	configParser := NewConfigParser()
	if err := configParser.parseLines([]string{"KEY_FILE=/path/to/key.pem"}); err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if err := configParser.options.Validate(); err == nil {
		t.Fatal("Expected error when KEY_FILE is set without CERT_FILE")
	}
}

func TestValidateCertFileAndKeyFile(t *testing.T) {
	configParser := NewConfigParser()
	if err := configParser.parseLines([]string{
		"CERT_FILE=/path/to/cert.pem",
		"KEY_FILE=/path/to/key.pem",
	}); err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if err := configParser.options.Validate(); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestValidateCertDomainAndCertFileMutuallyExclusive(t *testing.T) {
	configParser := NewConfigParser()
	if err := configParser.parseLines([]string{
		"CERT_DOMAIN=example.com",
		"CERT_FILE=/path/to/cert.pem",
		"KEY_FILE=/path/to/key.pem",
	}); err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if err := configParser.options.Validate(); err == nil {
		t.Fatal("Expected error when both CERT_DOMAIN and CERT_FILE are set")
	}
}

func TestValidateCertDomainAlone(t *testing.T) {
	configParser := NewConfigParser()
	if err := configParser.parseLines([]string{"CERT_DOMAIN=example.com"}); err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if err := configParser.options.Validate(); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestValidateMetricsUsernameWithoutPassword(t *testing.T) {
	configParser := NewConfigParser()
	if err := configParser.parseLines([]string{"METRICS_USERNAME=admin"}); err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if err := configParser.options.Validate(); err == nil {
		t.Fatal("Expected error when METRICS_USERNAME is set without METRICS_PASSWORD")
	}
}

func TestValidateMetricsPasswordWithoutUsername(t *testing.T) {
	configParser := NewConfigParser()
	if err := configParser.parseLines([]string{"METRICS_PASSWORD=secret"}); err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if err := configParser.options.Validate(); err == nil {
		t.Fatal("Expected error when METRICS_PASSWORD is set without METRICS_USERNAME")
	}
}

func TestValidateMetricsUsernameAndPassword(t *testing.T) {
	configParser := NewConfigParser()
	if err := configParser.parseLines([]string{
		"METRICS_USERNAME=admin",
		"METRICS_PASSWORD=secret",
	}); err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if err := configParser.options.Validate(); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestValidateDatabaseMinConnsGreaterThanMaxConns(t *testing.T) {
	configParser := NewConfigParser()
	if err := configParser.parseLines([]string{
		"DATABASE_MIN_CONNS=25",
		"DATABASE_MAX_CONNS=10",
	}); err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if err := configParser.options.Validate(); err == nil {
		t.Fatal("Expected error when DATABASE_MIN_CONNS > DATABASE_MAX_CONNS")
	}
}

func TestValidateDatabaseMinConnsEqualToMaxConns(t *testing.T) {
	configParser := NewConfigParser()
	if err := configParser.parseLines([]string{
		"DATABASE_MIN_CONNS=10",
		"DATABASE_MAX_CONNS=10",
	}); err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if err := configParser.options.Validate(); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestValidateSchedulerRoundRobinMinLessThanMax(t *testing.T) {
	configParser := NewConfigParser()
	if err := configParser.parseLines([]string{
		"SCHEDULER_ROUND_ROBIN_MIN_INTERVAL=60",
		"SCHEDULER_ROUND_ROBIN_MAX_INTERVAL=1440",
	}); err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if err := configParser.options.Validate(); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestValidateSchedulerEntryFrequencyMinLessThanMax(t *testing.T) {
	configParser := NewConfigParser()
	if err := configParser.parseLines([]string{
		"SCHEDULER_ENTRY_FREQUENCY_MIN_INTERVAL=5",
		"SCHEDULER_ENTRY_FREQUENCY_MAX_INTERVAL=1440",
	}); err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if err := configParser.options.Validate(); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
