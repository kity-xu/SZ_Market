package publish

import (
	"ProtocolBuffer/format/redis/pbdef/securitytable"

	. "haina.com/share/models"

	"github.com/golang/protobuf/proto"
	redigo "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
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

// 获取全局股票代码表
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

// 从应答缓存获取pb格式全局股票代码表(含应答码200)，如果没有，就缓存一份
func (rds *SecurityTable) GetSecurityTableReplyBytes() ([]byte, error) {
	key := rds.CacheKey + "-reply"
	bytes, err := redigo.Bytes(redis.Get(key))
	if err != nil {
		if err == redigo.ErrNil { // "redigo: nil returned"
			sc, err := rds.GetSecurityTable()
			if err != nil {
				return nil, err
			}

			reply := securitytable.ReplySecurityCodeTable{
				Code:   200,
				Stable: sc,
			}
			replypb, err := proto.Marshal(&reply)
			if err != nil {
				logging.Info("%v", err)
				return nil, err
			}

			if err := redis.Setex(key, REDISKEY_SECURITY_CODETABLE_REPLY_TTL, replypb); err != nil {
				logging.Error("Redis setex %s TTL %d: %s", key, REDISKEY_SECURITY_CODETABLE_REPLY_TTL, err)
				return nil, err
			}
			logging.Info("Redis setex %s TTL %d", key, REDISKEY_SECURITY_CODETABLE_REPLY_TTL)

			return replypb, nil
		}
		// 其他错误
		return nil, err
	}
	return bytes, nil
}
