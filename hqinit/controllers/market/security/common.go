package security

import (
	"errors"
)

/*********************finchina************************/

//证券类型
const (
	//股票
	SECURITY_STOCK   = '1'    //股票
	SECURITY_A       = "101"  //A股
	SECURITY_B       = "102"  //B股
	SECURITY_THIRD_A = "103"  //三板A股
	SECURITY_THIRD_B = "104"  //三板B股
	SECURITY_HK      = "105"  //港股
	SECURITY_SEAS_CN = "106"  //海外上市中国股
	SECURITY_OTHER   = "107"  //其他股票
	SECURITY_USA     = "1101" //美股

	//指数
	SECURITY_INDEX     = '7'   //指数
	STOCK_INDEX_CN     = "701" //国内指数
	STOCK_INDEX_HK     = "702" //香港指数
	STOCK_INDEX_GLOBAL = "703" //全球指数
	STOCK_INDEX_USER   = "704" //用户自定义
	STOCK_INDEX_OTHER  = "799" //其他

	//基金
	SECURITY_FUND = '3'  //基金
	SECURITY_CLOSEFUND = "301" //封闭式基金
	SECURITY_OPENFUND  = "302" //开放式基金

	//债券 '401','403','404','405','406','499','413'
	SECURITY_DEBT 			= '4'    //债券
	SECURITY_ZFDEBT 		= "401"  //	政府债券
	SECURITY_ZFZCJGDEBT 	= "403"	 //	政府支持机构债券
	SECURITY_GJKFJGDEBT 	= "404"  //国际开发机构债券
	SECURITY_JRDEBT 		= "405"  //金融债券
	SECURITY_QYDEBT 		= "406"  //企业(公司)债券
	SECURITY_UNKNOWDEBT 	= "499"   //	其他债券类型
	SECURITY_OTCDEBT 		= "413"		//OTC债券
)

//证券状态
const (
	SECURITY_READYLIST = '0' //未上市
	SECURITY_LIST      = '1' //上市
	SECURITY_DELIST    = '2' //退市
	SECURITY_NLIST     = '3' //非上市

)

/***********************haina**************************/
//证券类型
const (
	///第0字节(注意区分证券代码中的市场定义）
	SECURITY_TYPE_UNKNOW    = '-' ///< 未定义
	SECURITY_TYPE_MARKET_SH = '1' ///< 沪市
	SECURITY_TYPE_MARKET_SZ = '2' ///< 深市
	///证券类型第1字节
	SECURITY_TYPE_INDEX   = 'I' ///< 指数
	SECURITY_TYPE_STOCK   = 'S' ///< 股票
	SECURITY_TYPE_FUND    = 'F' ///< 基金
	SECURITY_TYPE_BOND    = 'B' ///< 债券
	SECURITY_TYPE_WARRANT = 'W' ///< 权证
	SECURITY_TYPE_OPTION  = 'O' ///< 期权
	SECURITY_TYPE_FUTURE  = 'U' ///< 期货
	SECURITY_TYPE_OPERATE = 'A' ///< 申购、收购、配号等证券业务
	SECURITY_TYPE_PS      = 'P' ///< 优先股
	///证券类型第2字节
	SECURITY_TYPE_AS         = 'A' ///< A股
	SECURITY_TYPE_BS         = 'B' ///< B股
	SECURITY_TYPE_HS         = 'X' ///< 港股			//徐晓东自定义
	SECURITY_TYPE_CS         = 'C' ///< 国际版
	SECURITY_TYPE_FUND_OPEN  = 'O' ///< 开放式基金类
	SECURITY_TYPE_FUND_CLOSE = 'C' ///< 封闭式基金类
	SECURITY_TYPE_GZ         = 'G' ///< 国债
	SECURITY_TYPE_QYZ        = 'Q' ///< 企债
	SECURITY_TYPE_JRZ        = 'J' ///< 金融债
	SECURITY_TYPE_KZZ        = 'K' ///< 可转债
	SECURITY_TYPE_ZQHG       = 'H' ///< 债券回购
	SECURITY_TYPE_ZY         = 'Z' ///< 质押
	SECURITY_TYPE_AR         = 'A' ///< A股权证
	SECURITY_TYPE_BR         = 'B' ///< B股权证
	SECURITY_TYPE_A_NEW      = '1' ///< A股申购
	SECURITY_TYPE_A_MON      = '2' ///< A股申购款
	SECURITY_TYPE_A_NUM      = '3' ///< A股申购配号
	SECURITY_TYPE_Z_NEW      = '4' ///< 债券申购
	SECURITY_TYPE_Z_MON      = '5' ///< 债券申购款
	SECURITY_TYPE_Z_NUM      = '6' ///< 债券申购配号
	SECURITY_TYPE_AZF        = '7' ///< A股增发
	SECURITY_TYPE_PG         = '8' ///< 配股
	///证券类型第3字节
	SECURITY_TYPE_ZB       = 'Z' ///< A股主板
	SECURITY_TYPE_ZXB      = 'X' ///< A股中小板
	SECURITY_TYPE_CYB      = 'C' ///< A股创业板
	SECURITY_TYPE_FUND_ETF = 'E' ///< ETF基金类
	SECURITY_TYPE_FUND_LOF = 'L' ///< LOF基金类
)

///证券状态 (Security::Status) char[4]
const (
	/// 第0位
	SECURITY_STATUS_UNKNOW = '-' ///< 未定义
	SECURITY_STATUS_NM     = '0' ///< 正常(normal)
	SECURITY_STATUS_FDL    = 'N' ///< 上市首日(First day of listing)
	SECURITY_STATUS_FDR    = 'R' ///< 恢复上市首日(First day of resumption of listing)
	SECURITY_STATUS_DL     = 'D' ///< 退市(Delisting)
	SECURITY_STATUS_SP     = 'S' ///< 停牌(suspended)
	SECURITY_STATUS_SL     = 'L' ///< 长期停牌
	SECURITY_STATUS_TSP    = 'T' ///< 临时停牌
	/// 第1位
	SECURITY_STATUS_ER  = 'R' ///< 除权(ex-rights)
	SECURITY_STATUS_ED  = 'D' ///< 除息(ex-divid)
	SECURITY_STATUS_EDR = 'C' ///< 除权除息
	/// 第2位
	SECURITY_STATUS_XST = '*' ///< *ST
	SECURITY_STATUS_ST  = 'S' ///< ST
	SECURITY_STATUS_DP  = 'P' ///< 退市整理期(Delisting finishing period)
	SECURITY_STATUS_TP  = 'T' ///< 暂停上市后协议转让(Transfer period)
	/// 第3位
	SECURITY_STATUS_VT = 'L' ///< 债券投资者适当性要求类
	SECURITY_STATUS_SR = 'G' ///< 未完成股改(share reform)
)

func HainaSecurityType(nsid, stype string) (string, error) {
	result := make([]byte, 4)

	//第0个字节
	if len(nsid) < 9 {
		return "", errors.New("Security nsid is null...")
	}
	if nsid[0] == SECURITY_TYPE_MARKET_SH {
		result[0] = SECURITY_TYPE_MARKET_SH

	} else if nsid[0] == SECURITY_TYPE_MARKET_SZ {
		result[0] = SECURITY_TYPE_MARKET_SZ
	} else {
		result[0] = SECURITY_TYPE_UNKNOW
	}

	//第1、2个字节
	if len(stype) < 3 {
		return "", errors.New("Security type is null...")
	}
	if stype[0] == SECURITY_STOCK { //股票
		result[1] = SECURITY_TYPE_STOCK
		if stype[2] == '1' { //A股
			result[2] = SECURITY_TYPE_AS

		} else if stype[2] == '2' { //B股
			result[2] = SECURITY_TYPE_BS
		} else {
			return "", errors.New("其他尚未实现的股票类型...")
		}

	} else if stype[0] == SECURITY_INDEX { //指数
		result[1] = SECURITY_TYPE_INDEX
		if stype[2] == '1' { //国内指数
			result[2] = SECURITY_TYPE_AS

		} else if stype[2] == '2' { //香港指数
			result[2] = SECURITY_TYPE_HS
		} else if stype[2] == '3' { //全球指数
			result[2] = SECURITY_TYPE_CS
		} else {
			return "", errors.New("其他尚未实现的指数类型...")
		}
	}else if stype[0] == SECURITY_FUND {
		result[1] = SECURITY_TYPE_FUND
		if stype[2] == '1'{
			result[2] = SECURITY_TYPE_FUND_OPEN
		}else if stype[2] == '2'{
			result[2] = SECURITY_TYPE_FUND_CLOSE
		}
	}else if stype[0] == SECURITY_DEBT{
		result[1] = SECURITY_TYPE_BOND
		if stype[2] == '1' && stype[1] == '0'{
			result[2] = SECURITY_TYPE_GZ
		}else if stype[2] == '3' && stype[1] == '0'{
			result[2] = SECURITY_TYPE_GZ
		}else if stype[2] == '4' && stype[1] == '0'{
			//暂无
			//result[2] = SECURITY_TYPE_GZ
		}else if stype[2] == '5' && stype[1] == '0'{
			result[2] = SECURITY_TYPE_JRZ
		}else if stype[2] == '6' && stype[1] == '0'{
			result[2] = SECURITY_TYPE_QYZ
		}
	}else {
		return "", errors.New("其他尚未实现的金融产品类型...")
	}

	if nsid[3] == '3' && nsid[4] == '0' && nsid[5] >= '0' && nsid[5] <= '9' { //创业板
		result[3] = SECURITY_TYPE_CYB
	} else {
		if nsid[3] == '0' && nsid[4] == '0' && nsid[5] >= '2' && nsid[5] <= '4' {
			result[3] = SECURITY_TYPE_ZXB
		} else {
			result[3] = SECURITY_TYPE_ZB
		}
	}
	return string(result), nil
}

func HainaSecurityStatus(status string) (string, error) {
	result := make([]byte, 4)
	if len(status) < 1 {
		return "", errors.New("Security status is null...")
	}
	sp := []byte(status)[0]

	//	if sp == SECURITY_READYLIST {
	//		result[0] = SECURITY_STATUS_SP
	//	} else if sp == SECURITY_LIST {
	//		result[0] = SECURITY_STATUS_NM

	//	} else if sp == SECURITY_DELIST {
	//		result[0] = SECURITY_STATUS_DL

	//	} else if sp == SECURITY_NLIST {
	if sp == '0' {
		result[0] = SECURITY_STATUS_SP
	} else if sp == '1' {
		result[0] = SECURITY_STATUS_NM
	} else {
		return "", errors.New("其他未知状态...")
	}
	result[1] = SECURITY_STATUS_UNKNOW // 考虑到用SECURITY_STATUS_UNKNOW（0）来代替的话，无占位符，没办法区分
	result[2] = SECURITY_STATUS_UNKNOW // "-" 表示尚不知该状态如何
	result[3] = SECURITY_STATUS_UNKNOW
	return string(result), nil
}
