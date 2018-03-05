package fcmysql


import (
_ "github.com/go-sql-driver/mysql"
"haina.com/share/gocraft/dbr"
. "haina.com/share/models"
)

// 数据对象名称：TQ_QT_INDEX    中文名称：指数行情表

type TQ_QT_BDQUOTE  struct {
	Model      `db:"-"`
	LCLOSENETPRICE dbr.NullFloat64 `db:"LCLOSENETPRICE"` // 昨收价
	SECODE dbr.NullString  `db:"SECODE"`
}

func NewTQ_QT_BDQUOTE() *TQ_QT_BDQUOTE {
	return &TQ_QT_BDQUOTE{
		Model: Model{
			TableName: TABLE_TQ_QT_FDQUOTE,
			Db:        MyCat,
		},
	}
}


func (this *TQ_QT_BDQUOTE) GetAllInfo() (map[dbr.NullString]TQ_QT_BDQUOTE, error) {
	var tss []TQ_QT_BDQUOTE
	var tssmap map[dbr.NullString]TQ_QT_BDQUOTE
	err := this.Db.Select("LCLOSENETPRICE,SECODE").From("TQ_QT_BDQUOTE").
		Where("TRADEDATE=date_format(curdate(),'%Y%m%d')").
	//	OrderBy("PUBLISHDATE  DESC").
	//Limit(1).
		LoadStruct(&tss)
	//转map
	tssmap = make(map[dbr.NullString]TQ_QT_BDQUOTE)
	for _, v := range tss{
		tssmap[v.SECODE] = v
	}
	return tssmap, err
}

