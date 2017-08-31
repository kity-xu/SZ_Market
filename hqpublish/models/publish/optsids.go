package publish

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"
	"fmt"
	"strconv"
	"strings"
	"time"

	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	. "haina.com/share/models"

	. "haina.com/market/hqpublish/models"
)

type MOptSids struct {
	Model    `db:"-"`
	MemberID int            // 会员ID
	OptStock dbr.NullString // sid列表
}

func NewMOptSids() *MOptSids {
	return &MOptSids{
		Model: Model{
			TableName: TABLE_HN_OPT_STOCK,
			Db:        DBmicrolink,
		},
	}
}

// get
func (this *MOptSids) SelectAllSidsByAccessToken(access_token string) (*protocol.PayloadOptstockGet, error) {

	mid, err := this.GetMemberIDByAccesstoken(access_token)
	if err != nil {
		return nil, err
	}
	nsids, err := this.selectSidsByMemberID(mid)
	if err != nil {
		return nil, err
	}

	paystocksids := protocol.PayloadOptstockGet{
		SidList: nsids,
	}
	return &paystocksids, err
}

// post
func (this *MOptSids) OperationStockSids(req *protocol.RequestOptstockPut, access_token string) error {
	var isMidExist bool = false

	mid, err := this.GetMemberIDByAccesstoken(access_token)
	if err != nil {
		return err
	}
	var sids string
	for _, v := range req.Sids {
		sids += ("," + strconv.Itoa(int(v)))
	}
	if len(sids) > 2 {
		sids = sids[1:]
	}
	params := map[string]interface{}{
		"MemberID":   mid,
		"OptStock":   sids,
		"UpdateDate": time.Now().Unix(),
	}

	mids, err := this.SelectMemberIDs()
	if err != nil {
		return err
	}
	if len(*mids) == 0 {
		if err = this.InsertStockSidList(params); err != nil {
			return err
		}
		return nil
	}

	for _, memberID := range *mids {
		if memberID == mid {
			isMidExist = true
			break
		}
	}

	var tp int
	var param_log map[string]interface{}
	if isMidExist {
		tp = 3
		if err = this.UpdateStockSidList(params); err != nil {
			return err
		}
		param_log = map[string]interface{}{
			"MemberID": mid,
			"Sids":     sids,
			"OpeType":  tp,
			"OpeTime":  time.Now().Unix(),
		}
	} else {
		tp = 4
		if err = this.InsertStockSidList(params); err != nil {
			return err
		}
		param_log = map[string]interface{}{
			"MemberID": mid,
			"Sids":     sids,
			"OpeType":  tp,
			"OpeTime":  time.Now().Unix(),
		}
	}
	if err = NewOptSidlog().insertOptLog(param_log); err != nil {
		logging.Debug("%v", err.Error())
		return err
	}
	return nil
}

// 单条删除(新增)自选股
func (this *MOptSids) OperateSids(req *protocol.RequestOptstockOperate, access_token string) error {
	mid, err := this.GetMemberIDByAccesstoken(access_token)
	if err != nil {
		return err
	}
	nsids, err := this.selectSidsByMemberID(mid)
	if err != nil {
		return err
	}

	var upsids string
	var tp int

	switch req.Type {
	case 1: //删除
		tp = 1
		for _, nsid := range nsids {
			if req.Sid != nsid {
				upsids += ("," + strconv.Itoa(int(nsid)))
			}
		}
	case 2: //新增
		tp = 2
		for _, nsid := range nsids {
			upsids += ("," + strconv.Itoa(int(nsid)))
		}
		upsids += ("," + strconv.Itoa(int(req.Sid)))
	}

	upsids = upsids[1:]
	params := map[string]interface{}{
		"MemberID":   mid,
		"OptStock":   upsids,
		"UpdateDate": time.Now().Unix(),
	}

	if err = this.UpdateStockSidList(params); err != nil {
		logging.Debug("%v", err)
		return err
	}

	// 操作入日志表
	param_log := map[string]interface{}{
		"MemberID": mid,
		"Sids":     strconv.Itoa(int(req.Sid)),
		"OpeType":  tp,
		"OpeTime":  time.Now().Unix(),
	}

	if err = NewOptSidlog().insertOptLog(param_log); err != nil {
		logging.Debug("%v", err.Error())
		return err
	}
	return nil
}

// 查询会员sids
func (this *MOptSids) selectSidsByMemberID(mid int) ([]int32, error) {
	var optstocks string
	builder := this.Db.Select("OptStock").From(this.TableName).Where("MemberID=" + strconv.Itoa(mid))
	if err := this.SelectWhere(builder, nil).LoadValue(&optstocks); err != nil {
		return nil, err
	}

	if len(optstocks) == 0 {
		logging.Debug("%v", MYSQL_NOT_FIND)
		return nil, MYSQL_NOT_FIND
	}

	var nsids []int32
	sids := strings.Split(optstocks, ",")

	for _, sid := range sids {
		nsid, _ := strconv.Atoi(sid)
		nsids = append(nsids, int32(nsid))
	}
	return nsids, nil
}

// Mysql 查询所有会员ID
func (this *MOptSids) SelectMemberIDs() (*[]int, error) {
	var memberids []int
	builder := this.Db.Select("MemberID").From(this.TableName)
	_, err := this.SelectWhere(builder, nil).LoadValues(&memberids)
	if err != nil {
		return nil, err
	}
	return &memberids, nil
}

// Mysql 插入一条新数据
func (this *MOptSids) InsertStockSidList(params map[string]interface{}) error {
	builder := this.Db.InsertInto(this.TableName)
	_, err := this.InsertParams(builder, params).Exec()
	return err
}

// Mysql 更新一条数据
func (this *MOptSids) UpdateStockSidList(params map[string]interface{}) error {
	var id int

	switch v := params["MemberID"].(type) {
	case int:
		id = v
	case int32:
		id = int(v)
	case int64:
		id = int(v)
	default:
		return fmt.Errorf("params is null")
	}

	builder := this.Db.Update(this.TableName).Where("MemberID=" + strconv.Itoa(id))
	_, err := this.UpdateParams(builder, params).Exec()
	return err
}

// Redis 查找会员ID（memberID）
func (this *MOptSids) GetMemberIDByAccesstoken(access_token string) (int, error) {
	key := fmt.Sprintf(REDIS_ACCESS_TOKEN_MEMBERID, access_token)
	id, err := RedisML.GetString(key)
	if err != nil || len(id) == 0 {
		return -1, REDIS_MEMBERID_NOT_FIND
	}

	memberID := IDDecrypt(id)
	return int(memberID), nil
}
