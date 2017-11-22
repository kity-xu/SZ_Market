package security

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"
	"fmt"

	. "haina.com/share/models"

	"github.com/golang/protobuf/proto"
	. "haina.com/market/hqpublish/models"

	"bytes"
	"encoding/binary"
	"io"

	"haina.com/share/logging"
)

const (
	REDISKEY_MARKET_SECURITY_TABLE_ASTOCK = "hq:market:sts:%s" ///A股市场
	REDISKEY_MARKET_SECURITY_TABLE        = "hq:market:sts:%d" ///<证券市场代码表(参数：MarketID)  (hq-init写入)
	REDISKEY_SECURITY_NAME_ID             = "hq:st:name:%d"    ///<证券代码(参数：sid) (hq-init写入)
	REDISKEY_SECURITY_NAME_CODE           = "hq:st:name:%s"    ///<证券代码(参数：scode) (hq-init写入)
	REDISKEY_SECURITY_STATIC              = "hq:st:static:%d"  ///<证券静态数据(参数：sid) (hq-init写入)
)

type SecurityInfo struct {
	Model `db:"-"`
}

func NewSecurityInfo() *SecurityInfo {
	return &SecurityInfo{
		Model: Model{
			CacheKey: REDISKEY_SECURITY_NAME_ID,
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
	secName := &TagSecurityName{}

	if err = binary.Read(bytes.NewBuffer(bs), binary.LittleEndian, secName); err != nil && err != io.EOF {
		return err
	}

	var kinfo = &protocol.SecurityName{
		NSID:        secName.NSID,
		NMarket:     secName.NMarket,
		SzSType:     ByteNToString(secName.SzSType),
		SzStatus:    ByteNToString(secName.SzStatus),
		SzSCode:     ByteNToString(secName.SzSCode),
		SzSymbol:    ByteNToString(secName.SzSymbol),
		SzISIN:      ByteNToString(secName.SzISIN),
		SzSName:     ByteNToString(secName.SzSName),
		SzSCName:    ByteNToString(secName.SzSCName),
		SzDESC:      ByteNToString(secName.SzDESC),
		SzPhonetic:  ByteNToString(secName.SzPhonetic),
		SzCUR:       ByteNToString(secName.SzCUR),
		SzIndusCode: ByteNToString(secName.SzIndusCode),
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
