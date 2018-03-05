package fcmysql

import (
	_ "github.com/go-sql-driver/mysql"
	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

type FcSecuNameTab struct {
	Model       `db:"-"`
	EXCHANGE    dbr.NullString `db:"EXCHANGE"`    // 交易市场
	SETYPE      dbr.NullString `db:"SETYPE"`      // 证券类别
	SECODE      dbr.NullString `db:"SECODE"`      // 证券内码
	COMPCODE    dbr.NullString `db:"COMPCODE"`    // 公司内码
	LISTSTATUS  dbr.NullString `db:"LISTSTATUS"`  // 上市状态
	SYMBOL      dbr.NullString `db:"SYMBOL"`      // 证券代码
	SECURITYID  dbr.NullString `db:"SECURITYID"`  // 证券合并内码
	SENAME      dbr.NullString `db:"SENAME"`      // 证券全称
	SESNAME     dbr.NullString `db:"SESNAME"`     // 证券简称
	SEENGNAME   dbr.NullString `db:"SEENGNAME"`   // 证券英文名
	SESPELL     dbr.NullString `db:"SESPELL"`     // 证券拼音
	CUR         dbr.NullString `db:"CUR"`         // 币种
	SzIndusCode dbr.NullString `db:"SzIndusCode"` // 行业代码

}

func NewFcSecuNameTab() *FcSecuNameTab {
	return &FcSecuNameTab{
		Model: Model{
			TableName: TABLE_TQ_OA_STCODE,
			Db:        MyCat,
		},
	}
}

// 查询沪深市场证券代码 个股
func (this *FcSecuNameTab) GetSecuNmList() ([]*FcSecuNameTab, error) {

	var data []*FcSecuNameTab
	_, err := this.Db.Select("*").
		From(this.TableName).
		Where("EXCHANGE in ('001002','001003') and SETYPE in('101') AND LISTSTATUS in (0,1) and ISVALID =1").
		OrderBy("SYMBOL").LoadStructs(&data)
	return data, err
}

// 查询沪深市场证券代码 个股
func (this *FcSecuNameTab) GetSecuAllNmList() ([]*FcSecuNameTab, error) {

	var data []*FcSecuNameTab
	_, err := this.Db.Select("*").
		From(this.TableName).
		Where("EXCHANGE in ('001002','001003') and SETYPE in('101','102','301','302','401','403','404','405','406','499','413','701') AND LISTSTATUS in (0,1) and ISVALID =1").
		OrderBy("SYMBOL").LoadStructs(&data)
	return data, err
}

// 查询沪深市场证券代码 指数
func (this *FcSecuNameTab) GetExponentList() ([]*FcSecuNameTab, error) {

	var data []*FcSecuNameTab
	_, err := this.Db.Select("*").
		From(this.TableName).
		Where("SETYPE ='701' AND (SYMBOL LIKE '399%' OR SYMBOL LIKE '000%') AND LISTSTATUS =1 and ISVALID =1").
		OrderBy("SYMBOL").LoadStructs(&data)
	return data, err
}

// 根据 str 查询证券信息
func (this *FcSecuNameTab) GetComCodeList(str string) ([]FcSecuNameTab, error) {
	var code []FcSecuNameTab
	_, err := this.Db.Select("*").
		From(this.TableName).
		Where("ISVALID =1").
		Where("SETYPE=101").
		Where("SECODE in (" + str + ")").
		OrderBy("SECODE").LoadStructs(&code)
	return code, err
}

// 查询沪深市场基金
func (this *FcSecuNameTab) GetFundList() ([]*FcSecuNameTab, error) {

	var data []*FcSecuNameTab
	_, err := this.Db.Select("*").
		From(this.TableName).
		Where("EXCHANGE in ('001002','001003') and SETYPE in('301','302') and LISTSTATUS in (0,1)and ISVALID =1").
		LoadStructs(&data)
	return data, err
}

// 查询沪深市场债券
func (this *FcSecuNameTab) GetDebtList() ([]*FcSecuNameTab, error) {

	var data []*FcSecuNameTab
	_, err := this.Db.Select("*").
		From(this.TableName).
		Where("EXCHANGE in ('001002','001003') and SETYPE in('401','403','404','405','406','499','413') and LISTSTATUS in (0,1)and ISVALID =1").
		LoadStructs(&data)
	return data, err
}
