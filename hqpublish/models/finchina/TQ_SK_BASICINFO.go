//股票基本信息表
package finchina

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/models"
)

type SecurityInfo struct {
	models.Model `db:"-" `
	LISTDATE     dbr.NullString  // 上市日期
	TOTALSHARE   dbr.NullFloat64 // 总股本
}

func NewSecurityInfo() *SecurityInfo {
	return &SecurityInfo{
		Model: models.Model{
			CacheKey:  "redis_key",
			TableName: TABLE_TQ_SK_BASICINFO,
			Db:        models.MyCat,
		},
	}
}

func (this *SecurityInfo) GetSecurityBasicInfo(scode string) (*SecurityInfo, error) {
	info := NewSecurityInfo()

	exps := map[string]interface{}{
		"COMPCODE=?": scode,
		"ISVALID=?":  1,
	}
	builder := info.Db.Select("*").From(info.TableName)
	err := info.SelectWhere(builder, exps).Limit(1).LoadStruct(info)
	if err != nil {
		logging.Error("%s", err.Error())
		return info, err
	}
	//logging.Debug("get company info success...")
	return info, err
}
