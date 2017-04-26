// 关键指标
// TQ_FIN_PROFINMAININDEX   主要财务指标（产品表）
//------------------------------------------------------------------------------
// 与以下3表组合中的数据构成关键指标数据
// TQ_FIN_PROINDICDATA      衍生财务指标（产品表）
// TQ_FIN_PROTTMINDIC       财务数据_TTM指标（产品表）
// TQ_FIN_PROCFSTTMSUBJECT	TTM现金科目产品表
package finchina

import (
	. "haina.com/share/models"

	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
)

type TQ_FIN_PROFINMAININDEX struct {
	Model      `db:"-"`
	ENDDATE    dbr.NullString //Date 	放置本次财报的截止日期
	CUR        dbr.NullString //币种 	放置本次财报的币种
	ACCSTACODE dbr.NullString //会计准则 代码 11001 中国会计准则(07年前版) 11002 中国会计准则(07年起版) 本表全为07年后的新会计准则

	//// 每股指标
	EPSFULLDILUTED dbr.NullFloat64 //DEPS      稀释每股收益
	EPSBASIC       dbr.NullFloat64 //EPS       基本每股收益
	NAPS           dbr.NullFloat64 //PSNA      每股净资产(值)

	//// 盈利能力
	//// 偿债能力
	//// 成长能力
	//// 营运能力
	//// 现金状况
	//// 分红能力
	//// 资本结构
	//// 收益质量
	NVALCHGITOTP dbr.NullFloat64 //ONIDTPTTM   价值变动净收益／利润总额_TTM
	NNONOPITOTP  dbr.NullFloat64 //NIVCDTPTTM  营业外收支净额／利润总额_TTM
	OPANITOTP    dbr.NullFloat64 //NPADNRGALNP 经营活动净收益／利润总额_TTM

	//// 杜邦分析

}

func NewTQ_FIN_PROFINMAININDEX() *TQ_FIN_PROFINMAININDEX {
	return &TQ_FIN_PROFINMAININDEX{
		Model: Model{
			TableName: TABLE_TQ_FIN_PROFINMAININDEX,
			Db:        MyCat,
		},
	}
}

type DateList []string

// 从 TQ_FIN_PROFINMAININDEX  主要财务指标（产品表）    取数据
func (this *TQ_FIN_PROFINMAININDEX) getListByCompcode(compcode string, report_type int, per_page int, page int) ([]TQ_FIN_PROFINMAININDEX, error) {
	var (
		sli_db []TQ_FIN_PROFINMAININDEX
		err    error
	)

	//表中 REPORTTYPE 的释义: 放置本次财报的类型（1、3为合并报表；2、4为母公司报表；1和2是第一次披露的期末值，3和4是最新一次披露的数值，结合是否实际披露字段，可得治是否发生过再次披露）
	//	主要财务指标表 REPORTTYPE 记录类型 有1,2,3,4 四种类型
	//  衍生财务指标表 财务数据_TTM指标表 TTM现金科目产品表 REPORTTYPE 记录类型只有 3,4, 没有1,2 类型
	//所以下面统一用REPORTTYPE=3(合并期末调整)
	builder := this.Db.Select("*").From(this.TableName)
	if report_type != 0 {
		builder.Where("REPORTDATETYPE=?", report_type)
	}
	err = builder.Where("COMPCODE=?", compcode).
		Where("ISVALID=1").
		Where("REPORTTYPE=?", 3).
		OrderBy("REPORTTYPE ASC, ENDDATE DESC").
		Paginate(uint64(page), uint64(per_page)).
		LoadStruct(&sli_db)

	if err != nil && err != dbr.ErrNotFound {
		logging.Error("%v", err)
		return nil, err
	}
	return sli_db, nil
}
func (this *TQ_FIN_PROFINMAININDEX) GetList(scode string, report_type int, per_page int, page int) ([]TQ_FIN_PROFINMAININDEX, error) {
	sc := NewTQ_OA_STCODE()
	if err := sc.getCompcode(scode); err != nil {
		logging.Error("%T GetListByEnddates error: %s", *this, err)
		return nil, err
	}
	return this.getListByCompcode(sc.COMPCODE.String, report_type, per_page, page)
}
