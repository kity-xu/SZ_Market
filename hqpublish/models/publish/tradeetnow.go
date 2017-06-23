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

type TradeEveryTimeNow struct {
	Model `db:"-"`
}

func NewTradeEveryTimeNow() *TradeEveryTimeNow {
	return &TradeEveryTimeNow{
		Model: Model{
			CacheKey: REDISKEY_SECURITY_TRADE,
		},
	}
}

// 未完善
func (this TradeEveryTimeNow) GetTradeEveryTimeNowObj2(req *pro.RequestTradeEveryTimeNow) (*pro.PayloadTradeEveryTimeNow, error) {
	curr := NewTraceRecord(req.SID)
	start := -req.Num
	stop := int32(-1)
	if err := curr.SyncAndGetTradeRecords(int(start), int(stop)); err != nil {
		return nil, err
	}
	return &pro.PayloadTradeEveryTimeNow{
		SID:     req.SID,
		Total:   curr.Total,
		Num:     int32(len(curr.List)),
		DTRList: curr.List,
	}, nil
}

func (this TradeEveryTimeNow) GetTradeEveryTimeNowObj(req *pro.RequestTradeEveryTimeNow) (*pro.PayloadTradeEveryTimeNow, error) {

	key := fmt.Sprintf(this.CacheKey, req.SID)

	slen, err := RedisStore.Llen(key)
	if err != nil {
		logging.Error("%v", err)
		return nil, err
	}
	if slen <= 0 {
		return &pro.PayloadTradeEveryTimeNow{
			SID:     req.SID,
			Total:   int32(slen),
			Num:     0,
			DTRList: nil,
		}, nil
	}

	bgn, end := 0, slen-1
	if req.Num > 0 {
		bgn = slen - int(req.Num)
	} else if req.Num <= 0 {
		bgn = 0
	}
	if bgn < 0 {
		bgn = 0
	}

	logging.Info("tradeetnow %d list range[%d,%d]", req.SID, bgn, end)

	ls, err := RedisStore.LRange(key, bgn, end)
	if err != nil {
		logging.Error("%v", err)
		return nil, err
	}

	rows := make([]*pro.TradeEveryTimeRecord, 0, 5000)

	for i := len(ls) - 1; i > -1; i-- {
		trade := &pro.TradeEveryTimeRecord{}
		bufer := bytes.NewBuffer([]byte(ls[i]))
		if err := binary.Read(bufer, binary.LittleEndian, trade); err != nil && err != io.EOF {
			logging.Error("%v", err)
			return nil, err
		}
		rows = append(rows, trade)
	}
	logging.Info("tradeetnow %d get list range[%d,%d] trade data done.", req.SID, bgn, end)

	return &pro.PayloadTradeEveryTimeNow{
		SID:     req.SID,
		Total:   int32(slen),
		Num:     int32(len(rows)),
		DTRList: rows,
	}, nil
}
