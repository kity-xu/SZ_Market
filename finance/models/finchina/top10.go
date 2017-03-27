package finchina

import (
	"errors"
	"fmt"

	"haina.com/share/gocraft/dbr"
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
	///////////下面数据是上面数据计算获得
	Sumh string // 前十大股东累计持有股份
	CR   string // 较上期变化
	Rate string // 累计占总股本比
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

type Top10Json struct {
	SCode string `json:"SCode"` // 股票代码
	Name  string `json:"Name"`  // 股东名称
	Ovwet string `json:"Ovwet"` // 增持股份
	Posi  string `json:"Posi"`  // 持股数
	Prop  string `json:"Prop"`  // 持股数量占总股本比例
	ISHIS int    `json:"ISHIS"` // 是否上一报告期存在股东
	//	ENDDATE string `json:"ENDDATE"` // 放置本次股东信息的截止日期
	//	Sum     string `json:"Sum"`     // 前十大股东累计持有股份
	//	CR      string `json:"CR"`      // 较上期变化
	//	Rate    string `json:"Rate"`    // 累计占总股本比
}

type TopList interface{}
type RetTopInfoJson struct {
	SCode   string      `json:SCode`
	Sum     string      `json:Sum`
	Rate    string      `json:Rate`
	CR      string      `json:CR`
	TopList interface{} `json:"TLSG"`
}

// 获单条数据
func (this *Top10) GetSingleByExps(enddate string, sCode string, limit int) error {
	builder := this.Db.Select(" SUM(a.Sumha) As Sumh,SUM(a.HOLDERRTO) As Rate ").
		From("(SELECT  HOLDERAMT As Sumha,HOLDERRTO FROM "+this.TableName+" As o").
		Join(TABLE_TQ_SK_LCPERSON+" As l ", " o.COMPCODE=l.COMPCODE ").
		Where("  l.SYMBOL='" + sCode + "' and ENDDATE= '" + enddate + "'").
		OrderBy(" HOLDERAMT  desc limit 10)a")
	err := this.SelectWhere(builder, nil).
		LoadStruct(this)

	return err
}

// 获取十大流通股东信息
func (this *Top10) GetTop10List(enddate string, sCode string, limit int) (RetTopInfoJson, error) {

	var data []*Top10
	var rij RetTopInfoJson
	bulid := this.Db.Select(" l.SYMBOL,s.ENDDATE ,s.SHHOLDERNAME,s.HOLDERAMT,s.HOLDERRTO,s.ISHIS,s.HOLDERSUMCHG  ").
		From(this.TableName+" As s ").
		Join(TABLE_TQ_SK_LCPERSON+" As l", "s.COMPCODE=l.COMPCODE").
		Where(" l.SYMBOL = '" + sCode + "' and ENDDATE= '" + enddate + "'").
		OrderBy(" s.HOLDERAMT  desc ")
	if limit > 0 {
		bulid = bulid.Limit(uint64(limit))
	}

	_, err := this.SelectWhere(bulid, nil).LoadStructs(&data)

	if err != nil {
		fmt.Println(err)
		//return nil, err
	}
	jsns := []*Top10Json{}

	for _, item := range data {

		jsn, err := this.GetJson(item)
		if err != nil {
			//return jsns, err
		}

		jsns = append(jsns, jsn)
	}

	this.GetSingleByExps(enddate, sCode, limit)

	rij.TopList = jsns
	rij.SCode = data[0].SYMBOL
	rij.Sum = this.Sumh
	rij.Rate = this.Rate

	return rij, nil
}

// 获取JSON
func (this *Top10) GetJson(top10 *Top10) (*Top10Json, error) {
	var jsn Top10Json
	if len(top10.SYMBOL) < 1 {
		return &jsn, errors.New("SYMBOL is nil")
	}

	return &Top10Json{
		Name:  top10.SHHOLDERNAME,
		Ovwet: top10.HOLDERSUMCHG.String,
		Posi:  top10.HOLDERAMT,
		Prop:  top10.HOLDERRTO,
		ISHIS: top10.ISHIS,
	}, nil
}
