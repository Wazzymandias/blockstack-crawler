package config

import "time"

const (
	defaultApiURL  = "core.blockstack.org"
	defaultTimeout = 30 * time.Second
	defaultApiURLScheme = "https"
)

var (
	ApiURL  = defaultApiURL
	Timeout = defaultTimeout
	ApiURLScheme = defaultApiURLScheme
)