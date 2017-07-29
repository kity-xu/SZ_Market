package control

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"haina.com/share/logging"
)

var basetime int32 = 20100101

type FactorGroup struct {
	Fa protocol.Factor
	Ls []*protocol.KInfo
}

// 将K线根据除权数据的日期进行分组, 一组K线属于一条除权数据
func FactorGroupOp(fs []*protocol.Factor, rows []*protocol.KInfo) ([]*FactorGroup, error) {
	// 从数据库取出来的除权因子数组是 下标小->大 时间小->大
	// 从Redis取出来的K线数据数组是  下标小->大 时间大->小
	// 这里反转一下K线数组, 使其符合：下标小->大 时间小->大
	ReverseKList(rows)

	fg, err := GroupRightRecoverKList(fs, rows)
	if err != nil {
		logging.Error("%v", err)
		return nil, err
	}
	CalcBeforeRightRecoverKList(fg)

	if fg == nil {
		return nil, nil
	}

	return fg, nil
}

// 把前复权K线和除权因子进行关联分组
func GroupRightRecoverKList(fs []*protocol.Factor, rows []*protocol.KInfo) ([]*FactorGroup, error) {
	if len(fs) == 0 || len(rows) == 0 {
		return nil, nil
	}

	//	fmt.Printf("k %03d %+v\n", 0, rows[0])
	//	fmt.Printf("k %03d %+v\n", len(rows)-1, rows[len(rows)-1])

	// 计算应使用的除权因子区间
	// 下标左值
	bgn := 0
	for n, v := range fs {
		k := rows[0]
		if int32(v.NBeginDate) <= k.NTime && int32(v.NEndDate) >= k.NTime {
			bgn = n
			break
		}
	}
	// 计算应使用的除权因子区间
	// 下标右值
	end := len(fs) - 1
	for n := len(fs) - 1; n > -1; n-- {
		k := rows[len(rows)-1]
		if int32(fs[n].NBeginDate) <= k.NTime && int32(fs[n].NEndDate) >= k.NTime {
			end = n
			break
		}
	}

	// 创建分组
	var fgs []*FactorGroup
	for n := bgn; n <= end; n++ {
		fg := &FactorGroup{
			Fa: *fs[n],
			Ls: make([]*protocol.KInfo, 0, 200),
		}
		fgs = append(fgs, fg)
	}

	s1 := 0
	for _, v := range fgs {
		m := s1
		for ; m < len(rows); m++ {
			k := rows[m]
			if k.NTime < int32(v.Fa.NBeginDate) || k.NTime > int32(v.Fa.NEndDate) {
				break
			}
		}
		v.Ls = append(v.Ls, rows[s1:m]...)
		s1 = m
		if m == len(rows) {
			break
		}
	}
	return fgs, nil
}

func CalcBeforeRightRecoverKList(fgs []*FactorGroup) {
	if fgs == nil || len(fgs) == 0 {
		return
	}
	for _, v := range fgs {
		for _, k := range v.Ls {
			CalcBeforeRightRecoverKLine(k, v.Fa.DfLTDXDY)
		}
	}
}

func CalcBeforeRightRecoverKLine(k *protocol.KInfo, factor float64) {
	//fmt.Println("Before Right calc origin", factor, k)
	k.NPreCPx = int32(float64(k.NPreCPx) / factor) // 昨收价
	k.NOpenPx = int32(float64(k.NOpenPx) / factor) // 开盘价
	k.NHighPx = int32(float64(k.NHighPx) / factor) // 最高价
	k.NLowPx = int32(float64(k.NLowPx) / factor)   // 最低价
	k.NLastPx = int32(float64(k.NLastPx) / factor) // 收盘价(最新价)

	//fmt.Println("Before Right calc result", factor, k)
}

func GetRangeKList(rows []*protocol.KInfo) ([]*protocol.KInfo, error) {
	n := LocationKLine(rows)
	if n == -1 {
		logging.Info("GetRangeKList base time %d is not found", basetime)
		return nil, nil
	}

	if n != -1 {
		return rows[:n+1], nil
	}

	return nil, nil
}

func LocationKLine(rows []*protocol.KInfo) int {
	// 找时间点K线
	index := LocationBinaryKLine(rows)
	if index != -1 {
		logging.Info("binary search basetime time %d is found with index %d: %+v", basetime, index, rows[index])
		return index
	}
	logging.Info("binary search basetime time %d is not found", basetime)

	for n, v := range rows {
		//fmt.Printf("1 1 n %d - %+v\n", n, *v)
		if basetime >= v.NTime {
			//fmt.Printf("1 2 n %d\n", n)
			if basetime == v.NTime {
				return n
			} else {
				return n - 1
			}
		}
	}

	return -1
}

// 二分查找
func LocationBinaryKLine(rows []*protocol.KInfo) int {
	// 下标小->大  时间大->小
	n0, ni := 0, len(rows)-1
	m := 0
	for n0 <= ni {
		m = (n0 + ni) >> 1
		if rows[m].NTime > basetime {
			n0 = m + 1
		} else if rows[m].NTime < basetime {
			ni = m - 1
		} else {
			return m
		}
	}
	return -1
}

//翻转
func ReverseKList(s []*protocol.KInfo) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
