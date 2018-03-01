package fcmysql

import (
	_ "github.com/go-sql-driver/mysql"
	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

// 数据对象名称：TQ_FIN_PROINCSTATEMENTNEW    中文名称：一般企业利润表(新准则产品表)

type TQ_FIN_PROINCSTATEMENTNEW struct {
	Model       `db:"-"`
	BASICEPS    dbr.NullFloat64 `db:"BASICEPS"`    // 基本每股收益
	TOTPROFIT   dbr.NullFloat64 `db:"TOTPROFIT"`   // 利润总额
	NETPROFIT   dbr.NullFloat64 `db:"NETPROFIT"`   // 净利润
	PUBLISHDATE dbr.NullInt64   `db:"PUBLISHDATE"` // 信息发布日期
	BIZINCO     dbr.NullFloat64 `db:"BIZINCO"`     // 主营业务收入
	PERPROFIT   dbr.NullFloat64 `db:"PERPROFIT"`   // 主营业务利润
	INVEINCO    dbr.NullFloat64 `db:"INVEINCO"`    // 投资收益
	ENDDATE     dbr.NullString  `db:"ENDDATE"`     // 截止日期
	PARENETP    dbr.NullFloat64 `db:"PARENETP"`    // 归属母公司净利润
	COMPCODE	dbr.NullString  `db:"COMPCODE"`     //
}

func NewTQ_FIN_PROINCSTATEMENTNEW() *TQ_FIN_PROINCSTATEMENTNEW {
	return &TQ_FIN_PROINCSTATEMENTNEW{
		Model: Model{
			TableName: TABLE_TQ_FIN_PROINCSTATEMENTNEW,
			Db:        MyCat,
		},
	}
}

func (this *TQ_FIN_PROINCSTATEMENTNEW) GetSingleInfo(comc string) (TQ_FIN_PROINCSTATEMENTNEW, error) {
	var tsp TQ_FIN_PROINCSTATEMENTNEW

	err := this.Db.Select("BASICEPS,TOTPROFIT,NETPROFIT,PUBLISHDATE,BIZINCO,PERPROFIT,INVEINCO,ENDDATE,PARENETP").
		From(this.TableName).
		Where("COMPCODE=" + comc + " and  ISVALID=1").
		Where("REPORTTYPE=1").
		OrderBy("ENDDATE DESC ").
		Limit(1).
		LoadStruct(&tsp)
	return tsp, err
}

func (this *TQ_FIN_PROINCSTATEMENTNEW) GetAllInfo() (map[dbr.NullString]TQ_FIN_PROINCSTATEMENTNEW, error) {
	var tsp []TQ_FIN_PROINCSTATEMENTNEW
	var tspmap map[dbr.NullString]TQ_FIN_PROINCSTATEMENTNEW
	err := this.Db.Select("BASICEPS,TOTPROFIT,NETPROFIT,PUBLISHDATE,BIZINCO,PERPROFIT,INVEINCO,ENDDATE,PARENETP,COMPCODE").
		From(this.TableName).
		Where("ID in (select max(ID) from tq_fin_proincstatementnew where ISVALID=1 and REPORTTYPE=1 group by COMPCODE)").
		LoadStruct(&tsp)

	//转map
	tspmap = make(map[dbr.NullString]TQ_FIN_PROINCSTATEMENTNEW)
	for _, v := range tsp{
		tspmap[v.COMPCODE] = v
	}
	return tspmap, err
}

// 倒序查询五期 净利润
func (this *TQ_FIN_PROINCSTATEMENTNEW) GetProinList(comc string) ([]*TQ_FIN_PROINCSTATEMENTNEW, error) {
	var tsp []*TQ_FIN_PROINCSTATEMENTNEW

	err := this.Db.Select("NETPROFIT,ENDDATE ").
		From(this.TableName).
		Where("COMPCODE=" + comc + " and  ISVALID=1").
		Where("REPORTTYPE=1").
		OrderBy("ENDDATE DESC ").
		Limit(5).
		LoadStruct(&tsp)
	return tsp, err
}
