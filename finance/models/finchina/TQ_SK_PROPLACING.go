//上市公司配股情况表（产品表）
package finchina

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/models"
)

//Rights Offering
type TQ_SK_PROPLACING struct {
	models.Model   `db:"-" `
	ALLOTCODE      dbr.NullString  //配股代码 			Code     //配股代码
	ALLOTPRICE     dbr.NullFloat64 //配股价格 			Price    //实际配股价格
	ALLOTSNAME     dbr.NullString  //配股简称 			Short  	 //配股简称
	ACTISSQTY      dbr.NullFloat64 //实际配售数量 		Vol   	 //实际配股数量
	PLANISSMAXQTY  dbr.NullFloat64 //拟配售数量上限		PVol 	 //计划配股数量
	ALLOTRT        dbr.NullFloat64 //实际配股比例 		Prop     //实际配股比例
	ACTNETRAISEAMT dbr.NullFloat64 //实际募集资金净额     NetRaise

	ACTTOTALLOTRT dbr.NullFloat64 //配股比例(10:X)  	PProp  	 //计划配股比例   ??

	PAYBEGDATE    dbr.NullString //配股缴款起始日
	PAYENDDATE    dbr.NullString //配股缴款截止日
	EXRIGHTDATE   dbr.NullString //除权日					ERDate   //配股除权日
	EQURECORDDATE dbr.NullString //股权登记日 				RegDate  //股权登记日
	LISTDATE      dbr.NullString //新增股份上市日 				LisDate  //配股上市日

	PUBLISHDATE dbr.NullString //首次公告日期 				PNDate   //预案公布日				??
	LISTPUBDATE dbr.NullString //上市公告日 					DNDate   //决案公布日				??
	UPDATEDATE  dbr.NullString //资料更新日期 				AGMD 	 //股东大会决议公告日		??

	ISFINSUC    dbr.NullInt64  //融资是否成功  1：是  0：否
	ISSUESTATUS dbr.NullString //发行状态
}

func (this *TQ_SK_PROPLACING) newTQ_SK_PROPLACING() *TQ_SK_PROPLACING {
	return &TQ_SK_PROPLACING{
		Model: models.Model{
			CacheKey:  "redis_key",
			TableName: TABLE_TQ_SK_PROPLACING,
			Db:        models.MyCat,
		},
	}
}

//获取配股RO数据
func (this *TQ_SK_PROPLACING) GetROListFromFC(scode string) ([]TQ_SK_PROPLACING, error) {
	var ros []TQ_SK_PROPLACING
	ro := this.newTQ_SK_PROPLACING()

	//根据股票代码获取公司内码
	sc := NewTQ_OA_STCODE()
	if err := sc.getCompcode(scode); err != nil {
		return ros, err
	}
	exps := map[string]interface{}{
		"COMPCODE=?": sc.COMPCODE,
		"ISVALID=?":  1,
	}
	builder := ro.Db.Select("*").From(ro.TableName)
	num, err := ro.SelectWhere(builder, exps).OrderBy("PUBLISHDATE desc").LoadStructs(&ros)
	if err != nil {
		logging.Error("%s", err.Error())
		return ros, err
	}
	logging.Debug("dataSize %d:", num)
	return ros, err
}
