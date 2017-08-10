package fcmysql

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	. "haina.com/share/models"
)

const (
	TABLE_TQ_SK_ANNOUNCEMT = "TQ_SK_ANNOUNCEMT" // 公告信息表
)

type TQ_SK_ANNOUNCEMT struct {
	Model        `db:"-" `
	ID           int32          // 公告ID
	ANNOUNCEMTID dbr.NullString // 对应目录id
	DECLAREDATE  int32          // 公告日期
	ANNTITLE     dbr.NullString // 公告标题
	ANNTYPE      dbr.NullString // 全文类型
	ANNTEXT      dbr.NullString // 公告内容
	LEVEL1       dbr.NullString // 一级分类

}

func NewTQ_SK_ANNOUNCEMT() *TQ_SK_ANNOUNCEMT {
	return &TQ_SK_ANNOUNCEMT{
		Model: Model{
			TableName: TABLE_TQ_SK_ANNOUNCEMT,
			Db:        MyCat,
		},
	}
}

func (this *TQ_SK_ANNOUNCEMT) GetNoticeInfo(ccode string, num int32, date string) ([]*TQ_SK_ANNOUNCEMT, error) {
	var tsa []*TQ_SK_ANNOUNCEMT

	var bulid *dbr.SelectBuilder
	if num != 0 {
		bulid = this.Db.Select("ID,ANNOUNCEMTID,ANNTYPE,DECLAREDATE,ANNTITLE,LEVEL1").
			From(this.TableName).
			Where("COMPCODE='" + ccode + "'").
			Where("DECLAREDATE >='" + date + "'").
			Where("ISVALID=1").
			OrderBy("DECLAREDATE desc").Limit(uint64(num))
	} else {
		bulid = this.Db.Select("ID,ANNOUNCEMTID,ANNTYPE,DECLAREDATE,ANNTITLE,LEVEL1").
			From(this.TableName).
			Where("COMPCODE='" + ccode + "'").
			Where("DECLAREDATE >='" + date + "'").
			Where("ISVALID=1").
			OrderBy("DECLAREDATE desc")
	}

	_, err := this.SelectWhere(bulid, nil).LoadStructs(&tsa)

	if err != nil {
		logging.Debug("%v", err)
		return tsa, err
	}
	return tsa, err
}

func (this *TQ_SK_ANNOUNCEMT) GetHisEvent(heid string) (*TQ_SK_ANNOUNCEMT, error) {
	var tsa TQ_SK_ANNOUNCEMT
	bulid := this.Db.Select("ID,ANNOUNCEMTID,ANNTYPE,DECLAREDATE,ANNTITLE,ANNTEXT,LEVEL1").
		From(this.TableName).
		Where("ID='" + heid + "'").
		Where("ISVALID=1")

	_, err := this.SelectWhere(bulid, nil).LoadStructs(&tsa)

	if err != nil {
		logging.Debug("%v", err)
		return &tsa, err
	}
	return &tsa, err
}
