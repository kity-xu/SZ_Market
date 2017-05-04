package stockcode

import (
	"haina.com/market/hqpost/models/tb_security"
)

//股票代码表
func GetSecurityTable() (*[]*tb_security.SecurityCode, error) {
	return tb_security.GetSecurityCodeTableFromMG()
}
