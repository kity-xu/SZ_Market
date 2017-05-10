//股票价格指数（stock index）
package security

import (
	"bytes"
	"encoding/binary"
	"fmt"

	sec "ProtocolBuffer/format/securitytable"

	. "haina.com/market/hqinit/controllers"

	"github.com/golang/protobuf/proto"
	"haina.com/market/hqinit/models/tb_security"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

func UpdateIndexTable() {
	stocks, err := tb_security.GetStockIndexTableFromMG()
	if err != nil {
		logging.Error("%v", err.Error())
		return
	}

	buffer := new(bytes.Buffer)

	for _, v := range *stocks {
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
		stock.SzStatus = v.SzStatus
		stock.SzSType = v.SzSType
		stock.SzSymbol = v.SzSymbol

		buf.NSID = v.NSID
		buf.NMarket = v.NMarket

		buf.SzSType = StringToByte_4(v.SzSType)
		buf.SzStatus = StringToByte_4(v.SzStatus)
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

	file, err := OpenFile("E:/security/index.dat")
	if err != nil {
		return
	}
	_, err1 := file.Write(buffer.Bytes())
	if err1 != nil {
		logging.Error("Write file error...")
	}

	defer file.Close()

}
