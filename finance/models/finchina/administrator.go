//公司高管信息
package finchina

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/models"
)

type Administrator struct {
}

//高管信息
type AdminInfo struct {
	models.Model   `db:"-"`
	ID             dbr.NullInt64  //流水号
	SECODE         dbr.NullString //股票内码
	SYMBOL         dbr.NullString //股票代码
	SNAME          dbr.NullString //股票名称
	COMPCODE       dbr.NullString //公司代码
	PSCODE         dbr.NullString //高管代码
	PSNAME         dbr.NullString //高管名称
	REMARK         dbr.NullString //简历
	GENDER         dbr.NullString //性别
	YEAR           dbr.NullString //年龄
	EDUCATIONLEVEL dbr.NullString //教育程度
	POST           dbr.NullString //现任职务
	PROVINCENAME   dbr.NullString //所属地区
	ORGTYPE        dbr.NullString //所有制形式
	CSRCNAME       dbr.NullString //证监会行业名称
	SWNAME         dbr.NullString //申万行业名称
	GICSNAME       dbr.NullString //GICS行业名称
	ISVALID        dbr.NullInt64  //是否有效 1:yes 0:no
	TMSTAMP        dbr.NullInt64  //时间戳
	ENTRYDATE      dbr.NullString //录入日期
	ENTRYTIME      dbr.NullString //录入时间
}

//高管持股变化
type AdminEquityChange struct {
	models.Model   `db:"-"`
	ID             dbr.NullInt64   //流水号
	PUBLISHDATE    dbr.NullString  //信息发布日期
	BEGINDATE      dbr.NullString  //起始日期
	ENDDATE        dbr.NullString  //截止日期
	COMPCODE       dbr.NullString  //公司代码
	SHHOLDERTYPE   dbr.NullString  //股东类型		5:高管
	SHHOLDERCODE   dbr.NullString  //股东代码
	SHHOLDERNAME   dbr.NullString  //股东名称
	SHHOLDERNATURE dbr.NullString  //股东性质		3:高管
	SHARESTYPE     dbr.NullString  //股份类型
	DUTY           dbr.NullString  //职务
	CHANGEDIRE     dbr.NullString  //变动方向	1:增持； 2:减持
	BIDCHGAMT      dbr.NullFloat64 //竞价交易变动数量
	BIDAVGPRICE    dbr.NullFloat64 //竞价交易交易均价
	BLOCKCHGAMT    dbr.NullFloat64 //大宗交易变动数量
	BLOCKAVGPRICE  dbr.NullFloat64 //大宗交易交易均价
	OTHCHGAMT      dbr.NullFloat64 //其他方式变动数量
	OTHAVGPRICE    dbr.NullFloat64 //其他方式交易均价
	TOTCHGAMT      dbr.NullFloat64 //变动数量总计
	TOTAVGPRICE    dbr.NullFloat64 //交易均价总计
	BFSHAREAMT     dbr.NullFloat64 //变动前持股数
	AFSHAREAMT     dbr.NullFloat64 //变动后持股数
	LIMSKAMT       dbr.NullFloat64 //有限售股
	CIRCSKAMT      dbr.NullFloat64 //流通股
	DATASOURCE     dbr.NullString  //信息来源  	1:临时公告	2:交易所网站
	MEMO           dbr.NullString  //备注
	ISVALID        dbr.NullInt64   //是否有效
	TMSTAMP        dbr.NullInt64   //时间戳
	ENTRYDATE      dbr.NullString  //录入日期
	ENTRYTIME      dbr.NullString  //录入时间
}

func (this *Administrator) NewAdminInfo() *AdminInfo {
	return &AdminInfo{
		Model: models.Model{
			CacheKey:  "redis_key",
			TableName: TABLE_TQ_SK_LCPERSON,
			Db:        models.MyCat,
		},
	}
}

func (this *Administrator) NewAdminEquityChange() *AdminEquityChange {
	return &AdminEquityChange{
		Model: models.Model{
			CacheKey:  "redis_key",
			TableName: TABLE_TQ_SK_SHAREHDCHG,
			Db:        models.MyCat,
		},
	}
}

func (this *Administrator) GetPSList(secode string) ([]AdminInfo, error) {
	ps := this.NewAdminInfo()
	var pss []AdminInfo
	//根据股票代码获取公司内码
	sc := NewSymbolToCompcode()
	if err := sc.getCompcode(secode); err != nil {
		return pss, err
	}
	if sc.COMPCODE.Valid == false {
		logging.Error("finchina db: select COMPCODE from %s where SYMBOL='%s'", TABLE_TQ_OA_STCODE, secode)
		return pss, ErrNullComp
	}

	exps := map[string]interface{}{
		"COMPCODE=?": sc.COMPCODE,
	}
	builder := ps.Db.Select("*").From(ps.TableName)
	num, err := ps.SelectWhere(builder, exps).LoadStructs(&pss)
	if err != nil {
		logging.Error("%s", err.Error())
		return pss, err
	}
	logging.Debug("dataSize %d:", num)
	return pss, err
}

func (this *Administrator) GetPSEquityChange(secode string) ([]AdminEquityChange, error) {
	eq := this.NewAdminEquityChange()
	var eqs []AdminEquityChange

	//根据股票代码获取公司内码
	sc := NewSymbolToCompcode()
	if err := sc.getCompcode(secode); err != nil {
		return eqs, err
	}
	if sc.COMPCODE.Valid == false {
		logging.Error("finchina db: select COMPCODE from %s where SYMBOL='%s'", TABLE_TQ_OA_STCODE, secode)
		return eqs, ErrNullComp
	}

	exps := map[string]interface{}{
		"COMPCODE=?":       sc.COMPCODE,
		"SHHOLDERTYPE=?":   5,
		"SHHOLDERNATURE=?": 3,
	}
	builder := eq.Db.Select("*").From(eq.TableName)
	num, err := eq.SelectWhere(builder, exps).OrderBy("ENDDATE desc").LoadStructs(&eqs)
	if err != nil {
		logging.Error("%s", err.Error())
		return eqs, err
	}
	logging.Debug("dataSize %d:", num)

	return eqs, nil
}
