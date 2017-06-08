package publish

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"

	ctrl "haina.com/market/hqpublish/controllers"
	. "haina.com/market/hqpublish/models"
	. "haina.com/share/models"

	pro "ProtocolBuffer/projects/hqpublish/go/protocol"

	"time"

	hsgrr "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
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

var tradeet sync.Mutex

func LockTradeet() {
	tradeet.Lock()
}
func UnlockTradeet() {
	tradeet.Unlock()
}

type TradeetSentinel struct {
	SyncFlag  bool
	Timestamp int64
	RWLock    *sync.RWMutex
}

func NewTradeetSentinel() *TradeetSentinel {
	return &TradeetSentinel{
		Timestamp: time.Now().Unix(),
		RWLock:    &sync.RWMutex{},
	}
}
func getNowTradeetSentinel(sid int32) *TradeetSentinel {
	LockTradeet()
	defer UnlockTradeet()

	t, ok := TradeetMap[sid]
	if ok {
		fmt.Printf("get trade every time sentinel %x\n", t)
		return t
	}

	fmt.Printf("new trade every time sentinel %x\n", t)
	t = NewTradeetSentinel()
	TradeetMap[sid] = t

	return t
}

var TradeetMap map[int32]*TradeetSentinel = make(map[int32]*TradeetSentinel, 5000)

// 将分笔成交数据同步到缓存
// 返回值 int
//   >0 本次同步到缓存的条数
//   =0 缓存和数据redis一致，本次没有进行同步
//   -1 redis缓存错误不可用
//   -2 redis数据错误，没有数据
func (this TradeEveryTime) SyncTradeEveryTimeRecord(sid int32) (int, error) {
	key := fmt.Sprintf(this.CacheKey, sid)

	t := getNowTradeetSentinel(sid)

	t.RWLock.RLock()
	defer t.RWLock.RUnlock()

	t.SyncFlag = false

	clen, err := RedisCache.Llen(key) // 出错时返回值 clen = -1
	if err != nil {
		if err != hsgrr.ErrNil {
			logging.Warning("RedisCache not available: %v", err)
			return -1, err
		} else {
			logging.Info("RedisCache no such key %s", key)
			clen = 0
		}
	}
	slen, err := RedisStore.Llen(key)
	if err != nil {
		logging.Error("RedisStore: %v", err)
		return -2, err
	}
	fmt.Printf("%s Store len %d, Cache len %d\n", key, slen, clen)

	if slen == clen {
		return 0, nil
	} else if slen < clen {
		RedisCache.Del(key)
	}

	ls, err := RedisStore.LRange(key, clen, slen-1)

	go func() {
		t.RWLock.Lock()
		defer t.RWLock.Unlock()

		t.SyncFlag = true
		for _, v := range ls {
			RedisCache.Rpush(key, []byte(v))
		}
		logging.Info("sid %d trade Sync TradeEveryTimeRecord done", sid)
		fmt.Printf("func %s Store len %d, Cache len %d\n", key, slen, clen)
		t.Timestamp = time.Now().Unix()
		t.SyncFlag = false
	}()

	fmt.Printf("---- %s Store len %d, Cache len %d\n", key, slen, clen)
	return slen - clen, nil
}

func (this TradeEveryTime) GetTradeEveryTimeJson(req *pro.RequestTradeEveryTime) ([]byte, error) {
	payload, err := this.GetTradeEveryTimeObj(req)
	if err != nil {
		return nil, err
	}
	return ctrl.MakeRespJson(200, payload)
}
func (this TradeEveryTime) GetTradeEveryTimePB(req *pro.RequestTradeEveryTime) ([]byte, error) {
	payload, err := this.GetTradeEveryTimeObj(req)
	if err != nil {
		return nil, err
	}
	return ctrl.MakeRespDataByPB(200, pro.HAINA_PUBLISH_CMD_ACK_TRADEET, payload)
}

func (this TradeEveryTime) GetTradeEveryTimeObj(req *pro.RequestTradeEveryTime) (*pro.PayloadTradeEveryTime, error) {

	key := fmt.Sprintf(this.CacheKey, req.SID)

	slen, err := RedisStore.Llen(key)
	if err != nil {
		logging.Error("%v", err)
		return nil, err
	}
	if slen == 0 {
		return &pro.PayloadTradeEveryTime{
			SID:     req.SID,
			Total:   int32(slen),
			Begin:   req.Begin,
			Num:     0,
			DTRList: nil,
		}, nil
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

	if end >= slen {
		end = slen - 1 // 本次最后一根 list 索引
	}

	logging.Info("tradeet %d list range[%d,%d]", req.SID, bgn, end)

	ls, err := RedisStore.LRange(key, bgn, end)
	if err != nil {
		logging.Error("%v", err)
		return nil, err
	}

	rows := make([]*pro.TradeEveryTimeRecord, 0, 5000)

	for _, v := range ls {
		trade := &pro.TradeEveryTimeRecord{}
		bufer := bytes.NewBuffer([]byte(v))
		if err := binary.Read(bufer, binary.LittleEndian, trade); err != nil && err != io.EOF {
			logging.Error("%v", err)
			return nil, err
		}
		rows = append(rows, trade)
	}
	logging.Info("tradeet %d get list range[%d,%d] trade data done.", req.SID, bgn, end)

	return &pro.PayloadTradeEveryTime{
		SID:     req.SID,
		Total:   int32(slen),
		Begin:   req.Begin,
		Num:     int32(len(rows)),
		DTRList: rows,
	}, nil

}
