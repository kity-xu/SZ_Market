package fcmysql

import (
	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

// 数据对象名称：TQ_IX_COMP    中文名称：指数成份股

type TQ_IX_COMP struct {
	Model       `db:"-" `
	SELECTECODE dbr.NullString `db:"SELECTECODE"` // 入选证券内码
	SAMPLECODE  dbr.NullString `db:"SAMPLECODE"`  // 样本券代码
}

func NewTQ_IX_COMP() *TQ_IX_COMP {
	return &TQ_IX_COMP{
		Model: Model{
			TableName: TABLE_TQ_IX_COMP,
			Db:        MyCat,
		},
	}
}

// 查询指数成分股信息
func (this *TQ_IX_COMP) GetIndexStockL(sec string) ([]TQ_IX_COMP, error) {

	var tss []TQ_IX_COMP
	_, err := this.Db.Select("SELECTECODE,SAMPLECODE").From("TQ_IX_COMP").
		Where("SECODE='" + sec + "'").
		Where("USESTATUS=1").
		Where("ISVALID=1").
		OrderBy("SAMPLECODE").LoadStructs(&tss)
	return tss, err
}
