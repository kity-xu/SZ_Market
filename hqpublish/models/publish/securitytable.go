package publish

import (
	"ProtocolBuffer/format/redis/pbdef/securitytable"

	. "haina.com/share/models"

	"github.com/golang/protobuf/proto"
	redigo "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/store/redis"
)

type SecurityTable struct {
	Model `db:"-"`
	securitytable.SecurityCodeTable
}

func NewSecurityTable() *SecurityTable {
	return &SecurityTable{
		Model: Model{
			CacheKey: REDISKEY_SECURITY_CODETABLE,
		},
	}
}

func (rds *SecurityTable) GetSecurityTable() (*securitytable.SecurityCodeTable, error) {
	key := rds.CacheKey
	bytes, err := redigo.Bytes(redis.Get(key))
	if err != nil {
		// 没找到
		if err == redigo.ErrNil {
			return nil, err
		}
		// 其他错误
		return nil, err
	}
	if err := proto.Unmarshal(bytes, &rds.SecurityCodeTable); err != nil {
		return nil, err
	}
	return &rds.SecurityCodeTable, nil
}
