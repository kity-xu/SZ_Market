package finchina

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	. "haina.com/share/models"
)

/**
  股东人数接口
  对应数据库表：TQ_SK_SHAREHOLDERNUM
  中文名称：股东户数统计
*/

type ShareHolder struct {
	Model              `db:"-" `
	PROPORTIONCHG      string // 户均持股较上期变化（%）
	TOTALSHAMT         string // 股东总户数（户）
	ENDDATE            string // 指标\日期
	KAVGSH             string // 户均持股数（股/户）
	HOLDPROPORTIONPACC string // 户均持股比例（%）
}

func NewShareHolder() *ShareHolder {
	return &ShareHolder{
		Model: Model{
			TableName: TABLE_TQ_SK_SHAREHOLDERNUM,
			Db:        MyCat,
		},
	}
}

func NewShareHolderTx(tx *dbr.Tx) *ShareHolder {
	return &ShareHolder{
		Model: Model{
			TableName: TABLE_TQ_SK_SHAREHOLDERNUM,
			Db:        MyCat,
			Tx:        tx,
		},
	}
}

// 获取多条股东人数信息
func (this *ShareHolder) GetListByExps(enddate string, sCode string, limit int) ([]*ShareHolder, error) {
	//根据证卷代码获取公司内码
	sc := NewSymbolToCompcode()
	if err := sc.getCompcode(sCode); err != nil {
		//return nil, err

	}

	if sc.COMPCODE.Valid == false {
		logging.Error("finchina db: select COMPCODE from %s where SYMBOL='%s'", TABLE_TQ_OA_STCODE, sc.COMPCODE)
		//return nil, ErrNullComp
	}

	var data []*ShareHolder

	bulid := this.Db.Select("COMPCODE,ENDDATE,TOTALSHAMT,KAVGSH,HOLDPROPORTIONPACC,PROPORTIONCHG ").
		From(this.TableName).
		Where(" COMPCODE = " + sc.COMPCODE.String).
		OrderBy(" ENDDATE  desc ")

	if limit > 0 {
		bulid = bulid.Limit(uint64(limit))
	}

	_, err := this.SelectWhere(bulid, nil).LoadStructs(&data)
	if err != nil {
		//return nil, err
	}

	return data, nil
}
