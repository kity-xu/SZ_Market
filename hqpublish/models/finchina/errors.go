package finchina

import "errors"

var (
	ErrIncData  = errors.New("finchina db: incomplete data")
	ErrNullComp = errors.New("finchina db: COMPCODE is NULL")
	ErrMarket   = errors.New("finchina db: Unknown Market type")
)
