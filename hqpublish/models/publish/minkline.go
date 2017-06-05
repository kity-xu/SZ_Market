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

type MinKLine struct {
	Model   `db:"-"`
	Compare compMinKline
}

type compMinKline func(k *pro.KInfo, req *pro.RequestMinK) bool

func compareMinKline(k *pro.KInfo, req *pro.RequestMinK) bool {
	if k.NTime%10000 >= req.BeginTime {
		return true
	} else {
		return false
	}
}

func NewMinKLine() *MinKLine {
	return &MinKLine{
		Model: Model{
			CacheKey: REDISKEY_SECURITY_MIN,
		},
		Compare: compareMinKline,
	}
}

////////////////////////////////////////////////////////////////////////////////
// 分钟K线 缓存 TTL 30秒
// BeginTime == 0 作为特例参数 取应答缓存 key.p key.j pb/json应答数据
// BeginTime != 0 非特例参数   取应答缓存 key 然后封装应答数据
////////////////////////////////////////////////////////////////////////////////
// 缓存规则：
//   key    从数据 redis 获取到的当前分钟K线完整列表 -> 编码 -> []byte
//   key.p  当前分钟K线的完整列表  pb  应答数据格式
//   key.j  当前分钟K线的完整列表 json 应答数据格式
////////////////////////////////////////////////////////////////////////////////

// 获取分钟K线JSON
// 由于涉及到应答结果缓存，这里直接返回用于应答给客户端的 json 数据
func (this MinKLine) GetMinKJson(req *pro.RequestMinK) ([]byte, error) {
	key := fmt.Sprintf(this.CacheKey, req.SID)

	// 全部
	if req.BeginTime == 0 {
		if bs, err := GetCacheJson(key); err == nil {
			return bs, nil
		}
	}

	cache, store, err := this.GetMinKObj(req)
	if cache != nil {
		j, err := ctrl.MakeRespJson(200, cache)
		if err != nil {
			return nil, err
		}
		return j, nil
	}
	if store != nil {
		go this.SaveToCache(key, store)
		j, err := ctrl.MakeRespJson(200, this.NewPayloadMinK(req, store.KList))
		if err != nil {
			return nil, err
		}
		return j, nil
	}

	go this.SaveToCache(key, nil)
	return nil, err
}

// 获取分钟K线PB
// 返回直接用于应答给客户端的Json
// 由于涉及到应答结果缓存，这里直接返回用于应答给客户端的 Header+Payload 数据
func (this MinKLine) GetMinKPB(req *pro.RequestMinK) ([]byte, error) {
	key := fmt.Sprintf(this.CacheKey, req.SID)

	// 全部
	if req.BeginTime == 0 {
		if bs, err := GetCachePB(key); err == nil {
			return bs, nil
		}
	}

	cache, store, err := this.GetMinKObj(req)
	if cache != nil {
		j, err := ctrl.MakeRespDataByPB(200, pro.HAINA_PUBLISH_CMD_ACK_MINKLINE, cache)
		if err != nil {
			return nil, err
		}
		return j, nil
	}
	if store != nil {
		go this.SaveToCache(key, store)
		p, err := ctrl.MakeRespDataByPB(200, pro.HAINA_PUBLISH_CMD_ACK_MINKLINE, this.NewPayloadMinK(req, store.KList))
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
func (this MinKLine) GetMinKObj(req *pro.RequestMinK) (*pro.PayloadMinK, *pro.PayloadMinK, error) {
	obj, err := this.GetCacheMinKObj(req)
	if err == nil {
		return obj, nil, nil
	}
	obj, err = this.GetStoreMinKObj(req)
	if err == nil {
		return nil, obj, nil
	}
	return nil, nil, err
}

// 从缓存中获取分钟K线 PayloadMinK 对象
func (this MinKLine) GetCacheMinKObj(req *pro.RequestMinK) (*pro.PayloadMinK, error) {
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

	kls := make([]*pro.KInfo, 0, 250)
	for _, k := range ls {
		if this.Compare(k, req) {
			kls = append(kls, k)
		}
	}
	return &pro.PayloadMinK{
		SID:       req.SID,
		BeginTime: req.BeginTime,
		Num:       int32(len(kls)),
		KList:     kls,
	}, nil
}

// 返回为当前全部分钟K线的 PayloadMinK 对象，用以后续做缓存
func (this MinKLine) GetStoreMinKObj(req *pro.RequestMinK) (*pro.PayloadMinK, error) {
	key := fmt.Sprintf(this.CacheKey, req.SID)

	kls := make([]*pro.KInfo, 0, 250)
	ls, err := RedisStore.LRange(key, 0, -1)
	if err != nil {
		logging.Warning("1 %v", err)
		return nil, err
	}
	if len(ls) == 0 {
		logging.Warning("redis no such %s", key)
		return nil, ERROR_KLINE_DATA_NULL
	}

	// 当日分钟K线每条为binary编码
	for _, v := range ls {
		k := &pro.KInfo{}
		buffer := bytes.NewBuffer([]byte(v))
		if err := binary.Read(buffer, binary.LittleEndian, k); err != nil && err != io.EOF {
			return nil, err
		}
		kls = append(kls, k)
	}

	return &pro.PayloadMinK{
		SID:       req.SID,
		BeginTime: 0,
		Num:       int32(len(kls)),
		KList:     kls,
	}, nil
}

// 存放到Cache前进行编码
func (this MinKLine) Encode(klist []*pro.KInfo) ([]byte, error) {
	if klist == nil {
		return nil, ERROR_REDIS_LIST_NULL
	}
	buffer := bytes.NewBuffer(nil)
	for _, k := range klist {
		if err := binary.Write(buffer, binary.LittleEndian, k); err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}

// 从Cache取出后进行解码
func (this MinKLine) Decode(bs []byte) ([]*pro.KInfo, error) {
	if bs == nil {
		return nil, ERROR_KLINE_DATA_NULL
	}
	one := pro.KInfo{}
	siz := binary.Size(&one)
	num := len(bs) / siz

	ls := make([]*pro.KInfo, 0, 250)

	buffer := bytes.NewBuffer(bs)
	for i := 0; i < num; i++ {
		obj := pro.KInfo{}
		if err := binary.Read(buffer, binary.LittleEndian, &obj); err != nil && err != io.EOF {
			return nil, err
		}
		ls = append(ls, &obj)
	}
	return ls, nil
}

func (this MinKLine) SaveToCache(key string, obj *pro.PayloadMinK) {
	if obj == nil {
		if p, err := ctrl.MakeRespDataByPB(40002, 0, nil); err == nil {
			SetCachePB(key, TTL_REDISKEY_SECURITY_MIN, p)
		}
		if j, err := ctrl.MakeRespJson(40002, nil); err == nil {
			SetCacheJson(key, TTL_REDISKEY_SECURITY_MIN, j)
		}
	} else {
		if p, err := ctrl.MakeRespDataByPB(200, pro.HAINA_PUBLISH_CMD_ACK_MINKLINE, obj); err == nil {
			SetCachePB(key, TTL_REDISKEY_SECURITY_MIN, p)
		}
		if j, err := ctrl.MakeRespJson(200, obj); err == nil {
			SetCacheJson(key, TTL_REDISKEY_SECURITY_MIN, j)
		}
	}
	// origin
	if o, err := this.Encode(obj.GetKList()); err == nil {
		SetCache(key, TTL_REDISKEY_SECURITY_MIN, o)
	}
}

func (this MinKLine) NewPayloadMinK(req *pro.RequestMinK, klist []*pro.KInfo) *pro.PayloadMinK {
	if req == nil || klist == nil {
		return nil
	}
	kls := make([]*pro.KInfo, 0, 250)
	for _, k := range klist {
		if this.Compare(k, req) {
			kls = append(kls, k)
		}
	}
	return &pro.PayloadMinK{
		SID:       req.SID,
		BeginTime: req.BeginTime,
		Num:       int32(len(kls)),
		KList:     kls,
	}
}
