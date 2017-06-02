package controllers

import (
	"errors"
)

var (
	ERROR_REQUEST_PARAM           = errors.New("invalid request parameter")
	ERROR_KLINE_BEGIN_TIME        = errors.New("invalid kline begin time")
	ERROR_INDEX_MAYBE_OUTOF_RANGE = errors.New("index maybe out of range")
)
