package fcmysql

import (
	_ "github.com/go-sql-driver/mysql"
	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

// TQ_FIN_PROFINMAININDEX    主要财务指标

type TQ_FIN_PROFINMAININDEX struct {
	Model   `db:"-"`
	NAPS    dbr.NullFloat64 `db:"NAPS"`    // 每股净值
	ENDDATE dbr.NullString  `db:"ENDDATE"` // 截止日期
}

func NewTQ_FIN_PROFINMAININDEX() *TQ_FIN_PROFINMAININDEX {
	return &TQ_FIN_PROFINMAININDEX{
		Model: Model{
			TableName: TABLE_TQ_FIN_PROFINMAININDEX,
			Db:        MyCat,
		},
	}
}

func (this *TQ_FIN_PROFINMAININDEX) GetSingleInfo(comc string) (TQ_FIN_PROFINMAININDEX, error) {
	var tsp TQ_FIN_PROFINMAININDEX

	err := this.Db.Select("NAPS,ENDDATE").
		From(this.TableName).
		Where("COMPCODE='" + comc + "' and  ISVALID=1").
		OrderBy("ENDDATE DESC ").
		Limit(1).
		LoadStruct(&tsp)
	return tsp, err
}
