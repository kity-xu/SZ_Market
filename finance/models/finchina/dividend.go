package finchina

//change test
import (
	"fmt"

	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

type Dividend struct {
}

type Symbol struct {
	Model    `db:"-" `
	COMPCODE string
}

func NewSymbol() *Symbol {
	return &Symbol{
		Model: Model{
			CacheKey:  "redis_key1",
			TableName: TABLE_TQ_OA_STCODE,
			Db:        MyCat,
		},
	}
}

//Dividend
type Div struct {
	Model `db:"-" `
	//ComCode  string
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

type DivJson struct {
	//ComCode  string  `json:"ComCode"`
	Bene     string  `json:"Bene"`     //分红对象
	Bonus    float64 `json:"Bonus"`    //送股（股）
	Date     string  `json:"Data"`     //年度
	Dividend float64 `json:"Dividend"` //分红（元，税前）
	DivDate  string  `json:"DivDate"`  //红利发放日
	DivRate  string  `json:"DivRate"`  //股利支付率（%）
	DNDate   string  `json:"DNDate"`   //决案公布日
	Evolve   string  `json:"Evolve"`   //事情进展
	ExDate   string  `json:"ExDate"`   //除权除息日
	INDate   string  `json:"INDate"`   //实施公告日
	LisDate  string  `json:"LisDate"`  //转股上市日
	PNDate   string  `json:"PNDate"`   //预案公布日
	RegDate  string  `json:"RegDate"`  //股权登记日
	Tran     float64 `json:"Tran"`     //转股（股）
}

//Seasoned Equity Offerings
type SEO struct {
	Model             `db:"-" `
	CSRCAPPRAGREEDATE dbr.NullString //发审委通过日期 		IECD 		//发审委公告日
	//LISTPUBDATE       dbr.NullString  //上市公告日 			LisDate 	//新股上市公告日
	UPDATEDATE dbr.NullString //资料更新日期		LisDate		//新股上市公告日
	//PUBLISHDATE     dbr.NullString  //首次公告日 			PNDate 		//预案公布日     ??
	ISSPRICE        dbr.NullFloat64 //					Price   	//实际发行价格
	PLANTOTRAISEAMT dbr.NullFloat64 //					PVal     	//预案募资金额
	ISSUEOBJECT     dbr.NullString  //					Range 		//发行对象类型
	//LISTDATE          dbr.NullString  //新增股份上市日 		SEOD 		//发行新股日
	SHARECORDDATE   dbr.NullString  //股份登记日（非公开增发） SEOD		//发行新股日
	CSRCAPPDPUBDATE dbr.NullString  //					SRCD   		//证监会核准公告日
	ISSUEMODEMEMO   dbr.NullString  //					Type  		//发行方式
	NEWTOTRAISEAMT  dbr.NullFloat64 //实际公司募集资金总额 	Val 		//实际募资金额
	ACTISSQTY       dbr.NullFloat64 //					Vol    		//实际发行数量
	PLANISSMAXQTY   dbr.NullFloat64 //拟发行数量上限 		PVol   		//预案发行数量
	ENQUMAXPRICE    dbr.NullFloat64 //询价发行价格上限 		PPrice  	//预案发行价格

	//AGMD string `json:"AGMD"` //股东大会决议公告日			??
	//Step string `json:"Step"` //事情进展
}

type SEOJson struct {
	AGMD    string  `json:"AGMD"`    //股东大会决议公告日
	IECD    string  `json:"IECD"`    //发审委公告日
	LisDate string  `json:"LisDate"` //新股上市日
	PNDate  string  `json:"PNDate"`  //预案公布日
	PPrice  float64 `json:"PPrice"`  //预案发行价格
	Price   float64 `json:"Price"`   //实际发行价格
	PVal    float64 `json:"PVal"`    //预案募资金额
	PVol    float64 `json:"PVol"`    //预案发行数量
	Range   string  `json:"Range"`   //发行对象类型
	SEOD    string  `json:"SEOD"`    //发行新股日
	SRCD    string  `json:"SRCD"`    //证监会核准公告日
	Step    string  `json:"Step"`    //事情进展
	Type    string  `json:"Type"`    //发行方式
	Val     float64 `json:"Val"`     //实际募资金额
	Vol     float64 `json:"Vol"`     //实际发行数量
}

//Rights Offering
type RO struct {
	Model         `db:"-" `
	ALLOTCODE     dbr.NullString  //配股代码 			Code     //配股代码
	ALLOTPRICE    dbr.NullFloat64 //配股价格 			Price    //实际配股价格
	ALLOTSNAME    dbr.NullString  //配股简称 			Short  	 //配股简称
	ACTISSQTY     dbr.NullFloat64 //实际配售数量 		Vol   	 //实际配股数量
	PLANISSMAXQTY dbr.NullFloat64 //拟配售数量上限		PVol 	 //计划配股数量
	ALLOTRT       dbr.NullFloat64 //实际配股比例 		Prop     //实际配股比例

	ACTTOTALLOTRT dbr.NullFloat64 //配股比例(10:X)  	PProp  	 //计划配股比例   ??

	PAYBEGDATE    dbr.NullString //配股缴款起始日
	PAYENDDATE    dbr.NullString //配股缴款截止日
	EXRIGHTDATE   dbr.NullString //除权日					ERDate   //配股除权日
	EQURECORDDATE dbr.NullString //股权登记日 				RegDate  //股权登记日
	LISTDATE      dbr.NullString //新增股份上市日 				LisDate  //配股上市日

	PUBLISHDATE dbr.NullString //首次公告日期 				PNDate   //预案公布日				??
	LISTPUBDATE dbr.NullString //上市公告日 					DNDate   //决案公布日				??
	UPDATEDATE  dbr.NullString //资料更新日期 				AGMD 	 //股东大会决议公告日		??
}

type ROJson struct {
	AGMD    string  `json:"AGMD"`    //股东大会决议公告日
	Code    string  `json:"Code"`    //配股代码
	DNDate  string  `json:"DNDate"`  //决案公布日
	ERDate  string  `json:"ERDate"`  //配股除权日
	LisDate string  `json:"LisDate"` //配股上市日
	PNDate  string  `json:"PNDate"`  //预案公布日
	PProp   float64 `json:"PProp"`   //计划配股比例
	Price   float64 `json:"Price"`   //实际配股价格
	Prop    float64 `json:"Prop"`    //实际配股比例
	PVol    float64 `json:"PVol"`    //计划配股数量
	RegDate string  `json:"RegDate"` //股权登记日
	ROPD    string  `json:"ROPD"`    //配股缴款起止日
	Short   string  `json:"Short"`   //配股简称
	Vol     float64 `json:"Vol"`     //实际配股数量
}

// Repo(repurchase agreement)
type RepoJson struct {
}

//New divJson
func (this *Dividend) NewDiv() *Div {
	return &Div{
		Model: Model{
			CacheKey:  "redis_key",
			TableName: TABLE_TQ_SK_DIVIDENTS,
			Db:        MyCat,
		},
	}
}

func (this *Dividend) NewSEO() *SEO {
	return &SEO{
		Model: Model{
			CacheKey:  "redis_key",
			TableName: TABLE_TQ_SK_PROADDISS,
			Db:        MyCat,
		},
	}
}
func (this *Dividend) NewRO() *RO {
	return &RO{
		Model: Model{
			CacheKey:  "redis_key",
			TableName: TABLE_TQ_SK_PROPLACING,
			Db:        MyCat,
		},
	}
}
func (this *Dividend) GetRepo() *RepoJson {
	return &RepoJson{}
}

//根据证券内码获取公司内码
func getCompcodeBySymbol(secode string, symbol *Symbol) error {
	exps := map[string]interface{}{
		"SYMBOL=?": secode,
	}
	builder := symbol.Db.Select("COMPCODE").From(symbol.TableName)
	err := symbol.SelectWhere(builder, exps).LoadStruct(symbol)
	return err
}

//获取Div数据
func (this *Dividend) GetDivList(sets uint64, secode string) ([]Div, error) {
	var divs []Div
	div := this.NewDiv()
	symbol := NewSymbol()
	if err := getCompcodeBySymbol(secode, symbol); err != nil {
		fmt.Println(err.Error())
	}
	exps := map[string]interface{}{
		"COMPCODE=?":    symbol.COMPCODE,
		"DATETYPE=?":    "4",
		"GRAOBJTYPE!=?": "99",
	}
	builder := div.Db.Select("*").From(div.TableName)
	num, err := div.SelectWhere(builder, exps).OrderBy("DIVIYEAR desc").Limit(sets).LoadStructs(&divs)
	if err != nil {
		fmt.Println(num, err.Error())
	}
	fmt.Println("dataSize::", num)
	return divs, err
}

//获取增发SEO数据
func (this *Dividend) GetSEOList(secode string) ([]SEO, error) {
	var seos []SEO
	seo := this.NewSEO()
	symbol := NewSymbol()
	if err := getCompcodeBySymbol(secode, symbol); err != nil {
		fmt.Println(err.Error())
	}
	exps := map[string]interface{}{
		"COMPCODE=?": symbol.COMPCODE,
	}
	builder := seo.Db.Select("*").From(seo.TableName)
	num, err := seo.SelectWhere(builder, exps).OrderBy("PUBLISHDATE desc").LoadStructs(&seos)
	if err != nil {
		fmt.Println(num, err.Error())
	}
	fmt.Println("dataSize::", num)
	return seos, err
}

//获取配股RO数据
func (this *Dividend) GetROList(secode string) ([]RO, error) {
	var ros []RO
	ro := this.NewRO()
	symbol := NewSymbol()
	if err := getCompcodeBySymbol(secode, symbol); err != nil {
		fmt.Println(err.Error())
	}
	exps := map[string]interface{}{
		"COMPCODE=?": symbol.COMPCODE,
	}
	builder := ro.Db.Select("*").From(ro.TableName)
	num, err := ro.SelectWhere(builder, exps).OrderBy("PUBLISHDATE desc").LoadStructs(&ros)
	if err != nil {
		fmt.Println(num, err.Error())
	}
	fmt.Println("dataSize::", num)
	return ros, err
}
