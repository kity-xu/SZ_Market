//高管及联系人持股变动情况表
package finchina

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/models"
)

type TQ_COMP_SKHOLDERCHG struct {
	models.Model `db:"-" `
	PUBLISHDATE  dbr.NullString  //信息发布日期
	COMPCODE     dbr.NullString  //公司内码
	BEGINDATE    dbr.NullString  //变动起始日
	ENDDATE      dbr.NullString  //变动截止日
	PSCODE       dbr.NullString  //高管代码
	PSCNAME      dbr.NullString  //高管姓名
	HOLDCODE     dbr.NullString  //持股人代码
	HOLDNAME     dbr.NullString  //持股人
	REALTETYPE   dbr.NullString  //持股人与高管的关系	1:本人  2:父母 	3:配偶	4:子女	5:兄弟姐妹	6:受控法人	9:其他
	DUTY         dbr.NullString  //职务名称
	HOLDBEAMT    dbr.NullFloat64 //变动前持股数
	HOLDAFAMT    dbr.NullFloat64 //变动后持股数
	CHGAMT       dbr.NullFloat64 //期间变动股数
	CHGAVGPRICE  dbr.NullFloat64 //变动均价
	CHGRSN       dbr.NullString  //变动原因类型 		1:出售 	3:配股	4:增发	5:购买	6:送、转	7:债转股	8:股权奖励	9:股权投资	10:股权分置	11:司法划转	12:二级市场买卖	13:大宗交易	14:竞价交易	90:未披露	99:其它
}

type HolderChange struct {
	PUBLISHDATE dbr.NullString  //信息发布日期
	PSCODE      dbr.NullString  //高管代码
	PSCNAME     dbr.NullString  //高管姓名
	HOLDCODE    dbr.NullString  //持股人代码
	HOLDNAME    dbr.NullString  //持股人姓名
	HOLDBEAMT   dbr.NullFloat64 //变动前持股数
	HOLDAFAMT   dbr.NullFloat64 //变动后持股数
}

func (this *TQ_COMP_SKHOLDERCHG) newTQ_COMP_SKHOLDERCHG() *TQ_COMP_SKHOLDERCHG {
	return &TQ_COMP_SKHOLDERCHG{
		Model: models.Model{
			CacheKey:  "redis_key",
			TableName: TABLE_TQ_COMP_SKHOLDERCHG,
			Db:        models.MyCat,
		},
	}
}

func (this *TQ_COMP_SKHOLDERCHG) GetHoldAMTlist(scode string) ([]HolderChange, error) {
	var holders []HolderChange
	amt := this.newTQ_COMP_SKHOLDERCHG()
	sc := NewSymbolToCompcode()
	if err := sc.getCompcode(scode); err != nil {
		return holders, err
	}

	exps := map[string]interface{}{
		"COMPCODE=?": sc.COMPCODE,
	}
	builder := amt.Db.Select("*").From(amt.TableName)
	num, err := amt.SelectWhere(builder, exps).OrderBy("PUBLISHDATE desc").LoadStructs(&holders)
	if err != nil {
		logging.Error("%s", err.Error())
		return holders, err
	}
	logging.Debug("dataSize %d:", num)
	return holders, err
}
