package company

import (
	"errors"

	"haina.com/market/finance/models/finchina"
)

//总股本
type StructureEquity struct {
	ID     int64  `json:"-"`      // ID
	OuSh   string `json:"OuSh"`   // 流通股份
	OuShTO string `json:"OuShTO"` // 流通股份所占比例
	NOS    string `json:"NOS"`    // 未流通股份
	NOSTO  string `json:"Prop"`   // 未流通股份所占比例
	ROS    string `json:"ROS"`    // 限售流通股份
	ROSTO  string `json:"ROSTO"`  // 限售流通股份所占比例
	//}

	//流通A股本
	//type SharestruchgA struct {
	//流通A股

	CAMT   string `json:"CAMT"`   // 已上市流通A股
	CAMTTO string `json:"CAMTTO"` // 已上市流通A股所占比例
	OAMT   string `json:"OAMT"`   // 其他流通股
	OAMTTO string `json:"OAMTTO"` // 其他流通股所占比例
	RAMT   string `json:"RAMT"`   // 限售流通A股
	RAMTTO string `json:"RAMTTO"` // 限售流通A股所占比例
}

//股本变动
type ChangesEquity struct {
	ID   int64  `json:"-"`    // ID
	CDCV string `json:"CDCV"` // 变动日期对应值
	CCCV string `json:"CCCV"` // 变动原因对应值
	NSCV string `json:"NSCV"` // 流通A股数及变化比例对应值
	SPCV string `json:"SPCV"` // 限售A股数及变动比例对应值
	TPCV string `json:"TPCV"` // 总股本及变化比例对应值
}

////////////股本结构
type TrucList interface{}
type TrucAList interface{}
type RetTrucInfoJson struct {
	SCode     string      `json:"scode"`
	TrucList  interface{} `json:"TSC"`
	TrucAList interface{} `json:"CAS"`
}

//////////////股本变动
type ShaChaList interface{}
type RetShaInfoJson struct {
	SCode      string      `json:"scode"`
	ShaChaList interface{} `json:"ChEq"`
}

/**
  获取股本结构信息
*/
func GetStructure(scode string) (RetTrucInfoJson, error) {
	data, err := finchina.NewTQ_SK_SHARESTRUCHG().GetSingleBySCode(scode)
	var js RetTrucInfoJson
	jsn, err := GetStruInfo(data)
	jsna, err := GetAInfo(data)

	js.SCode = scode
	js.TrucList = jsn
	js.TrucAList = jsna
	return js, err
}

// 获取JSON
func GetStruInfo(sharestruchg *finchina.TQ_SK_SHARESTRUCHG) (*StructureEquity, error) {
	var jsn StructureEquity
	if len(sharestruchg.CIRCSKAMT) < 1 {
		return &jsn, errors.New("obj is nil")
	}

	return &StructureEquity{
		OuSh:   sharestruchg.CIRCSKAMT,   // 流通股份
		OuShTO: sharestruchg.CIRCSKRTO,   // 流通股份所占比例
		ROS:    sharestruchg.LIMSKAMT,    // 限售流通股份
		ROSTO:  sharestruchg.LIMSKRTO,    // 限售流通股份所占比例
		NOS:    sharestruchg.NCIRCAMT,    // 未流通股份
		NOSTO:  sharestruchg.NONNEGSKRTO, // 未流通股份所占比例
	}, nil
}

// 获取流通A股JSON
func GetAInfo(sharestruchg *finchina.TQ_SK_SHARESTRUCHG) (*StructureEquity, error) {
	var jsn StructureEquity
	if len(sharestruchg.CIRCAAMT.String) < 1 {
		return &jsn, errors.New("obj is nil")
	}

	return &StructureEquity{
		//流通A股
		CAMT:   sharestruchg.CIRCAAMT.String,     // 已上市流通A股
		CAMTTO: sharestruchg.CIRCAAMTTO,          // 已上市流通A股所占比例
		OAMT:   sharestruchg.OTHERCIRCAMT.String, // 其他流通股
		OAMTTO: sharestruchg.OTHERCIRCAMTTO,      // 其他流通股所占比例
		RAMT:   sharestruchg.RECIRCAAMT.String,   // 限售流通A股
		RAMTTO: sharestruchg.RECIRCAAMTTO,        // 限售流通A股所占比例
	}, nil
}

///////////////////////////////////////////////////////////////////股本变动
/**
  获取股本变动信息
*/
func GetChangesStrInfo(enddate string, scode string, limit int) (RetShaInfoJson, error) {
	data, err := finchina.NewTQ_SK_SHARESTRUCHG().GetChangesStrGroup(enddate, scode, limit)
	var rij RetShaInfoJson
	jsns := []*ChangesEquity{}

	for _, item := range data {
		jsn, err := GetChaEquInfo(item)
		if err != nil {
			return rij, err
		}

		jsns = append(jsns, jsn)
	}
	rij.SCode = scode
	rij.ShaChaList = jsns
	return rij, err
}

// 获取JSON
func GetChaEquInfo(ce *finchina.TQ_SK_SHARESTRUCHG) (*ChangesEquity, error) {
	var jsn ChangesEquity
	if len(ce.ENDDATEV) < 1 {
		return &jsn, errors.New("obj is nil")
	}

	return &ChangesEquity{
		CDCV: ce.ENDDATEV,    // 变动日期对应值
		CCCV: ce.SHCHGRSNV,   // 变动原因对应值
		NSCV: ce.CIRCAAMTV,   // 流通A股数及变化比例对应值
		SPCV: ce.RECIRCAAMTV, // 限售A股数及变动比例对应值
		TPCV: ce.TOTALSHAREV, // 总股本及变化比例对应值
	}, nil
}
