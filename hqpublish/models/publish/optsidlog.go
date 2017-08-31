package publish

import (
	"haina.com/share/gocraft/dbr"
	. "haina.com/share/models"

	. "haina.com/market/hqpublish/models"
)

type OptSidlog struct {
	Model    `db:"-"`
	MemberID int            // 会员ID
	SidList  dbr.NullString // 自选股ID
	OptStock int            // 操作类型 1：删除  2：新增 3：批量更新	4：批量上传
	OpeTime  int            // 操作时间
}

func NewOptSidlog() *OptSidlog {
	return &OptSidlog{
		Model: Model{
			TableName: TABLE_HN_OPT_STOCKLOG,
			Db:        DBmicrolink,
		},
	}
}

func (this *OptSidlog) insertOptLog(params map[string]interface{}) error {
	builder := this.Db.InsertInto(this.TableName)
	_, err := this.InsertParams(builder, params).Exec()
	return err
}
