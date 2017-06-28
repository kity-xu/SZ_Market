package publish

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	//	"strconv"

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
			//fmt.Printf("0 1 n %d - %+v\n", n, *rows[n])
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

	fgs, err := this.FactorGroupOp(req, fcs, ls)
	if err != nil {
		return nil, err
	}
	if fgs == nil {
		return &pro.PayloadXRXD{
			SID:   req.SID,
			Total: int32(len(ls)),
			Begin: req.TimeIndex,
			Num:   0,
			KList: nil,
		}, nil
	}

	result_ls := make([]*pro.KInfo, 0, 1024)
	for _, v := range fgs {
		result_ls = append(result_ls, v.Ls[:]...)
	}

	begin := int32(0)
	if len(result_ls) > 0 {
		begin = result_ls[0].NTime
	}

	return &pro.PayloadXRXD{
		SID:   req.SID,
		Total: int32(len(ls)),
		Begin: begin,
		Num:   int32(len(result_ls)),
		KList: result_ls,
	}, nil
}
