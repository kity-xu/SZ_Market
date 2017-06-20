package publish

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	ctrl "haina.com/market/hqpublish/controllers"
	. "haina.com/market/hqpublish/models"
	. "haina.com/share/models"

	pro "ProtocolBuffer/projects/hqpublish/go/protocol"

	hsgrr "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
	//	"haina.com/share/store/redis"
)

var (
	_ = GetCache
	_ = ctrl.MakeRespDataByBytes
	_ = errors.New
	_ = fmt.Println
	_ = hsgrr.Bytes
	_ = logging.Info
	_ = bytes.NewBuffer
	_ = binary.Read
	_ = io.ReadFull
)

type TradeEveryTime struct {
	Model `db:"-"`
}

func NewTradeEveryTime() *TradeEveryTime {
	return &TradeEveryTime{
		Model: Model{
			CacheKey: REDISKEY_SECURITY_TRADE,
		},
	}
}

func (this TradeEveryTime) GetTradeEveryTimeObj(req *pro.RequestTradeEveryTime) (*pro.PayloadTradeEveryTime, error) {
	curr := NewTraceRecord(req.SID)
	start := req.Begin
	stop := int32(-1)
	if req.Num > 0 {
		stop = req.Begin + req.Num - 1
	}
	if err := curr.SyncAndGetTradeRecords(int(start), int(stop)); err != nil {
		return nil, err
	}
	return &pro.PayloadTradeEveryTime{
		SID:     req.SID,
		Total:   curr.Total,
		Begin:   start,
		Num:     int32(len(curr.List)),
		DTRList: curr.List,
	}, nil
}
