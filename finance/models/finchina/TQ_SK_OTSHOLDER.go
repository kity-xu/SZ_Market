package finchina

import (
	"haina.com/share/logging"

	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

/**
  十大流通股东接口
  对应数据库表：TQ_SK_OTSHOLDER
  中文名称：流通股东名单
*/

type TQ_SK_OTSHOLDER struct {
	Model        `db:"-" `
	ID           int64          // ID
	ENDDATE      string         // 放置本次股东信息的截止日期
	HOLDERSUMCHG dbr.NullString // 增持股份 (?大于1是增持小于是减少)
	HOLDERAMT    string         // 持股数
	HOLDERRTO    string         // 持股数量占总股本比例
	ISHIS        int            // 是否上一报告期存在股东
	SYMBOL       string         // 股票代码
	SHHOLDERNAME string         // 股东名称
}

func NewTQ_SK_OTSHOLDER() *TQ_SK_OTSHOLDER {
	return &TQ_SK_OTSHOLDER{
		Model: Model{
			TableName: TABLE_TQ_SK_OTSHOLDER,
			Db:        MyCat,
		},
	}
}

func NewTQ_SK_OTSHOLDERTx(tx *dbr.Tx) *TQ_SK_OTSHOLDER {
	return &TQ_SK_OTSHOLDER{
		Model: Model{
			TableName: TABLE_TQ_SK_OTSHOLDER,
			Db:        MyCat,
			Tx:        tx,
		},
	}
}
func NewCalculate() *Calculate {
	return &Calculate{
		Model: Model{
			TableName: TABLE_TQ_SK_OTSHOLDER,
			Db:        MyCat,
		},
	}
}

type Calculate struct {
	Model `db:"-" `
	///////////下面数据计算获得
	CR   string // 较上期变化
	Rate string // 累计占总股本比
	Sumh string // 前十大股东累计持有股份
}

/**
  获取结算时间列表
*/
func (this *TQ_SK_OTSHOLDER) GetEndDate(sCode string) ([]*TQ_SK_OTSHOLDER, error) {
	var dataTop10 []*TQ_SK_OTSHOLDER
	//根据证卷代码获取公司内码
	sc := NewSymbolToCompcode()
	if err := sc.getCompcode(sCode); err != nil {
		return dataTop10, err

	}

	bulid := this.Db.Select("DISTINCT(ENDDATE)").
		From(this.TableName).
		Where("COMPCODE=" + sc.COMPCODE.String).
		OrderBy("ENDDATE desc").Limit(8)

	_, err := this.SelectWhere(bulid, nil).LoadStructs(&dataTop10)

	if err != nil {
		logging.Debug("%v", err)
		return dataTop10, err
	}
	return dataTop10, err
}

// 获单条数据
func (this *Calculate) GetSingleCalculate(enddate string, scode string) *Calculate {
	//根据证卷代码获取公司内码
	sc := NewSymbolToCompcode()
	if err := sc.getCompcode(scode); err != nil {
		return this

	}

	builder := this.Db.Select("SUM(a.Sumh) As Sumh,SUM(a.HOLDERRTO) As Rate").
		From("(SELECT  HOLDERAMT As Sumh ,HOLDERRTO FROM " + this.TableName).
		Where("COMPCODE='" + sc.COMPCODE.String + "' and ENDDATE= '" + enddate + "'").
		OrderBy("HOLDERAMT desc limit 10)a")
	err := this.SelectWhere(builder, nil).
		LoadStruct(this)
	if err != nil {
		logging.Debug("%v", err)
	}
	return this
}

// 获取十大流通股东信息
func (this *TQ_SK_OTSHOLDER) GetTop10Group(enddate string, scode string, limit int) ([]*TQ_SK_OTSHOLDER, error) {
	var data []*TQ_SK_OTSHOLDER
	//根据证卷代码获取公司内码
	sc := NewSymbolToCompcode()
	if err := sc.getCompcode(scode); err != nil {
		return data, err

	}

	bulid := this.Db.Select(" * ").
		From(this.TableName).
		Where("COMPCODE = '" + sc.COMPCODE.String + "' and ENDDATE= '" + enddate + "'").
		OrderBy("HOLDERAMT  desc ")

	bulid = bulid.Limit(uint64(limit))

	_, err := this.SelectWhere(bulid, nil).LoadStructs(&data)

	if err != nil {
		logging.Debug("%v", err)
		return data, err
	}

	return data, nil
}
