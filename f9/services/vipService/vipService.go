package vipService

type vipJson struct {
	Id           int64       `json:"id"`
	MemberId     int64       `json:"memberId"`
	ProductId    int64       `json:"productId"`
	Order        int64       `json:"order"`
	ProductName  string      `json:"productName"`
	Description  string      `json:"description"`
	FriendlyName string      `json:"abviserName"`
	Price        interface{} `json:"Price"`
}

func GetVipList() ([]*vipJson, error) {
	data, err := NewVip().GetVipList()
	var json []*vipJson
	for _, v := range data {
		var j vipJson
		j.Id = v.Id
		j.MemberId = v.MemberId
		j.ProductId = v.ProductId
		j.Order = v.Order
		j.ProductName = v.ProductName.String
		j.Description = v.Description.String
		j.FriendlyName = v.FriendlyName.String
		price, _ := GetVipPrice(v.ProductId)
		j.Price = price
		json = append(json, &j)
	}
	return json, err
}

type priceJson struct {
	Id           int64   `db:"id"`
	Price        float64 `db:"price"`
	Prices       float64 `db:"prices"`
	Num          int64   `db:"num"`
	DiscountProb float64 `db:"discountProb"`
	CreateTime   int64   `db:"createTime"`
}

func GetVipPrice(productId int64) ([]*priceJson, error) {
	data, err := NewVipPrice().GetVipPrice(productId)
	var js []*priceJson
	//var chu float64 = 100
	for _, v := range data {
		var j priceJson
		j.Id = v.Id
		j.Num = v.Num
		j.CreateTime = v.CreateTime
		j.Price = float64(v.Price) / float64(100)
		j.Prices = v.Price
		j.DiscountProb = v.DiscountProb * 100
		js = append(js, &j)
	}
	return js, err
}
