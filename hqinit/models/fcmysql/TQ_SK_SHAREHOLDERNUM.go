package fcmysql

import (
	_ "github.com/go-sql-driver/mysql"
	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

// 数据对象名称：TQ_SK_SHAREHOLDERNUM    中文名称：股东户数统计

type TQ_SK_SHAREHOLDERNUM struct {
	Model      `db:"-"`
	TOTALSHAMT dbr.NullString `db:"TOTALSHAMT"` // 股东总户数
}

func NewTQ_SK_SHAREHOLDERNUM() *TQ_SK_SHAREHOLDERNUM {
	return &TQ_SK_SHAREHOLDERNUM{
		Model: Model{
			TableName: TABLE_TQ_SK_SHAREHOLDERNUM,
			Db:        MyCat,
		},
	}
}

// 查询证券信息
func (this *TQ_SK_SHAREHOLDERNUM) GetSingleInfo(comc string) (TQ_SK_SHAREHOLDERNUM, error) {
	var tss TQ_SK_SHAREHOLDERNUM
	err := this.Db.Select("TOTALSHAMT").From(this.TableName).
		Where("COMPCODE='" + comc + "'").
		Where("ISVALID=1").
		OrderBy("ENDDATE DESC").
		Limit(1).
		LoadStruct(&tss)
	return tss, err
}
