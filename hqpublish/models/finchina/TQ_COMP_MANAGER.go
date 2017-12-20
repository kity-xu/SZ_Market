//公司高管表
package finchina

import (
	"fmt"

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

type TQ_COM_PERSONRECORD struct {
	CNAME        dbr.NullString // 名称
	BIRTHDAY     dbr.NullString // 出生日期
	DEGREE       dbr.NullString // 最高学历
	PERSONALCODE dbr.NullString // 员工唯一码
	MEMO         dbr.NullString // 高管简介

}

func NewTQ_COMP_MANAGER() *TQ_COMP_MANAGER {
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

	exps := map[string]interface{}{
		"COMPCODE=?":  scode,
		"NOWSTATUS=?": 2,
		"ISVALID=?":   1,
	}
	builder := this.Db.Select("*").From(this.TableName).Where("DUTYCODE in ('0200299','0200101','0100101') ")
	_, err := this.SelectWhere(builder, exps).OrderBy("DUTYCODE desc").LoadStructs(&primal)
	if err != nil {
		logging.Error("%s", err.Error())
		return primal, err
	}
	//logging.Debug("dataSize %d:", num)
	return primal, err
}

func (this *TQ_COMP_MANAGER) GetPersonRecordInfo(percode string) ([]*TQ_COM_PERSONRECORD, error) {
	var cper []*TQ_COM_PERSONRECORD
	builder := this.Db.Select("CNAME,BIRTHDAY,DEGREE,PERSONALCODE,MEMO").
		From(TABLE_TQ_COMP_PERSONRECORD).
		Where(fmt.Sprintf("PERSONALCODE in (%v) ", percode))
	_, err := this.SelectWhere(builder, nil).
		//OrderBy("PERSONALCODE").
		LoadStructs(&cper)
	if err != nil {
		logging.Error("%s", err.Error())
		return nil, err
	}
	//logging.Debug("dataSize %d:", num)
	return cper, err
}
