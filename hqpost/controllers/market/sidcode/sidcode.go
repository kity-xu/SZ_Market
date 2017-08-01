package sidcode

import (
	"haina.com/market/hqpost/models/redistore"
	"haina.com/share/logging"
)

func GetSecurityTable() (*[]int32, error) {
	codes, err := redistore.NewGlobalSid("hq:st:name:*").GetGlobalSidFromRedis()
	if err != nil {
		logging.Error("%v", err.Error())
		return nil, err
	}
	return codes, nil
}
