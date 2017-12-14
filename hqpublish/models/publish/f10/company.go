package f10

import (
	"strconv"

	"time"

	"haina.com/market/hqpublish/models/finchina"
	"haina.com/share/logging"
)

type HN_F10_Company struct {
}

func NewHN_F10_Company() *HN_F10_Company {
	return &HN_F10_Company{}
}

//1.公司详细资料
type Compinfo struct {
	Name       string              `json:"name"`       // 公司名称
	ListTime   int32               `json:"listTime"`   // 上市日期
	IssueVue   float32             `json:"issueVue"`   // 发行价格
	IssueVol   float32             `json:"issueVol"`   // 发行数量
	RegCap     float32             `json:"regCap"`     // 注册资金
	Indus      string              `json:"industry"`   // 公司所属证监会行业（聚源）
	Prov       string              `json:"area"`       // 省份
	Secretary  string              `json:"secretary"`  // 董秘
	Director   string              `json:"director"`   // 董事长
	RegAddress string              `json:"regAddress"` // 注册地址
	MainBus    string              `json:"mainBus"`    // 主营业务
	PTime      string              `json:"pTime"`      // 主营收入构成 日期
	Constitute []*BusiinfoKeyValue `json:"constitute"` // 主营收入构成
	GgLeader   []*Leader           `json:"leader"`     // 高管数据包
}

// 高管信息
type Leader struct {
	Name      string `json:"name"`      // 高管姓名
	Age       int32  `json:"age"`       // 年龄
	Education string `json:"education"` // 学历
	Duty      string `json:"duty"`      // 职位
	Intro     string `json:"intro"`     // 简介
	BeginDate string `json:"beginDate"` // 在职起始日期
}

// 获取公司详细信息
func GetF10Company(scode string) (*Compinfo, error) {

	sc := finchina.NewTQ_OA_STCODE()
	if err := sc.GetCompcode(scode); err != nil {
		return nil, err
	}
	scode = sc.COMPCODE.String

	comp := finchina.NewCompInfo()
	cinfo, err := comp.GetCompInfo(scode)     // 查询公司资料
	industry, err := comp.GetCompTrade(scode) // 查询行业
	if err != nil {
		logging.Debug("%v", err.Error())
		return nil, err
	}
	// 查询上市日期 总股本
	sercui := finchina.NewSecurityInfo()
	securdate, err := sercui.GetSecurityBasicInfo(scode)
	if err != nil {
		logging.Debug("%v", err.Error())
		return nil, err
	}
	listdate, err := strconv.Atoi(securdate.LISTDATE.String) // 上市日期转int
	if err != nil {
		logging.Debug("%v", err.Error())
		return nil, err
	}
	// 主营收入构成
	busiinfo := finchina.NewTQ_SK_BUSIINFO()
	busilist, err := busiinfo.GetBusiInfo(scode)
	if err != nil {
		logging.Debug("%v", err.Error())
		return nil, err
	}
	fbdata := ""
	var busil []*BusiinfoKeyValue
	for i, v := range busilist {
		if i == 0 {
			fbdata = v.ENTRYDATE.String
		}
		var kv BusiinfoKeyValue
		kv.KeyName = v.CLASSNAME.String
		kv.Value = v.TCOREBIZINCOME.Float64
		kv.Ratio = v.COREBIZINCRTO.Float64
		busil = append(busil, &kv)
	}
	// 查询高管信息
	mangdate, err := finchina.NewTQ_COMP_MANAGER().GetManagersFromFC(scode)
	if err != nil {
		logging.Debug("%v", err.Error())
		return nil, err
	}

	hcg := managersToOnly(mangdate)
	str := ""
	for i, v := range hcg {
		if i == 0 {
			str = "'" + v.PERSONALCODE.String + "'"
		} else {
			str += ",'" + v.PERSONALCODE.String + "'"
		}
	}
	// 查询高管详细信息
	person, err := finchina.NewTQ_COMP_MANAGER().GetPersonRecordInfo(str)
	if err != nil {
		logging.Debug("%v", err.Error())
		return nil, err
	}
	// 获取当前年
	year := time.Now().Year()

	var ld []*Leader
	for _, v := range hcg {
		var l Leader
		for _, k := range person {
			if v.PERSONALCODE.String == k.PERSONALCODE.String {
				l.Name = v.CNAME.String
				i, err := strconv.Atoi(k.BIRTHDAY.String[:4])
				if err != nil {
					logging.Error("%v", err)
					continue
				}
				l.Age = int32(year - i)
				l.Education = degreeTransform(k.DEGREE.String)
				l.Duty = v.ACTDUTYNAME.String
				l.Intro = k.MEMO.String
				l.BeginDate = v.BEGINDATE.String
			}
		}
		ld = append(ld, &l)
	}
	// 查询发行价格和数量
	ail, err := finchina.NewTQ_SK_ALLISSUE().GetAllissueL(scode)
	if err != nil {
		logging.Debug("%v", err.Error())
		return nil, err
	}
	t1 := Compinfo{
		Name:       cinfo.COMPNAME.String,
		ListTime:   int32(listdate),
		IssueVue:   float32(ail.ISSPRICE.Float64),
		IssueVol:   float32(ail.ACTISSQTY.Float64),
		RegCap:     float32(cinfo.REGCAPITAL.Float64),
		Indus:      industry,
		Prov:       getProvince(cinfo.REGION.String),
		Secretary:  cinfo.BSECRETARY.String,
		Director:   cinfo.CHAIRMAN.String,
		RegAddress: cinfo.REGADDR.String,
		MainBus:    cinfo.MAJORBIZ.String,
		PTime:      fbdata,
		Constitute: busil,
		GgLeader:   ld,
	}

	return &t1, nil
}

//高管信息表数据去重，取UPDATEDATE最新
func managersToOnly(primal []finchina.TQ_COMP_MANAGER) []finchina.TQ_COMP_MANAGER {
	swap := make(map[string]finchina.TQ_COMP_MANAGER)

	var managers []finchina.TQ_COMP_MANAGER
	for _, v := range primal {
		if _, ok := swap[v.PERSONALCODE.String]; !ok {
			swap[v.PERSONALCODE.String] = v //不存在
		} else { //存在
			update := swap[v.PERSONALCODE.String]
			update.ACTDUTYNAME.String = v.ACTDUTYNAME.String + ", " + update.ACTDUTYNAME.String
			swap[v.PERSONALCODE.String] = update
		}
	}
	var count int = 0
	for _, v := range swap {
		managers = append(managers, v)
		count++
	}

	return managers
}

// 学历转换
func degreeTransform(istr string) string {
	degr := ""
	switch istr {
	case "1":
		degr = "学士"
	case "2":
		degr = "硕士"
	case "3":
		degr = "博士"
	case "99":
		degr = "其他"
	}
	return degr
}
