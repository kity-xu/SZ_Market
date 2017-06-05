//股票价格指数（stock index）
package security

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"

	"haina.com/market/hqinit/config"
	"haina.com/share/lib"

	sec "ProtocolBuffer/format/securitytable"

	. "haina.com/market/hqinit/controllers"

	"github.com/golang/protobuf/proto"
	"haina.com/market/hqinit/models/tb_security"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

func UpdateIndexTable(cfg *config.AppConfig) {
	var stype, status string

	stocks := tb_security.GetStockIndexTableFromMG()
	buffer := new(bytes.Buffer)
	var err error
	for _, v := range *stocks {
		stype, err = HainaSecurityType(strconv.Itoa(int(v.NSID)), v.SzSType)
		if err != nil {
			logging.Error("%v", err.Error())
		}
		status, err = HainaSecurityStatus(v.SzStatus)
		if err != nil {
			logging.Error("%v", err.Error())
		}

		buf := TagSecurityName{}

		stock := sec.SecurityInfo{}
		stock.NMarket = v.NMarket
		stock.NSID = v.NSID
		stock.SzCUR = v.SzCUR
		stock.SzDESC = v.SzDESC
		stock.SzIndusCode = v.SzIndusCode
		stock.SzISIN = v.SzISIN
		stock.SzPhonetic = v.SzPhonetic
		stock.SzSCName = v.SzSCName
		stock.SzSCode = v.SzSCode
		stock.SzSName = v.SzSName
		stock.SzStatus = status
		stock.SzSType = stype
		stock.SzSymbol = v.SzSymbol

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

		if err := binary.Write(buffer, binary.LittleEndian, buf); err != nil {
			logging.Fatal(err)
		}

		//转PB
		data, err := proto.Marshal(&stock)
		if err != nil {
			logging.Error("Encode protocbuf of week Line error...%v", err.Error())
			return
		}

		key := fmt.Sprintf(REDISKEY_SECURITY_NAME_ID, stock.NSID)
		if err := redis.Set(key, data); err != nil {
			logging.Fatal("%v", err)
		}

	}
	lib.CheckDir(cfg.File.Path)
	file, err := OpenFile(cfg.File.Path + cfg.File.IndexName)
	if err != nil {
		return
	}
	_, err1 := file.Write(buffer.Bytes())
	if err1 != nil {
		logging.Error("Write file error...")
	}

	defer file.Close()

}
