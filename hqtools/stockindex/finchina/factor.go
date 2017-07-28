//根据sid 获取历史日线的复权因子
package finchina

import (
	pro "ProtocolBuffer/projects/hqpublish/go/protocol"

	"errors"
	"fmt"

	"haina.com/share/logging"
)

func ErrDataInvalid(fields string, sid int32, secode string) error {
	return errors.New(fmt.Sprintf("finchina TQ_SK_XDRY fields[%s] invalid by sid[%d], secode[%s]", fields, sid, secode))
}

// 从财汇数据库获取 *股票除权因子*
func GetReferFactors(sid int32) ([]*pro.Factor, error) {
	real_sid := sid % 1000000
	secode, err := NewTQ_OA_STCODE().GetSecode(fmt.Sprintf("%06d", real_sid))
	if err != nil {
		logging.Error("%v", err)
		return nil, err
	}
	fc, err := NewTQ_SK_XDRY().GetFactorsBySecode(secode)
	if err != nil {
		logging.Error("%v", err)
		return nil, err
	}
	frows := make([]*pro.Factor, 0, 100)
	for _, v := range fc {
		switch {
		case !v.BEGINDATE.Valid:
			return nil, ErrDataInvalid("BEGINDATE", sid, secode)
		case !v.ENDDATE.Valid:
			return nil, ErrDataInvalid("ENDDATE", sid, secode)
		case !v.XDY.Valid:
			return nil, ErrDataInvalid("XDY", sid, secode)
		case !v.LTDXDY.Valid:
			return nil, ErrDataInvalid("LTDXDY", sid, secode)
		case !v.THELTDXDY.Valid:
			return nil, ErrDataInvalid("THELTDXDY", sid, secode)
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
		logging.Error("finchina TQ_SK_XDRY no data by sid %d", sid)
		return frows, nil
	}
	// 当最后一条没有截止日期时：ENDDATE 数据库值为 19000101，这样的话就出现了 起始时间<截止时间
	// 这里改为 99999999 更符合主观认知和方便后续逻辑处理
	if frows[len(frows)-1].NEndDate == 19000101 {
		frows[len(frows)-1].NEndDate = 99999999
	}
	return frows, nil
}
