package valueModel

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/models"
)

type Num struct {
	models.Model `db:"-"`
	COMPCODE     dbr.NullInt64 `db:"COMPCODE"`   //公司内码
	ENDDATE      dbr.NullInt64 `db:"ENDDATE"`    //截止日期
	TOTALSHAMT   dbr.NullInt64 `db:"TOTALSHAMT"` //股东人数
}

func NewNum() *Num {
	return &Num{
		Model: models.Model{
			TableName: "",
			Db:        models.MyCat,
		},
	}
}
func (this *Num) GetNumChartData(compcode string) ([]*Num, error) {
	var data []*Num
	builder := this.Db.SelectBySql(`
	          select
                b.COMPCODE,
                b.ENDDATE,
                b.PUBLISHDATE,
                b.UPDATEDATE,
                b.TOTALSHAMT
               from (select
                a.COMPCODE,
                a.ENDDATE,
                a.PUBLISHDATE,
                a.UPDATEDATE,
                a.TOTALSHAMT
                from TQ_SK_SHAREHOLDERNUM as a
                where  a.COMPCODE = ?  and a.ISVALID = 1
                order by a.ENDDATE desc limit 20) as b
                group by b.PUBLISHDATE order by b.ENDDATE desc`, compcode)
	_, err := this.SelectWhere(builder, nil).LoadStructs(&data)
	if err != nil {
		logging.Info("===", err.Error())
		return data, err
	}
	return data, err
}
func (this *Num) GetNumTextData(swlevelcode string) ([]*Num, error) {
	data := []*Num{}
	builder := this.Db.SelectBySql(`
	        select
                c.COMPCODE,
                c.ENDDATE,
                c.PUBLISHDATE,
                c.UPDATEDATE,
                c.TOTALSHAMT
                from  (select
                        b.COMPCODE,
                        b.ENDDATE,
                        b.PUBLISHDATE,
                        b.UPDATEDATE,
                        b.TOTALSHAMT
            from  TQ_SK_BASICINFO as a
            LEFT JOIN TQ_SK_SHAREHOLDERNUM as b on a.COMPCODE=b.COMPCODE
            WHERE a.SWLEVEL1CODE=? and a.ISVALID = 1 and a.LISTSTATUS = 1 and a.SETYPE = 101 and   a.EXCHANGE in ( '001002' , '001003')
            and b.ISVALID = 1  order by b.TOTALSHAMT desc) c
            group by c.COMPCODE
	`, swlevelcode)
	_, err := this.SelectWhere(builder, nil).LoadStructs(&data)
	if err != nil {
		logging.Info("==", err.Error())
		return data, err
	}
	return data, err
}
