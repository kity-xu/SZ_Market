package publish

import (
	"bytes"
	"encoding/binary"
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

//const REDISKEY_SECURITY_TRADE = "hq:st:trade:%d" ///<证券分笔成交数据(参数：sid) (calc写入)

func NewTradeEveryTime() *TradeEveryTime {
	return &TradeEveryTime{
		Model: Model{
			CacheKey: REDISKEY_SECURITY_TRADE,
		},
	}
}

func (this TradeEveryTime) GetTradeEveryTimeJson(req *pro.RequestTradeEveryTime) ([]byte, error) {
	key := fmt.Sprintf(this.CacheKey, req.SID)

	// 全部
	if req.Begin == 0 && req.Num == 0 {
		if bs, err := GetCacheJson(key); err == nil {
			return bs, nil
		}
	}

	cache, store, err := this.GetTradeEveryTimeObj(req)
	if cache != nil {
		j, err := ctrl.MakeRespJson(200, cache)
		if err != nil {
			return nil, err
		}
		return j, nil
	}
	if store != nil {
		go this.SaveToCache(key, store)
		j, err := ctrl.MakeRespJson(200, this.NewPayloadTradeEveryTime(req, nil))
		if err != nil {
			return nil, err
		}
		return j, nil
	}

	go this.SaveToCache(key, nil)
	return nil, err
}

func (this TradeEveryTime) GetTradeEveryTimePB(req *pro.RequestTradeEveryTime) ([]byte, error) {
	key := fmt.Sprintf(this.CacheKey, req.SID)

	// 全部
	if req.Num == 0 {
		if bs, err := GetCachePB(key); err == nil {
			return bs, nil
		}
	}

	cache, store, err := this.GetTradeEveryTimeObj(req)
	if cache != nil {
		j, err := ctrl.MakeRespDataByPB(200, pro.HAINA_PUBLISH_CMD_ACK_TRADEDT, cache)
		if err != nil {
			return nil, err
		}
		return j, nil
	}
	if store != nil {
		go this.SaveToCache(key, store)
		p, err := ctrl.MakeRespDataByPB(200, pro.HAINA_PUBLISH_CMD_ACK_TRADEDT, this.NewPayloadTradeEveryTime(req, nil))
		if err != nil {
			return nil, err
		}
		return p, nil
	}

	go this.SaveToCache(key, nil)
	return nil, err
}

// 第一个返回参数：从缓存Redis里拿到的符合条件的应答Payload对象
// 第二个返回参数：从数据Redis里拿到的所有分钟K线Payload对象(后续直接用于缓存)
func (this TradeEveryTime) GetTradeEveryTimeObj(req *pro.RequestTradeEveryTime) (*pro.PayloadTradeEveryTime, *pro.PayloadTradeEveryTime, error) {
	obj, err := this.GetCacheTradeEveryTimeObj(req)
	if err == nil {
		return obj, nil, nil
	}
	obj, err = this.GetStoreTradeEveryTimeObj(req)
	if err == nil {
		return nil, obj, nil
	}
	return nil, nil, err
}

// 从缓存中获取 PayloadTradeEveryTime 对象
func (this TradeEveryTime) GetCacheTradeEveryTimeObj(req *pro.RequestTradeEveryTime) (*pro.PayloadTradeEveryTime, error) {
	key := fmt.Sprintf(this.CacheKey, req.SID)
	bs, err := GetCache(key)
	if err != nil {
		return nil, err
	}

	// cache hit
	ls, err := this.Decode(bs)
	if err != nil {
		return nil, err
	}

	kls := make([]*pro.TradeEveryTimeRecord, 0, 250)
	for _, k := range ls {
		if k.NSn >= uint32(req.Begin) {
			kls = append(kls, k)
		}
	}
	return &pro.PayloadTradeEveryTime{
		Num: 0,
	}, nil
}

func (this TradeEveryTime) GetStoreTradeEveryTimeObj(req *pro.RequestTradeEveryTime) (*pro.PayloadTradeEveryTime, error) {
	key := fmt.Sprintf(this.CacheKey, req.SID)

	var ls []string

	kls := make([]*pro.TradeEveryTimeRecord, 0, 5000)
	slen, err := RedisStore.Llen(key)
	clen, err := RedisCache.Llen(key)
	fmt.Printf("%s Store len %d, Cache len %d\n", key, slen, clen)

	if clen < 0 {
		ls, err = RedisStore.LRange(key, 0, -1)
		if err == nil {
			for _, v := range ls {
				RedisCache.Rpush(key, []byte(v))
			}
		}
	} else if clen < slen {
		ls, err = RedisStore.LRange(key, clen, slen)
		if err == nil {
			for _, v := range ls {
				RedisCache.Rpush(key, []byte(v))
			}
		}
	}
	slen, err = RedisStore.Llen(key)
	clen, err = RedisCache.Llen(key)
	fmt.Printf("%s Store len %d, Cache len %d\n", key, slen, clen)

	//	RedisCache.Do("expire", key, 10000)
	if err != nil {
		logging.Warning("1 %v", err)
		return nil, err
	}
	if len(ls) == 0 {
		logging.Warning("redis no such %s", key)
		return nil, ERROR_INVALID_DATA
	}

	for _, v := range ls {
		k := &pro.TradeEveryTimeRecord{}
		buffer := bytes.NewBuffer([]byte(v))
		if err := binary.Read(buffer, binary.LittleEndian, k); err != nil && err != io.EOF {
			return nil, err
		}
		kls = append(kls, k)
	}

	return &pro.PayloadTradeEveryTime{
		SID:   req.SID,
		Begin: 0,
		Num:   int32(len(kls)),
	}, nil
}

// 存放到Cache前进行编码
func (this TradeEveryTime) Encode(klist []*pro.TradeEveryTimeRecord) ([]byte, error) {
	return nil, nil
}

// 从Cache取出后进行解码
func (this TradeEveryTime) Decode(bs []byte) ([]*pro.TradeEveryTimeRecord, error) {
	return nil, nil
}

func (this TradeEveryTime) SaveToCache(key string, obj *pro.PayloadTradeEveryTime) {
}

func (this TradeEveryTime) NewPayloadTradeEveryTime(req *pro.RequestTradeEveryTime, klist []*pro.TradeEveryTimeRecord) *pro.PayloadTradeEveryTime {
	if req == nil || klist == nil {
		return nil
	}
	kls := make([]*pro.TradeEveryTimeRecord, 0, 5000)
	for _, k := range klist {
		if k.NSn >= uint32(req.Begin) {
			kls = append(kls, k)
		}
	}
	return &pro.PayloadTradeEveryTime{
		Num: 0,
	}
}
