package models

import "errors"

// App Setting
//---------------------------------------------------------------------------------
const (
	APP_NAME    = "hqpost"
	APP_VERSION = "0.0.1.0"
	APP_PID     = "hqpost"
)

// global_table
//---------------------------------------------------------------------------------
const (
	MOGON_SECURITY_TABLE = "basic_securityinfo_table"
	MOGON_MARKET_TABLE   = "basic_securityinfo_table"
)

var (
	ERROR_REDIS_LIST_NULL = errors.New("redis list is null")
)
