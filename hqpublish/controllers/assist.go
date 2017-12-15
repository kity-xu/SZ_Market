package controllers

import (
	"fmt"
	"strconv"
)

type Sid struct {
	Sid    int
	Symbol string
	Market string
}

func NewSid(sid int) *Sid {
	s := &Sid{
		Sid: sid,
	}
	s.Parse()
	return s
}

func (s *Sid) Parse() {
	s.Symbol = strconv.Itoa(int(s.Sid % 1000000))
	s.Market = strconv.Itoa(int(s.Sid / 1000000))
}

func (s Sid) String() string {
	return fmt.Sprintf("{sid:%v market:\"%s\" symbol:\"%s\"}", s.Sid, s.Market, s.Symbol)
}

func PackAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
