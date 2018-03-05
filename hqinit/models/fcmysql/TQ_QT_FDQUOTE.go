package fcmysql


import (
_ "github.com/go-sql-driver/mysql"
"haina.com/share/gocraft/dbr"
. "haina.com/share/models"
)

// 数据对象名称：TQ_QT_INDEX    中文名称：指数行情表

type TQ_QT_FDQUOTE  struct {
	Model      `db:"-"`
	LCLOSE dbr.NullFloat64 `db:"LCLOSE"` // 昨收价
	SECODE dbr.NullString  `db:"SECODE"`
}

func NewTQ_QT_FDQUOTE() *TQ_QT_FDQUOTE {
	return &TQ_QT_FDQUOTE{
		Model: Model{
			TableName: TABLE_TQ_QT_FDQUOTE,
			Db:        MyCat,
		},
	}
}


func (this *TQ_QT_FDQUOTE) GetAllInfo() (map[dbr.NullString]TQ_QT_FDQUOTE, error) {
	var tss []TQ_QT_FDQUOTE
	var tssmap map[dbr.NullString]TQ_QT_FDQUOTE
	err := this.Db.Select("LCLOSE,SECODE").From("TQ_QT_FDQUOTE").
		Where("TRADEDATE=date_format(curdate(),'%Y%m%d')").
	//	OrderBy("PUBLISHDATE  DESC").
	//Limit(1).
		LoadStruct(&tss)
	//转map
	tssmap = make(map[dbr.NullString]TQ_QT_FDQUOTE)
	for _, v := range tss{
		tssmap[v.SECODE] = v
	}
	return tssmap, err
}

