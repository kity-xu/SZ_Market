package publish

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"haina.com/share/kityxu/utils"

	ctrl "haina.com/market/hqpublish/controllers"
	. "haina.com/market/hqpublish/models"
	. "haina.com/share/models"

	pro "ProtocolBuffer/projects/hqpublish/go/protocol"

	"github.com/golang/protobuf/proto"

	hsgrr "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

var (
	_ = redis.Init
	_ = GetCache
	_ = ctrl.MakeRespDataByBytes
	_ = errors.New
	_ = fmt.Println
	_ = hsgrr.Bytes
	_ = logging.Info
	_ = bytes.NewBuffer
	_ = binary.Read
	_ = io.ReadFull
)

type XRXD struct {
	Model `db:"-"`
}

func NewXRXD() *XRXD {
	return &XRXD{
		Model: Model{
			CacheKey: REDISKEY_SECURITY_HDAY,
		},
	}
}

// redis list 中单根日K线存储格式为protobuf(徐晓东存入)
// 日K线每一根都进行了PB编码，这里需要对所有K线进行解码
func (this XRXD) SingleDecode(bin []byte) (*pro.KInfo, error) {
	obj := pro.KInfo{}
	if err := proto.Unmarshal(bin, &obj); err != nil {
		logging.Error("%v", err)
		return nil, err
	}
	return &obj, nil
}

// 从 redis 取出后进行解码
func (this XRXD) Decode(lsbin []string) ([]*pro.KInfo, error) {
	ls := make([]*pro.KInfo, 0, len(lsbin))
	for _, v := range lsbin {
		if obj, err := this.SingleDecode([]byte(v)); err != nil {
			logging.Error("%v", err)
			return nil, err
		} else {
			ls = append(ls, obj)
		}
	}
	return ls, nil
}

type FactorGroup struct {
	Fa pro.Factor
	Ls []*pro.KInfo
}

func (this XRXD) ErrDataInvalid(fields string, sid int32, secode string) error {
	return errors.New(fmt.Sprintf("finchina TQ_SK_XDRY fields[%s] invalid by sid[%d], sid[%d]", "NBeginDate", sid, secode))
}

// 二分查找
func (this XRXD) LocationBinaryKLine(req *pro.RequestXRXD, rows []*pro.KInfo) int {
	// 下标小->大  时间大->小
	n0, ni := 0, len(rows)-1
	m := 0
	for n0 <= ni {
		m = (n0 + ni) >> 1
		if rows[m].NTime > req.TimeIndex {
			n0 = m + 1
		} else if rows[m].NTime < req.TimeIndex {
			ni = m - 1
		} else {
			return m
		}
	}
	fmt.Printf("n0 %d ni %d m %d\n", n0, ni, m)
	return -1
}
func (this XRXD) LocationKLine(req *pro.RequestXRXD, rows []*pro.KInfo) int {
	// 找时间点K线
	index := this.LocationBinaryKLine(req, rows)
	if index != -1 {
		fmt.Printf("req time %d is found with index %d: %+v\n", req.TimeIndex, index, rows[index])
		return index
	}
	fmt.Printf("req time %d is no found\n", req.TimeIndex)

	// 按条件找范围内第一根时间点K线
	if req.Direct == 0 {
		// 时间轴减小<-方向向左, but K线的存储模式是 下标小->大 时间大->小
		for n := len(rows) - 1; n > -1; n-- {
			if req.TimeIndex <= rows[n].NTime {
				//fmt.Printf("0 2 n %d\n", n)
				if req.TimeIndex == rows[n].NTime {
					return n
				} else {
					return n + 1
				}
			}
		}

	} else {
		// 向右->时间轴增大
		for n, v := range rows {
			//fmt.Printf("1 1 n %d - %+v\n", n, *v)
			if req.TimeIndex >= v.NTime {
				//fmt.Printf("1 2 n %d\n", n)
				if req.TimeIndex == v.NTime {
					return n
				} else {
					return n - 1
				}
			}
		}

	}
	return -1
}

func (this XRXD) GetRangeKList(req *pro.RequestXRXD, rows []*pro.KInfo) ([]*pro.KInfo, error) {
	n := this.LocationKLine(req, rows)
	if n == -1 {
		return nil, nil
		fmt.Printf("GetRangeKList req time %d no found\n", req.TimeIndex)
	}

	if req.Direct == 0 {
		if n != -1 {
			if req.Num > 0 {
				m := n + int(req.Num)
				if m > len(rows) {
					m = len(rows)
				}
				return rows[n:m], nil
			}
			return rows[n:], nil
		}
	} else {
		if n != -1 {
			if req.Num > 0 {
				m := n + 1 - int(req.Num)
				if m < 0 {
					m = 0
				}
				return rows[m : n+1], nil
			}
			return rows[:n+1], nil
		}
	}
	return nil, nil
}

func (this XRXD) ReverseFactors(s []*pro.Factor) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
func (this XRXD) ReverseKList(s []*pro.KInfo) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func (this XRXD) CalcBeforeRightRecoverKLine(k *pro.KInfo, factor float64) {
	//fmt.Println("Before Right calc origin", factor, k)
	k.NOpenPx = int32(float64(k.NOpenPx) / factor) // 开盘价
	k.NHighPx = int32(float64(k.NHighPx) / factor) // 最高价
	k.NLowPx = int32(float64(k.NLowPx) / factor)   // 最低价
	k.NLastPx = int32(float64(k.NLastPx) / factor) // 收盘价(最新价)
	//fmt.Println("Before Right calc result", factor, k)
}
func (this XRXD) CalcAfterRightRecoverKLine(k *pro.KInfo, factor float64) {
	//fmt.Println("Before Right calc origin", factor, k)
	k.NOpenPx = int32(float64(k.NOpenPx) * factor) // 开盘价
	k.NHighPx = int32(float64(k.NHighPx) * factor) // 最高价
	k.NLowPx = int32(float64(k.NLowPx) * factor)   // 最低价
	k.NLastPx = int32(float64(k.NLastPx) * factor) // 收盘价(最新价)
	//fmt.Println("Before Right calc result", factor, k)
}

// 把前复权K线和除权因子进行关联分组
func (this XRXD) GroupRightRecoverKList(fs []*pro.Factor, rows []*pro.KInfo) ([]*FactorGroup, error) {
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
			Ls: make([]*pro.KInfo, 0, 200),
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
func (this XRXD) CalcBeforeRightRecoverKList(fgs []*FactorGroup) {
	if fgs == nil || len(fgs) == 0 {
		return
	}
	for _, v := range fgs {
		for _, k := range v.Ls {
			this.CalcBeforeRightRecoverKLine(k, v.Fa.DfLTDXDY)
		}
	}
}
func (this XRXD) CalcAfterRightRecoverKList(fgs []*FactorGroup) {
	if fgs == nil || len(fgs) == 0 {
		return
	}
	for _, v := range fgs {
		for _, k := range v.Ls {
			this.CalcAfterRightRecoverKLine(k, v.Fa.DfTHELTDXDY)
		}
	}
}
func (this XRXD) FactorGroupTotal(req *pro.RequestXRXD, fs []*pro.Factor, rows []*pro.KInfo) ([]*FactorGroup, error) {
	var fg []*FactorGroup // *Factor([]*KInfo), 即每一个复权因子关联一组K线
	var err error

	this.ReverseKList(rows)

	if req.Method == 1 {
		if fg, err = this.GroupRightRecoverKList(fs, rows); err != nil {
			logging.Error("%v", err)
			return nil, err
		}
		this.CalcBeforeRightRecoverKList(fg)
	} else if req.Method == 2 {
		if fg, err = this.GroupRightRecoverKList(fs, rows); err != nil {
			logging.Error("%v", err)
			return nil, err
		}
		this.CalcAfterRightRecoverKList(fg)
	} else {
		g1 := &FactorGroup{
			Ls: make([]*pro.KInfo, 0, 200),
		}
		g1.Ls = append(g1.Ls, rows...)
		fg = append(fg, g1)
	}
	if fg == nil {
		return nil, nil
	}

	return fg, nil
}

// 将K线根据除权数据的日期进行分组, 一组K线属于一条除权数据
func (this XRXD) FactorGroupOp(req *pro.RequestXRXD, fs []*pro.Factor, rows []*pro.KInfo) ([]*FactorGroup, error) {
	var fg []*FactorGroup // *Factor([]*KInfo), 即每一个复权因子关联一组K线
	var kg []*pro.KInfo   // 根据条件计算出来的需要参与除权除息计算的合法K线切片
	//var err error

	kg, err := this.GetRangeKList(req, rows)
	if err != nil {
		logging.Error("%v", err)
	}
	// 从数据库取出来的除权因子数组是 下标小->大 时间小->大
	// 从Redis取出来的K线数据数组是  下标小->大 时间大->小
	// 这里反转一下K线数组, 使其符合：下标小->大 时间小->大
	this.ReverseKList(kg)

	//  // debug show
	//	for n, v := range kg {
	//		fmt.Printf("for GetRange %02d %+v\n", n, v)
	//	}

	if req.Method == 1 {
		if fg, err = this.GroupRightRecoverKList(fs, kg); err != nil {
			logging.Error("%v", err)
			return nil, err
		}
		this.CalcBeforeRightRecoverKList(fg)
	} else if req.Method == 2 {
		if fg, err = this.GroupRightRecoverKList(fs, kg); err != nil {
			logging.Error("%v", err)
			return nil, err
		}
		this.CalcAfterRightRecoverKList(fg)
	} else {
		g1 := &FactorGroup{
			Ls: make([]*pro.KInfo, 0, 200),
		}
		g1.Ls = append(g1.Ls, kg...)
		fg = append(fg, g1)
	}
	if fg == nil {
		return nil, nil
	}

	return fg, nil
}

func (this XRXD) GetXRXDObj(req *pro.RequestXRXD) (*pro.PayloadXRXD, error) {
	key := fmt.Sprintf(this.CacheKey, req.SID)
	var fgs []*FactorGroup

	lsbin, err := RedisStore.LRange(key, 0, -1)
	if err != nil {
		logging.Error("%v", err)
		return nil, err
	}

	// 获取除权因子列表
	fcs, err := NewFactor().GetReferFactors(req.SID)
	if err != nil {
		return nil, err
	}

	// 解码
	ls, err := this.Decode(lsbin)
	if err != nil {
		return nil, err
	}

	switch req.Type {
	case 1: //日线
		fgs, err = this.FactorGroupOp(req, fcs, ls)
		if err != nil {
			return nil, err
		}
	case 2, 3, 4: // 其他K线
		fgs, err = this.FactorGroupTotal(req, fcs, ls)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Invalid request parameter - 'req.Type'")
	}

	if fgs == nil {
		return &pro.PayloadXRXD{
			SID:   req.SID,
			Type:  req.Type,
			Total: int32(len(ls)),
			Begin: req.TimeIndex,
			Num:   0,
			KList: nil,
		}, nil
	}

	result_ls := make([]*pro.KInfo, 0, 1024)
	var kline *[]*pro.KInfo

	for _, v := range fgs {
		result_ls = append(result_ls, v.Ls[:]...)
	}
	if len(result_ls) == 0 {
		return nil, fmt.Errorf("The klineData of xrxd is null")
	}

	switch req.Type {
	case 1:
		break
	case 2:
		if kline, err = ToWeekLine(&result_ls); err != nil {
			return nil, err
		}
	case 3:
		if kline, err = ToMonthLine(&result_ls); err != nil {
			return nil, err
		}
	case 4:
		if kline, err = ToYearLine(&result_ls); err != nil {
			return nil, err
		}
	}

	if req.Type == 1 {
		return &pro.PayloadXRXD{
			SID:   req.SID,
			Type:  req.Type,
			Total: int32(len(ls)),
			Begin: result_ls[0].NTime,
			Num:   int32(len(result_ls)),
			KList: result_ls,
		}, nil
	} else {
		res, err := this.getKlistByRequest(req, *kline)
		if err != nil {
			logging.Error("%v", err)
		}

		return &pro.PayloadXRXD{
			SID:   req.SID,
			Type:  req.Type,
			Total: int32(len(*kline)),
			Begin: res[0].NTime,
			Num:   int32(len(res)),
			KList: res,
		}, nil
	}
}

func (this XRXD) getKlistByRequest(req *pro.RequestXRXD, rows []*pro.KInfo) ([]*pro.KInfo, error) {
	var tmptime int32
	var n int
	var v *pro.KInfo
	for n, v = range rows {

		if tmptime < req.TimeIndex && req.TimeIndex < v.NTime {
			if tmptime == 0 {
				tmptime = v.NTime
			}
			break
		}
		tmptime = v.NTime
	}

	if req.Direct == 1 {
		if req.Num > 0 {
			m := n + int(req.Num)
			if m > len(rows) {
				m = len(rows)
			}
			return rows[n:m], nil
		}
		return rows[n:], nil

	} else {
		if req.Num > 0 {
			m := n - int(req.Num)
			if m < 0 {
				m = 0
			}
			return rows[m:n], nil
		}
		return rows[:n], nil
	}
}

//复权后的日K线转周K
func ToWeekLine(ksrc *[]*pro.KInfo) (*[]*pro.KInfo, error) {
	var aweek, weeks []*pro.KInfo

	sat, _ := utils.DateAdd((*ksrc)[0].NTime) //该股票第一个交易日所在周的周日（周六可能会有交易）
	for i, kl := range *ksrc {
		var wk *pro.KInfo
		var lengh int

		if utils.IntToTime(kl.NTime).Before(sat) {
			aweek = append(aweek, kl)
			if i == int(len(*ksrc)-1) { //执行到最后一个
				wk = daysToAKline(aweek) //一周形成
				if lengh = len(weeks); lengh > 0 {
					wk.NPreCPx = weeks[len(weeks)-1].NLastPx //昨收价取前一周最新价
				}
				weeks = append(weeks, wk)
			}
		} else {
			wk = daysToAKline(aweek) //一周形成
			if lengh = len(weeks); lengh > 0 {
				wk.NPreCPx = weeks[lengh-1].NLastPx //昨收价取前一周最新价
			}
			weeks = append(weeks, wk)

			sat, _ = utils.DateAdd(kl.NTime)
			aweek = nil
			aweek = append(aweek, kl)
		}
	}
	return &weeks, nil
}

//复权后的日K线转月K
func ToMonthLine(ksrc *[]*pro.KInfo) (*[]*pro.KInfo, error) {
	var amonth, months []*pro.KInfo
	var yesterday int32 = 0

	for i, kl := range *ksrc {
		var monk *pro.KInfo
		var lengh int

		if i == 0 {
			amonth = append(amonth, kl)
			yesterday = kl.NTime / 100
			continue
		}
		if yesterday == kl.NTime/100 {
			amonth = append(amonth, kl)
			if i == int(len(*ksrc)-1) { //执行到最后一个
				monk = daysToAKline(amonth) //一月形成
				if lengh = len(months); lengh > 0 {
					monk.NPreCPx = months[lengh-1].NLastPx //昨收价取前一月最新价
				}
				months = append(months, monk)
			}
		} else {
			monk = daysToAKline(amonth) //一月形成
			if lengh = len(months); lengh > 0 {
				monk.NPreCPx = months[lengh-1].NLastPx //昨收价取前一月最新价
			}
			months = append(months, monk)

			amonth = nil
			amonth = append(amonth, kl)
		}
		yesterday = kl.NTime / 100
	}
	return &months, nil
}

//复权后的日K线转年K
func ToYearLine(ksrc *[]*pro.KInfo) (*[]*pro.KInfo, error) {
	var ayear, years []*pro.KInfo
	var yesterday int32 = 0

	for i, kl := range *ksrc {
		var yk *pro.KInfo
		var lengh int

		if i == 0 {
			ayear = append(ayear, kl)
			yesterday = kl.NTime / 10000
			continue
		}
		if yesterday == kl.NTime/10000 {
			ayear = append(ayear, kl)
			if i == int(len(*ksrc)-1) { //执行到最后一个
				yk = daysToAKline(ayear) //一年形成
				if lengh = len(years); lengh > 0 {
					yk.NPreCPx = years[lengh-1].NLastPx //昨收价取前一年最新价
				}
				years = append(years, yk)
			}
		} else {
			yk = daysToAKline(ayear) //一年形成
			if lengh = len(years); lengh > 0 {
				yk.NPreCPx = years[lengh-1].NLastPx //昨收价取前一年最新价
			}
			years = append(years, yk)

			ayear = nil
			ayear = append(ayear, kl)
		}
		yesterday = kl.NTime / 10000
	}
	return &years, nil
}

//将一组K线合成一根
func daysToAKline(days []*pro.KInfo) *pro.KInfo {
	var (
		i          int
		AvgPxTotal uint32
		tmp        pro.KInfo
		dk         *pro.KInfo
	)

	for i, dk = range days {
		if tmp.NHighPx < dk.NHighPx || tmp.NHighPx == 0 { //最高价
			tmp.NHighPx = dk.NHighPx
		}
		if tmp.NLowPx > dk.NLowPx || tmp.NLowPx == 0 { //最低价
			tmp.NLowPx = dk.NLowPx
		}
		tmp.LlVolume += dk.LlVolume //成交量
		tmp.LlValue += dk.LlValue   //成交额
		AvgPxTotal += dk.NAvgPx
	}

	tmp.NSID = days[0].NSID
	tmp.NTime = days[0].NTime               //时间（取每周第一天）
	tmp.NOpenPx = days[0].NOpenPx           //开盘价（每周第一天的开盘价）
	tmp.NPreCPx = days[len(days)-1].NPreCPx //取每星期最后一天（之后再替换，防止第一周的昨收为零）
	tmp.NLastPx = days[i].NLastPx           //最新价
	tmp.NAvgPx = AvgPxTotal / uint32(i+1)   //平均价
	return &tmp
}
