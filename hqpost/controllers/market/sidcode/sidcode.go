package sidcode

import (
	"haina.com/market/hqpost/models/tb_security"
)

func GetSecurityTable() (*[]*tb_security.SecurityCode, error) {
	return tb_security.GetSecurityCodeTableFromMG()
}
