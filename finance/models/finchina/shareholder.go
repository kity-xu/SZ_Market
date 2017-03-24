package finchina

import (
	"errors"

	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

/**
  股东人数接口
  对应数据库表：TQ_SK_SHAREHOLDERNUM
  中文名称：股东户数统计
*/

type ShareHolder struct {
	Model              `db:"-" `
	SYMBOL             string // 股票代码
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

type ShareHolderJson struct {
	//GUID  string `json:"_id"`   // GUID
	//SCode string `json:"SCode"` // 股票代码
	CRPS string `json:"CRPS"` // 户均持股较上期变化（%）
	TNS  string `json:"TNS"`  // 股东总户数（户）
	Date string `json:"Date"` // 指标\日期
	ANS  string `json:"ANS"`  // 户均持股数（股/户）
	APS  string `json:"APS"`  // 户均持股比例（%）
}

type SharList interface{}
type RetInfoJson struct {
	SCode    string      `json:SCode`
	SharList interface{} `json:"Shareholders"`
}

// 获取多条股东人数信息
func (this *ShareHolder) GetListByExps(enddate string, sCode string, limit int) (RetInfoJson, error) {
	var data []*ShareHolder
	var js RetInfoJson

	bulid := this.Db.Select("l.COMPCODE,s.COMPCODE,l.SYMBOL,s.ENDDATE,s.TOTALSHAMT,s.KAVGSH,s.HOLDPROPORTIONPACC,s.PROPORTIONCHG ").
		From(this.TableName+" As s").
		Join(TABLE_TQ_SK_LCPERSON+" As l", "s.COMPCODE=l.COMPCODE").
		Where(" l.SYMBOL = " + sCode).
		OrderBy(" s.ENDDATE  desc ")

	if limit > 0 {
		bulid = bulid.Limit(uint64(limit))
	}

	_, err := this.SelectWhere(bulid, nil).LoadStructs(&data)
	if err != nil {
		//return nil, err
	}

	jsns := []*ShareHolderJson{}

	for _, item := range data {
		jsn, err := this.GetJson(item)
		if err != nil {
			//	return jsns, err
		}

		jsns = append(jsns, jsn)

		js.SharList = jsns
		js.SCode = data[0].SYMBOL
	}

	return js, nil
}

// 获取JSON
func (this *ShareHolder) GetJson(shareHolder *ShareHolder) (*ShareHolderJson, error) {
	var jsn ShareHolderJson
	if len(shareHolder.SYMBOL) < 1 {
		return &jsn, errors.New("SYMBOL is nil")
	}

	return &ShareHolderJson{
		Date: shareHolder.ENDDATE,
		TNS:  shareHolder.TOTALSHAMT,
		ANS:  shareHolder.KAVGSH,
		APS:  shareHolder.HOLDPROPORTIONPACC,
		CRPS: shareHolder.PROPORTIONCHG,
	}, nil
}
