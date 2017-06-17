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
// 读写锁版
// RedisCache 缓存成交记录list(和RedisStore一致)
type TradeSentinel struct {
	Timestamp time.Time     // 最后同步时间戳
	SyncFlag  bool          // true/false:正在同步/完成同步
	RWLock    *sync.RWMutex // 读写锁
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
// 将分笔成交数据同步到缓存
// 返回参数
//  第一组 int err  缓存 list 长度 和 Redis 错误
//  第二组 int err  数据 list 长度 和 Redis 错误
// 返回值说明 int
//   -1  错误（如果缓存redis不可用，后续操作跳过缓存redis)
//    0  this.Key-list长度为0 或 不存在该 this.Key
//   >0  this.Key-list长度值
func SyncTradeEveryTimeRecord(sid int32) (int, error, int, error) {
	key := fmt.Sprintf(TradeCacheKey, sid)
	t := UseTradeSentinel(sid)
	slen, serr := RedisStore.Llen(key)
	if t.SyncFlag {
		fmt.Printf("%s in sync\n", key)
		return 0, nil, slen, serr
	} else {
		fmt.Printf("%s no sync\n", key)
	}

	t.RWLock.RLock()
	defer t.RWLock.RUnlock()

	clen, cerr := RedisCache.Llen(key) // Llen 出错时 返回 clen = -1
	switch {
	case cerr != nil && cerr != hsgrr.ErrNil: // 因缓存Redis不可用，跳过同步
		logging.Error("RedisCache: %v", cerr)
		return clen, cerr, slen, serr
	case serr != nil: // 没有数据或数据源Redis不可用
		logging.Error("RedisStore: %v", serr)
		return clen, cerr, slen, serr
	}
	fmt.Printf("%s Store len %d, Cache len %d\n", key, slen, clen)

	if slen == clen {
		return clen, nil, slen, nil
	}

	// 同步 RedisStore 到 RedisCache
	now := time.Now()
	del := false // 缓存删除标志
	nowhm := now.Hour()*100 + now.Minute()

	switch {
	case slen < clen:
		fallthrough
	case 900 < nowhm && nowhm < 930: // 从9:30以后开始缓存
		del = true // 清除缓存
		clen = 0   // 全部同步
		break
	}

	// 9:00 ~ 9:25 之间，以数据Redis为准
	if nowhm > 900 && nowhm < 925 {
		return 0, nil, slen, nil
	}

	if clen < 0 {
		clen = 0
	}
	ls, serr := RedisStore.LRange(key, clen, slen-1)
	if serr != nil {
		return clen, nil, slen, serr
	}

	t.SyncFlag = true

	go func() {
		logging.Info(">>> ready t. sync %s ... ", key)
		t.RWLock.Lock()
		defer func() {
			t.SyncFlag = false
			t.RWLock.Unlock()
		}()

		if del {
			RedisCache.Del(key)
		}

		start := time.Now()
		logging.Info(">>> sync %s start... ", key)
		RedisCache.Do("MULTI", "")
		for _, v := range ls {
			RedisCache.Rpush(key, []byte(v))
		}
		RedisCache.Do("EXPIREAT", key, ExpireAt(9, 25, 0).Unix()) // 缓存Redis key设置9:25自动删除
		RedisCache.Do("EXEC", "")
		t.Timestamp = time.Now()
		logging.Info(">>> sync %s finish, rows %d took %v", key, slen-clen, t.Timestamp.Sub(start))
	}()

	return clen, nil, slen, nil
}

//------------------------------------------------------------------------------
// 老函数
func (this *TraceRecord) SyncTradeEveryTimeRecord(sid int32) (int, error, int, error) {
	slen, serr := RedisStore.Llen(this.Key)
	if this.Sentinel.SyncFlag {
		fmt.Printf("%s in sync\n", this.Key)
		return 0, nil, slen, serr
	} else {
		fmt.Printf("%s no sync\n", this.Key)
	}

	this.Sentinel.RWLock.RLock()
	defer this.Sentinel.RWLock.RUnlock()

	// Llen 出错时 返回 clen = -1
	clen, cerr := RedisCache.Llen(this.Key)
	switch {
	case cerr != nil && cerr != hsgrr.ErrNil:
		// 因缓存Redis不可用，跳过同步
		logging.Error("RedisCache: %v", cerr)
		return clen, cerr, slen, serr
	case serr != nil:
		// 没有数据或数据源Redis不可用
		logging.Error("RedisStore: %v", serr)
		return clen, cerr, slen, serr
	}
	fmt.Printf("%s Store len %d, Cache len %d\n", this.Key, slen, clen)

	if slen == clen {
		return clen, nil, slen, nil
	}

	// 同步 RedisStore 到 RedisCache
	now := time.Now()
	del := false                           // 缓存删除标志
	nowhm := now.Hour()*100 + now.Minute() // 当前时间时分形式 例如 9:25 925

	switch {
	case slen < clen:
		fallthrough
	case 900 < nowhm && nowhm < 930: // 从9:30以后开始缓存
		del = true // 清除缓存
		clen = 0   // 同步所有
		break
	}

	// 9:00 ~ 9:25 之间，以数据Redis为准
	if nowhm > 900 && nowhm < 925 {
		return 0, nil, slen, nil
	}

	if clen < 0 {
		clen = 0
	}
	ls, serr := RedisStore.LRange(this.Key, clen, slen-1)
	if serr != nil {
		return clen, nil, slen, serr
	}

	this.Sentinel.SyncFlag = true

	go func() {
		logging.Info(">>> ready this.Sentinel. sync %s ... ", this.Key)
		this.Sentinel.RWLock.Lock()
		defer func() {
			this.Sentinel.SyncFlag = false
			this.Sentinel.RWLock.Unlock()
		}()

		if del {
			RedisCache.Del(this.Key)
		}

		start := time.Now()
		logging.Info(">>> sync %s start... ", this.Key)
		RedisCache.Do("MULTI", "")
		for _, v := range ls {
			RedisCache.Rpush(this.Key, []byte(v))
		}
		RedisCache.Do("EXPIREAT", this.Key, ExpireAt(9, 25, 0).Unix()) // 缓存Redis key设置9:25自动删除
		RedisCache.Do("EXEC", "")
		this.Sentinel.Timestamp = time.Now()
		logging.Info(">>> sync %s finish, rows %d took %v", this.Key, slen-clen, this.Sentinel.Timestamp.Sub(start))
	}()

	return clen, nil, slen, nil
}

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
	this.Sentinel.RWLock.RLock()
	defer this.Sentinel.RWLock.RUnlock()

	clen, err := RedisCache.Llen(this.Key)
	if err != nil && err == hsgrr.ErrNil {
		// 缓存出错不可用
		return this.GetTradeRecordsFrom(RedisStore, start, stop)
	}

	slen, err := RedisStore.Llen(this.Key)
	if err != nil {
		return err
	}

	if slen == clen {
		return this.GetTradeRecordsFrom(RedisCache, start, stop)
	} else if slen < clen {
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
	fmt.Println(len(ls))

	if slen < clen {
		go this.SafeSyncTraceRecord(true, ls)
	} else {
		go this.SafeSyncTraceRecord(false, ls)
	}

	logging.Info("mark !!!!!!!!!!!!!!!!!!!!!!!!")
	if stop == -1 { // 到结束
		if start <= clen {
			logging.Info("mark !!!!!!!!!!!!!!!!!!!!!!!!")
			goto merge
		} else {
			logging.Info("mark !!!!!!!!!!!!!!!!!!!!!!!!")
			goto only_from_cache
		}
	} else {
		if stop <= clen {
			logging.Info("mark !!!!!!!!!!!!!!!!!!!!!!!!")
			goto only_from_cache
		} else if start < clen && clen < stop {
			logging.Info("mark !!!!!!!!!!!!!!!!!!!!!!!!")
			goto merge
		} else if start > clen {
			logging.Info("mark !!!!!!!!!!!!!!!!!!!!!!!!")
			goto only_from_store
		} else {
			// 此处如何写，思索中
		}
	}
	logging.Info("mark !!!!!!!!!!!!!!!!!!!!!!!!")

	return nil

only_from_cache:
	{
		cls, err := RedisCache.LRange(this.Key, start, stop)
		if err != nil {
			return err
		}
		this.Total = int32(clen)
		return this.StringsToTradeRecordss(cls)
	}
only_from_store:
	{
		//cls, err := RedisStore.LRange(this.Key, start, stop)
		cls := ls[start:stop]
		this.Total = int32(slen)
		return this.StringsToTradeRecordss(cls)
	}
merge:
	{
		// 需要从数据Redis和缓存Redis各取一部分数据
		cls, err := RedisCache.LRange(this.Key, start, clen)
		if err != nil {
			return err
		}
		sls, err := RedisStore.LRange(this.Key, clen+1, stop)
		if err != nil {
			return err
		}
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

func (this TraceRecord) SyncTradeRecords(del bool, ls []string) error {
	if del {
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
		return err
	}
	return nil
}
func (this TraceRecord) SafeSyncTraceRecord(del bool, ls []string) error {
	this.Sentinel.RWLock.Lock()
	defer this.Sentinel.RWLock.Unlock()
	return this.SyncTradeRecords(del, ls)
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
