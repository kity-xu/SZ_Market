package fcmysql

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

type TQ_OA_NTRDSCHEDULE struct {
	Model  `db:"-"`
	SECODE dbr.NullString `db:"SECODE"` // 证券内码
}

func NewTQ_OA_NTRDSCHEDULE() *TQ_OA_NTRDSCHEDULE {
	return &TQ_OA_NTRDSCHEDULE{
		Model: Model{
			TableName: TABLE_TQ_OA_NTRDSCHEDULE,
			Db:        MyCat,
		},
	}
}

// 查询沪深市场证券代码 个股
func (this *TQ_OA_NTRDSCHEDULE) GetNtrdsList() ([]*TQ_OA_NTRDSCHEDULE, error) {

	time := time.Now().Format("2006-01-02 00:00:00")
	var data []*TQ_OA_NTRDSCHEDULE
	_, err := this.Db.Select("SECODE").
		From(this.TableName).
		Where("SETYPE=101 and ISVALID=1").
		Where("NTRADEBEGDATE  <='" + time + "' and ( (RESUMEDATE='1900-01-01 00:00:00:000' and  NTRADEENDDATE ='1900-01-01 00:00:00:000') or ( NTRADEENDDATE >= '" + time + "') )").
		OrderBy("SECODE").LoadStructs(&data)
	return data, err
}
