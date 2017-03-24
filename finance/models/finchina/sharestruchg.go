package finchina

import (
	"errors"

	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

/**
  股本结构接口
  对应数据库表：TQ_SK_SHARESTRUCHG
  中文名称：股本结构变化
*/

//总股本
type Sharestruchg struct {
	Model       `db:"-" `
	SYMBOL      string // 股票代码
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

}

//股本变动
type ChangesEquity struct {
	Model      `db:"-" `
	ENDDATE    string // 变动日期对应值
	SHCHGRSN   string // 变动原因对应值
	CIRCAAMT   string // 流通A股数及变化比例对应值
	RECIRCAAMT string // 限售A股数及变动比例对应值
	TOTALSHARE string // 总股本及变化比例对应值
}

func NewSharestruchg() *Sharestruchg {
	return &Sharestruchg{
		Model: Model{
			TableName: TABLE_TQ_SK_SHARESTRUCHG,
			Db:        MyCat,
		},
	}
}

func NewSharestruchgTx(tx *dbr.Tx) *Sharestruchg {
	return &Sharestruchg{
		Model: Model{
			TableName: TABLE_TQ_SK_SHARESTRUCHG,
			Db:        MyCat,
			Tx:        tx,
		},
	}
}

func NewChangesEquity() *ChangesEquity {
	return &ChangesEquity{
		Model: Model{
			TableName: TABLE_TQ_SK_SHARESTRUCHG,
			Db:        MyCat,
		},
	}
}

//总股本
type SharestruchgJson struct {
	//SCode  string `json:"SCode"`  // 股票代码
	OuSh   string `json:"OuSh"`   // 流通股份
	OuShTO string `json:"OuShTO"` // 流通股份所占比例
	NOS    string `json:"NOS"`    // 未流通股份
	NOSTO  string `json:"Prop"`   // 未流通股份所占比例
	ROS    string `json:"ROS"`    // 限售流通股份
	ROSTO  string `json:"ROSTO"`  // 限售流通股份所占比例
}

//流通A股本
type SharestruchgAJson struct {
	//流通A股

	CAMT   string `json:"CAMT"`   // 已上市流通A股
	CAMTTO string `json:"CAMTTO"` // 已上市流通A股所占比例
	OAMT   string `json:"OAMT"`   // 其他流通股
	OAMTTO string `json:"OAMTTO"` // 其他流通股所占比例
	RAMT   string `json:"RAMT"`   // 限售流通A股
	RAMTTO string `json:"RAMTTO"` // 限售流通A股所占比例
}

//股本变动
type ChangesEquityJson struct {
	CDCV string `json:"CDCV"` // 变动日期对应值
	CCCV string `json:"CCCV"` // 变动原因对应值
	NSCV string `json:"NSCV"` // 流通A股数及变化比例对应值
	SPCV string `json:"SPCV"` // 限售A股数及变动比例对应值
	TPCV string `json:"TPCV"` // 总股本及变化比例对应值
}

type TrucList interface{}
type TrucAList interface{}
type RetTrucInfoJson struct {
	SCode     string      `json:SCode`
	TrucList  interface{} `json:"TSC"`
	TrucAList interface{} `json:"CAS"`
}

//获取股本结构信息
func (this *Sharestruchg) GetSingleByExps(sCode string) (RetTrucInfoJson, error) {
	var str = ""
	var js RetTrucInfoJson
	str += " ENDDATE=(select ENDDATE from TQ_SK_SHARESTRUCHG "
	str += " where COMPCODE=(select COMPCODE from TQ_SK_LCPERSON where SYMBOL='" + sCode + "')"
	str += " ORDER BY ENDDATE desc LIMIT 1)and"
	str += " COMPCODE=(select COMPCODE from TQ_SK_LCPERSON where SYMBOL='" + sCode + "')"
	var strs = ""
	strs += "ENDDATE, CIRCSKAMT,CIRCSKRTO , LIMSKAMT, LIMSKRTO,	NCIRCAMT ,NONNEGSKRTO,	TOTALSHARE ,"
	strs += " CIRCAAMT ,(CIRCAAMT/TOTALSHARE)As CIRCAAMTTO,"
	strs += " OTHERCIRCAMT,(OTHERCIRCAMT/TOTALSHARE)As OTHERCIRCAMTTO,"
	strs += " RECIRCAAMT,(RECIRCAAMT/TOTALSHARE)As RECIRCAAMTTO"
	bulid := this.Db.Select(strs).
		From(this.TableName).
		Where(str)
	err := this.SelectWhere(bulid, nil).
		Limit(1).
		LoadStruct(this)

	jsn, err := this.GetJson(this)
	jsna, err := this.GetAJson(this)
	if err != nil {
		//return jsn, err
	}
	js.TrucList = jsn
	js.TrucAList = jsna

	return js, err
}

// 获取JSON
func (this *Sharestruchg) GetJson(sharestruchg *Sharestruchg) (*SharestruchgJson, error) {
	var jsn SharestruchgJson
	if len(sharestruchg.CIRCSKAMT) < 1 {
		return &jsn, errors.New("obj is nil")
	}

	return &SharestruchgJson{
		//SCode:  sharestruchg.SYMBOL,      // 股票代码
		OuSh:   sharestruchg.CIRCSKAMT,   // 流通股份
		OuShTO: sharestruchg.CIRCSKRTO,   // 流通股份所占比例
		ROS:    sharestruchg.LIMSKAMT,    // 限售流通股份
		ROSTO:  sharestruchg.LIMSKRTO,    // 限售流通股份所占比例
		NOS:    sharestruchg.NCIRCAMT,    // 未流通股份
		NOSTO:  sharestruchg.NONNEGSKRTO, // 未流通股份所占比例
	}, nil
}

// 获取流通A股JSON
func (this *Sharestruchg) GetAJson(sharestruchg *Sharestruchg) (*SharestruchgAJson, error) {
	var jsn SharestruchgAJson
	if len(sharestruchg.CIRCAAMT.String) < 1 {
		return &jsn, errors.New("obj is nil")
	}

	return &SharestruchgAJson{
		//流通A股
		CAMT:   sharestruchg.CIRCAAMT.String,     // 已上市流通A股
		CAMTTO: sharestruchg.CIRCAAMTTO,          // 已上市流通A股所占比例
		OAMT:   sharestruchg.OTHERCIRCAMT.String, // 其他流通股
		OAMTTO: sharestruchg.OTHERCIRCAMTTO,      // 其他流通股所占比例
		RAMT:   sharestruchg.RECIRCAAMT.String,   // 限售流通A股
		RAMTTO: sharestruchg.RECIRCAAMTTO,        // 限售流通A股所占比例
	}, nil
}

/////////////////////////股本变动
type ShaChaList interface{}
type RetShaInfoJson struct {
	ShaChaList interface{} `json:"ChEq"`
}

func (this *ChangesEquity) GetChangesStrJson(enddate string, sCode string, limit int) (RetShaInfoJson, error) {
	var data []*ChangesEquity
	var rij RetShaInfoJson
	bulid := this.Db.Select(" ENDDATE,SHCHGRSN,TOTALSHARE,CIRCAAMT, RECIRCAAMT ").
		From(this.TableName).
		Where(" COMPCODE=(select COMPCODE from TQ_SK_LCPERSON where SYMBOL='" + sCode + "')").
		OrderBy("ENDDATE  desc ")
	if limit > 0 {
		bulid = bulid.Limit(uint64(limit))
	}

	_, err := this.SelectWhere(bulid, nil).LoadStructs(&data)

	if err != nil {
		//return nil, err
	}
	jsns := []*ChangesEquityJson{}

	for _, item := range data {
		jsn, err := this.GetChaEquJson(item)
		if err != nil {
			//return jsns, err
		}

		jsns = append(jsns, jsn)
	}
	rij.ShaChaList = jsns
	return rij, nil
}

// 获取JSON
func (this *ChangesEquity) GetChaEquJson(ce *ChangesEquity) (*ChangesEquityJson, error) {
	var jsn ChangesEquityJson
	if len(ce.ENDDATE) < 1 {
		return &jsn, errors.New("obj is nil")
	}

	return &ChangesEquityJson{
		CDCV: ce.ENDDATE,    // 变动日期对应值
		CCCV: ce.SHCHGRSN,   // 变动原因对应值
		NSCV: ce.CIRCAAMT,   // 流通A股数及变化比例对应值
		SPCV: ce.RECIRCAAMT, // 限售A股数及变动比例对应值
		TPCV: ce.TOTALSHARE, // 总股本及变化比例对应值
	}, nil
}
