package publish

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	//	"strconv"

	ctrl "haina.com/market/hqpublish/controllers"
	. "haina.com/market/hqpublish/models"
	. "haina.com/share/models"

	pro "ProtocolBuffer/projects/hqpublish/go/protocol"

	"github.com/golang/protobuf/proto"

	"haina.com/market/hqpublish/models/fcmysql"
	hsgrr "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

var (
	_ = proto.Marshal
	_ = redis.Init
	_ = GetCache
	_ = ctrl.MakeRespDataByBytes
	_ = errors.New
	_ = fmt.Println
	_ = hsgrr.Bytes
	_ = logging.Info
	_ = bytes.NewBuffer
	_ = binary.Read
	_ = io.ReadFull
)

type Factor struct {
	Model `db:"-"`
}

func NewFactor() *Factor {
	return &Factor{
		Model: Model{
			CacheKey: REDISKEY_SECURITY_HDAY,
		},
	}
}

func (this Factor) ErrDataInvalid(fields string, sid int32, secode string) error {
	return errors.New(fmt.Sprintf("finchina TQ_SK_XDRY fields[%s] invalid by sid[%d], secode[%d]", "BEGINDATE", sid, secode))
}

// 从财汇数据库获取 *股票除权因子*
func (this Factor) GetReferFactors(sid int32) ([]*pro.Factor, error) {
	real_sid := sid % 1000000
	secode, err := fcmysql.NewTQ_OA_STCODE().GetSecode(fmt.Sprintf("%06d", real_sid))
	if err != nil {
		logging.Error("%v", err)
		return nil, err
	}
	fc, err := fcmysql.NewTQ_SK_XDRY().GetFactorsBySecode(secode)
	if err != nil {
		logging.Error("%v", err)
		return nil, err
	}
	frows := make([]*pro.Factor, 0, 100)
	for _, v := range fc {
		switch {
		case !v.BEGINDATE.Valid:
			return nil, this.ErrDataInvalid("BEGINDATE", sid, secode)
		case !v.ENDDATE.Valid:
			return nil, this.ErrDataInvalid("ENDDATE", sid, secode)
		case !v.XDY.Valid:
			return nil, this.ErrDataInvalid("XDY", sid, secode)
		case !v.LTDXDY.Valid:
			return nil, this.ErrDataInvalid("LTDXDY", sid, secode)
		case !v.THELTDXDY.Valid:
			return nil, this.ErrDataInvalid("THELTDXDY", sid, secode)
		}
		newf := pro.Factor{
			NBeginDate:  v.BEGINDATE.Int64,
			NEndDate:    v.ENDDATE.Int64,
			DfXDY:       v.XDY.Float64,
			DfLTDXDY:    v.LTDXDY.Float64,
			DfTHELTDXDY: v.THELTDXDY.Float64,
		}
		frows = append(frows, &newf)
	}
	if len(frows) == 0 {
		return nil, errors.New(fmt.Sprintf("finchina TQ_SK_XDRY no data by sid %d", sid))
	}
	// 当最后一条没有截止日期时：ENDDATE 数据库值为 19000101，这样的话就出现了 起始时间<截止时间
	// 这里改为 99999999 更符合主观认知和方便后续逻辑处理
	if frows[len(frows)-1].NEndDate == 19000101 {
		frows[len(frows)-1].NEndDate = 99999999
	}
	return frows, nil
}

func (this Factor) GetFactorObj(req *pro.RequestFactor) (*pro.PayloadFactor, error) {
	// 获取除权因子列表
	factors, err := this.GetReferFactors(req.SID)
	if err != nil {
		return nil, err
	}

	return &pro.PayloadFactor{
		SID:   req.SID,
		Total: int32(len(factors)),
		FList: factors,
	}, nil
}
