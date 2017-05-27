package security

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"
	"fmt"

	. "haina.com/share/models"

	"github.com/golang/protobuf/proto"
	. "haina.com/market/hqpublish/models"

	"haina.com/market/hqpublish/models/publish"
	"haina.com/share/logging"
)

type SecurityStatic struct {
	Model `db:"-"`
}

func NewSecurityStatic() *SecurityStatic {
	return &SecurityStatic{
		Model: Model{
			CacheKey: publish.REDISKEY_SECURITY_STATIC,
		},
	}
}

func (this *SecurityStatic) GetSecurityStatic(req *protocol.RequestSecurityStatic) (*protocol.PayloadSecurityStatic, error) {
	return this.getSecurityStaticFromeCache(req.SID)
}

func (this *SecurityStatic) getSecurityStaticFromeCache(sid int32) (*protocol.PayloadSecurityStatic, error) {
	key := fmt.Sprintf(this.CacheKey, sid)
	var single = &protocol.PayloadSecurityStatic{}

	bs, err := GetCache(key)
	if err != nil {
		if err = getSecurityStaticFromeStore(key, single); err != nil {
			return nil, err
		}

		if err = setSecurityStaticToCache(key, single); err != nil {
			logging.Error("%v", err.Error())
		}

	} else {
		if err = proto.Unmarshal(bs, single); err != nil {
			return nil, err
		}
	}
	return single, nil
}

func getSecurityStaticFromeStore(key string, single *protocol.PayloadSecurityStatic) error {
	bs, err := RedisStore.GetBytes(key)
	if err != nil {
		return err
	}

	if err = proto.Unmarshal(bs, single); err != nil {
		return err
	}
	return nil
}

func setSecurityStaticToCache(key string, single *protocol.PayloadSecurityStatic) error {
	bs, err := proto.Marshal(single)
	if err != nil {
		return err
	}

	if err = SetCache(key, 60*5, bs); err != nil {
		return err
	}
	return nil
}
