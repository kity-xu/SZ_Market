// 包含如下两个模块
// 0.2:按市场分类的股票代码表
// 0.3:证券基本信息
package security

import (
	"ProtocolBuffer/projects/hqinit/go/protocol"
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	"haina.com/market/hqinit/config"

	. "haina.com/market/hqinit/controllers"
	"haina.com/share/logging"

	"github.com/golang/protobuf/proto"
	"haina.com/market/hqinit/models/tb_security"
	"haina.com/share/store/redis"
)

type Market struct {
	MarketID int32
	Info     []TagSecurityName
}

type TagSecurityName struct {
	//SID int32 `bson:"nSID"`
	NSID        int32
	NMarket     int32                   // 市场类型
	SzSType     [4]byte                 // 证券类型										len:4
	SzStatus    [4]byte                 // 证券状态										len:4
	SzSCode     [SECURITY_CODE_LEN]byte // 证券代码: 600036.SH							len:SECURITY_CODE_LEN
	SzSymbol    [SECURITY_CODE_LEN]byte // 证券原始: 600036								len:SECURITY_CODE_LEN
	SzISIN      [SECURITY_ISIN_LEN]byte // 证券国际代码信息								len:SECURITY_ISIN_LEN
	SzSName     [SECURITY_NAME_LEN]byte // 证券名称 (超过24字节部分被省略)					len:SECURITY_NAME_LEN
	SzSCName    [SECURITY_NAME_LEN]byte // 证券简体中文名称 (美股、港股超过40字节部分被省略		len:SECURITY_NAME_LEN
	SzDESC      [SECURITY_DESC_LEN]byte // 英文简称										len:SECURITY_DESC_LEN
	SzPhonetic  [SECURITY_CODE_LEN]byte // 拼音											len:SECURITY_CODE_LEN
	SzCUR       [4]byte                 // 币种											len:4
	SzIndusCode [INDUSTRY_CODE_LEN]byte // 行业代码										len:INDUSTRY_CODE_LEN
}

//市场代码表
func MarketTable() (*[]tb_security.TagSecurityInfo, error) {
	return tb_security.GetSecurityInfoTableFromMG()
}

//证券基本信息和单市场的证券代码表的实现
func UpdateSecurityTable(cfg *config.AppConfig) {
	var (
		stype, status string
		//sec_sh, sec_sz sec.MarketSecurityCodeTable
		sec_sh, sec_sz protocol.PayloadMarketSecurityNameTable
	)

	table, err := MarketTable()
	if err != nil {
		logging.Error("%v", err)
		return
	}

	//入文件
	buffer := new(bytes.Buffer)

	for _, v := range *table {
		stype, err = HainaSecurityType(strconv.Itoa(int(v.NSID)), v.SzSType)
		if err != nil {
			logging.Error("%v", err.Error())
		}
		status, err = HainaSecurityStatus(v.SzStatus)
		if err != nil {
			logging.Error("%v", err.Error())
		}

		buf := TagSecurityName{}

		single := protocol.SecurityName{}
		single.NMarket = int32(v.NMarket)
		single.NSID = int32(v.NSID)
		single.SzSType = stype
		single.SzStatus = status
		single.SzSCode = v.SzSCode
		single.SzSymbol = v.SzSymbol
		single.SzISIN = v.SzISIN
		single.SzSName = v.SzSName
		single.SzSCName = v.SzSCName
		single.SzDESC = v.SzDESC
		single.SzPhonetic = v.SzPhonetic
		single.SzCUR = v.SzCUR
		single.SzIndusCode = v.SzIndusCode

		buf.NSID = v.NSID
		buf.NMarket = v.NMarket

		buf.SzSType = StringToByte_4(stype)
		buf.SzStatus = StringToByte_4(status)
		buf.SzSCode = StringToByte_SECURITY_CODE_LEN(v.SzSCode)
		buf.SzSymbol = StringToByte_SECURITY_CODE_LEN(v.SzSymbol)
		buf.SzISIN = StringToByte_SECURITY_ISIN_LEN(v.SzISIN)
		buf.SzSName = StringToByte_SECURITY_NAME_LEN(v.SzSName)
		buf.SzSCName = StringToByte_SECURITY_NAME_LEN(v.SzSCName)
		buf.SzDESC = StringToByte_SECURITY_DESC_LEN(v.SzDESC)
		buf.SzPhonetic = StringToByte_SECURITY_CODE_LEN(v.SzPhonetic)
		buf.SzCUR = StringToByte_4(v.SzCUR)
		buf.SzIndusCode = StringToByte_INDUSTRY_CODE_LEN(v.SzIndusCode)

		//入文件
		if err := binary.Write(buffer, binary.LittleEndian, buf); err != nil {
			logging.Fatal(err)
		}

		/*********************证券基本信息************************/
		//转PB
		data, err := proto.Marshal(&single)
		if err != nil {
			logging.Error("Encode protocbuf of week Line error...%v", err.Error())
			return
		}

		//入redis
		key := fmt.Sprintf(REDISKEY_SECURITY_NAME_ID, single.NSID)
		if err := redis.Set(key, data); err != nil {
			logging.Fatal("%v", err)
		}

		/*************************OVER******************************/

		if v.NMarket == protocol.HAINA_PUBLISH_MARKET_value["SH"] {
			sec_sh.SNList = append(sec_sh.SNList, &single)
		} else if v.NMarket == protocol.HAINA_PUBLISH_MARKET_value["SZ"] {
			sec_sz.SNList = append(sec_sz.SNList, &single)
		} else {
			logging.Error("security info nMarket ID error ...")
			return
		}
	}
	/*************************START******************************/
	file, err := OpenFile(cfg.File.Path + cfg.File.StockName)
	if err != nil {
		return
	}

	_, err1 := file.Write(buffer.Bytes())
	if err1 != nil {
		logging.Error("Write file error...")
	}

	defer file.Close()
	/*************************OVER******************************/

	sec_sh.MarketID = protocol.HAINA_PUBLISH_MARKET_value["SH"]
	sec_sh.Num = int32(len(sec_sh.SNList))
	sec_sh.TimeStamp = int32(time.Now().Unix())

	sec_sz.MarketID = protocol.HAINA_PUBLISH_MARKET_value["SZ"]
	sec_sz.Num = int32(len(sec_sz.SNList))
	sec_sz.TimeStamp = int32(time.Now().Unix())

	//上海入redis
	data_sh, err := proto.Marshal(&sec_sh)
	if err != nil {
		logging.Error("Encode protocbuf of week Line error...%v", err.Error())
		return
	}
	logging.Info("Lengh of SH security table:%v", len(sec_sh.SNList))

	key_sh := fmt.Sprintf(REDISKEY_MARKET_SECURITY_TABLE, sec_sh.MarketID)
	if err := redis.Set(key_sh, data_sh); err != nil {
		logging.Fatal("%v", err)
	}

	//深圳入redis
	data_sz, err := proto.Marshal(&sec_sz)
	if err != nil {
		logging.Error("Encode protocbuf of week Line error...%v", err.Error())
		return
	}
	logging.Info("Lengh of SZ security table:%v", len(sec_sz.SNList))

	key_sz := fmt.Sprintf(REDISKEY_MARKET_SECURITY_TABLE, sec_sz.MarketID)
	if err := redis.Set(key_sz, data_sz); err != nil {
		logging.Fatal("%v", err)
	}

}
