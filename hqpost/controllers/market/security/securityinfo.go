package security

import (
	sec "ProtocolBuffer/format/securitytable"
	"fmt"
	"time"

	. "haina.com/market/hqpost/controllers"
	"haina.com/share/logging"

	"github.com/golang/protobuf/proto"
	"haina.com/market/hqpost/models/tb_security"
	"haina.com/share/store/redis"
)

type Market struct {
	MarketID int32
	Info     []TagSecurityInfo
}

type TagSecurityInfo struct {
	//SID int32 `bson:"nSID"`
	NSID        int32
	NMarket     int32  // 市场类型
	SzSType     string // 证券类型
	SzStatus    string // 证券状态
	SzSCode     string // 证券代码: 600036.SH
	SzSymbol    string // 证券原始: 600036
	SzISIN      string // 证券国际代码信息
	SzSName     string // 证券名称 (超过24字节部分被省略)
	SzSCName    string // 证券简体中文名称 (美股、港股超过40字节部分被省略
	SzDESC      string // 英文简称
	SzPhonetic  string // 拼音
	SzCUR       string // 币种
	SzIndusCode string // 行业代码
}

//市场代码表
func sh_MarketTable() (*[]tb_security.TagSecurityInfo, error) {
	return tb_security.SH_SecurityInfoTableFromMG()
}

func sz_MarketTable() (*[]tb_security.TagSecurityInfo, error) {
	return tb_security.SZ_SecurityInfoTableFromMG()
}

func UpdateSecurityInfo() error {
	/*******************************SH**************************************/
	sh, err := sh_MarketTable()
	if err != nil {
		return err
	}

	var sec_sh sec.MarketSecurityCodeTable

	var count int
	for _, v := range *sh {
		single := sec.SecurityInfo{}
		single.NMarket = v.NMarket
		single.NSID = v.NSID
		single.SzSType = v.SzSType
		single.SzStatus = v.SzStatus
		single.SzSCode = v.SzSCode
		single.SzSymbol = v.SzSymbol
		single.SzISIN = v.SzISIN
		single.SzSName = v.SzSName
		single.SzSCName = v.SzSCName
		single.SzDESC = v.SzDESC
		single.SzPhonetic = v.SzPhonetic
		single.SzCUR = v.SzCUR
		single.SzIndusCode = v.SzIndusCode

		sec_sh.List = append(sec_sh.List, &single)
		count++
	}
	sec_sh.MarketID = sec.Market_value["SH"]
	sec_sh.Num = int32(count)
	sec_sh.TimeStamp = int32(time.Now().Unix())

	//入redis
	data_sh, err := proto.Marshal(&sec_sh)
	if err != nil {
		logging.Error("Encode protocbuf of week Line error...%v", err.Error())
		return err
	}

	key_sh := fmt.Sprintf(REDISKEY_MARKET_SECURITY_TABLE, sec_sh.MarketID)
	if err := redis.Set(key_sh, data_sh); err != nil {
		logging.Fatal("%v", err)
	}

	/*******************************SZ***************************************/

	sz, err := sz_MarketTable()
	if err != nil {
		return err
	}

	var sec_sz sec.MarketSecurityCodeTable

	count = 0
	for _, v := range *sz {
		single := sec.SecurityInfo{}

		single.NMarket = v.NMarket
		single.NSID = v.NSID
		single.SzSType = v.SzSType
		single.SzStatus = v.SzStatus
		single.SzSCode = v.SzSCode
		single.SzSymbol = v.SzSymbol
		single.SzISIN = v.SzISIN
		single.SzSName = v.SzSName
		single.SzSCName = v.SzSCName
		single.SzDESC = v.SzDESC
		single.SzPhonetic = v.SzPhonetic
		single.SzCUR = v.SzCUR
		single.SzIndusCode = v.SzIndusCode

		sec_sz.List = append(sec_sz.List, &single)
		count++
	}
	sec_sz.MarketID = sec.Market_value["SZ"]
	sec_sz.Num = int32(count)
	sec_sz.TimeStamp = int32(time.Now().Unix())

	//入redis
	data_sz, err := proto.Marshal(&sec_sz)
	if err != nil {
		logging.Error("Encode protocbuf of week Line error...%v", err.Error())
		return err
	}

	key_sz := fmt.Sprintf(REDISKEY_MARKET_SECURITY_TABLE, sec_sz.MarketID)
	if err := redis.Set(key_sz, data_sz); err != nil {
		logging.Fatal("%v", err)
	}

	return nil

}
