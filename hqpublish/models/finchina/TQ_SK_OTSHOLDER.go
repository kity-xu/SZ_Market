package finchina

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"

	"fmt"

	. "haina.com/share/models"
)

/**
  十大流通股东接口
  对应数据库表：TQ_SK_OTSHOLDER
  中文名称：流通股东名单
*/

type TQ_SK_OTSHOLDER struct {
	Model           `db:"-" `
	SHHOLDERCODE    dbr.NullString  // 股东代码
	SHHOLDERNAME    string          // 股东名称
	HOLDERAMT       dbr.NullFloat64 // 持股数
	RANK            int             //股东排名
	PCTOFFLOTSHARES dbr.NullFloat64 // 占流通股比
	ENDDATE         int             // 放置本次股东信息的截止日期
	HOLDERSUMCHG    dbr.NullFloat64 //持股数量增减
	ISHIS           int             //上一期股东是否存在
}

func NewTQ_SK_OTSHOLDER() *TQ_SK_OTSHOLDER {
	return &TQ_SK_OTSHOLDER{
		Model: Model{
			TableName: TABLE_TQ_SK_OTSHOLDER,
			Db:        MyCat,
		},
	}
}

/**
  获取结算时间列表
*/
func (this *TQ_SK_OTSHOLDER) GetEndDate(sCode string, edata string, limit int, market string) ([]*TQ_SK_OTSHOLDER, error) {
	var dataTop10 []*TQ_SK_OTSHOLDER
	//根据证券代码获取公司内码
	sc := NewTQ_OA_STCODE()
	if err := sc.getCompcode(sCode); err != nil {
		return dataTop10, err

	}

	bulid := this.Db.Select("DISTINCT(ENDDATE)").
		From(this.TableName).
		Where("COMPCODE=" + sc.COMPCODE.String + edata).
		Where("ISVALID=1").
		OrderBy("ENDDATE desc").Limit(uint64(limit))

	_, err := this.SelectWhere(bulid, nil).LoadStructs(&dataTop10)

	if err != nil {
		logging.Debug("%v", err)
		return dataTop10, err
	}
	return dataTop10, err
}

// 获取十大流通股东 zxw
func (this *TQ_SK_OTSHOLDER) GetOtshTop10L(scode string, limit int32, enddate int) ([]*TQ_SK_OTSHOLDER, error) {
	var data []*TQ_SK_OTSHOLDER

	if enddate < 1 {
		bulid := this.Db.Select("SHHOLDERCODE, SHHOLDERNAME, HOLDERAMT, RANK, PCTOFFLOTSHARES, ENDDATE, HOLDERSUMCHG, ISHIS").
			From(this.TableName).
			Where(fmt.Sprintf("COMPCODE ='%v'", scode)).
			Where("ISVALID='1'").
			OrderBy("ENDDATE DESC, RANK ASC").Limit(uint64(limit))

		_, err := this.SelectWhere(bulid, nil).
			LoadStructs(&data)

		if err != nil {
			logging.Debug("%v", err)
			return data, err
		}
		return data, err
	} else {
		bulid := this.Db.Select("SHHOLDERCODE, SHHOLDERNAME, HOLDERAMT, RANK, PCTOFFLOTSHARES, ENDDATE, HOLDERSUMCHG, ISHIS").
			From(this.TableName).
			Where(fmt.Sprintf("COMPCODE ='%v'", scode)).
			Where("ISVALID='1'").
			Where(fmt.Sprintf("ENDDATE='%v'", enddate)).
			OrderBy("ENDDATE DESC,RANK ASC").Limit(uint64(limit))

		_, err := this.SelectWhere(bulid, nil).
			LoadStructs(&data)

		if err != nil {
			logging.Debug("%v", err)
			return data, err
		}
		return data, err
	}

}

// 获取发布日期
func (this *TQ_SK_OTSHOLDER) GetOtshEndDate(scode string) ([]int, error) {
	var times []int
	bulid := this.Db.Select("ENDDATE").
		From(this.TableName).
		Where(fmt.Sprintf("COMPCODE ='%v'", scode)).
		Where("ISVALID='1'").
		GroupBy("ENDDATE").
		OrderBy("ENDDATE DESC").
		Limit(uint64(5))

	_, err := this.SelectWhere(bulid, nil).
		LoadStructs(&times)

	if err != nil {
		logging.Debug("%v", err)
		return times, err
	}
	return times, err
}
