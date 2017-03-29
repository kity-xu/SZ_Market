package finchina

import (
	"fmt"

	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	. "haina.com/share/models"
)

/**
  十大流通股东接口
  对应数据库表：TQ_SK_OTSHOLDER
  中文名称：流通股东名单
*/

type Top10 struct {
	Model        `db:"-" `
	SYMBOL       string         // 股票代码
	SHHOLDERNAME string         // 股东名称
	HOLDERSUMCHG dbr.NullString // 增持股份 (?大于1是增持小于是减少)
	HOLDERAMT    string         // 持股数
	HOLDERRTO    string         // 持股数量占总股本比例
	ISHIS        int            // 是否上一报告期存在股东
	ENDDATE      string         // 放置本次股东信息的截止日期

}

func NewTop10() *Top10 {
	return &Top10{
		Model: Model{
			TableName: TABLE_TQ_SK_OTSHOLDER,
			Db:        MyCat,
		},
	}
}

func NewTop10Tx(tx *dbr.Tx) *Top10 {
	return &Top10{
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
	Sumh string // 前十大股东累计持有股份
	CR   string // 较上期变化
	Rate string // 累计占总股本比

}

/**
  获取结算时间列表
*/
func (this *Top10) GetEndDate(sCode string) ([]*Top10, error) {
	//根据证卷代码获取公司内码
	sc := NewSymbolToCompcode()
	if err := sc.getCompcode(sCode); err != nil {
		//return nil, err

	}

	if sc.COMPCODE.Valid == false {
		logging.Error("finchina db: select COMPCODE from %s where SYMBOL='%s'", TABLE_TQ_OA_STCODE, sc.COMPCODE)
		//return nil, ErrNullComp
	}

	var dataTop10 []*Top10
	bulid := this.Db.Select(" DISTINCT(ENDDATE)  ").
		From(this.TableName).
		Where(" COMPCODE = " + sc.COMPCODE.String).
		OrderBy(" ENDDATE  desc ").Limit(8)

	_, err := this.SelectWhere(bulid, nil).LoadStructs(&dataTop10)

	if err != nil {
		fmt.Println(err)
		//return nil, err
	}
	return dataTop10, err
}

// 获单条数据
func (this *Calculate) GetSingleByExps(enddate string, sCode string) *Calculate {
	//根据证卷代码获取公司内码
	sc := NewSymbolToCompcode()
	if err := sc.getCompcode(sCode); err != nil {
		//return nil, err

	}

	if sc.COMPCODE.Valid == false {
		logging.Error("finchina db: select COMPCODE from %s where SYMBOL='%s'", TABLE_TQ_OA_STCODE, sc.COMPCODE)
		//return nil, ErrNullComp
	}

	builder := this.Db.Select(" SUM(a.Sumh) As Sumh,SUM(a.HOLDERRTO) As Rate").
		From("(SELECT  HOLDERAMT As Sumh ,HOLDERRTO FROM " + this.TableName).
		Where("  COMPCODE='" + sc.COMPCODE.String + "' and ENDDATE= '" + enddate + "'").
		OrderBy(" HOLDERAMT  desc limit 10)a")
	err := this.SelectWhere(builder, nil).
		LoadStruct(this)
	fmt.Println(err)
	return this
}

// 获取十大流通股东信息
func (this *Top10) GetTop10List(enddate string, sCode string, limit int) ([]*Top10, error) {

	//根据证卷代码获取公司内码
	sc := NewSymbolToCompcode()
	if err := sc.getCompcode(sCode); err != nil {
		//return nil, err

	}

	if sc.COMPCODE.Valid == false {
		logging.Error("finchina db: select COMPCODE from %s where SYMBOL='%s'", TABLE_TQ_OA_STCODE, sc.COMPCODE)
		//return nil, ErrNullComp
	}

	var data []*Top10

	bulid := this.Db.Select(" ENDDATE ,SHHOLDERNAME,HOLDERAMT,HOLDERRTO,ISHIS,HOLDERSUMCHG   ").
		From(this.TableName).
		Where(" COMPCODE = '" + sc.COMPCODE.String + "' and ENDDATE= '" + enddate + "'").
		OrderBy(" HOLDERAMT  desc ")
	if limit > 0 {
		bulid = bulid.Limit(uint64(limit))
	}

	_, err := this.SelectWhere(bulid, nil).LoadStructs(&data)

	if err != nil {
		fmt.Println(err)
		//return nil, err
	}

	return data, nil
}
