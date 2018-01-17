package finchina

import (
	"fmt"

	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

/**
  机构持股接口
  对应数据库表：TQ_SK_SHAREHOLDER
  中文名称：股东名单
*/

type TQ_SK_SHAREHOLDER struct {
	Model        `db:"-" `
	SHHOLDERCODE dbr.NullString  //股东代码
	SHHOLDERNAME string          // 股东名称
	RANK         int             // 股东排名
	HOLDERAMT    dbr.NullFloat64 // 持股数量
	HOLDERRTO    dbr.NullFloat64 // 持股数量占总股本比例
	ENDDATE      int             // 截止日期
	CURCHG       dbr.NullFloat64 // 本期变动数量
	ISHIS        int             //是否上一期存在的股东
}

func NewTQ_SK_SHAREHOLDER() *TQ_SK_SHAREHOLDER {
	return &TQ_SK_SHAREHOLDER{
		Model: Model{
			TableName: TABLE_TQ_SK_SHAREHOLDER,
			Db:        MyCat,
		},
	}
}

/***************************以下是移动端f10页面*****************************************/
// 该处是十大股东信息

type ShareHoldersTop10 struct {
	Model        `db:"-" `
	SHHOLDERNAME dbr.NullString  //股东名称
	RANK         dbr.NullInt64   //股东排名
	HOLDERRTO    dbr.NullFloat64 //持股数量占总股本比例
	ENDDATE      dbr.NullString  //截止日期
}

func NewShareHoldersTop10() *ShareHoldersTop10 {
	return &ShareHoldersTop10{
		Model: Model{
			TableName: TABLE_TQ_SK_SHAREHOLDER,
			Db:        MyCat,
		},
	}
}

func (this *ShareHoldersTop10) GetShareHoldersTop10(compCode, diviTime string) (*[]ShareHoldersTop10, error) {
	var top10 []ShareHoldersTop10
	exps := map[string]interface{}{
		"COMPCODE=?": compCode,
		"ENDDATE=?":  diviTime,
		"ISVALID=?":  1,
	}
	builder := this.Db.Select("SHHOLDERNAME,RANK,HOLDERRTO,ENDDATE").From(this.TableName).OrderBy("ENDDATE desc") //变动起始日
	_, err := this.SelectWhere(builder, exps).LoadStructs(&top10)
	if err != nil {
		return nil, err
	}
	return &top10, nil
}

// 十大股东发布日期
func (this *TQ_SK_SHAREHOLDER) GetSharEndDate(comcode string) ([]int, error) {
	var dates []int

	Bulid := this.Db.Select("ENDDATE").
		From(this.TableName).
		Where(fmt.Sprintf("COMPCODE ='%v'", comcode)).
		Where("ISVALID=1").
		GroupBy("ENDDATE").
		OrderBy("ENDDATE desc ").
		Limit(5)

	_, err := this.SelectWhere(Bulid, nil).
		LoadStructs(&dates)

	return dates, err
}

// 查询十大股东信息
func (this *TQ_SK_SHAREHOLDER) GetSharBaseL(comcode string, limit int32, enddate int) ([]*TQ_SK_SHAREHOLDER, error) {
	var shar []*TQ_SK_SHAREHOLDER

	if enddate < 1 {
		Bulid := this.Db.Select("SHHOLDERCODE, SHHOLDERNAME, RANK, HOLDERAMT, HOLDERRTO, CURCHG, ENDDATE, ISHIS").
			From(this.TableName).
			Where(fmt.Sprintf("COMPCODE ='%v'", comcode)).
			Where("ISVALID=1").
			OrderBy("ENDDATE desc, RANK ASC").
			Limit(uint64(limit))

		_, err := this.SelectWhere(Bulid, nil).
			LoadStructs(&shar)

		return shar, err
	}
	Bulid := this.Db.Select("SHHOLDERCODE, SHHOLDERNAME, RANK, HOLDERAMT, HOLDERRTO, CURCHG, ENDDATE").
		From(this.TableName).
		Where(fmt.Sprintf("COMPCODE ='%v'", comcode)).
		Where(fmt.Sprintf("ENDDATE='%v'", enddate)).
		Where("ISVALID=1").
		OrderBy("ENDDATE DESC, RANK ASC").
		Limit(uint64(limit))

	_, err := this.SelectWhere(Bulid, nil).
		LoadStructs(&shar)

	return shar, err

}
