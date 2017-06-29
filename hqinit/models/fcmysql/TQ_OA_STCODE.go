package fcmysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
)

type FcSecuNameTab struct {
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

// 查询沪深市场证券代码 个股
func (this *FcSecuNameTab) GetSecuNmList(sess *dbr.Session) ([]*FcSecuNameTab, error) {

	var data []*FcSecuNameTab
	_, err := sess.Select("*").From("TQ_OA_STCODE").
		Where("EXCHANGE in ('001002','001003') and SETYPE in('101') AND LISTSTATUS =1 and ISVALID =1").
		OrderBy("SYMBOL").Limit(20).LoadStructs(&data)
	return data, err
}

// 查询沪深市场证券代码 指数
func (this *FcSecuNameTab) GetExponentList(sess *dbr.Session) ([]*FcSecuNameTab, error) {

	var data []*FcSecuNameTab
	_, err := sess.Select("*").From("TQ_OA_STCODE").
		Where("SETYPE ='701' AND (SYMBOL LIKE '399%' OR SYMBOL LIKE '000%') AND LISTSTATUS =1 and ISVALID =1").
		OrderBy("SYMBOL").LoadStructs(&data)
	return data, err
}

// 根据 str 查询证券信息
func (this *FcSecuNameTab) GetComCodeList(sess *dbr.Session, str string) ([]FcSecuNameTab, error) {
	var code []FcSecuNameTab
	_, err := sess.Select("*").From("TQ_OA_STCODE").
		Where("ISVALID =1").
		Where("SETYPE=101").
		Where("SECODE in (" + str + ")").
		OrderBy("SECODE").LoadStructs(&code)
	return code, err
}
