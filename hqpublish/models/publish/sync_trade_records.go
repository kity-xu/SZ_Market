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

	//	pro "ProtocolBuffer/projects/hqpublish/go/protocol"

	"time"

	hsgrr "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
	//	"haina.com/share/store/redis"
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

//------------------------------------------------------------------------------
// 读写锁版
// RedisCache 缓存成交记录list(和RedisStore一致)
type TradeetSentinel struct {
	Timestamp time.Time     // 最后同步时间戳
	SyncFlag  bool          // true/false:正在同步/完成同步
	RWLock    *sync.RWMutex // 读写锁
}

type TradeetMap map[int32]*TradeetSentinel

//------------------------------------------------------------------------------
var (
	tradeetMap   TradeetMap = make(TradeetMap, 5000)
	tradeetMutex sync.Mutex
)

func LockTradeet() { // for tradeetMap
	tradeetMutex.Lock()
}
func UnlockTradeet() {
	tradeetMutex.Unlock()
}

//------------------------------------------------------------------------------
func NewTradeetSentinel() *TradeetSentinel {
	return &TradeetSentinel{
		SyncFlag:  false,
		Timestamp: time.Now(),
		RWLock:    &sync.RWMutex{},
	}
}

func UseTradeetSentinel(sid int32) *TradeetSentinel {
	LockTradeet()
	defer UnlockTradeet()

	t, ok := tradeetMap[sid]
	if ok {
		return t
	}
	t = NewTradeetSentinel()
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
//    0  key-list长度为0 或 不存在该 key
//   >0  key-list长度值
func SyncTradeEveryTimeRecord(sid int32) (int, error, int, error) {
	key := fmt.Sprintf(TradeCacheKey, sid)

	t := UseTradeetSentinel(sid)
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

	// 9:00 ~ 9:25 之间，以走数据Redis为准
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
		t.RWLock.Lock()
		defer func() {
			t.SyncFlag = false
			t.RWLock.Unlock()
		}()

		if del {
			RedisCache.Del(key)
		}

		stop := now

		// 计算TTL: 当前时间到下一个9:25之间的秒数
		if nowhm >= 925 {
			stop = stop.AddDate(0, 0, 1)
		}
		stopstr := fmt.Sprintf("%04d-%02d-%02d %02d:%02d", stop.Year(), int(stop.Month()), stop.Day(), 9, 25)
		ttlstop, _ := time.ParseInLocation("2006-01-02 15:04", stopstr, local)
		ttls := ttlstop.Unix()

		begin := time.Now()
		logging.Info(">>> sync %s start... ", key)
		RedisCache.Do("MULTI", "")
		for _, v := range ls {
			RedisCache.Rpush(key, []byte(v))
		}
		RedisCache.Do("EXPIREAT", key, ttls) // 缓存Redis TTL设置 下一个9:25自动删除
		RedisCache.Do("EXEC", "")
		t.Timestamp = time.Now()
		logging.Info(">>> sync %s finish, rows %d took %v", key, slen-clen, t.Timestamp.Sub(begin))
	}()

	return clen, nil, slen, nil
}
