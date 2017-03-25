// F10 财务分析公用
package company

import (
	"strconv"
	"strings"

	. "haina.com/market/finance/models/company"
)

func CheckAndNewRequestParam(scode string, stype string, perPage string, spage string) *RequestParam {
	if len(scode) == 0 {
		return nil
	}

	var itype, iperPage, ipage int

	scodePrefix := strings.Split(scode, ".")[0]
	if scodePrefix == "" {
		return nil
	}

	if len(stype) == 0 {
		itype = 0
	} else {
		i, err := strconv.Atoi(stype)
		if err != nil || i < 1 || 4 < i {
			return nil
		}
		itype = i
	}

	if len(perPage) == 0 {
		iperPage = 4
	} else {
		i, err := strconv.Atoi(perPage)
		if err != nil || i < 1 {
			return nil
		}
		iperPage = i
	}

	if len(spage) == 0 {
		ipage = 1
	} else {
		i, err := strconv.Atoi(spage)
		if err != nil || i < 1 {
			return nil
		}
		ipage = i
	}

	return &RequestParam{
		SCodeOrigin: scode,
		SCode:       scodePrefix,
		Type:        itype,
		PerPage:     iperPage,
		Page:        ipage,
	}
}
