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

	"haina.com/market/hqpublish/models/fcmysql"
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

type Factor struct {
	BEGINDATE int     //起始日期VARCHAR(8) 	           本次除权因子的有效起始日期（即实际上的除权除息日）
	ENDDATE   int     //截止日期VARCHAR(8)              本次除权因子的有效截止日期（当尚无下一次除权因素的具体日期前，为19000101）
	XDY       float64 //当次除权因子NUMERIC(32,19)	   本次除权日，因分红送股转增等因素，依照除权前后价值不变动的原则计算的当次除权的折价因子
	LTDXDY    float64 //逆推累积除权因子NUMERIC(29,16)   以当前最新一天交易价格为标准，计算每次时间区间的累积除权因子，既每天实际交易价格与逆推复权价格之间的比值关系
	THELTDXDY float64 //顺推累计除权因子NUMERIC(29,16)   以上市第一天为标准，计算每次时间区间的累积除权因子，既顺推复权价格与每天实际交易价格与之间的比值关系
}

type FactorGroup struct {
	Fa Factor
	Ls []*pro.KInfo
}

func (this XRXD) ErrDataInvalid(fields string, sid int32, secode string) error {
	return errors.New(fmt.Sprintf("finchina TQ_SK_XDRY fields[%s] invalid by sid[%d], sid[%d]", "BEGINDATE", sid, secode))
}

// 从财汇数据库获取 *股票除权因子*
func (this XRXD) GetReferFactors(sid int32) ([]*Factor, error) {
	real_sid := sid % 1000000
	secode, err := fcmysql.NewTQ_OA_STCODE().GetSecode(fmt.Sprintf("%d", real_sid))
	if err != nil {
		logging.Error("%v", err)
		return nil, err
	}
	fc, err := fcmysql.NewTQ_SK_XDRY().GetFactorsBySecode(secode)
	if err != nil {
		logging.Error("%v", err)
		return nil, err
	}
	frows := make([]*Factor, 0, 100)
	for _, v := range fc {
		switch {
		case !v.BEGINDATE.Valid:
			return nil, this.ErrDataInvalid("BEGINDATE", sid, secode)
		case !v.ENDDATE.Valid:
			return nil, this.ErrDataInvalid("ENDDATE", sid, secode)
		case !v.XDY.Valid:
			return nil, this.ErrDataInvalid("XDY", sid, secode)
		case !v.LTDXDY.Valid:
			return nil, this.ErrDataInvalid("LTDXDY", sid, secode)
		case !v.THELTDXDY.Valid:
			return nil, this.ErrDataInvalid("THELTDXDY", sid, secode)
		}
		newf := Factor{
			BEGINDATE: int(v.BEGINDATE.Int64),
			ENDDATE:   int(v.ENDDATE.Int64),
			XDY:       v.XDY.Float64,
			LTDXDY:    v.LTDXDY.Float64,
			THELTDXDY: v.THELTDXDY.Float64,
		}
		frows = append(frows, &newf)
	}
	if len(frows) == 0 {
		return nil, errors.New(fmt.Sprintf("finchina TQ_SK_XDRY no datas by %d", sid))
	}
	if frows[len(frows)-1].ENDDATE == 19000101 {
		frows[len(frows)-1].ENDDATE = 99999999
	}
	return frows, nil
}

// After the right to recover
// 后复权计算前准备
// 获取后复权范围内K线列表
// 徐晓东写入的日K线使用的是redis list lpush的插入模式，时间大->小 下标小->大
func (this XRXD) GetBeforeRightRecoverKList(req *pro.RequestXRXD, rows []*pro.KInfo) ([]*pro.KInfo, error) {
	g := make([]*pro.KInfo, 0, 500)
	index := -1
	for n, v := range rows {
		if req.TimeIndex >= v.NTime {
			fmt.Println("find Before point", v)
			index = n
			break
		}
	}

	if index != -1 {
		n := index
		m := index + int(req.Num)
		if m > len(rows) {
			m = len(rows)
		}
		g = append(g, rows[n:m]...)
	}
	return g, nil
}

// Before the right to recover
// 前复权计算前准备
// 获取前复权范围内K线列表
func (this XRXD) GetAfterRightRecoverKList(req *pro.RequestXRXD, rows []*pro.KInfo) ([]*pro.KInfo, error) {
	g := make([]*pro.KInfo, 0, 500)
	index := -1
	for n := len(rows) - 1; n > -1; n-- {
		if req.TimeIndex <= rows[n].NTime {
			index = n
			break
		}
	}

	if index != -1 {
		for c, n := int32(0), index; n > -1 && c < req.Num; c, n = c+1, n-1 {
			g = append(g, rows[n])
		}
	}
	return g, nil
}

func (this XRXD) ReverseFactors(s []*Factor) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
func (this XRXD) ReverseKList(s []*pro.KInfo) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// 把前复权K线和除权因子进行关联分组
func (this XRXD) GroupBeforeRightRecoverKList(fs []*Factor, rows []*pro.KInfo) ([]*FactorGroup, error) {
	if len(fs) == 0 || len(rows) == 0 {
		return nil, nil
	}

	fmt.Printf("k %03d %+v\n", 0, rows[0])
	fmt.Printf("k %03d %+v\n", len(rows)-1, rows[len(rows)-1])

	// 从数据库取出来的除权因子数组是 下标小->大 时间小->大
	// 这里反转一下, 符合rows的方向  下标小->大 时间大->小
	this.ReverseFactors(fs)

	// 计算应使用的除权因子区间 // 下标左值
	bgn := 0
	for n, v := range fs {
		k := rows[0]
		if int32(v.BEGINDATE) <= k.NTime && int32(v.ENDDATE) >= k.NTime {
			bgn = n
			break
		}
	}
	// 计算应使用的除权因子区间 // 下标右值
	end := len(fs) - 1
	for n := len(fs) - 1; n > -1; n-- {
		k := rows[len(rows)-1]
		if int32(fs[n].BEGINDATE) <= k.NTime && int32(fs[n].ENDDATE) >= k.NTime {
			end = n
			break
		}
	}
	fmt.Println("----------------原有除权因子")
	for n, v := range fs {
		fmt.Printf("f %03d %+v\n", n, *v)
	}
	fmt.Printf("bgn end %d %d\n", bgn, end)
	fmt.Println("----------------范围内除权因子")
	nowfs := fs[bgn : end+1]
	for n, v := range nowfs {
		fmt.Printf("f %03d %+v\n", n, *v)
	}
	fmt.Println("----------------")

	// //debug show k line list
	//	for n, v := range rows {
	//		fmt.Println("k", n, v)
	//	}

	fmt.Println("-------创建除权分组")
	// 创建分组
	var fgs []*FactorGroup
	for n := bgn; n <= end; n++ {
		fg := &FactorGroup{
			Fa: *fs[n],
			Ls: make([]*pro.KInfo, 0, 200),
		}
		fgs = append(fgs, fg)
		fmt.Println("f", n, fg)
	}

	fmt.Println("-------除权分组 <- K线分组")
	s1 := 0
	for _, v := range fgs {
		m := s1
		for ; m < len(rows); m++ {
			k := rows[m]
			if k.NTime < int32(v.Fa.BEGINDATE) || k.NTime > int32(v.Fa.ENDDATE) {
				fmt.Println("k find", m, k)
				fmt.Println("f", *v, "append rows[", s1, ":", m, "]")
				break
			}
		}
		v.Ls = append(v.Ls, rows[s1:m]...)
		s1 = m
		if m == len(rows) {
			break
		}
	}

	//debug show factor group data
	for n, v := range fgs {
		fmt.Println("f", n, v.Fa, "k len", len(v.Ls))
		if len(v.Ls) > 0 {
			for n, v := range v.Ls {
				fmt.Println("  k", n, *v)
			}
		}
		fmt.Println("---------------")
	}

	return fgs, nil
}

func (this XRXD) CalcBeforeRightRecoverKList(fgs []*FactorGroup) {
	for _, v := range fgs {
		factor := v.Fa
		for _, k := range v.Ls {
			a := k.NPreCPx
			k.NPreCPx = int32(float64(k.NPreCPx) / 10000 / factor.LTDXDY * 10000)
			fmt.Printf("%d .... %d, factor %f\n", a, k.NPreCPx, factor.LTDXDY)
		}
	}
}
func (this XRXD) CalcAfterRightRecoverKList(fgs []*FactorGroup) {
}

// 把后复权K线和除权因子进行关联分组
func (this XRXD) GroupAfterRightRecoverKList(fs []*Factor, rows []*pro.KInfo) ([]*FactorGroup, error) {
	if len(fs) == 0 || len(rows) == 0 {
		return nil, nil
	}
	fmt.Printf("k %03d %+v\n", 0, rows[0])
	fmt.Printf("k %03d %+v\n", len(rows)-1, rows[len(rows)-1])
	for n, v := range fs {
		fmt.Printf("%d %+v\n", n, v)
	}
	return nil, nil
}

// 将K线根据除权数据的日期进行分组, 一组K线属于一条除权数据
func (this XRXD) FactorGroupOp(req *pro.RequestXRXD, fs []*Factor, rows []*pro.KInfo) ([]*FactorGroup, error) {
	var fg []*FactorGroup // *Factor([]*KInfo), 即每一个除权数据包含一组K线
	var kg []*pro.KInfo   // 根据条件计算出来的需要参与除权除息计算的合法K线切片
	var err error
	if req.Direct == 0 {
		if kg, err = this.GetBeforeRightRecoverKList(req, rows); err != nil {
			logging.Error("%v", err)
			return nil, err
		}
		if fg, err = this.GroupBeforeRightRecoverKList(fs, kg); err != nil {
			logging.Error("%v", err)
			return nil, err
		}

		if len(fg) > 0 {
			this.CalcBeforeRightRecoverKList(fg)
		}
	} else {
		if kg, err = this.GetAfterRightRecoverKList(req, rows); err != nil {
			logging.Error("%v", err)
			return nil, err
		}
		if fg, err = this.GroupAfterRightRecoverKList(fs, kg); err != nil {
			logging.Error("%v", err)
			return nil, err
		}
	}
	if fg == nil {
		return nil, nil
	}

	//	//debug show
	//	for _, v := range kg {
	//		fmt.Println(v)
	//	}

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
	fcs, err := this.GetReferFactors(req.SID)
	if err != nil {
		return nil, err
	}

	//	for _, v := range fcs {
	//		fmt.Println(*v)
	//	}

	// 解码
	ls, err := this.Decode(lsbin)
	if err != nil {
		return nil, err
	}

	//	for n, v := range ls {
	//		fmt.Println(n, v)
	//	}

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

	//return nil, errors.New(fmt.Sprintf("For debug error suspend"))

	return &pro.PayloadXRXD{
		SID:   req.SID,
		Total: int32(len(ls)),
		Begin: req.TimeIndex,
		Num:   int32(len(ls)),
		KList: ls,
	}, nil
}
