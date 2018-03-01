package fcmysql

import (
	_ "github.com/go-sql-driver/mysql"
	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

// 数据对象名称：TQ_FIN_PROBALSHEETNEW    中文名称：一般企业资产负债表(新准则产品表)

type TQ_FIN_PROBALSHEETNEW struct {
	Model         `db:"-"`
	TOTASSET      dbr.NullFloat64 `db:"TOTASSET"`      // 资产总计
	TOTALCURRLIAB dbr.NullFloat64 `db:"TOTALCURRLIAB"` // 流动负债合计
	TOTLIAB       dbr.NullFloat64 `db:"TOTLIAB"`       // 负债合计
	CAPISURP      dbr.NullFloat64 `db:"CAPISURP"`      // 资本公积
	TOTCURRASSET  dbr.NullFloat64 `db:"TOTCURRASSET"`  // 流动资产合计
	RIGHAGGR      dbr.NullFloat64 `db:"RIGHAGGR"`      // 所有股东权益合计
	PARESHARRIGH  dbr.NullFloat64 `db:"PARESHARRIGH"`  // 归属于母公司股东权益合计(元)
	COMPCODE	  dbr.NullString  `db:"COMPCODE"`
}

func NewTQ_FIN_PROBALSHEETNEW() *TQ_FIN_PROBALSHEETNEW {
	return &TQ_FIN_PROBALSHEETNEW{
		Model: Model{
			TableName: TABLE_TQ_FIN_PROBALSHEETNEW,
			Db:        MyCat,
		},
	}
}

//
func (this *TQ_FIN_PROBALSHEETNEW) GetSingleInfo(comc string) (TQ_FIN_PROBALSHEETNEW, error) {
	var tsp TQ_FIN_PROBALSHEETNEW

	err := this.Db.Select("TOTASSET,TOTALCURRLIAB,TOTLIAB,CAPISURP,TOTCURRASSET,RIGHAGGR,PARESHARRIGH").
		From(this.TableName).
		Where("COMPCODE=" + comc).
		Where("REPORTTYPE=1").
		Where("ISVALID=1").
		OrderBy("ENDDATE DESC ").
		Limit(1).
		LoadStruct(&tsp)
	return tsp, err
}


//查询全部企业负债信息
func (this *TQ_FIN_PROBALSHEETNEW) GetAllInfo() (map[dbr.NullString]TQ_FIN_PROBALSHEETNEW, error) {
	var tsp []TQ_FIN_PROBALSHEETNEW
	var tspmap map[dbr.NullString]TQ_FIN_PROBALSHEETNEW
	err := this.Db.Select("TOTASSET,TOTALCURRLIAB,TOTLIAB,CAPISURP,TOTCURRASSET,RIGHAGGR,PARESHARRIGH,COMPCODE").
		From(this.TableName).
		Where("ID in (select max(ID) from tq_fin_probalsheetnew where ISVALID=1 and REPORTTYPE=1  group by COMPCODE)").
		LoadStruct(&tsp)

	//转map
	tspmap = make(map[dbr.NullString]TQ_FIN_PROBALSHEETNEW)
	for _, v := range tsp{
		tspmap[v.COMPCODE] = v
	}
	return tspmap, err
}
