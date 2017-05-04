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

type SecurityInfo struct {
	Model `db:"-"`
}

func NewSecurityInfo() *SecurityInfo {
	return &SecurityInfo{
		Model: Model{
			CacheKey: publish.REDISKEY_MARKET_SECURITY_TABLE,
		},
	}
}

func (this *SecurityInfo) GetSecurityTable(MarketID int32) (*securitytable.SecurityCodeTable, int, error) {
	var table securitytable.SecurityCodeTable

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

	if err := proto.Unmarshal(data, &table); err != nil {
		logging.Error("%v", err.Error())
		return nil, 0, err
	}
	logging.Debug("---%+v", table)

	//	for _, v := range table {
	//		k := &securitytable.SecurityInfo{}
	//		buffer := bytes.NewBuffer([]byte(v))
	//		if err := binary.Read(buffer, binary.LittleEndian, k); err != nil && err != io.EOF {
	//			return table, 0, err
	//		}

	//		info = append(info, k)
	//		logging.Debug("---%v", k)

	//	}
	//	table.List = info
	return &table, len(table.List), nil
}
