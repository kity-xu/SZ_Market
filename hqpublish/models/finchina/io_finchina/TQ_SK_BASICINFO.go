// 证券基本信息表
package io_finchina

import (
	. "haina.com/share/models"

	"haina.com/share/gocraft/dbr"

	"haina.com/market/hqpublish/models/finchina"
)

// TQ_OA_STCODE    证券内码表
// ---------------------------------------------------------------------
type TQ_SK_BASICINFO struct {
	Model    `db:"-"`
	LISTDATE dbr.NullString //公司代码(公司内码) 通过 SYMBOL 得到
}

func NewTQ_SK_BASICINFO() *TQ_SK_BASICINFO {
	return &TQ_SK_BASICINFO{
		Model: Model{
			TableName: finchina.TABLE_TQ_SK_BASICINFO,
			Db:        MyCat,
		},
	}
}

func (this *TQ_SK_BASICINFO) GetBaseinfo(secode string) (*TQ_SK_BASICINFO, error) {
	builder := this.Db.Select("LISTDATE").From(this.TableName)

	err := builder.Where("SECODE=?", secode).
		Where("ISVALID=1").
		LoadStruct(this)
	if err != nil { //&& err != dbr.ErrNotFound
		return nil, err
	}
	return this, nil
}
