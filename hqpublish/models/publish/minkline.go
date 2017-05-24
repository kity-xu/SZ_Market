package publish

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"

	. "haina.com/market/hqpublish/models"
	. "haina.com/share/models"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	redigo "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

var (
	_ = fmt.Println
	_ = redigo.Bytes
	_ = logging.Info
	_ = bytes.NewBuffer
	_ = binary.Read
	_ = io.ReadFull
)

type MinKLine struct {
	Model `db:"-"`
}

const TTL_REDISKEY_SECURITY_MIN = 30 // 暂时放在，避免多人合并冲突 constants.go

func NewMinKLine() *MinKLine {
	return &MinKLine{
		Model: Model{
			CacheKey: REDISKEY_SECURITY_MIN,
		},
	}
}

// 分钟K线 缓存分钟K线 TTL 30秒

// 获取分钟K线
func (this MinKLine) GetMinKObj(req *protocol.RequestMinK) (*protocol.PayloadMinK, error) {
	key := fmt.Sprintf(this.CacheKey, req.SID)

	kls := make([]*protocol.KInfo, 0, 250)
	bs, err := GetCache(key)
	if err == nil {
		logging.Info("GetCache %s hit", key)
		ls, err := this.Decode(bs)
		if err != nil {
			logging.Error("----- %v", err)
			return nil, err
		}
		for _, k := range ls {
			fmt.Printf("--------- %+v\n", k)
			if k.NTime >= req.BeginTime {
				kls = append(kls, k)
			}
		}
		ret := protocol.PayloadMinK{
			SID:       req.SID,
			BeginTime: req.BeginTime,
			Num:       int32(len(kls)),
			KList:     kls,
		}
		logging.Info("GetCache %s hit pass", key)
		return &ret, nil
	}

	ls, err := redis.LRange(key, 0, -1)
	if err != nil {
		logging.Warning("1 %v", err)
		return nil, err
	}
	if len(ls) == 0 {
		logging.Warning("2")
		return nil, ERROR_KLINE_DATA_NULL
	}

	for _, v := range ls {
		k := &protocol.KInfo{}
		buffer := bytes.NewBuffer([]byte(v))
		if err := binary.Read(buffer, binary.LittleEndian, k); err != nil && err != io.EOF {
			return nil, err
		}
		if k.NTime >= req.BeginTime {
			kls = append(kls, k)
		}
	}
	cache, err := this.Encode(ls)
	if err != nil {
		return nil, err
	}
	SetCache(key, TTL_REDISKEY_SECURITY_MIN, cache)

	ret := protocol.PayloadMinK{
		SID:       req.SID,
		BeginTime: req.BeginTime,
		Num:       int32(len(kls)),
		KList:     kls,
	}

	return &ret, nil
}

func (this MinKLine) Encode(ss []string) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	count := 0
	for _, v := range ss {
		count += len(v)
		if err := binary.Write(buffer, binary.LittleEndian, []byte(v)); err != nil {
			return nil, err
		}
	}

	fmt.Println("count", count, "len", len(buffer.Bytes()))

	return buffer.Bytes(), nil
}

func (this MinKLine) Decode(bs []byte) ([]*protocol.KInfo, error) {
	one := protocol.KInfo{}
	siz := binary.Size(&one)
	num := len(bs) / siz

	fmt.Println("size", siz, "num", num)

	ls := make([]*protocol.KInfo, 0, 250)

	buffer := bytes.NewBuffer(bs)
	for i := 0; i < num; i++ {
		if err := binary.Read(buffer, binary.LittleEndian, &one); err != nil && err != io.EOF {
			return nil, err
		}
		ls = append(ls, &one)
		fmt.Printf("decode %+v\n", one)
	}

	return ls, nil
}

func (this MinKLine) GetMinKJson(req *protocol.RequestMinK) ([]byte, error) {
	key := fmt.Sprintf(this.CacheKey, req.SID)
	if req.BeginTime == 0 {
		bs, err := GetCacheJson(key)
		if err == nil {
			logging.Info("GetCacheJson %s hit", key)
			return bs, nil
		}
		logging.Info("GetCacheJson %s: %v", key, err)
	}

	obj, err := this.GetMinKObj(req)
	if err != nil {
		logging.Warning("3 %v", err)
		res := map[string]interface{}{"code": 40002}
		js, jerr := json.Marshal(&res)
		if jerr != nil {
			return nil, jerr
		}
		SetCacheJson(key, TTL_REDISKEY_SECURITY_MIN, js)
		return nil, err
	}
	res := map[string]interface{}{"code": 200}
	res["data"] = obj

	js, err := json.Marshal(&res)
	if err != nil {
		return nil, err
	}
	SetCacheJson(key, TTL_REDISKEY_SECURITY_MIN, js)
	return js, nil
}

// 获取分钟K线PB
func (this MinKLine) GetMinKPB(req *protocol.RequestMinK) ([]byte, error) {
	return nil, nil
}
