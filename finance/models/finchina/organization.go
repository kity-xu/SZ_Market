package finchina

import (
	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"
)

/**
  机构持股接口
  对应数据库表：TQ_SK_SHAREHOLDER
  中文名称：股东名单
*/

type Organization struct {
	Model  `db:"-" `
	SYMBOL string // 股票代码
	Name   string // 机构名称
	Count  string // 持股股份
	PEI    string // 所占比例
}

func NewOrganization() *Organization {
	return &Organization{
		Model: Model{
			TableName: TABLE_TQ_SK_SHAREHOLDER,
			Db:        MyCat,
		},
	}
}

func NewOrganizationTx(tx *dbr.Tx) *Organization {
	return &Organization{
		Model: Model{
			TableName: TABLE_TQ_SK_SHAREHOLDER,
			Db:        MyCat,
			Tx:        tx,
		},
	}
}

type OrganizationJson struct {
	SCode string `json:"SCode"` // 股票代码
	Name  string `json:"Name"`  // 机构名称
	Count string `json:"Count"` // 持股股份
	PEI   string `json:"PEI"`   // 所占比例
}

// 获取机构持股信息
func (this *Organization) GetListByExps(exps map[string]interface{}, limit uint64) ([]*Organization, error) {
	var data []*Organization
	bulid := this.Db.Select("*").From(this.TableName)
	if limit > 0 {
		bulid = bulid.Limit(limit)
	}
	_, err := this.SelectWhere(bulid, exps).LoadStructs(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}
