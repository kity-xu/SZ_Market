package fcmysql

import (
	_ "github.com/go-sql-driver/mysql"
	. "haina.com/share/models"
	"haina.com/share/gocraft/dbr"
)

// 数据对象名称：TQ_SK_INTERVALQT    中文名称：股票固定区间行情表

type TQ_SK_INTERVALQT struct {
	Model `db:"-"`
	VOL5D dbr.NullInt64 `db:"VOL5D"` // 五日交易量
	SECODE dbr.NullString  `db:"SECODE"`
}

func NewTQ_SK_INTERVALQT() *TQ_SK_INTERVALQT {
	return &TQ_SK_INTERVALQT{
		Model: Model{
			TableName: TABLE_TQ_SK_INTERVALQT,
			Db:        MyCat,
		},
	}
}

// 查询证券信息
func (this *TQ_SK_INTERVALQT) GetSingleInfo(sec string) (TQ_SK_INTERVALQT, error) {
	var tss TQ_SK_INTERVALQT
	err := this.Db.Select("VOL5D").From(this.TableName).
		Where("SECODE='" + sec + "' and  ISVALID=1").
		Where("VOL5D>0").
		OrderBy("TRADEDATE DESC").
		Limit(1).
		LoadStruct(&tss)
	return tss, err
}

// 查询所有证券5日成交量信息
func (this *TQ_SK_INTERVALQT) GetALLInfo() (map[dbr.NullString]TQ_SK_INTERVALQT, error) {
	var tss []TQ_SK_INTERVALQT
	var tssmap map[dbr.NullString]TQ_SK_INTERVALQT
	err := this.Db.Select("VOL5D,SECODE").From(this.TableName).
		Where(" ID in(select max(ID) from finchina.TQ_SK_INTERVALQT where ISVALID=1 group by SECODE)").
		LoadStruct(&tss)

	//转map
	tssmap = make(map[dbr.NullString]TQ_SK_INTERVALQT)
	for _, v := range tss{
		tssmap[v.SECODE] = v
	}
	return tssmap, err
}
