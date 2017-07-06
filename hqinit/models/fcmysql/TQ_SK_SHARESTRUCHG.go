package fcmysql

import (
	_ "github.com/go-sql-driver/mysql"
	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

// 数据对象名称：TQ_SK_SHARESTRUCHG    中文名称：股本结构变化

type TQ_SK_SHARESTRUCHG struct {
	Model      `db:"-"`
	TOTALSHARE dbr.NullFloat64 `db:"TOTALSHARE"` // 总股本
	CIRCAAMT   dbr.NullFloat64 `db:"CIRCAAMT"`   // 流通A股
	CIRCBAMT   dbr.NullFloat64 `db:"CIRCBAMT"`   // 流通B股
	CIRCHAMT   dbr.NullFloat64 `db:"CIRCHAMT"`   // 流通H股
	CIRCSKRTO  dbr.NullFloat64 `db:"CIRCSKRTO"`  // 流通股合计占总股本比例
	CIRCSKAMT  dbr.NullFloat64 `db:"CIRCSKAMT"`  // 流通股
}

func NewTQ_SK_SHARESTRUCHG() *TQ_SK_SHARESTRUCHG {
	return &TQ_SK_SHARESTRUCHG{
		Model: Model{
			TableName: TABLE_TQ_SK_SHARESTRUCHG,
			Db:        MyCat,
		},
	}
}

// 查询证券信息
func (this *TQ_SK_SHARESTRUCHG) GetSingleInfo(comc string) (TQ_SK_SHARESTRUCHG, error) {
	var tss TQ_SK_SHARESTRUCHG
	err := this.Db.Select("TOTALSHARE,CIRCAAMT,CIRCBAMT,CIRCHAMT,CIRCSKRTO,CIRCSKAMT").From("TQ_SK_SHARESTRUCHG").
		Where("COMPCODE='" + comc + "' and  ISVALID=1").
		OrderBy("PUBLISHDATE  DESC").
		Limit(1).
		LoadStruct(&tss)
	return tss, err
}
