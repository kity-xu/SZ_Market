package finchina

import (
	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

/**
  机构持股接口
  对应数据库表：TQ_SK_SHAREHOLDER
  中文名称：股东名单
*/

type TQ_SK_SHAREHOLDER struct {
	Model          `db:"-" `
	ID             int64   // ID
	COMPCODE       string  // 公司内码
	ENDDATE        string  // 截止日期
	HOLDERAMT      float64 // 持股数量
	HOLDERRTO      float64 // 持股数量占总股本比例
	SHHOLDERNAME   string  // 股东名称
	SHHOLDERTYPE   int64   // 股东机构类型
	SHARESTYPE     string  // 股份类型
	UNLIMHOLDERAMT float64 // 其中:无限售股份数量
}

func NewTQ_SK_SHAREHOLDER() *TQ_SK_SHAREHOLDER {
	return &TQ_SK_SHAREHOLDER{
		Model: Model{
			TableName: TABLE_TQ_SK_SHAREHOLDER,
			Db:        MyCat,
		},
	}
}

func NewTQ_SK_SHAREHOLDERTx(tx *dbr.Tx) *TQ_SK_SHAREHOLDER {
	return &TQ_SK_SHAREHOLDER{
		Model: Model{
			TableName: TABLE_TQ_SK_SHAREHOLDER,
			Db:        MyCat,
			Tx:        tx,
		},
	}
}

func (this *TQ_SK_SHAREHOLDER) GetSingleByScode(scode string) (*TQ_SK_SHAREHOLDER, string, error) {
	//根据证券代码获取公司内码
	sc := NewSymbolToCompcode()
	if err := sc.getCompcode(scode); err != nil {
		return this, "", err

	}

	Bulid := this.Db.Select("ENDDATE").
		From(this.TableName).
		Where("COMPCODE=" + sc.COMPCODE.String).
		OrderBy("ENDDATE desc ")

	Bulid = Bulid.Limit(1)
	_, err := this.SelectWhere(Bulid, nil).LoadStructs(this)

	return this, sc.COMPCODE.String, err
}

// 获取机构持股信息
func (this *TQ_SK_SHAREHOLDER) GetListByExps(exps map[string]interface{}) ([]*TQ_SK_SHAREHOLDER, error) {
	var data []*TQ_SK_SHAREHOLDER
	bulid := this.Db.Select("*").
		From(this.TableName)
	_, err := this.SelectWhere(bulid, exps).LoadStructs(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}
