package finchina

import (
	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

/**
  股东人数接口
  对应数据库表：TQ_SK_SHAREHOLDERNUM
  中文名称：股东户数统计
*/

type TQ_SK_SHAREHOLDERNUM struct {
	Model              `db:"-" `
	ID                 int64  // ID
	ENDDATE            string // 指标\日期
	HOLDPROPORTIONPACC string // 户均持股比例（%）
	KAVGSH             string // 户均持股数（股/户）
	PROPORTIONCHG      string // 户均持股较上期变化（%）
	TOTALSHAMT         string // 股东总户数（户）
}

func NewTQ_SK_SHAREHOLDERNUM() *TQ_SK_SHAREHOLDERNUM {
	return &TQ_SK_SHAREHOLDERNUM{
		Model: Model{
			TableName: TABLE_TQ_SK_SHAREHOLDERNUM,
			Db:        MyCat,
		},
	}
}

func NewTQ_SK_SHAREHOLDERNUMTx(tx *dbr.Tx) *TQ_SK_SHAREHOLDERNUM {
	return &TQ_SK_SHAREHOLDERNUM{
		Model: Model{
			TableName: TABLE_TQ_SK_SHAREHOLDERNUM,
			Db:        MyCat,
			Tx:        tx,
		},
	}
}

// 获取多条股东人数信息
func (this *TQ_SK_SHAREHOLDERNUM) GetListByExps(scode string, limit int, strdate string) ([]*TQ_SK_SHAREHOLDERNUM, error) {
	var data []*TQ_SK_SHAREHOLDERNUM
	//根据证券代码获取公司内码
	sc := NewTQ_OA_STCODE()
	if err := sc.getCompcode(scode); err != nil {
		return data, err

	}

	bulid := this.Db.Select("*").
		From(this.TableName).
		Where("COMPCODE = " + sc.COMPCODE.String + strdate + " and ISVALID =1").
		OrderBy("ENDDATE  desc ")

	if limit > 0 {
		bulid = bulid.Limit(uint64(limit))
	}

	_, err := this.SelectWhere(bulid, nil).LoadStructs(&data)
	if err != nil {
		return data, err
	}

	return data, nil
}
