package fcmysql

import (
	_ "github.com/go-sql-driver/mysql"
	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

// 数据对象名称：TQ_QT_SKDAILYPRICE    中文名称：股票历史日交易

type TQ_QT_SKDAILYPRICE struct {
	Model      `db:"-"`
	LCLOSE dbr.NullFloat64 `db:"LCLOSE"` // 昨收价
	SECODE dbr.NullString  `db:"SECODE"`
}

func NewTQ_QT_SKDAILYPRICE() *TQ_QT_SKDAILYPRICE {
	return &TQ_QT_SKDAILYPRICE{
		Model: Model{
			TableName: TABLE_TQ_QT_SKDAILYPRICE,
			Db:        MyCat,
		},
	}
}


func (this *TQ_QT_SKDAILYPRICE) GetAllInfo() (map[dbr.NullString]TQ_QT_SKDAILYPRICE, error) {
	var tss []TQ_QT_SKDAILYPRICE
	var tssmap map[dbr.NullString]TQ_QT_SKDAILYPRICE
	err := this.Db.Select("LCLOSE,SECODE").From("TQ_QT_SKDAILYPRICE").
		Where("TRADEDATE=date_format(curdate(),'%Y%m%d')").
	//	OrderBy("PUBLISHDATE  DESC").
	//Limit(1).
		LoadStruct(&tss)
	//转map
	tssmap = make(map[dbr.NullString]TQ_QT_SKDAILYPRICE)
	for _, v := range tss{
		tssmap[v.SECODE] = v
	}
	return tssmap, err
}
