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

//A股市场代码表
func (this *SecurityNameTable) GetSecurityTableAStock() (*protocol.PayloadMarketSecurityNameTable, error) {
	return this.getSecurityTableAStockFromeCache()
}

func (this *SecurityNameTable) getSecurityTableAStockFromeCache() (*protocol.PayloadMarketSecurityNameTable, error) {
	key := fmt.Sprintf(publish.REDISKEY_MARKET_SECURITY_TABLE_ASTOCK, "astock")
	var market = &protocol.PayloadMarketSecurityNameTable{}

	bs, err := GetCache(key)
	if err != nil {
		key_sh := fmt.Sprintf(this.CacheKey, 100000000)
		key_sz := fmt.Sprintf(this.CacheKey, 200000000)

		var market_100 = &protocol.PayloadMarketSecurityNameTable{}
		var market_200 = &protocol.PayloadMarketSecurityNameTable{}

		if err = getSecurityTableFromeStore(key_sh, market_100); err != nil {
			return nil, err
		}

		if err = getSecurityTableFromeStore(key_sz, market_200); err != nil {
			return nil, err
		}

		market_100.MarketID = 0
		market_100.TimeStamp = market_100.TimeStamp
		market_100.Num = market_100.Num + market_200.Num

		for _, v := range market_200.SNList {
			market_100.SNList = append(market_100.SNList, v)
		}

		var marketNew *protocol.PayloadMarketSecurityNameTable
		marketNew = market_100

		if err = setSecurityTableToCache(key, marketNew); err != nil {
			logging.Error("%v", err.Error())
		}
		return marketNew, nil
	} else {
		if err = proto.Unmarshal(bs, market); err != nil {
			return nil, err
		}
	}
	return market, nil
}

//单市场股票代码表
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
