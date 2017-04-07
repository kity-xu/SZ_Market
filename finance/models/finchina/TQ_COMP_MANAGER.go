//公司高管表
package finchina

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/models"
)

type TQ_COMP_MANAGER struct {
	models.Model `db:"-" `
	ID           dbr.NullInt64
	UPDATEDATE   dbr.NullString //更新日期
	COMPCODE     dbr.NullString //公司代码
	POSTYPE      dbr.NullString //职位属性  			注：	1	:经营层  2:董事会  3:董秘  4:证券事务代表  5:监事会  6:其他
	DUTYCODE     dbr.NullString //统一职位编码
	DUTYMOD      dbr.NullString //职位性质  			注：	1	:正式  2:暂代
	ACTDUTYNAME  dbr.NullString //实际职位名称
	PERSONALCODE dbr.NullString //人物代码
	CNAME        dbr.NullString //人名
	MGENTRYS     dbr.NullInt64  //经营层入职时董事会届次				放置经营层在入职时董事会的届次，标识这经营层人员是哪届董事会聘请的
	DENTRYS      dbr.NullInt64  //董事会、监事会人员任职所在届...		放置董事会和监事会人员本身任职所在届次
	NOWSTATUS    dbr.NullString //当前的任职状态		注：	2:在职	5:申请离职	6:已离职
	BEGINDATE    dbr.NullString //在职起始日期
	ENDDATE      dbr.NullString //离职日期/当前届次截止日
	DIMREASON    dbr.NullString //离职原因			注：	1:辞职	2:免去	3:换届	4:死亡	5:工作调动	6:出国学习	7:个人原因	8:退休	9:其它	10:股权变动	11:身体健康
	ISRELDIM     dbr.NullInt64  //是否实际离职		注：	0:否		1:是
	MEMO         dbr.NullString //备注信息
}

func (this *TQ_COMP_MANAGER) newTQ_COMP_MANAGER() *TQ_COMP_MANAGER {
	return &TQ_COMP_MANAGER{
		Model: models.Model{
			CacheKey:  "redis_key",
			TableName: TABLE_TQ_COMP_MANAGER,
			Db:        models.MyCat,
		},
	}
}

func (this *TQ_COMP_MANAGER) GetManagersFromFC(scode string) ([]TQ_COMP_MANAGER, error) {
	var primal []TQ_COMP_MANAGER
	mg := this.newTQ_COMP_MANAGER()

	//根据股票代码获取公司内码
	sc := NewTQ_OA_STCODE()
	if err := sc.getCompcode(scode); err != nil {
		return primal, err
	}

	exps := map[string]interface{}{
		"COMPCODE=?":  sc.COMPCODE,
		"NOWSTATUS=?": 2,
	}
	builder := mg.Db.Select("*").From(mg.TableName)
	num, err := mg.SelectWhere(builder, exps).OrderBy("UPDATEDATE desc").LoadStructs(&primal)
	if err != nil {
		logging.Error("%s", err.Error())
		return primal, err
	}
	logging.Debug("dataSize %d:", num)
	return primal, err
}
