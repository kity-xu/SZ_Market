package finchina

import (
	"fmt"

	"github.com/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/models"
)

type TQ_SK_BUSIINFO struct {
	models.Model   `db:"-" `
	CLASSNAME      dbr.NullString  // 分类值
	TCOREBIZINCOME dbr.NullFloat64 // 本期主营业务收入(万元)
	COREBIZINCRTO  dbr.NullFloat64 // 占主营业务收入比例(%)
	ENTRYDATE      dbr.NullString  // 日期
}

func NewTQ_SK_BUSIINFO() *TQ_SK_BUSIINFO {
	return &TQ_SK_BUSIINFO{
		Model: models.Model{
			CacheKey:  "redis_key",
			TableName: TABLE_TQ_SK_BUSIINFO,
			Db:        models.MyCat,
		},
	}
}

func (this *TQ_SK_BUSIINFO) GetBusiInfo(scode string) ([]*TQ_SK_BUSIINFO, error) {
	var info []*TQ_SK_BUSIINFO

	builder := this.Db.Select("CLASSNAME,TCOREBIZINCOME,COREBIZINCRTO ,ENTRYDATE").
		From(this.TableName).
		Where(fmt.Sprintf("COMPCODE ='%v'", scode)).
		Where("typestyle = '2' AND ISVALID ='1'").
		Where(fmt.Sprintf("publishdate  =(select publishdate  FROM tq_sk_busiinfo WHERE COMPCODE ='%v' AND typestyle = '2' AND ISVALID ='1' ORDER BY publishdate DESC LIMIT 1)", scode)).
		OrderBy("TCOREBIZINCOME DESC")

	_, err := this.SelectWhere(builder, nil).
		LoadStructs(&info)
	if err != nil {
		logging.Error("%s", err.Error())
		return info, err
	}
	//logging.Debug("get Busi info success...")
	return info, err
}
