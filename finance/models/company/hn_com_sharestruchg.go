package company

import (
	"errors"

	"haina.com/market/finance/models/finchina"
)

//总股本
type StructureEquity struct {
	// ------------------------------------------------原接口
	//	ID     int64  `json:"-"`      // ID
	//	OuSh   string `json:"OuSh"`   // 流通股份
	//	OuShTO string `json:"OuShTO"` // 流通股份所占比例
	//	NOS    string `json:"NOS"`    // 未流通股份
	//	NOSTO  string `json:"Prop"`   // 未流通股份所占比例
	//	ROS    string `json:"ROS"`    // 限售流通股份
	//	ROSTO  string `json:"ROSTO"`  // 限售流通股份所占比例

	//	CAMT   string `json:"CAMT"`   // 已上市流通A股
	//	CAMTTO string `json:"CAMTTO"` // 已上市流通A股所占比例
	//	OAMT   string `json:"OAMT"`   // 其他流通股
	//	OAMTTO string `json:"OAMTTO"` // 其他流通股所占比例
	//	RAMT   string `json:"RAMT"`   // 限售流通A股
	//	RAMTTO string `json:"RAMTTO"` // 限售流通A股所占比例
	// ------------------------------------------------原接口
	GenCap     float32 `json:"GenCap"`  // 总股本
	CIRCAAMT   float32 `json:"CirA"`    // 流通A股
	RECIRCAAMT float32 `json:"RecA"`    // 限售流通A股
	AGenCap    float32 `json:"AGenCap"` // A股总股本
	Edate      string  `json:"Edate"`   // 截止日期
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
type TrucsList interface{}

//type TrucAList interface{}
type RetTrucsInfoJson struct {
	SCode     string `json:"scode"`
	TrucsList interface{}
}

//////////////股本变动
type ShaChaList interface{}
type RetShaInfoJson struct {
	SCode      string      `json:"scode"`
	ShaChaList interface{} `json:"ChEq"`
}

///**
//  获取股本结构信息
//*/
//func _GetStructure(scode string, selwhe string, limit int) (RetTrucsInfoJson, error) {
//	data, err := finchina.NewTQ_SK_SHARESTRUCHG().GetSingleBySCode(scode, selwhe, limit)
//	//	var js RetTrucInfoJson
//	//	//jsn, err := GetStruInfo(data)

//	//	//jsna, err := GetAInfo(data)

//	//	js.SCode = scode
//	//	js.TrucList = jsn
//	//	js.TrucAList = jsna
//	//	return js, err
//	var rtj RetTrucInfoJson
//	jsnse := []*StructureEquity{}

//	for _, item := range data {

//		//if len(item.ENDDATE) > 6 {
//		//str := item.ENDDATE[4:]

//		//}
//		jsn, err := GetAInfo(item)
//		if err != nil {
//			return rtj, err
//		}
//		jsnse = append(jsnse, jsn)
//	}
//	rtj.SCode = scode
//	rtj.TrucList = jsnse
//	return rtj, err
//}
func GetStructure(scode string, selwhe string, limit int) ([]*StructureEquity, error) {
	data, err := finchina.NewTQ_SK_SHARESTRUCHG().GetSingleBySCode(scode, selwhe, limit)

	//var rtj RetTrucInfoJson
	jsnse := []*StructureEquity{}

	for _, item := range data {

		//if len(item.ENDDATE) > 6 {
		//str := item.ENDDATE[4:]

		//}
		jsn, err := GetAInfo(item)
		if err != nil {
			return nil, err
		}
		jsnse = append(jsnse, jsn)
	}
	return jsnse, err
}

//// 获取JSON
//func GetStruInfo(sharestruchg *finchina.TQ_SK_SHARESTRUCHG) (*StructureEquity, error) {
//	var jsn StructureEquity
//	if len(sharestruchg.TOTALSHARE) < 1 {
//		return &jsn, errors.New("obj is nil")
//	}

//	return &StructureEquity{
//		OuSh:   sharestruchg.CIRCSKAMT,   // 流通股份
//		OuShTO: sharestruchg.CIRCSKRTO,   // 流通股份所占比例
//		ROS:    sharestruchg.LIMSKAMT,    // 限售流通股份
//		ROSTO:  sharestruchg.LIMSKRTO,    // 限售流通股份所占比例
//		NOS:    sharestruchg.NCIRCAMT,    // 未流通股份
//		NOSTO:  sharestruchg.NONNEGSKRTO, // 未流通股份所占比例
//	}, nil
//}

// 获取流通A股JSON
func GetAInfo(st *finchina.TQ_SK_SHARESTRUCHG) (*StructureEquity, error) {
	var jsn StructureEquity
	if st.TOTALSHARE < 1 {
		return &jsn, errors.New("obj is nil")
	}

	return &StructureEquity{
		//		//流通A股
		//		CAMT:   sharestruchg.CIRCAAMT.String,     // 已上市流通A股
		//		CAMTTO: sharestruchg.CIRCAAMTTO,          // 已上市流通A股所占比例
		//		OAMT:   sharestruchg.OTHERCIRCAMT.String, // 其他流通股
		//		OAMTTO: sharestruchg.OTHERCIRCAMTTO,      // 其他流通股所占比例
		//		RAMT:   sharestruchg.RECIRCAAMT.String,   // 限售流通A股
		//		RAMTTO: sharestruchg.RECIRCAAMTTO,        // 限售流通A股所占比例

		GenCap:     st.TOTALSHARE,               // 总股本
		CIRCAAMT:   st.CIRCAAMT,                 // 流通A股
		RECIRCAAMT: st.RECIRCAAMT,               // 限售流通A股
		AGenCap:    st.CIRCAAMT + st.RECIRCAAMT, // A股总股本
		Edate:      st.ENDDATE,                  // 截止日期
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
