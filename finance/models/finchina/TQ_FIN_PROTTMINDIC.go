// 关键指标
// TQ_FIN_PROTTMINDIC       财务数据_TTM指标（产品表）
//------------------------------------------------------------------------------
// 与以下3表组合中的数据构成关键指标数据
// TQ_FIN_PROINDICDATA      衍生财务指标（产品表）
// TQ_FIN_PROCFSTTMSUBJECT	TTM现金科目产品表
// TQ_FIN_PROFINMAININDEX   主要财务指标（产品表）
package finchina

import (
	. "haina.com/share/models"

	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
)

type TQ_FIN_PROTTMINDIC struct {
	Model   `db:"-"`
	ENDDATE dbr.NullString //Date 	放置本次财报的截止日期

	//// 每股指标
	OPNCFPS dbr.NullFloat64 //PSNBCFTTM 每股经营活动产生的现金流量净额_TTM
	NCFPS   dbr.NullFloat64 //PSNCFTTM  每股现金流量净额_TTM

	//// 盈利能力
	EBITTOTOPI dbr.NullFloat64 //EBITDRTTM 息税前利润／营业总收入_TTM
	SGPMARGIN  dbr.NullFloat64 //SGPMTTM   销售毛利率_TTM

	//// 偿债能力
	//// 成长能力
	//// 营运能力
	//// 现金状况
	OPANCFTOOPNI    dbr.NullFloat64 //NBAGCFDNETTM 经营活动产生的现金流量净额／经营活动净收益_TTM
	SCASHREVTOOPIRT dbr.NullFloat64 //SGPCRSDRTTM  销售商品提供劳务收到的现金／营业收入_TTM

	//// 分红能力
	//// 资本结构
	//// 收益质量
	//// 杜邦分析
}

func NewTQ_FIN_PROTTMINDIC() *TQ_FIN_PROTTMINDIC {
	return &TQ_FIN_PROTTMINDIC{
		Model: Model{
			TableName: TABLE_TQ_FIN_PROTTMINDIC,
			Db:        MyCat,
		},
	}
}

// 从 TQ_FIN_PROTTMINDIC      财务数据_TTM指标（产品表）取数据
func (this *TQ_FIN_PROTTMINDIC) getListByCompcode(compcode string, report_type int, per_page int, page int, date DateList) ([]TQ_FIN_PROTTMINDIC, error) {
	var (
		sli_db []TQ_FIN_PROTTMINDIC
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
	err = builder.Where("COMPCODE=?", compcode).
		Where("ISVALID=1").
		Where("REPORTTYPE=?", 3).
		OrderBy("ENDDATE DESC").
		Limit(uint64(per_page)).
		LoadStruct(&sli_db)
	if err != nil && err != dbr.ErrNotFound {
		return nil, err
	}
	return sli_db, nil
}

func (this *TQ_FIN_PROTTMINDIC) GetListByEnddates(scode string, market string, report_type int, per_page int, page int, date DateList) ([]TQ_FIN_PROTTMINDIC, error) {
	sc := NewTQ_OA_STCODE()
	if err := sc.getCompcode(scode, market); err != nil {
		logging.Error("%T GetListByEnddates error: %s", *this, err)
		return nil, err
	}
	return this.getListByCompcode(sc.COMPCODE.String, report_type, per_page, page, date)
}
