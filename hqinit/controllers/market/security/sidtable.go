//股票代码表（只是股票代码sid这东西）
package security

import (
	//"strconv"

	//. "haina.com/market/hqinit/controllers"

	"haina.com/market/hqinit/models/tb_security"
	//"haina.com/share/store/redis"
)

func getSecurityTable() *[]*tb_security.SecurityCode {
	return tb_security.GetSecurityCodeTableFromMG()
}

func UpdateSecurityCodeTable() {
	//sids := getSecurityTable()
	//redis.Del(REDISKEY_SECURITY_NSID_TABLE)
	//for _, sid := range *sids {
	//	redis.Lpush(REDISKEY_SECURITY_NSID_TABLE, []byte(strconv.Itoa(int(sid.SID))))
	//}
}
