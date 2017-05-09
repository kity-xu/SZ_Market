//股票基本信息表
package finchina

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/models"
)

type SecurityInfo struct {
	models.Model `db:"-" `
	LISTDATE     dbr.NullString //上市日期
}

func newSecurityInfo() *SecurityInfo {
	return &SecurityInfo{
		Model: models.Model{
			CacheKey:  "redis_key",
			TableName: TABLE_TQ_SK_BASICINFO,
			Db:        models.MyCat,
		},
	}
}

func (this *SecurityInfo) GetSecurityBasicInfo(scode string, market string) (*SecurityInfo, error) {
	info := newSecurityInfo()

	//根据股票代码获取公司内码
	sc := NewTQ_OA_STCODE()
	if err := sc.getCompcode(scode, market); err != nil {
		return info, err
	}

	exps := map[string]interface{}{
		"COMPCODE=?": sc.COMPCODE,
		"ISVALID=?":  1,
	}
	builder := info.Db.Select("*").From(info.TableName)
	err := info.SelectWhere(builder, exps).Limit(1).LoadStruct(info)
	if err != nil {
		logging.Error("%s", err.Error())
		return info, err
	}
	logging.Debug("get company info success...")
	return info, err
}
