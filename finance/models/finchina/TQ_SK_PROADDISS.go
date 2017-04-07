//上市公司增发情况表（产品表）
package finchina

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/models"
)

//Seasoned Equity Offerings
type TQ_SK_PROADDISS struct {
	models.Model      `db:"-" `
	CSRCAPPRAGREEDATE dbr.NullString //发审委通过日期 			IECD 		//发审委公告日
	//LISTPUBDATE       dbr.NullString  //上市公告日 			LisDate 	//新股上市公告日
	UPDATEDATE dbr.NullString //资料更新日期		LisDate		//新股上市公告日
	//PUBLISHDATE     dbr.NullString  //首次公告日 			PNDate 		//预案公布日     ??
	ISSPRICE        dbr.NullFloat64 //						Price   	//实际发行价格
	PLANTOTRAISEAMT dbr.NullFloat64 //						PVal     	//预案募资金额
	ISSUEOBJECT     dbr.NullString  //						Range 		//发行对象类型
	//LISTDATE          dbr.NullString  //新增股份上市日 		SEOD 		//发行新股日
	SHARECORDDATE   dbr.NullString  //股份登记日（非公开增发）	SEOD		//发行新股日
	CSRCAPPDPUBDATE dbr.NullString  //						SRCD   		//证监会核准公告日
	ISSUEMODEMEMO   dbr.NullString  //						Type  		//发行方式
	ACTNETRAISEAMT  dbr.NullFloat64 //实际本次发行股份资金净额 	Val 		//实际募资金额
	ACTISSQTY       dbr.NullFloat64 //						Vol    		//实际发行数量
	PLANISSMAXQTY   dbr.NullFloat64 //拟发行数量上限 			PVol   		//预案发行数量
	ENQUMAXPRICE    dbr.NullFloat64 //询价发行价格上限 		PPrice  	//预案发行价格
	ISFINSUC        dbr.NullInt64   //融资是否成功  1：是  0：否

	//AGMD string `json:"AGMD"` //股东大会决议公告日			??
	//Step string `json:"Step"` //事情进展
}

func (this *TQ_SK_PROADDISS) newTQ_SK_PROADDISS() *TQ_SK_PROADDISS {
	return &TQ_SK_PROADDISS{
		Model: models.Model{
			CacheKey:  "redis_key",
			TableName: TABLE_TQ_SK_PROADDISS,
			Db:        models.MyCat,
		},
	}
}

//获取增发SEO数据
func (this *TQ_SK_PROADDISS) GetSEOListFromFC(scode string) ([]TQ_SK_PROADDISS, error) {
	var seos []TQ_SK_PROADDISS
	seo := this.newTQ_SK_PROADDISS()

	//根据股票代码获取公司内码
	sc := NewTQ_OA_STCODE()
	if err := sc.getCompcode(scode); err != nil {
		return seos, err
	}

	exps := map[string]interface{}{
		"COMPCODE=?": sc.COMPCODE,
	}
	builder := seo.Db.Select("*").From(seo.TableName)
	num, err := seo.SelectWhere(builder, exps).OrderBy("PUBLISHDATE desc").LoadStructs(&seos)
	if err != nil {
		logging.Error("%s", err.Error())
		return seos, err
	}
	logging.Debug("dataSize %d:", num)
	return seos, err
}
