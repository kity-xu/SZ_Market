package sidcode

import (
	"haina.com/market/hqpost/models/mogstore"
	"haina.com/share/logging"
)

func GetSecurityTable() (*[]int32, error) {
	var sids []int32

	codes, err := mogstore.GetSecurityCodeTableFromMG()
	if err != nil {
		logging.Error("%v", err.Error())
		return nil, err
	}

	for _, v := range *codes {
		sids = append(sids, v.SID)
	}
	return &sids, nil
}
