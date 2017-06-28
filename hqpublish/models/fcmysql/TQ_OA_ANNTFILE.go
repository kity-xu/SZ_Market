package fcmysql

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	. "haina.com/share/models"
)

const (
	TABLE_TQ_OA_ANNTFILE = "TQ_OA_ANNTFILE" // 公告信息表
)

type TQ_OA_ANNTFILE struct {
	Model              `db:"-" `
	ANNOUNCEMENTFILEID int32          // 公告文件目录表ID
	RELAID             int32          // 对应编码
	FILELINK           dbr.NullString // 路径相对目录
	FILEEXTNAME        dbr.NullString // 附件扩展名

}

func NewTQ_OA_ANNTFILE() *TQ_OA_ANNTFILE {
	return &TQ_OA_ANNTFILE{
		Model: Model{
			TableName: TABLE_TQ_OA_ANNTFILE,
			Db:        MyCat,
		},
	}
}

func (this *TQ_OA_ANNTFILE) GetAnntfile(fileid string) (*TQ_OA_ANNTFILE, error) {
	var annfile TQ_OA_ANNTFILE
	bulid := this.Db.Select("ANNOUNCEMENTFILEID,RELAID,FILELINK,FILEEXTNAME").
		From(this.TableName).
		Where("RELAID='" + fileid + "'").
		Where("RELATABLE=1").
		Where("ISVALID=1")

	_, err := this.SelectWhere(bulid, nil).LoadStructs(&annfile)

	if err != nil {
		logging.Debug("%v", err)
		return &annfile, err
	}
	return &annfile, err
}
