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
	//	. "haina.com/share/models"

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

var local *time.Location

func init() {
	local, _ = time.LoadLocation("Local")
}

var TradeCacheKey = REDISKEY_SECURITY_TRADE

type TraceRecord struct {
	Sid      int32
	Key      string
	Sentinel *TradeSentinel
	Total    int32                       // 返回的当前交易数量
	List     []*pro.TradeEveryTimeRecord // 返回的成交记录
}

func NewTraceRecord(sid int32) *TraceRecord {
	t := UseTradeSentinel(sid)
	return &TraceRecord{
		Sid:      sid,
		Key:      fmt.Sprintf(TradeCacheKey, sid),
		Sentinel: t,
	}
}

//------------------------------------------------------------------------------
// RedisCache 缓存成交记录list(和RedisStore一致)
type TradeSentinel struct {
	Timestamp time.Time     // 最后同步时间戳
	SyncFlag  bool          // true/false:正在同步/完成同步
	RWLock    *sync.RWMutex // 保护RedisCache读写
}

type TradeMap map[int32]*TradeSentinel

//------------------------------------------------------------------------------
var (
	tradeetMap   TradeMap = make(TradeMap, 5000)
	tradeetMutex sync.Mutex
)

func LockTrade() { // for tradeetMap
	tradeetMutex.Lock()
}
func UnlockTrade() {
	tradeetMutex.Unlock()
}

//------------------------------------------------------------------------------
func NewTradeSentinel() *TradeSentinel {
	return &TradeSentinel{
		SyncFlag:  false,
		Timestamp: time.Now(),
		RWLock:    &sync.RWMutex{},
	}
}

func UseTradeSentinel(sid int32) *TradeSentinel {
	LockTrade()
	defer UnlockTrade()

	t, ok := tradeetMap[sid]
	if ok {
		return t
	}
	t = NewTradeSentinel()
	tradeetMap[sid] = t
	return t
}

//------------------------------------------------------------------------------

//------------------------------------------------------------------------------
func GetLocalTime(year int, month int, day int, hour int, min int, sec int) time.Time {
	local, _ := time.LoadLocation("Local")
	v := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, min, sec)
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", v, local)
	return t
}

func ExpireAt(hour int, min int, sec int) time.Time {
	now := time.Now()
	nowhms := now.Hour()*10000 + now.Minute()*100 + now.Second()
	ttlhms := hour*10000 + min*100 + sec
	stop := now

	if nowhms >= ttlhms {
		stop = stop.AddDate(0, 0, 1)
	}

	local, _ := time.LoadLocation("Local")
	v := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", stop.Year(), int(stop.Month()), stop.Day(), hour, min, sec)
	stop, _ = time.ParseInLocation("2006-01-02 15:04:05", v, local)
	return stop
}

func (this *TraceRecord) SyncAndGetTradeRecords(start int, stop int) error {
	if this.Sentinel.SyncFlag {
		logging.Info("Running in sync, use RedisStore")
		return this.GetTradeRecordsFrom(RedisStore, start, stop)
	}

	this.Sentinel.RWLock.RLock()
	defer this.Sentinel.RWLock.RUnlock()

	clen, err := RedisCache.Llen(this.Key)
	if err != nil && err == hsgrr.ErrNil {
		// 缓存出错不可用
		logging.Info("RedisCache error, use RedisStore")
		return this.GetTradeRecordsFrom(RedisStore, start, stop)
	}

	slen, err := RedisStore.Llen(this.Key)
	if err != nil {
		return err
	}

	clear := false
	if slen == clen {
		logging.Info("RedisStore and RedisCache is the same, use RedisCache")
		return this.GetTradeRecordsFrom(RedisCache, start, stop)
	} else if slen < clen {
		clear = true
		clen = 0
	}

	/// 获取数据Redis剩余的用以同步
	if clen < 0 {
		clen = 0
	}
	ls, serr := RedisStore.LRange(this.Key, clen, slen-1)
	if serr != nil {
		return serr
	}

	go this.SafeSyncTraceRecord(clear, ls)

	if stop == -1 {
		// 到结束
		stop = slen // redis 包含stop, slice 不包含stop
	}
	n := start
	m := stop
	c := clen
	s := slen

	switch {
	case n < c && m < c:
		goto from_cache
	case n >= c && m >= c: // from ls
		goto from_store
	case n < c && m >= c:
		goto merge_cast
	default:
		logging.Fatal("exception condition...")
		return this.GetTradeRecordsFrom(RedisStore, start, stop)
	}

	return nil

from_cache:
	{
		// 从缓存Redis取数据
		logging.Info("Get from Cache list[n m]=[%d %d], clen=%d, slen=%d\n", n, m, c, s)
		cls, err := RedisCache.LRange(this.Key, start, stop) //redis list: lrange key n m包含m,区别于go slice[n:m]不包含m
		if err != nil {
			return err
		}
		this.Total = int32(clen)
		return this.StringsToTradeRecordss(cls)
	}
from_store:
	{
		// 从数据Redis取数据, 复用ls切片
		logging.Info("Get from Store list[n m]=[%d %d], clen=%d, slen=%d\n", n, m, c, s)
		n = n - c
		m = m - c + 1 // 切片的后置下标不包含m,故+1
		if m > len(ls) {
			m = len(ls)
		}
		cls := ls[n:m]           // 切片的后置下标不包含m
		this.Total = int32(slen) // 本次请求时交易总条数
		return this.StringsToTradeRecordss(cls)
	}
merge_cast:
	{
		// 需要从数据Redis和缓存Redis各取一部分数据
		logging.Info("Get merge CaSt list[n m]=[%d %d], clen=%d, slen=%d\n", n, m, c, s)
		cls, err := RedisCache.LRange(this.Key, n, c-1)
		if err != nil {
			return err
		}
		n = 0
		m = m - c + 1
		if m > len(ls) {
			m = len(ls)
		}
		sls := ls[n:m]
		cls = append(cls, sls...)
		this.Total = int32(slen)
		return this.StringsToTradeRecordss(cls)
	}
	return nil
}

func (this *TraceRecord) StringsToTradeRecordss(ls []string) error {
	rows := make([]*pro.TradeEveryTimeRecord, 0, 5000)
	for _, v := range ls {
		trade, err := this.StringToTradeRecords(&v)
		if err != nil {
			return err
		}
		rows = append(rows, trade)
	}
	this.List = rows
	return nil
}

func (this *TraceRecord) StringToTradeRecords(v *string) (*pro.TradeEveryTimeRecord, error) {
	trade := &pro.TradeEveryTimeRecord{}
	bufer := bytes.NewBuffer([]byte(*v))
	if err := binary.Read(bufer, binary.LittleEndian, trade); err != nil && err != io.EOF {
		logging.Error("%v", err)
		return nil, err
	}
	return trade, nil
}

func (this TraceRecord) SyncTradeRecords(clear bool, ls []string) error {
	this.Sentinel.SyncFlag = true
	defer func() {
		this.Sentinel.SyncFlag = false
	}()
	if clear {
		return RedisCache.Del(this.Key)
	}
	start := time.Now()
	logging.Info("%T >>> sync %s start... ", this, this.Key)
	RedisCache.Do("MULTI", "")
	for _, v := range ls {
		RedisCache.Rpush(this.Key, []byte(v))
	}
	RedisCache.Do("EXPIREAT", this.Key, ExpireAt(9, 25, 0).Unix()) // 缓存Redis this.Key设置9:25自动删除
	RedisCache.Do("EXEC", "")
	this.Sentinel.Timestamp = time.Now()
	logging.Info("%T >>> sync %s finish, rows %d took %v", this, this.Key, len(ls), this.Sentinel.Timestamp.Sub(start))

	_, err := RedisCache.Llen(this.Key)
	if err != nil {
		logging.Warning("sync trade: %v", err)
		return err
	}
	return nil
}
func (this TraceRecord) SafeSyncTraceRecord(clear bool, ls []string) error {
	this.Sentinel.RWLock.Lock()
	defer this.Sentinel.RWLock.Unlock()
	return this.SyncTradeRecords(clear, ls)
}

func (this *TraceRecord) GetTradeRecordsFrom(r *redis.RedisPool, start int, stop int) error {
	total, err := r.Llen(this.Key)
	if err != nil {
		logging.Error("%v", err)
		return err
	}
	ls, err := r.LRange(this.Key, start, stop)
	if err != nil {
		logging.Error("%v", err)
		return err
	}
	rows := make([]*pro.TradeEveryTimeRecord, 0, 5000)
	for _, v := range ls {
		trade := &pro.TradeEveryTimeRecord{}
		bufer := bytes.NewBuffer([]byte(v))
		if err := binary.Read(bufer, binary.LittleEndian, trade); err != nil && err != io.EOF {
			logging.Error("%v", err)
			return err
		}
		rows = append(rows, trade)
	}
	this.List = rows
	this.Total = int32(total)
	return nil
}
func (this *TraceRecord) SafeGetTradeRecordsFrom(r *redis.RedisPool, start int, stop int) error {
	this.Sentinel.RWLock.RLock()
	defer this.Sentinel.RWLock.RUnlock()
	return this.SafeGetTradeRecordsFrom(r, start, stop)
}

//==============================================================================
// 存PB版
type TraceRecordObj struct {
	Sid      int32
	Key      string
	Sentinel *TradeSentinel
	Total    int32                      // 返回的当前交易数量
	Obj      *pro.PayloadTradeEveryTime // 返回的成交记录
}

func NewTraceRecordObj(sid int32) *TraceRecordObj {
	t := UseTradeSentinel(sid)
	return &TraceRecordObj{
		Sid:      sid,
		Key:      fmt.Sprintf(TradeCacheKey, sid),
		Sentinel: t,
	}
}
