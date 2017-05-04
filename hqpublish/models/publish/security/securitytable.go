package security

import (
	"fmt"

	. "haina.com/share/models"

	"ProtocolBuffer/format/securitytable"

	"github.com/golang/protobuf/proto"

	"haina.com/market/hqpublish/models/publish"

	redigo "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

var _ = fmt.Println
var _ = redigo.Bytes
var _ = logging.Info

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

func (this *SecurityNameTable) GetSecurityTable(MarketID int32) (*securitytable.MarketSecurityCodeTable, int, error) {
	var market securitytable.MarketSecurityCodeTable

	key := fmt.Sprintf(this.CacheKey, MarketID)
	logging.Debug("key %v", key)

	str, err := redis.Get(key)
	if err != nil {
		return nil, 0, err
	}

	if str == "" {
		return nil, 0, publish.ERROR_REDIS_LIST_NULL
	}
	data := []byte(str)

	if err := proto.Unmarshal(data, &market); err != nil {
		logging.Error("%v", err.Error())
		return nil, 0, err
	}

	return &market, len(market.List), nil
}
