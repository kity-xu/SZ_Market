//分红情况表
package finchina

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/models"
)

type TQ_SK_DIVIDENTS struct {
	models.Model       `db:"-" `
	TOTCASHDV          dbr.NullFloat64 //年度分红金额合计  	ToCash
	GRAOBJ             dbr.NullString  //分红对象				Bene   		 //分红对象
	PROBONUSRT         dbr.NullFloat64 //送股（股）			Bonus     	 //送股（股）
	DIVIYEAR           dbr.NullString  //年度				Date      	 //年度
	PRETAXCASHMAXDVCNY dbr.NullFloat64 //分红（元，税前）		TQ_SK_DIVIDENTS	 //分红（元，税前）
	XDRDATE            dbr.NullString  //除权除息日			ExDate    	 //除权除息日
	EQURECORDDATE      dbr.NullString  //股权登记日			RegDate  	 //股权登记日
	CASHDVARRBEGDATE   dbr.NullString  //红利发放日			DivDate 	 //红利发放日
	PUBLISHDATE        dbr.NullString  //公告日期				INDate    	 //实施公告日
	SHHDMEETRESPUBDATE dbr.NullString  //股东大会决议公告日	DNDate    	 //决案公布日
	LISTDATE           dbr.NullString  //转股上市日			LisDate   	 //送转股上市日
	TRANADDRT          dbr.NullFloat64 //转股（股）			Tran     	 //转股（股）

	//				预案公布日nil								PNDate		 //预案公布日
	// 						nil								DivRate  	 //股利支付率（%）
	// 						nil								Evolve   	 //事情进展
}

//New divJson
func (this *TQ_SK_DIVIDENTS) newTQ_SK_DIVIDENTS() *TQ_SK_DIVIDENTS {
	return &TQ_SK_DIVIDENTS{
		Model: models.Model{
			CacheKey:  "redis_key",
			TableName: TABLE_TQ_SK_DIVIDENTS,
			Db:        models.MyCat,
		},
	}
}

//获取Div数据
func (this *TQ_SK_DIVIDENTS) GetDivListFromFC(sets uint64, scode string, market string) ([]TQ_SK_DIVIDENTS, error) {
	var divs []TQ_SK_DIVIDENTS
	div := this.newTQ_SK_DIVIDENTS()

	//根据股票代码获取公司内码
	sc := NewTQ_OA_STCODE()
	if err := sc.getCompcode(scode, market); err != nil {
		return divs, err
	}

	exps := map[string]interface{}{
		"COMPCODE=?":    sc.COMPCODE,
		"DATETYPE=?":    "4",
		"GRAOBJTYPE!=?": "99",
		"ISVALID=?":     1,
	}
	builder := div.Db.Select("*").From(div.TableName)
	num, err := div.SelectWhere(builder, exps).OrderBy("DIVIYEAR desc").Limit(sets).LoadStructs(&divs)
	if err != nil {
		logging.Error("%s", err.Error())
		return divs, err
	}
	logging.Debug("dataSize %d:", num)
	return divs, err
}

//给陈亮宇用的
func (this *TQ_SK_DIVIDENTS) GetDivListFromDB(sets uint64, scode string, market string) ([]TQ_SK_DIVIDENTS, error) {
	var divs []TQ_SK_DIVIDENTS
	div := this.newTQ_SK_DIVIDENTS()

	//根据股票代码获取公司内码
	sc := NewTQ_OA_STCODE()
	if err := sc.getCompcode(scode, market); err != nil {
		return divs, err
	}

	exps := map[string]interface{}{
		"COMPCODE=?": sc.COMPCODE,
		"ISVALID=?":  1,
	}
	builder := div.Db.Select("*").From(div.TableName)
	num, err := div.SelectWhere(builder, exps).OrderBy("DIVIYEAR desc").Limit(sets).LoadStructs(&divs)
	if err != nil {
		logging.Error("%s", err.Error())
		return divs, err
	}
	logging.Debug("dataSize %d:", num)
	return divs, err
}

/***********************************以下是移动端f10页面******************************************/
// 该处实现分红配股

type DividendRO struct {
	models.Model       `db:"-" `
	DIVIYEAR           dbr.NullString  //年度
	PRETAXCASHMAXDVCNY dbr.NullFloat64 //分红
	PROBONUSRT         dbr.NullFloat64 //送股比例(10:X)
	TRANADDRT          dbr.NullFloat64 //转增比例(10:X)
	BONUSRT            dbr.NullFloat64 //赠股比例(10:X)
	EQURECORDDATE      dbr.NullString  //股权登记日
}

func NewDividendRO() *DividendRO {
	return &DividendRO{
		Model: models.Model{
			TableName: TABLE_TQ_SK_DIVIDENTS,
			Db:        models.MyCat,
		},
	}
}

func (this *DividendRO) GetDividendRO(compCode string) (*[]DividendRO, error) {
	var divs []DividendRO
	exps := map[string]interface{}{
		"COMPCODE=?":    compCode,
		"DATETYPE=?":    "4",
		"GRAOBJTYPE!=?": "99",
		"CUR=?":         "CNY",
		"ISVALID=?":     1,
	}
	builder := this.Db.Select("DIVIYEAR,PRETAXCASHMAXDVCNY,PROBONUSRT,TRANADDRT,BONUSRT,EQURECORDDATE").From(this.TableName)
	_, err := this.SelectWhere(builder, exps).OrderBy("DIVIYEAR desc").Limit(10).LoadStructs(&divs)
	if err != nil {
		return nil, err
	}
	return &divs, nil
}
