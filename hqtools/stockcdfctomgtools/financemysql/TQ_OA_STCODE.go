package financemysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"haina.com/share/logging"
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

// 查询沪深市场证券代码
func (this *FcSecuNameTab) GetSecuNmList() ([]*FcSecuNameTab, error) {
	conn, err := dbr.Open("mysql", "finchina:finchina@tcp(114.55.105.11:3306)/finchina?charset=utf8", nil)
	if err != nil {
		logging.Debug("mysql onn", err)
	}
	sess := conn.NewSession(nil)

	var data []*FcSecuNameTab
	_, err = sess.Select("*").From("TQ_OA_STCODE").
		Where("EXCHANGE in ('001002','001003') and SETYPE='101' AND LISTSTATUS =1 and ISVALID =1").
		OrderBy("SYMBOL").LoadStructs(&data)
	return data, err
}
