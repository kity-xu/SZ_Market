//股票价格指数（stock index）
package security

import (
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

	for _, v := range *stocks {
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

}
