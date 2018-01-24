package vipService

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/models"
)

type Vip struct {
	models.Model `db:"-"`
	Id           int64          `db:"id"`
	MemberId     int64          `db:"memberId"`
	ProductId    int64          `db:"productId"`
	Order        int64          `db:"order"`
	ProductName  dbr.NullString `db:"ProductName"`
	Description  dbr.NullString `db:"Description"`
	FriendlyName dbr.NullString `db:"FriendlyName"`
}

func NewVip() *Vip {
	return &Vip{
		Model: models.Model{
			TableName: "hn_v_advance_service",
			Db:        models.MyCat,
		},
	}
}

func (this *Vip) GetVipList() ([]*Vip, error) {
	var data []*Vip
	exps := map[string]interface{}{
		"v.runStatus=?": 2,
	}
	builder := this.Db.Select("v.id,v.productId,v.memberId,v.order,p.ProductName,p.Description,m.FriendlyName").
		From(this.TableName+" as v").
		LeftJoin("hn_products as p", "v.productId=p.ID").
		LeftJoin("hn_members as m", "v.memberId=m.ID").
		OrderBy("v.id desc")
	_, err := this.SelectWhere(builder, exps).LoadStructs(&data)
	if err != nil {
		logging.Debug("%v", err)
		return data, err
	}
	return data, err
}

type VipPrice struct {
	models.Model `db:"-"`
	Id           int64   `db:"id"`
	Price        float64 `db:"price"`
	Num          int64   `db:"num"`
	DiscountProb float64 `db:"discountProb"`
	CreateTime   int64   `db:"createTime"`
}

func NewVipPrice() *VipPrice {
	return &VipPrice{
		Model: models.Model{
			TableName: "hn_v_product_price",
			Db:        models.MyCat,
		},
	}
}

func (this *VipPrice) GetVipPrice(productId int64) ([]*VipPrice, error) {
	var data []*VipPrice
	exps := map[string]interface{}{
		"productId=?": productId,
	}
	builder := this.Db.Select("*").From(this.TableName)
	_, err := this.SelectWhere(builder, exps).LoadStructs(&data)
	if err != nil {
		logging.Debug("%v", err)
		return data, err
	}
	return data, err
}
