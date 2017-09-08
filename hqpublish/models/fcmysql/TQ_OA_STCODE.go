package fcmysql

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	. "haina.com/share/models"
)

const (
	TABLE_TQ_OA_STCODE = "TQ_OA_STCODE" // 证券内码表
)

type TQ_OA_STCODE struct {
	Model    `db:"-" `
	COMPCODE dbr.NullString // 公司内码
	SECODE   dbr.NullString // 证券内码
}

func NewTQ_OA_STCODE() *TQ_OA_STCODE {
	return &TQ_OA_STCODE{
		Model: Model{
			TableName: TABLE_TQ_OA_STCODE,
			Db:        MyCat,
		},
	}
}

// 获取公司内码	 COMPCODE
func (this *TQ_OA_STCODE) GetStcodeInfo(sid string) (TQ_OA_STCODE, error) {
	var tsa TQ_OA_STCODE

	bulid := this.Db.Select("COMPCODE").
		From(this.TableName).
		Where("EXCHANGE in('001002','001003')").
		Where("SETYPE='101'").
		Where("SYMBOL='" + sid + "'").Limit(1)

	_, err := this.SelectWhere(bulid, nil).LoadStructs(&tsa)

	if err != nil {
		logging.Debug("%v", err)
		return tsa, err
	}
	return tsa, err
}

// 获取证券内码	 SECODE
func (this *TQ_OA_STCODE) GetStockSecode(sid string) (string, error) {
	return this.getSecode(sid, "101")
}

// 获取指数内码	 SECODE
func (this *TQ_OA_STCODE) GetIndexSecode(sid string) (string, error) {
	return this.getSecode(sid, "701")
}

// 获取内码	 SECODE
func (this *TQ_OA_STCODE) getSecode(sid string, setype string) (string, error) {
	var v TQ_OA_STCODE

	bulid := this.Db.Select("SECODE").
		From(this.TableName).
		Where("EXCHANGE in('001002','001003')").
		Where("SETYPE=?", setype).
		Where("SYMBOL='" + sid + "'").Limit(1)

	_, err := this.SelectWhere(bulid, nil).LoadStructs(&v)

	if err != nil {
		logging.Error("%v", err)
		return "", err
	}
	return v.SECODE.String, err
}
