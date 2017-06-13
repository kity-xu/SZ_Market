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

var TradeCacheKey = REDISKEY_SECURITY_TRADE

//------------------------------------------------------------------------------
// 读写锁版实现
// RedisCache 缓存成交记录list(和RedisStore一致)
type TradeetSentinel struct {
	Timestamp int64         // 最后同步时间戳
	SyncFlag  bool          // true/false:正在同步/完成同步
	RWLock    *sync.RWMutex // 读写锁
}

//------------------------------------------------------------------------------
var (
	tradeetMap map[int32]*TradeetSentinel = make(map[int32]*TradeetSentinel, 5000)
	tradeet    sync.Mutex
)

func LockTradeet() { // for tradeetMap
	tradeet.Lock()
}
func UnlockTradeet() {
	tradeet.Unlock()
}

//------------------------------------------------------------------------------
func NewTradeetSentinel() *TradeetSentinel {
	return &TradeetSentinel{
		SyncFlag:  false,
		Timestamp: time.Now().Unix(),
		RWLock:    &sync.RWMutex{},
	}
}

func UseTradeetSentinel(sid int32) *TradeetSentinel {
	LockTradeet()
	defer UnlockTradeet()

	t, ok := tradeetMap[sid]
	if ok {
		fmt.Printf("get trade every time sentinel %+v\n", t)
		return t
	}

	fmt.Printf("new trade every time sentinel %+v\n", t)
	t = NewTradeetSentinel()
	tradeetMap[sid] = t

	return t
}

//------------------------------------------------------------------------------

//------------------------------------------------------------------------------
// 将分笔成交数据同步到缓存
// 返回值 int
//   -1  错误，缓存redis不可用(本次后续操作跳过缓存redis)
//   -2  错误，数据redis不可用或没有数据(返回客户端 code 40002)
//   >=0 缓存key-list长度(llen key)，目前缓存redis和数据redis已同步一致
func SyncTradeEveryTimeRecord(sid int32) (int, error) {
	key := fmt.Sprintf(TradeCacheKey, sid)

	t := UseTradeetSentinel(sid)
	t.RWLock.RLock()
	if t.SyncFlag {
		t.RWLock.RUnlock()
		fmt.Printf("%d In the synchronous\n", sid)
		return 0, nil
	}
	t.RWLock.RUnlock()

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
		return slen, nil
	} else if slen < clen {
		RedisCache.Del(key)
	}

	ls, err := RedisStore.LRange(key, clen, slen-1)

	t.RWLock.Lock()
	t.SyncFlag = true
	t.RWLock.Unlock()
	{
		for _, v := range ls {
			RedisCache.Rpush(key, []byte(v))
		}
		logging.Info("sid %d trade Sync TradeEveryTimeRecord done", sid)

		t.Timestamp = time.Now().Unix()
	}
	t.RWLock.Lock()
	t.SyncFlag = false
	t.RWLock.Unlock()

	return slen, nil
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

	go SyncTradeEveryTimeRecord(req.SID)

	t := UseTradeetSentinel(req.SID)
	t.RWLock.RLock()
	defer t.RWLock.RUnlock()

	if t.SyncFlag {
		return this.GetPayloadTradeEveryTimeObj(RedisStore, req)
	}

	payload, err := this.GetPayloadTradeEveryTimeObj(RedisCache, req)
	if err == nil {
		return payload, nil
	}

	return this.GetPayloadTradeEveryTimeObj(RedisStore, req)
}

func (this TradeEveryTime) GetTradeEveryTimeRecordList(rds *redis.RedisPool, req *pro.RequestTradeEveryTime) ([]*pro.TradeEveryTimeRecord, int, error) {

	key := fmt.Sprintf(this.CacheKey, req.SID)

	slen, err := rds.Llen(key)
	if err != nil {
		logging.Error("%v", err)
		return nil, -1, err
	}
	if slen == 0 {
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

	if end >= slen {
		end = slen - 1 // 本次最后一根 list 索引
	}

	logging.Info("tradeet %d list range[%d,%d]", req.SID, bgn, end)

	ls, err := rds.LRange(key, bgn, end)
	if err != nil {
		logging.Error("%v", err)
		return nil, slen, err
	}

	rows := make([]*pro.TradeEveryTimeRecord, 0, 5000)

	for _, v := range ls {
		trade := &pro.TradeEveryTimeRecord{}
		bufer := bytes.NewBuffer([]byte(v))
		if err := binary.Read(bufer, binary.LittleEndian, trade); err != nil && err != io.EOF {
			logging.Error("%v", err)
			return nil, slen, err
		}
		rows = append(rows, trade)
	}
	logging.Info("tradeet %d get list range[%d,%d] trade data done.", req.SID, bgn, end)

	return rows, slen, nil
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
