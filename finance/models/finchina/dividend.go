package finchina

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/models"
)

type Dividend struct {
}

//Dividend
type Div struct {
	models.Model       `db:"-" `
	TOTCASHDV          dbr.NullFloat64 //年度分红金额合计  	ToCash
	GRAOBJ             dbr.NullString  //分红对象				Bene   		 //分红对象
	PROBONUSRT         dbr.NullFloat64 //送股（股）			Bonus     	 //送股（股）
	DIVIYEAR           dbr.NullString  //年度				Date      	 //年度
	PRETAXCASHMAXDVCNY dbr.NullFloat64 //分红（元，税前）		Dividend	 //分红（元，税前）
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

//Seasoned Equity Offerings
type SEO struct {
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

//Rights Offering
type RO struct {
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

	ISFINSUC dbr.NullInt64 //融资是否成功  1：是  0：否
}

// Repo(repurchase agreement)
type Repo struct {
}

//New divJson
func (this *Dividend) NewDiv() *Div {
	return &Div{
		Model: models.Model{
			CacheKey:  "redis_key",
			TableName: TABLE_TQ_SK_DIVIDENTS,
			Db:        models.MyCat,
		},
	}
}

func (this *Dividend) NewSEO() *SEO {
	return &SEO{
		Model: models.Model{
			CacheKey:  "redis_key",
			TableName: TABLE_TQ_SK_PROADDISS,
			Db:        models.MyCat,
		},
	}
}
func (this *Dividend) NewRO() *RO {
	return &RO{
		Model: models.Model{
			CacheKey:  "redis_key",
			TableName: TABLE_TQ_SK_PROPLACING,
			Db:        models.MyCat,
		},
	}
}
func (this *Dividend) GetRepo() *Repo {
	return &Repo{}
}

//获取Div数据
func (this *Dividend) GetDivList(sets uint64, secode string) ([]Div, error) {
	var divs []Div
	div := this.NewDiv()

	//根据股票代码获取公司内码
	sc := NewSymbolToCompcode()
	if err := sc.getCompcode(secode); err != nil {
		return divs, err
	}
	if sc.COMPCODE.Valid == false {
		logging.Error("finchina db: select COMPCODE from %s where SYMBOL='%s'", TABLE_TQ_OA_STCODE, secode)
		return divs, ErrNullComp
	}

	exps := map[string]interface{}{
		"COMPCODE=?":    sc.COMPCODE,
		"DATETYPE=?":    "4",
		"GRAOBJTYPE!=?": "99",
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

//获取增发SEO数据
func (this *Dividend) GetSEOList(secode string) ([]SEO, error) {
	var seos []SEO
	seo := this.NewSEO()

	//根据股票代码获取公司内码
	sc := NewSymbolToCompcode()
	if err := sc.getCompcode(secode); err != nil {
		return seos, err
	}
	if sc.COMPCODE.Valid == false {
		logging.Error("finchina db: select COMPCODE from %s where SYMBOL='%s'", TABLE_TQ_OA_STCODE, secode)
		return seos, ErrNullComp
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

//获取配股RO数据
func (this *Dividend) GetROList(secode string) ([]RO, error) {
	var ros []RO
	ro := this.NewRO()

	//根据股票代码获取公司内码
	sc := NewSymbolToCompcode()
	if err := sc.getCompcode(secode); err != nil {
		return ros, err
	}
	if sc.COMPCODE.Valid == false {
		logging.Error("finchina db: select COMPCODE from %s where SYMBOL='%s'", TABLE_TQ_OA_STCODE, secode)
		return ros, ErrNullComp
	}

	exps := map[string]interface{}{
		"COMPCODE=?": sc.COMPCODE,
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
