package publish

import (
	protocol "ProtocolBuffer/projects/hqpublish/go/protocol"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	. "haina.com/market/hqpublish/models"
	. "haina.com/share/models"

	"github.com/golang/protobuf/proto"
	hsgrr "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

var (
	_ = fmt.Print
	_ = logging.Error
	_ = proto.Marshal
	_ = protocol.QuoteRecord{}
	_ = hsgrr.Dial
	_ = redis.NewRedisPool
	_ = GetStore
	_ = bytes.NewBuffer
	_ = io.ReadFull
)

type MarketStatus struct {
	Model `db:"-"`
}

func NewMarketStatus() *MarketStatus {
	return &MarketStatus{
		Model: Model{
			CacheKey: REDISKEY_MARKET_STATUS,
		},
	}
}

////////////////////////////////////////////////////////////////////////////////
// 市场状态
// 有多个市场，根据请求的参数不同，返回的内容组合也不同，如果将应答缓存下来,其内容不具有唯一性
// 所以合理做法是将数据的每个原始市场状态信息缓存下来，使用时就近取用，灵活组合
////////////////////////////////////////////////////////////////////////////////

// 返回的 PayloadXXX 类型是返回客户端的 json -> data item
func (this MarketStatus) GetPayloadObj(req *protocol.RequestMarketStatus) (*protocol.PayloadMarketStatus, error) {
	var ret = protocol.PayloadMarketStatus{
		Num:    req.Num,
		MSList: make([]*protocol.MarketStatus, 0, req.Num),
	}
	for _, v := range req.MarketIDList {
		single, err := this.GetSingle(v)
		if err != nil {
			return nil, err
		}
		ret.MSList = append(ret.MSList, single)
	}
	return &ret, nil
}

// 返回的[]byte是用于返回给客户端的 Payload(PB)
func (this MarketStatus) GetPayloadPB(req *protocol.RequestMarketStatus) ([]byte, error) {
	obj, err := this.GetPayloadObj(req)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(obj)
}

// MarketStatus在Redis里是以二进制形式存储
func (this MarketStatus) Decode(bin []byte) (*protocol.MarketStatus, error) {
	var obj protocol.MarketStatus
	buffer := bytes.NewBuffer([]byte(bin))
	if err := binary.Read(buffer, binary.LittleEndian, &obj); err != nil && err != io.EOF {
		return nil, err
	}
	return &obj, nil
}

//------------------------------------------------------------------------------
// 取单个市场状态
func (this MarketStatus) GetSingle(mid int32) (*protocol.MarketStatus, error) {
	key := fmt.Sprintf(this.CacheKey, mid)
	bin, err := GetCache(key)
	if err == nil {
		return this.Decode(bin)
	}
	logging.Info("GetCache %s: %v", key, err)

	// 如果请求者提供非法的市场ID，那么数据Redis里一定查不到该ID的数据
	// 如果将没找到也当成错误返回，那么已查到的市场状态数据也将会被丢弃  如
	//   if err != nil { return nil, err }
	// 如果不将其没找到当成错误处理，那么客户端会得到能查询到的那部分市场状态信息 如
	//   if err != nil && err != hsgrr.ErrNil { return nil, err }
	bin, err = GetStore(key)
	if err != nil && err != hsgrr.ErrNil {
		return nil, err
	}

	if err == nil {
		SetCache(key, TTL_REDISKEY_MARKETSTATUS, bin)
	}
	return this.Decode(bin)
}
