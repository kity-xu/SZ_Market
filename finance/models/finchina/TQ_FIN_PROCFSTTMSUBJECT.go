// 关键指标
// TQ_FIN_PROCFSTTMSUBJECT	TTM现金科目产品表
//------------------------------------------------------------------------------
// 与以下3表组合中的数据构成关键指标数据
// TQ_FIN_PROFINMAININDEX   主要财务指标（产品表）
// TQ_FIN_PROINDICDATA      衍生财务指标（产品表）
// TQ_FIN_PROTTMINDIC       财务数据_TTM指标（产品表）
package finchina

import (
	. "haina.com/share/models"

	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
)

// TQ_FIN_PROCFSTTMSUBJECT	  TTM现金科目产品表
type TQ_FIN_PROCFSTTMSUBJECT struct {
	Model    `db:"-"`
	ENDDATE  dbr.NullString  //Date 	放置本次财报的截止日期
	CASHNETR dbr.NullFloat64 //CACENI       现金及现金等价物 净增加额
}

func NewTQ_FIN_PROCFSTTMSUBJECT() *TQ_FIN_PROCFSTTMSUBJECT {
	return &TQ_FIN_PROCFSTTMSUBJECT{
		Model: Model{
			TableName: TABLE_TQ_FIN_PROCFSTTMSUBJECT,
			Db:        MyCat,
		},
	}
}

// 从 TQ_FIN_PROCFSTTMSUBJECT TTM现金科目产品表        取数据
func (this *TQ_FIN_PROCFSTTMSUBJECT) getListByCompcode(compcode string, report_type int, per_page int, page int, date DateList) ([]TQ_FIN_PROCFSTTMSUBJECT, error) {
	var (
		sli_db []TQ_FIN_PROCFSTTMSUBJECT
		err    error
	)

	builder := this.Db.Select("*").From(this.TableName)
	if report_type != 0 {
		builder.Where("REPORTDATETYPE=?", report_type)
	}
	if len(date) > 0 {
		sets := ""
		for _, v := range date {
			sets += v + ","
		}
		sets = sets[:len(sets)-1]
		builder.Where("ENDDATE in (" + sets + ")")
	}
	err = builder.Where("COMPCODE = ?", compcode).
		Where("REPORTTYPE = 3").
		OrderBy("ENDDATE DESC").
		Limit(uint64(per_page)).
		LoadStruct(&sli_db)
	if err != nil && err != dbr.ErrNotFound {
		return nil, err
	}
	return sli_db, nil
}

func (this *TQ_FIN_PROCFSTTMSUBJECT) GetListByEnddates(scode string, report_type int, per_page int, page int, date DateList) ([]TQ_FIN_PROCFSTTMSUBJECT, error) {
	sc := NewSymbolToCompcode()
	if err := sc.getCompcode(scode); err != nil {
		logging.Error("%T GetListByEnddates error: %s", *this, err)
		return nil, err
	}
	return this.getListByCompcode(sc.COMPCODE.String, report_type, per_page, page, date)
}
