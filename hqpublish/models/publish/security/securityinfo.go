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

type SecurityInfo struct {
	Model `db:"-"`
}

func NewSecurityInfo() *SecurityInfo {
	return &SecurityInfo{
		Model: Model{
			CacheKey: publish.REDISKEY_SECURITY_NAME_ID,
		},
	}
}

func (this *SecurityInfo) GetSecurityBasicInfo(req *protocol.RequestSingleSecurity) (*protocol.PayloadSingleSecurity, error) {
	return this.getSecurityInfoFromeCache(req.SID)
}

func (this *SecurityInfo) getSecurityInfoFromeCache(sid int32) (*protocol.PayloadSingleSecurity, error) {
	key := fmt.Sprintf(this.CacheKey, sid)
	var single = &protocol.PayloadSingleSecurity{}

	bs, err := GetCache(key)
	if err != nil {
		if err = getSecurityInfoFromeStore(key, single); err != nil {
			return nil, err
		}

		if err = setSecurityInfoToCache(key, single); err != nil {
			logging.Error("%v", err.Error())
		}

	} else {
		if err = proto.Unmarshal(bs, single); err != nil {
			return nil, err
		}
	}
	return single, nil
}

func getSecurityInfoFromeStore(key string, single *protocol.PayloadSingleSecurity) error {
	bs, err := RedisStore.GetBytes(key)
	if err != nil {
		return err
	}

	var kinfo = &protocol.SecurityName{}

	if err = proto.Unmarshal(bs, kinfo); err != nil {
		return err
	}
	single.SID = kinfo.NSID
	single.SNInfo = kinfo
	return nil
}

func setSecurityInfoToCache(key string, single *protocol.PayloadSingleSecurity) error {
	bs, err := proto.Marshal(single)
	if err != nil {
		return err
	}

	if err = SetCache(key, 60*5, bs); err != nil {
		return err
	}
	return nil
}
