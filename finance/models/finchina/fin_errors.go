package finchina

import "errors"

var (
	ErrIncData  = errors.New("finchina db: incomplete data")
	ErrNullComp = errors.New("finchina db: COMPCODE is NULL")
)
