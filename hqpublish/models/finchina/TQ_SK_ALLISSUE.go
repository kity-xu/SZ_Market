package finchina

import (
	"fmt"

	"github.com/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/models"
)

type TQ_SK_ALLISSUE struct {
	models.Model `db:"-" `
	ISSPRICE     dbr.NullFloat64 // 发行价格（元）
	ACTISSQTY    dbr.NullFloat64 // 实际发行数量(万股)
	LISTDATE     dbr.NullString  // 新增股份上市日
}

func NewTQ_SK_ALLISSUE() *TQ_SK_ALLISSUE {
	return &TQ_SK_ALLISSUE{
		Model: models.Model{
			CacheKey:  "redis_key",
			TableName: TABLE_TQ_SK_ALLISSUE,
			Db:        models.MyCat,
		},
	}
}

func (this *TQ_SK_ALLISSUE) GetAllissueL(scode string) (*TQ_SK_ALLISSUE, error) {
	var info *TQ_SK_ALLISSUE

	builder := this.Db.Select("ISSPRICE,ACTISSQTY,LISTDATE").
		From(this.TableName).
		Where(fmt.Sprintf("COMPCODE ='%v'", scode)).
		Where("ISSUETYPE='01' and ISVALID=1 AND SETYPE='101'").
		Where("CUR='CNY' AND EXCHANGE in('001002','001003')").
		Limit(1)

	err := this.SelectWhere(builder, nil).
		LoadStruct(&info)
	if err != nil {
		logging.Error("%s", err.Error())
		return info, err
	}
	logging.Debug("get AllissueL success...")
	return info, err
}
