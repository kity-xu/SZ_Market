package financemysql

import (
	"github.com/gocraft/dbr"
)

// 数据对象名称：TQ_SK_SHAREHOLDERNUM    中文名称：股东户数统计

type TQ_SK_SHAREHOLDERNUM struct {
	//	CIRCSKAAMT dbr.NullString `db:"CIRCSKAAMT"` // 无限售流通A股股本
	//	TOTALSHARE dbr.NullString `db:"TOTALSHARE"` // 总股本
	TOTALSHAMT dbr.NullString `db:"TOTALSHAMT"` // 股东总户数
}

// 查询证券信息
func (this *TQ_SK_SHAREHOLDERNUM) GetSingleInfo(sess *dbr.Session, comc string) (TQ_SK_SHAREHOLDERNUM, error) {
	var tss TQ_SK_SHAREHOLDERNUM
	err := sess.Select("TOTALSHAMT").From("TQ_SK_SHAREHOLDERNUM").
		Where("COMPCODE='" + comc + "' and  ISVALID=1").OrderBy("ENDDATE DESC").Limit(1).LoadStruct(&tss)
	return tss, err
}
