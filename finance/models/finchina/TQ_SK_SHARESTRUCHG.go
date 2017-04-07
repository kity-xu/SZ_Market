package finchina

import (
	"haina.com/share/logging"

	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

/**
  股本结构接口
  对应数据库表：TQ_SK_SHARESTRUCHG
  中文名称：股本结构变化
*/

//总股本
type TQ_SK_SHARESTRUCHG struct {
	Model       `db:"-" `
	ID          int64  // ID
	CIRCSKAMT   string // 流通股份
	CIRCSKRTO   string // 流通股份所占比例
	LIMSKAMT    string // 限售流通股份
	LIMSKRTO    string // 限售流通股份所占比例
	NCIRCAMT    string // 未流通股份
	NONNEGSKRTO string // 未流通股份所占比例
	//流通A股
	CIRCAAMT   dbr.NullString // 已上市流通A股
	CIRCAAMTTO string         // 已上市流通A股所占比例
	//未找到对应字段默认为空
	//---------------
	SIPS   dbr.NullString // 战略投资者配售持股
	SIPSTO string         // 战略投资者配售持股所占比例
	GCPS   dbr.NullString // 一般法人配售持股
	GCPSTO string         // 一般法人配售持股所占比例
	FPS    dbr.NullString // 基金配售持股
	FPSTO  string         // 基金配售持股所占比例
	ARIU   dbr.NullString // 增发未上市
	ARIUTO string         // 增发未上市所占比例
	ASIU   dbr.NullString // 配股未上市
	ASIUTO string         // 配股未上市所占比例
	//----------------
	OTHERCIRCAMT   dbr.NullString // 其他流通股
	OTHERCIRCAMTTO string         // 其他流通股所占比例
	RECIRCAAMT     dbr.NullString // 限售流通A股
	RECIRCAAMTTO   string         // 限售流通A股所占比例

	//股本变动

	ENDDATEV    string // 变动日期对应值
	SHCHGRSNV   string // 变动原因对应值
	CIRCAAMTV   string // 流通A股数及变化比例对应值
	RECIRCAAMTV string // 限售A股数及变动比例对应值
	TOTALSHAREV string // 总股本及变化比例对应值
}

func NewTQ_SK_SHARESTRUCHG() *TQ_SK_SHARESTRUCHG {
	return &TQ_SK_SHARESTRUCHG{
		Model: Model{
			TableName: TABLE_TQ_SK_SHARESTRUCHG,
			Db:        MyCat,
		},
	}
}

func NewTQ_SK_SHARESTRUCHGTx(tx *dbr.Tx) *TQ_SK_SHARESTRUCHG {
	return &TQ_SK_SHARESTRUCHG{
		Model: Model{
			TableName: TABLE_TQ_SK_SHARESTRUCHG,
			Db:        MyCat,
			Tx:        tx,
		},
	}
}

//获取股本结构信息
func (this *TQ_SK_SHARESTRUCHG) GetSingleBySCode(scode string) (*TQ_SK_SHARESTRUCHG, error) {
	var cheq *TQ_SK_SHARESTRUCHG
	//根据证卷代码获取公司内码
	sc := NewSymbolToCompcode()
	if err := sc.getCompcode(scode); err != nil {
		return this, err
	}
	shBulid := this.Db.Select("ENDDATE AS ENDDATEV ").
		From(this.TableName).
		Where("COMPCODE=" + sc.COMPCODE.String).OrderBy(" ENDDATE desc ")
	err1 := this.SelectWhere(shBulid, nil).Limit(1).LoadStruct(&cheq)
	if err1 != nil {
		logging.Debug("%v", err1)
	}
	var strs = ""
	strs += "ENDDATE, CIRCSKAMT,CIRCSKRTO , LIMSKAMT, LIMSKRTO,	NCIRCAMT ,NONNEGSKRTO,	TOTALSHARE ,"
	strs += " CIRCAAMT ,(CIRCAAMT/TOTALSHARE)As CIRCAAMTTO,"
	strs += " OTHERCIRCAMT,(OTHERCIRCAMT/TOTALSHARE)As OTHERCIRCAMTTO,"
	strs += " RECIRCAAMT,(RECIRCAAMT/TOTALSHARE)As RECIRCAAMTTO"
	bulid := this.Db.Select(strs).
		From(this.TableName).
		Where("ENDDATE = " + cheq.ENDDATEV + " and COMPCODE=" + sc.COMPCODE.String)
	err := this.SelectWhere(bulid, nil).
		Limit(1).
		LoadStruct(this)
	if err != nil {
		return this, err
	}

	return this, err
}

/////////////////////////股本变动

func (this *TQ_SK_SHARESTRUCHG) GetChangesStrGroup(enddate string, scode string, limit int) ([]*TQ_SK_SHARESTRUCHG, error) {
	var data []*TQ_SK_SHARESTRUCHG
	//根据证卷代码获取公司内码
	sc := NewSymbolToCompcode()
	if err := sc.getCompcode(scode); err != nil {
		return data, err
	}

	var enddateDx = ""
	if enddate != "" {
		enddateDx = " and ENDDATE < " + enddate
	}
	bulid := this.Db.Select("ENDDATE AS ENDDATEV,SHCHGRSN AS SHCHGRSNV,TOTALSHARE AS TOTALSHAREV,CIRCAAMT AS CIRCAAMTV, RECIRCAAMT AS RECIRCAAMTV").
		From(this.TableName).
		Where("COMPCODE=" + sc.COMPCODE.String + enddateDx).
		OrderBy("ENDDATE  desc ")

	bulid = bulid.Limit(uint64(limit))

	_, err := this.SelectWhere(bulid, nil).LoadStructs(&data)

	if err != nil {
		return data, err
	}

	return data, nil
}
