package newService

import "haina.com/market/f9/services"

func GetBigData(sid string) (*commonService.SecBasicInfo, error) {
	data, err := commonService.GetCommonData(sid)
	if err != nil {
		return nil, err
	}

	if data.Status == 0 { // 未上市
		return nil, err
	}

	return data, nil
}
