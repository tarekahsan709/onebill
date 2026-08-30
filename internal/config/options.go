// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package config // import "onebill.app/internal/config"

import (
	"errors"
	"maps"
	"net"
	"net/url"
	"slices"
	"strings"
	"time"
)

type optionPair struct {
	Key   string
	Value string
}

type configValueType int

const (
	stringType configValueType = iota
	stringListType
	boolType
	intType
	int64Type
	urlType
	secondType
	minuteType
	hourType
	dayType
	secretFileType
	bytesType
)

type configValue struct {
	parsedStringValue string
	parsedIntValue    int
	parsedInt64Value  int64
	parsedDuration    time.Duration
	parsedStringList  []string
	parsedURLValue    *url.URL
	parsedBytesValue  []byte
	parsedBoolValue   bool

	secret    bool
	rawValue  string
	valueType configValueType
	targetKey string

	validator func(string) error
}

type configOptions struct {
	rootURL  string
	basePath string
	options  map[string]*configValue
}

// NewConfigOptions creates a new instance of ConfigOptions with default values.
func NewConfigOptions() *configOptions {
	return &configOptions{
		rootURL:  "http://localhost",
		basePath: "",
		options: map[string]*configValue{
			"ADMIN_PASSWORD": {
				parsedStringValue: "",
				rawValue:          "",
				valueType:         stringType,
				secret:            true,
			},
			"ADMIN_PASSWORD_FILE": {
				parsedStringValue: "",
				rawValue:          "",
				valueType:         secretFileType,
				targetKey:         "ADMIN_PASSWORD",
			},
			"ADMIN_USERNAME": {
				parsedStringValue: "",
				rawValue:          "",
				valueType:         stringType,
			},
			"ADMIN_USERNAME_FILE": {
				parsedStringValue: "",
				rawValue:          "",
				valueType:         secretFileType,
				targetKey:         "ADMIN_USERNAME",
			},
			"AUTH_PROXY_HEADER": {
				parsedStringValue: "",
				rawValue:          "",
				valueType:         stringType,
			},
			"AUTH_PROXY_USER_CREATION": {
				parsedBoolValue: false,
				rawValue:        "0",
				valueType:       boolType,
			},
			"BASE_URL": {
				parsedStringValue: "http://localhost",
				rawValue:          "http://localhost",
				valueType:         stringType,
			},
			"CERT_DOMAIN": {
				parsedStringValue: "",
				rawValue:          "",
				valueType:         stringType,
			},
			"CERT_FILE": {
				parsedStringValue: "",
				rawValue:          "",
				valueType:         stringType,
			},
			"CLEANUP_ARCHIVE_BATCH_SIZE": {
				parsedIntValue: 10000,
				rawValue:       "10000",
				valueType:      intType,
				validator: func(rawValue string) error {
					return validateGreaterOrEqualThan(rawValue, 1)
				},
			},
			"CLEANUP_ARCHIVE_READ_DAYS": {
				parsedDuration: time.Hour * 24 * 60,
				rawValue:       "60",
				valueType:      dayType,
			},
			"CLEANUP_ARCHIVE_UNREAD_DAYS": {
				parsedDuration: time.Hour * 24 * 180,
				rawValue:       "180",
				valueType:      dayType,
			},
			"CLEANUP_FREQUENCY_HOURS": {
				parsedDuration: time.Hour * 24,
				rawValue:       "24",
				valueType:      hourType,
				validator: func(rawValue string) error {
					return validateGreaterOrEqualThan(rawValue, 1)
				},
			},
			"CLEANUP_REMOVE_SESSIONS_DAYS": {
				parsedDuration: time.Hour * 24 * 30,
				rawValue:       "30",
				valueType:      dayType,
			},
			"CREATE_ADMIN": {
				parsedBoolValue: false,
				rawValue:        "0",
				valueType:       boolType,
			},
			"DATABASE_CONNECTION_LIFETIME": {
				parsedDuration: time.Minute * 5,
				rawValue:       "5",
				valueType:      minuteType,
				validator: func(rawValue string) error {
					return validateGreaterThan(rawValue, 0)
				},
			},
			"DATABASE_MAX_CONNS": {
				parsedIntValue: 20,
				rawValue:       "20",
				valueType:      intType,
				validator: func(rawValue string) error {
					return validateGreaterOrEqualThan(rawValue, 1)
				},
			},
			"DATABASE_MIN_CONNS": {
				parsedIntValue: 1,
				rawValue:       "1",
				valueType:      intType,
				validator: func(rawValue string) error {
					return validateGreaterOrEqualThan(rawValue, 0)
				},
			},
			"DATABASE_URL": {
				parsedStringValue: "user=postgres password=postgres dbname=onebill sslmode=disable",
				rawValue:          "user=postgres password=postgres dbname=onebill sslmode=disable",
				valueType:         stringType,
				secret:            true,
			},
			"DATABASE_URL_FILE": {
				parsedStringValue: "",
				rawValue:          "",
				valueType:         secretFileType,
				targetKey:         "DATABASE_URL",
			},
			"DISABLE_API": {
				parsedBoolValue: false,
				rawValue:        "0",
				valueType:       boolType,
			},
			"DISABLE_HSTS": {
				parsedBoolValue: false,
				rawValue:        "0",
				valueType:       boolType,
			},
			"DISABLE_HTTP_SERVICE": {
				parsedBoolValue: false,
				rawValue:        "0",
				valueType:       boolType,
			},
			"DISABLE_LOCAL_AUTH": {
				parsedBoolValue: false,
				rawValue:        "0",
				valueType:       boolType,
			},
			"DISABLE_SCHEDULER_SERVICE": {
				parsedBoolValue: false,
				rawValue:        "0",
				valueType:       boolType,
			},
			"HTTP_SERVER_TIMEOUT": {
				parsedDuration: 300 * time.Second,
				rawValue:       "300",
				valueType:      secondType,
				validator: func(rawValue string) error {
					return validateGreaterOrEqualThan(rawValue, 1)
				},
			},
			"HTTPS": {
				parsedBoolValue: false,
				rawValue:        "0",
				valueType:       boolType,
			},
			"INTEGRATION_ALLOW_PRIVATE_NETWORKS": {
				parsedBoolValue: false,
				rawValue:        "0",
				valueType:       boolType,
			},
			"KEY_FILE": {
				parsedStringValue: "",
				rawValue:          "",
				valueType:         stringType,
			},
			"LISTEN_ADDR": {
				parsedStringList: []string{"127.0.0.1:8080"},
				rawValue:         "127.0.0.1:8080",
				valueType:        stringListType,
			},
			"LOG_DATE_TIME": {
				parsedBoolValue: false,
				rawValue:        "0",
				valueType:       boolType,
			},
			"LOG_FILE": {
				parsedStringValue: "stderr",
				rawValue:          "stderr",
				valueType:         stringType,
			},
			"LOG_FORMAT": {
				parsedStringValue: "text",
				rawValue:          "text",
				valueType:         stringType,
				validator: func(rawValue string) error {
					return validateChoices(rawValue, []string{"text", "json"})
				},
			},
			"LOG_LEVEL": {
				parsedStringValue: "info",
				rawValue:          "info",
				valueType:         stringType,
				validator: func(rawValue string) error {
					return validateChoices(rawValue, []string{"debug", "info", "warning", "error"})
				},
			},
			"MAINTENANCE_MESSAGE": {
				parsedStringValue: "OneBill is currently under maintenance",
				rawValue:          "OneBill is currently under maintenance",
				valueType:         stringType,
			},
			"MAINTENANCE_MODE": {
				parsedBoolValue: false,
				rawValue:        "0",
				valueType:       boolType,
			},
			"METRICS_ALLOWED_NETWORKS": {
				parsedStringList: []string{"127.0.0.1/8"},
				rawValue:         "127.0.0.1/8",
				valueType:        stringListType,
			},
			"METRICS_COLLECTOR": {
				parsedBoolValue: false,
				rawValue:        "0",
				valueType:       boolType,
			},
			"METRICS_PASSWORD": {
				parsedStringValue: "",
				rawValue:          "",
				valueType:         stringType,
				secret:            true,
			},
			"METRICS_PASSWORD_FILE": {
				parsedStringValue: "",
				rawValue:          "",
				valueType:         secretFileType,
				targetKey:         "METRICS_PASSWORD",
			},
			"METRICS_REFRESH_INTERVAL": {
				parsedDuration: 60 * time.Second,
				rawValue:       "60",
				valueType:      secondType,
				validator: func(rawValue string) error {
					return validateGreaterOrEqualThan(rawValue, 1)
				},
			},
			"METRICS_USERNAME": {
				parsedStringValue: "",
				rawValue:          "",
				valueType:         stringType,
			},
			"METRICS_USERNAME_FILE": {
				parsedStringValue: "",
				rawValue:          "",
				valueType:         secretFileType,
				targetKey:         "METRICS_USERNAME",
			},
			"OAUTH2_CLIENT_ID": {
				parsedStringValue: "",
				rawValue:          "",
				valueType:         stringType,
				secret:            true,
			},
			"OAUTH2_CLIENT_ID_FILE": {
				parsedStringValue: "",
				rawValue:          "",
				valueType:         secretFileType,
				targetKey:         "OAUTH2_CLIENT_ID",
			},
			"OAUTH2_CLIENT_SECRET": {
				parsedStringValue: "",
				rawValue:          "",
				valueType:         stringType,
				secret:            true,
			},
			"OAUTH2_CLIENT_SECRET_FILE": {
				parsedStringValue: "",
				rawValue:          "",
				valueType:         secretFileType,
				targetKey:         "OAUTH2_CLIENT_SECRET",
			},
			"OAUTH2_OIDC_DISCOVERY_ENDPOINT": {
				parsedStringValue: "",
				rawValue:          "",
				valueType:         stringType,
			},
			"OAUTH2_OIDC_PROVIDER_NAME": {
				parsedStringValue: "OpenID Connect",
				rawValue:          "OpenID Connect",
				valueType:         stringType,
			},
			"OAUTH2_PROVIDER": {
				parsedStringValue: "",
				rawValue:          "",
				valueType:         stringType,
				validator: func(rawValue string) error {
					return validateChoices(rawValue, []string{"oidc", "google"})
				},
			},
			"OAUTH2_REDIRECT_URL": {
				parsedStringValue: "",
				rawValue:          "",
				valueType:         stringType,
			},
			"OAUTH2_USER_CREATION": {
				parsedBoolValue: false,
				rawValue:        "0",
				valueType:       boolType,
			},
			"POLLING_PARSING_ERROR_LIMIT": {
				parsedIntValue: 3,
				rawValue:       "3",
				valueType:      intType,
				validator: func(rawValue string) error {
					return validateGreaterOrEqualThan(rawValue, 0)
				},
			},
			"PORT": {
				parsedStringValue: "",
				rawValue:          "",
				valueType:         stringType,
				validator: func(rawValue string) error {
					return validateRange(rawValue, 1, 65535)
				},
			},
			"RUN_MIGRATIONS": {
				parsedBoolValue: false,
				rawValue:        "0",
				valueType:       boolType,
			},
			"TRUSTED_REVERSE_PROXY_NETWORKS": {
				parsedStringList: []string{},
				rawValue:         "",
				valueType:        stringListType,
				validator: func(rawValue string) error {
					networks := parseStringListValue(rawValue, nil)
					if len(networks) == 0 {
						return errors.New("at least one CIDR notation network is required")
					}
					for _, ip := range networks {
						if _, _, err := net.ParseCIDR(ip); err != nil {
							return err
						}
					}

					return nil
				},
			},
			"WATCHDOG": {
				parsedBoolValue: true,
				rawValue:        "1",
				valueType:       boolType,
			},
			"WEBAUTHN": {
				parsedBoolValue: false,
				rawValue:        "0",
				valueType:       boolType,
			},
		},
	}
}

func (c *configOptions) AdminPassword() string {
	return c.options["ADMIN_PASSWORD"].parsedStringValue
}

func (c *configOptions) AdminUsername() string {
	return c.options["ADMIN_USERNAME"].parsedStringValue
}

func (c *configOptions) AuthProxyHeader() string {
	return c.options["AUTH_PROXY_HEADER"].parsedStringValue
}

func (c *configOptions) AuthProxyUserCreation() bool {
	return c.options["AUTH_PROXY_USER_CREATION"].parsedBoolValue
}

func (c *configOptions) BasePath() string {
	return c.basePath
}

func (c *configOptions) BaseURL() string {
	return c.options["BASE_URL"].parsedStringValue
}

func (c *configOptions) RootURL() string {
	return c.rootURL
}

func (c *configOptions) CertDomain() string {
	return c.options["CERT_DOMAIN"].parsedStringValue
}

func (c *configOptions) CertFile() string {
	return c.options["CERT_FILE"].parsedStringValue
}

func (c *configOptions) CleanupArchiveBatchSize() int {
	return c.options["CLEANUP_ARCHIVE_BATCH_SIZE"].parsedIntValue
}

func (c *configOptions) CleanupArchiveReadInterval() time.Duration {
	return c.options["CLEANUP_ARCHIVE_READ_DAYS"].parsedDuration
}

func (c *configOptions) CleanupArchiveUnreadInterval() time.Duration {
	return c.options["CLEANUP_ARCHIVE_UNREAD_DAYS"].parsedDuration
}

func (c *configOptions) CleanupFrequency() time.Duration {
	return c.options["CLEANUP_FREQUENCY_HOURS"].parsedDuration
}

func (c *configOptions) CleanupRemoveSessionsInterval() time.Duration {
	return c.options["CLEANUP_REMOVE_SESSIONS_DAYS"].parsedDuration
}

func (c *configOptions) CreateAdmin() bool {
	return c.options["CREATE_ADMIN"].parsedBoolValue
}

func (c *configOptions) DatabaseConnectionLifetime() time.Duration {
	return c.options["DATABASE_CONNECTION_LIFETIME"].parsedDuration
}

func (c *configOptions) DatabaseMaxConns() int {
	return c.options["DATABASE_MAX_CONNS"].parsedIntValue
}

func (c *configOptions) DatabaseMinConns() int {
	return c.options["DATABASE_MIN_CONNS"].parsedIntValue
}

func (c *configOptions) DatabaseURL() string {
	return c.options["DATABASE_URL"].parsedStringValue
}

func (c *configOptions) DisableHSTS() bool {
	return c.options["DISABLE_HSTS"].parsedBoolValue
}

func (c *configOptions) DisableHTTPService() bool {
	return c.options["DISABLE_HTTP_SERVICE"].parsedBoolValue
}

func (c *configOptions) DisableLocalAuth() bool {
	return c.options["DISABLE_LOCAL_AUTH"].parsedBoolValue
}

func (c *configOptions) DisableSchedulerService() bool {
	return c.options["DISABLE_SCHEDULER_SERVICE"].parsedBoolValue
}

func (c *configOptions) HasAPI() bool {
	return !c.options["DISABLE_API"].parsedBoolValue
}

func (c *configOptions) HasHTTPService() bool {
	return !c.options["DISABLE_HTTP_SERVICE"].parsedBoolValue
}

func (c *configOptions) HasHSTS() bool {
	return !c.options["DISABLE_HSTS"].parsedBoolValue
}

func (c *configOptions) HasMaintenanceMode() bool {
	return c.options["MAINTENANCE_MODE"].parsedBoolValue
}

func (c *configOptions) HasMetricsCollector() bool {
	return c.options["METRICS_COLLECTOR"].parsedBoolValue
}

func (c *configOptions) HasSchedulerService() bool {
	return !c.options["DISABLE_SCHEDULER_SERVICE"].parsedBoolValue
}

func (c *configOptions) HasWatchdog() bool {
	return c.options["WATCHDOG"].parsedBoolValue
}

func (c *configOptions) HTTPServerTimeout() time.Duration {
	return c.options["HTTP_SERVER_TIMEOUT"].parsedDuration
}

func (c *configOptions) HTTPS() bool {
	return c.options["HTTPS"].parsedBoolValue
}

func (c *configOptions) IntegrationAllowPrivateNetworks() bool {
	if c == nil {
		return false
	}
	return c.options["INTEGRATION_ALLOW_PRIVATE_NETWORKS"].parsedBoolValue
}

func (c *configOptions) IsAuthProxyUserCreationAllowed() bool {
	return c.options["AUTH_PROXY_USER_CREATION"].parsedBoolValue
}

func (c *configOptions) IsDefaultDatabaseURL() bool {
	return c.options["DATABASE_URL"].rawValue == "user=postgres password=postgres dbname=onebill sslmode=disable"
}

func (c *configOptions) IsOAuth2UserCreationAllowed() bool {
	return c.options["OAUTH2_USER_CREATION"].parsedBoolValue
}

func (c *configOptions) CertKeyFile() string {
	return c.options["KEY_FILE"].parsedStringValue
}

func (c *configOptions) ListenAddr() []string {
	return c.options["LISTEN_ADDR"].parsedStringList
}

func (c *configOptions) LogFile() string {
	return c.options["LOG_FILE"].parsedStringValue
}

func (c *configOptions) LogDateTime() bool {
	return c.options["LOG_DATE_TIME"].parsedBoolValue
}

func (c *configOptions) LogFormat() string {
	return c.options["LOG_FORMAT"].parsedStringValue
}

func (c *configOptions) LogLevel() string {
	return c.options["LOG_LEVEL"].parsedStringValue
}

func (c *configOptions) MaintenanceMessage() string {
	return c.options["MAINTENANCE_MESSAGE"].parsedStringValue
}

func (c *configOptions) MaintenanceMode() bool {
	return c.options["MAINTENANCE_MODE"].parsedBoolValue
}

func (c *configOptions) MetricsAllowedNetworks() []string {
	return c.options["METRICS_ALLOWED_NETWORKS"].parsedStringList
}

func (c *configOptions) MetricsCollector() bool {
	return c.options["METRICS_COLLECTOR"].parsedBoolValue
}

func (c *configOptions) MetricsPassword() string {
	return c.options["METRICS_PASSWORD"].parsedStringValue
}

func (c *configOptions) MetricsRefreshInterval() time.Duration {
	return c.options["METRICS_REFRESH_INTERVAL"].parsedDuration
}

func (c *configOptions) MetricsUsername() string {
	return c.options["METRICS_USERNAME"].parsedStringValue
}

func (c *configOptions) OAuth2ClientID() string {
	return c.options["OAUTH2_CLIENT_ID"].parsedStringValue
}

func (c *configOptions) OAuth2ClientSecret() string {
	return c.options["OAUTH2_CLIENT_SECRET"].parsedStringValue
}

func (c *configOptions) OAuth2OIDCDiscoveryEndpoint() string {
	return c.options["OAUTH2_OIDC_DISCOVERY_ENDPOINT"].parsedStringValue
}

func (c *configOptions) OAuth2OIDCProviderName() string {
	return c.options["OAUTH2_OIDC_PROVIDER_NAME"].parsedStringValue
}

func (c *configOptions) OAuth2Provider() string {
	return c.options["OAUTH2_PROVIDER"].parsedStringValue
}

func (c *configOptions) OAuth2RedirectURL() string {
	return c.options["OAUTH2_REDIRECT_URL"].parsedStringValue
}

func (c *configOptions) OAuth2UserCreation() bool {
	return c.options["OAUTH2_USER_CREATION"].parsedBoolValue
}

func (c *configOptions) PollingParsingErrorLimit() int {
	return c.options["POLLING_PARSING_ERROR_LIMIT"].parsedIntValue
}

func (c *configOptions) Port() string {
	return c.options["PORT"].parsedStringValue
}

func (c *configOptions) RunMigrations() bool {
	return c.options["RUN_MIGRATIONS"].parsedBoolValue
}

func (c *configOptions) SetLogLevel(level string) {
	c.options["LOG_LEVEL"].parsedStringValue = level
	c.options["LOG_LEVEL"].rawValue = level
}

func (c *configOptions) SetHTTPSValue(value bool) {
	c.options["HTTPS"].parsedBoolValue = value
	if value {
		c.options["HTTPS"].rawValue = "1"
	} else {
		c.options["HTTPS"].rawValue = "0"
	}
}

func (c *configOptions) TrustedReverseProxyNetworks() []string {
	return c.options["TRUSTED_REVERSE_PROXY_NETWORKS"].parsedStringList
}

func (c *configOptions) Watchdog() bool {
	return c.options["WATCHDOG"].parsedBoolValue
}

func (c *configOptions) WebAuthn() bool {
	return c.options["WEBAUTHN"].parsedBoolValue
}

func (c *configOptions) ConfigMap(redactSecret bool) []*optionPair {
	sortedKeys := slices.Sorted(maps.Keys(c.options))
	sortedOptions := make([]*optionPair, 0, len(sortedKeys))
	for _, key := range sortedKeys {
		value := c.options[key]
		displayValue := value.rawValue
		if displayValue != "" && redactSecret && value.secret {
			displayValue = "<redacted>"
		}
		sortedOptions = append(sortedOptions, &optionPair{Key: key, Value: displayValue})
	}
	return sortedOptions
}

func (c *configOptions) String() string {
	var builder strings.Builder

	for _, option := range c.ConfigMap(false) {
		builder.WriteString(option.Key)
		builder.WriteByte('=')
		builder.WriteString(option.Value)
		builder.WriteByte('\n')
	}

	return builder.String()
}
