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
	COMPCODE  dbr.NullString `db:"COMPCODE"` //
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

// 查询所有证券股东信息
func (this *TQ_SK_SHAREHOLDERNUM) GetAllInfo() (map[dbr.NullString]TQ_SK_SHAREHOLDERNUM, error) {
	var tss []TQ_SK_SHAREHOLDERNUM
	var tssmap map[dbr.NullString]TQ_SK_SHAREHOLDERNUM
	err := this.Db.Select("TOTALSHAMT,COMPCODE").From(this.TableName).
		Where(" ID in (select max(ID) from tq_sk_shareholdernum where ISVALID=1  group by COMPCODE)").
		//Where("ISVALID=1").
		//OrderBy("ENDDATE DESC").
		//Limit(1).
		LoadStruct(&tss)
	//转map
	tssmap = make(map[dbr.NullString]TQ_SK_SHAREHOLDERNUM)
	for _, v := range tss{
		tssmap[v.COMPCODE] = v
	}
	return tssmap, err
}
