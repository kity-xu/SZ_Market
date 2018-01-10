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
	ss := strconv.Itoa(int(s.Sid))
	s.Symbol = ss[3:]
	s.Market = ss[:3]
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
