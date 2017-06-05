// 证券代码表
package tb_security

import (
	"haina.com/market/hqinit/servers"
	"haina.com/share/logging"
)

type SecurityCode struct {
	SID int32 `bson:"nSID"`
}

// 证券代码表
func GetSecurityCodeTableFromMG() *[]*SecurityCode {
	var codes []*SecurityCode
	tagS := new(servers.TagSecurityInfo).GetStockInfo("s1")
	for _, item := range tagS {
		var sc SecurityCode
		sc.SID = item.NSID
		codes = append(codes, &sc)
	}
	logging.Debug("lenght of sidcode tables:%v", len(codes))

	return &codes
}
