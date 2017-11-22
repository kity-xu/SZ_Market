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

type SecurityStatic struct {
	Model `db:"-"`
}

func NewSecurityStatic() *SecurityStatic {
	return &SecurityStatic{
		Model: Model{
			CacheKey: REDISKEY_SECURITY_STATIC,
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
			logging.Error("%v", err.Error())
			return nil, err
		}

		if err = setSecurityStaticToCache(key, single); err != nil {
			logging.Error("%v", err.Error())
		}

	} else {
		if err = proto.Unmarshal(bs, single); err != nil {
			logging.Error("%v", err.Error())
			return nil, err
		}
	}
	return single, nil
}

func getSecurityStaticFromeStore(key string, single *protocol.PayloadSecurityStatic) error {
	static := &StockStatic{}

	bs, err := RedisStore.GetBytes(key)
	if err != nil {
		return err
	}

	if err = binary.Read(bytes.NewBuffer(bs), binary.LittleEndian, static); err != nil && err != io.EOF {
		return err
	}
	pStatic := &protocol.StockStatic{
		NSID:               static.NSID,
		SzSType:            ByteNToString(static.SzSType),
		SzStatus:           ByteNToString(static.SzStatus),
		NListDate:          static.NListDate,
		NLastTradeDate:     static.NLastTradeDate,
		NDelistDate:        static.NDelistDate,
		LlCircuShare:       static.LlCircuShare,
		LlTotalShare:       static.LlTotalShare,
		LlLast5Volume:      static.LlLast5Volume,
		NEPS:               static.NEPS,
		LlTotalProperty:    static.LlTotalProperty,
		LlFlowProperty:     static.LlFlowProperty,
		NAVPS:              static.NAVPS,
		LlMainIncoming:     static.LlMainIncoming,
		LlMainProfit:       static.LlMainProfit,
		LlTotalProfit:      static.LlTotalProfit,
		LlNetProfit:        static.LlNetProfit,
		NHolders:           static.NHolders,
		NReportDate:        static.NReportDate,
		NQuickMovingRatio:  static.NQuickMovingRatio,
		NCurrentRatio:      static.NCurrentRatio,
		NEUndisProfit:      static.NEUndisProfit,
		NFlowLiab:          static.NFlowLiab,
		NTotalLiabilities:  static.NTotalLiabilities,
		NTotalHolderEquity: static.NTotalHolderEquity,
		NCapitalReserve:    static.NCapitalReserve,
		NIncomeInvestments: static.NIncomeInvestments,
	}

	single.SSInfo = pStatic

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
