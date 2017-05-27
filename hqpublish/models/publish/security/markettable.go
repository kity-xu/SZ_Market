package security

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"
	"fmt"

	"github.com/golang/protobuf/proto"
	. "haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish"
	"haina.com/share/logging"
	. "haina.com/share/models"
)

type SecurityNameTable struct {
	Model `db:"-"`
}

func NewSecurityNameTable() *SecurityNameTable {
	return &SecurityNameTable{
		Model: Model{
			CacheKey: publish.REDISKEY_MARKET_SECURITY_TABLE,
		},
	}
}

func (this *SecurityNameTable) GetSecurityTable(req *protocol.RequestMarketSecurityNameTable) (*protocol.PayloadMarketSecurityNameTable, error) {
	return this.getSecurityTableFromeCache(req.MarketID)
}

func (this *SecurityNameTable) getSecurityTableFromeCache(MarketID int32) (*protocol.PayloadMarketSecurityNameTable, error) {
	key := fmt.Sprintf(this.CacheKey, MarketID)
	var market = &protocol.PayloadMarketSecurityNameTable{}

	bs, err := GetCache(key)
	if err != nil {
		if err = getSecurityTableFromeStore(key, market); err != nil {
			return nil, err
		}

		if err = setSecurityTableToCache(key, market); err != nil {
			logging.Error("%v", err.Error())
		}

	} else {
		if err = proto.Unmarshal(bs, market); err != nil {
			return nil, err
		}
	}
	return market, nil
}

func getSecurityTableFromeStore(key string, market *protocol.PayloadMarketSecurityNameTable) error {
	logging.Debug("security code table key:%v", key)
	bs, err := RedisStore.GetBytes(key)
	if err != nil {
		return err
	}
	logging.Debug("redis store len:%v", len(bs))

	if err = proto.Unmarshal(bs, market); err != nil {
		return err
	}
	return nil
}

func setSecurityTableToCache(key string, market *protocol.PayloadMarketSecurityNameTable) error {
	bs, err := proto.Marshal(market)
	if err != nil {
		return err
	}

	if err = SetCache(key, 60*5, bs); err != nil {
		return err
	}
	return nil
}
