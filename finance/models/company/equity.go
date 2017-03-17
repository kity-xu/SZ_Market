package company

import (
	_ "encoding/json"
	"fmt"

	. "haina.com/share/models"

	"haina.com/share/gocraft/dbr"
	"haina.com/share/lib"
	"haina.com/share/lib/crypto"
	"haina.com/share/models/common"
	"haina.com/share/store/redis"
)

type Equity struct {
	Model `db:"-" `
	ID    int64  // GUID
	SCode string // 股票代码
	Date  string // 指标\日期
	TNS   string // 股东总户数（户）
	ANS   string // 户均持股数（股/户）
	APS   string // 户均持股比例（%）
	CRPS  string // 户均持股较上期变化（%）
}

func NewEquity() *Equity {
	return &Equity{
		Model: Model{
			CacheKey:  REDIS_ADVISER,
			TableName: TABLE_EMPLOYEE,
			Db:        MyCat,
		},
	}
}

func NewEquityTx(tx *dbr.Tx) *Adviser {
	return &Equity{
		Model: Model{
			CacheKey:  REDIS_ADVISER,
			TableName: TABLE_EMPLOYEE,
			Db:        MyCat,
			Tx:        tx,
		},
	}
}

type EquityJson struct {
	GUID            string `json:"_id"`
	Avatar          string `json:"avatar"`
	Department      string `json:"department"`        // 部门
	Description     string `json:"description"`       // 简介/说明
	Email           string `json:"email"`             // 邮箱
	Intro           string `json:"intro"`             // 简介
	JobNumber       int    `json:"job_number"`        // 工号
	JobNumberString string `json:"job_number_string"` // 加密工号
	Level           int    `json:"level"`             // 投顾等级
	MaxMember       int    `json:"max_number"`        // 服务人数限制
	MemberID        string `json:"member_id"`         // hn_members.ID
	OfficePhone     string `json:"office_phone"`      // 办公电话
	Pinyin          string `json:"name_pinyin"`       // 姓名拼音
	Position        string `json:"position"`          // 职位
	QCer            string `json:"qcer"`              // 证券资格号
	Recommend       int    `json:"recommend"`         // 推荐标志
	Tags            string `json:"tags"`              // 标签
	TrueName        string `json:"name"`              // 姓名
}

// 获取单条数据
func (this *Equity) GetSingleByMemberID(id int64) error {
	// 在为PHP后端提供相应的刷新redis数据的接口前，请直接从mysql里读取数据
	//	cacheKey := fmt.Sprintf(this.CacheKey, id)

	//	rec, err := redis.Hgetall(cacheKey)
	//	if err == nil && len(rec) > 0 {
	//		if err := MapToStruct(this, rec); err != nil {
	//			redis.Del(cacheKey)
	//		} else {
	//			return nil
	//		}
	//	}

	exps := map[string]interface{}{
		"MemberID=?": id,
	}

	return this.GetSingleByExps(exps)

	//	if err := this.GetSingleByExps(exps); err != nil {
	//		return err
	//	}

	//	return redis.Hmset(cacheKey, StructToMap(this))
}

// 获取单条数据
func (this *Equity) GetSingleByExps(exps map[string]interface{}) error {
	builder := this.Db.Select("*").From(this.TableName)
	err := this.SelectWhere(builder, exps).
		Limit(1).
		LoadStruct(this)
	return err
}

// 获取单条数据
func (this *Equity) GetSingleJsonByMemberID(id int64) (*AdviserJson, error) {
	if err := this.GetSingleByMemberID(id); err != nil {
		return nil, err
	}

	return this.GetJson(this)
}

// 获取单条数据--通过JobNumber
func (this *Equity) GetSingleJsonByJobNumber(jobNumber int) (*AdviserJson, error) {
	exps := map[string]interface{}{
		"JobNumber=?": jobNumber,
	}
	if err := this.GetSingleByExps(exps); err != nil {
		return nil, err
	}

	return this.GetJson(this)
}

// 获取多条数据
func (this *Equity) GetListByExps(exps map[string]interface{}, limit uint64) ([]*Adviser, error) {
	var data []*Adviser
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

// 获取多条JSON数据
func (this *Equity) GetJsonListByExps(exps map[string]interface{}, limit uint64) ([]*AdviserJson, error) {
	data, err := this.GetListByExps(exps, limit)
	if err != nil {
		return nil, err
	}

	jsns := []*AdviserJson{}
	for _, item := range data {
		jsn, err := this.GetJson(item)
		if err != nil {
			return jsns, err
		}
		jsns = append(jsns, jsn)
	}

	return jsns, nil
}

// 获取所有投顾推荐标记的成员列表
func (this *Equity) GetWeeklyRecommend() ([]*Adviser, error) {
	exps := map[string]interface{}{
		"`Recommend`=?": 1,
	}
	bulider := this.Db.Select("*").From(this.TableName)

	var data []*Adviser
	_, err := this.SelectWhere(bulider, exps).LoadStructs(&data)
	if err != nil {
		return nil, err
	}

	return data, err
}

// 获取所有投顾推荐标记的成员列表
func (this *Equity) GetWeeklyRecommendJsonList() ([]*AdviserJson, error) {
	data, err := this.GetWeeklyRecommend()
	if err != nil {
		return nil, err
	}

	jsns := []*AdviserJson{}
	for _, item := range data {
		jsn, err := this.GetJson(item)
		if err != nil {
			return jsns, err
		}
		jsns = append(jsns, jsn)
	}

	return jsns, nil
}

// 获取JSON
func (this *Equity) GetJson(adviser *Adviser) (*AdviserJson, error) {
	var jsn AdviserJson
	if adviser.ID < 1 {
		return &jsn, ErrUndefinedMemberID
	}

	return &AdviserJson{
		GUID:            IDEncrypt(adviser.ID),
		Avatar:          BASE_URL_AVATAR + crypto.GetMD5(fmt.Sprintf("%v", adviser.MemberID), false),
		Department:      adviser.Department,
		Description:     adviser.Description,
		Email:           adviser.Email,
		Intro:           adviser.Intro.String,
		JobNumber:       adviser.JobNumber,
		JobNumberString: IDEncrypt(int64(adviser.JobNumber)),
		Level:           adviser.Level,
		MaxMember:       adviser.MaxMember,
		MemberID:        IDEncrypt(adviser.MemberID),
		OfficePhone:     adviser.OfficePhone,
		Pinyin:          lib.Pinyin(adviser.TrueName),
		Position:        adviser.Position,
		QCer:            adviser.QCer,
		Recommend:       adviser.Recommend,
		Tags:            adviser.Tags,
		TrueName:        adviser.TrueName,
	}, nil
}

func (this *Equity) UpdateOwnerByMemberID(memberId int64, jobNumber int) error {
	exps := map[string]interface{}{
		"ID=?": memberId,
	}
	params := map[string]interface{}{
		"Owner": jobNumber,
	}

	builder := this.Db.Update(TABLE_MEMBERS)
	this.UpdateParams(builder, params)
	_, err := this.UpdateWhere(builder, exps).Exec()

	return err
}

func (this *Equity) ResetMemberCache(id int64) error {
	redis.Del(fmt.Sprintf(REDIS_MEMBERS, id))

	// 使用此方法重新生成缓存数据
	return common.NewMember().GetSingle(id)
}

func (this *Equity) GetAdviserType(id int64) (int, error) {
	var adviserType int

	exps := map[string]interface{}{
		"e.`ID`=?": id,
	}

	builder := this.Db.Select("m.`AdvisorType`").
		From(TABLE_MEMBERS+" AS m").
		Join(this.TableName+" AS e", "e.`MemberID`=m.`ID`")

	err := this.SelectWhere(builder, exps).Limit(1).LoadValue(&adviserType)

	return adviserType, err
}
