package fcmysql

import (
	"github.com/gocraft/dbr"
)

type TQ_COMP_BOARDMAP struct {
	BOARDCODE dbr.NullString `db:"BOARDCODE"` // 对照板块编码
	BOARDNAME dbr.NullString `db:"BOARDNAME"` // 对照板块名称
	COMPCODE  dbr.NullString `db:"COMPCODE"`  // 公司内码
	KEYCODE   dbr.NullString `db:"KEYCODE"`   // 对照键值代码
	KEYNAME   dbr.NullString `db:"KEYNAME"`   // 对照键值名称
	SECODE    dbr.NullString `db:"SECODE"`    // 证券内码
}

// 查询一级板块
func (this *TQ_COMP_BOARDMAP) GetBoardmapList(sess *dbr.Session) ([]TQ_COMP_BOARDMAP, error) {
	var boa []TQ_COMP_BOARDMAP
	_, err := sess.Select("*").From("TQ_COMP_BOARDMAP").
		Where("BOARDCODE in('1102','1105','1109') and ISVALID =1").
		//BOARDCODE in('1102','1103','1105','1106','1107','1108','1109') and ISVALID =1
		GroupBy("BOARDCODE").
		LoadStructs(&boa)
	return boa, err
}

// 二级板块
func (this *TQ_COMP_BOARDMAP) GetBoardmap2List(sess *dbr.Session, bocode string) ([]TQ_COMP_BOARDMAP, error) {
	var boa []TQ_COMP_BOARDMAP
	_, err := sess.Select("*").From("TQ_COMP_BOARDMAP").
		Where("BOARDCODE in('1102','1105','1109') and ISVALID =1").
		Where("BOARDCODE='" + bocode + "'").
		GroupBy("KEYCODE").
		OrderBy("BOARDCODE").
		LoadStructs(&boa)
	return boa, err
}

// 根据 KEYCODE 查询所属板块证券
func (this *TQ_COMP_BOARDMAP) GetBoardmapInfoList(sess *dbr.Session, str string) ([]TQ_COMP_BOARDMAP, error) {
	var boa []TQ_COMP_BOARDMAP
	_, err := sess.Select("*").From("TQ_COMP_BOARDMAP").
		Where("KEYCODE='" + str + "'").
		OrderBy("COMPCODE").
		LoadStructs(&boa)
	return boa, err
}
