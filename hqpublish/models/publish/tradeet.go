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
	"haina.com/share/store/redis"
)

var (
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

	clen, cerr, slen, serr := SyncTradeEveryTimeRecord(req.SID)
	if serr != nil {
		return nil, serr
	}

	nothing := pro.PayloadTradeEveryTime{
		SID:     req.SID,
		Total:   0,
		Begin:   req.Begin,
		Num:     0,
		DTRList: nil,
	}
	t := UseTradeetSentinel(req.SID)

	// 当缓存Redis正在同步时,使用数据Redis
	if t.SyncFlag {
		goto useRedisStore
	}

	// 缓存可用时使用缓存Redis
	if cerr == nil {
		t.RWLock.RLock()
		defer t.RWLock.RUnlock()
		if clen <= 0 {
			return &nothing, nil
		} else {
			return this.GetPayloadTradeEveryTimeObj(RedisCache, req)
		}
	}

useRedisStore:
	// 缓存不可用时使用数据Redis
	if serr != nil {
		return nil, serr
	}
	if slen <= 0 {
		return &nothing, nil
	} else {
		return this.GetPayloadTradeEveryTimeObj(RedisStore, req)
	}
}

func (this TradeEveryTime) GetTradeEveryTimeRecordList(rds *redis.RedisPool, req *pro.RequestTradeEveryTime) ([]*pro.TradeEveryTimeRecord, int, error) {

	key := fmt.Sprintf(this.CacheKey, req.SID)

	lslen, err := rds.Llen(key)
	if err != nil {
		logging.Error("%v", err)
		return nil, 0, err
	}
	if lslen == 0 {
		return nil, 0, hsgrr.ErrNil
	}

	bgn, end := 0, -1
	if req.Begin > 0 {
		bgn = int(req.Begin)
	}
	if req.Num > 0 {
		end = int(req.Begin + req.Num - 1)
	} else if req.Num <= 0 {
		end = -1
	}

	if end >= lslen {
		end = lslen - 1 // 最后一根 list 索引
	}

	logging.Info("tradeet %d list range[%d,%d]", req.SID, bgn, end)

	ls, err := rds.LRange(key, bgn, end)
	if err != nil {
		logging.Error("%v", err)
		return nil, lslen, err
	}

	rows := make([]*pro.TradeEveryTimeRecord, 0, 5000)

	for _, v := range ls {
		trade := &pro.TradeEveryTimeRecord{}
		bufer := bytes.NewBuffer([]byte(v))
		if err := binary.Read(bufer, binary.LittleEndian, trade); err != nil && err != io.EOF {
			logging.Error("%v", err)
			return nil, lslen, err
		}
		rows = append(rows, trade)
	}
	logging.Info("tradeet %d get list range[%d,%d] trade data done.", req.SID, bgn, end)

	return rows, lslen, nil
}

func (this TradeEveryTime) GetPayloadTradeEveryTimeObj(rds *redis.RedisPool, req *pro.RequestTradeEveryTime) (*pro.PayloadTradeEveryTime, error) {

	rows, total, err := this.GetTradeEveryTimeRecordList(rds, req)
	if err != nil {
		return nil, err
	}

	return &pro.PayloadTradeEveryTime{
		SID:     req.SID,
		Total:   int32(total),
		Begin:   req.Begin,
		Num:     int32(len(rows)),
		DTRList: rows,
	}, nil
}
