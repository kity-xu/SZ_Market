package f10

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"haina.com/market/hqpublish/models"
	"haina.com/share/logging"
)

var (
	ERR_URL_PARAM_FORMAT = errors.New("Error URL Parameter format")
	ERR_URL_PARAM_RANGE  = errors.New("Error URL Parameter range")
	ERR_URL_PARAM_VALUE  = errors.New("Error URL Parameter value")
)

func NEW_ERR_URL_PARAM_FORMAT(i interface{}) error {
	return errors.New(fmt.Sprintf("URL Parameter error: %v", i))
}

var _ = fmt.Print

type RequestParam struct {
	SCode   string // 截取股票代码数字部分后的scode
	Market  string // 市场
	Type    int    // 类型(1:一季报 2:中报 3:三季报 4:年报)
	PerPage int    // 每页条数,默认100
	Page    int    // 第几页的页码,默认1
}

func CheckAndNewRequestParam(scode string, stype string, perPage string, spage string) *RequestParam {
	if len(scode) == 0 {
		logging.Error("%v: %s value is %s", ERR_URL_PARAM_FORMAT, models.CONTEXT_SCODE, "null")
		return nil
	}

	var itype, iperPage, ipage int

	scodePrefix, market, err := ParseSCode(scode)
	if err != nil {
		logging.Error("%v: %s=%s", err, models.CONTEXT_SCODE, scode)
		return nil
	}

	if len(stype) == 0 {
		itype = 0
	} else {
		i, err := strconv.Atoi(stype)
		if err != nil {
			logging.Error("%v", err)
			return nil
		}
		if i < 0 || 4 < i {
			logging.Error("%v: %s=%s", ERR_URL_PARAM_RANGE, models.CONTEXT_TYPE, stype)
			return nil
		}
		itype = i
	}

	if len(perPage) == 0 {
		iperPage = 100
	} else {
		i, err := strconv.Atoi(perPage)
		if err != nil {
			logging.Error("%v", err)
			return nil
		}
		if i < 1 {
			logging.Error("%v: %s=%s", ERR_URL_PARAM_RANGE, models.CONTEXT_PERPAGE, perPage)
			return nil
		}
		iperPage = i
	}

	if len(spage) == 0 {
		ipage = 1
	} else {
		i, err := strconv.Atoi(spage)
		if err != nil {
			logging.Error("%v", err)
			return nil
		}
		if i < 1 {
			logging.Error("%v: %s=%s", ERR_URL_PARAM_RANGE, models.CONTEXT_PAGE, spage)
			return nil
		}
		ipage = i
	}

	return &RequestParam{
		SCode:   scodePrefix,
		Market:  market,
		Type:    itype,
		PerPage: iperPage,
		Page:    ipage,
	}
}

// 成功返回 code, market, nil
// 失败返回   "",     "", err
func ParseSCode(scode string) (string, string, error) {
	// scode A股合法格式判断: 代码.市场
	sli := strings.Split(scode, ".")
	if len(sli) != 2 {
		return "", "", ERR_URL_PARAM_FORMAT
	}
	code := sli[0]
	if code == "" {
		return "", "", ERR_URL_PARAM_FORMAT
	}
	market := strings.ToUpper(sli[1])
	switch market {
	case "SH", "SZ":
	default:
		logging.Error("Unknown Stock Market %s", market)
		return "", "", ERR_URL_PARAM_FORMAT
	}

	return code, market, nil
}

