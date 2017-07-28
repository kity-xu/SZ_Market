package fcmysql

import (
	_ "github.com/go-sql-driver/mysql"
	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

type TQ_COMP_BOARDMAP struct {
	Model     `db:"-"`
	BOARDCODE dbr.NullString `db:"BOARDCODE"` // 对照板块编码
	BOARDNAME dbr.NullString `db:"BOARDNAME"` // 对照板块名称
	COMPCODE  dbr.NullString `db:"COMPCODE"`  // 公司内码
	KEYCODE   dbr.NullString `db:"KEYCODE"`   // 对照键值代码
	KEYNAME   dbr.NullString `db:"KEYNAME"`   // 对照键值名称
	SECODE    dbr.NullString `db:"SECODE"`    // 证券内码
}

func NewTQ_COMP_BOARDMAP() *TQ_COMP_BOARDMAP {
	return &TQ_COMP_BOARDMAP{
		Model: Model{
			TableName: TABLE_TQ_COMP_BOARDMAP,
			Db:        MyCat,
		},
	}
}

// 查询一级板块
func (this *TQ_COMP_BOARDMAP) GetBoardmapList() ([]TQ_COMP_BOARDMAP, error) {
	var boa []TQ_COMP_BOARDMAP
	_, err := this.Db.Select("*").From("TQ_COMP_BOARDMAP").
		Where("BOARDCODE in('1102','1105','1109') and ISVALID =1").
		GroupBy("BOARDCODE").
		LoadStructs(&boa)
	return boa, err
}

// 查询一级板块 入 redis
func (this *TQ_COMP_BOARDMAP) GetBoardmapRedis() ([]TQ_COMP_BOARDMAP, error) {
	var boa []TQ_COMP_BOARDMAP
	_, err := this.Db.Select("*").From("TQ_COMP_BOARDMAP").
		Where("BOARDCODE in('1102','1105','1109') and ISVALID =1").
		LoadStructs(&boa)
	return boa, err
}

// 二级板块
func (this *TQ_COMP_BOARDMAP) GetBoardmap2List(bocode string) ([]TQ_COMP_BOARDMAP, error) {
	var boa []TQ_COMP_BOARDMAP
	_, err := this.Db.Select("*").From("TQ_COMP_BOARDMAP").
		Where("BOARDCODE in('1102','1105','1109') and ISVALID =1").
		Where("BOARDCODE='" + bocode + "'").
		GroupBy("KEYCODE").
		OrderBy("BOARDCODE").
		LoadStructs(&boa)
	return boa, err
}

// 根据 KEYCODE 查询所属板块证券
func (this *TQ_COMP_BOARDMAP) GetBoardmapInfoList(str string) ([]TQ_COMP_BOARDMAP, error) {
	var boa []TQ_COMP_BOARDMAP
	_, err := this.Db.Select("*").From("TQ_COMP_BOARDMAP").
		Where("KEYCODE='" + str + "'").
		OrderBy("COMPCODE").
		LoadStructs(&boa)
	return boa, err
}

// 根据sid查所属板块
func (this *TQ_COMP_BOARDMAP) GetBoadBySID(sid string) ([]TQ_COMP_BOARDMAP, error) {
	var boa []TQ_COMP_BOARDMAP
	_, err := this.Db.Select("*").From("TQ_COMP_BOARDMAP").
		Where("SECODE='" + sid + "'").
		Where("BOARDCODE in('1102','1105','1109')").
		Where("ISVALID=1").
		LoadStructs(&boa)
	return boa, err
}
