//机构资料表
package finchina

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/models"
)

type TQ_COMP_INFO struct {
	models.Model `db:"-" `
	ACCFIRM      dbr.NullString //会计师事务所
	OFFICEADDR   dbr.NullString //公司办公地址
	//A股证券代码
	COMPNAME  dbr.NullString //机构全称
	COMPINTRO dbr.NullString //公司简介
	FOUNDDATE dbr.NullString //公司成立日期
	COMPEMAIL dbr.NullString //联系人电子邮箱
	//公司所属证监会行业（聚源）
	LEGREP        dbr.NullString //法人代表
	BIZLICENSENO  dbr.NullString //企业法人营业执照注册号
	MAJORBIZ      dbr.NullString //主营业务
	MANAGER       dbr.NullString //总经理
	BIZSCOPE      dbr.NullString //经营范围
	OFFICEZIPCODE dbr.NullString //公司办公地址邮编
	REGION        dbr.NullString //省份
	COMPSNAME     dbr.NullString //机构简称
	REGADDR       dbr.NullString //注册地址
	COMPTEL       dbr.NullString //公司电话

	REGCAPITAL dbr.NullFloat64 //注册资本
	ORGCODE    dbr.NullString  //机构组织代码
}

func (this *TQ_COMP_INFO) newTQ_COMP_INFO() *TQ_COMP_INFO {
	return &TQ_COMP_INFO{
		Model: models.Model{
			CacheKey:  "redis_key",
			TableName: TABLE_TQ_COMP_INFO,
			Db:        models.MyCat,
		},
	}
}

//获取公司信息数据
func (this *TQ_COMP_INFO) GetCompInfoFromFC(scode string, market string) (*TQ_COMP_INFO, error) {
	comp := this.newTQ_COMP_INFO()

	//根据股票代码获取公司内码
	sc := NewTQ_OA_STCODE()
	if err := sc.getCompcode(scode, market); err != nil {
		return comp, err
	}

	exps := map[string]interface{}{
		"COMPCODE=?": sc.COMPCODE,
		"ISVALID=?":  1,
	}
	builder := comp.Db.Select("*").From(comp.TableName)
	err := comp.SelectWhere(builder, exps).Limit(1).LoadStruct(comp)
	if err != nil {
		logging.Error("%s", err.Error())
		return comp, err
	}
	logging.Debug("get compinfo success...")
	return comp, err
}
