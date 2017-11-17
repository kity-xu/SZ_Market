package filestore

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"

	"ProtocolBuffer/projects/hqpost/go/protocol"
	"errors"

	"haina.com/share/lib"

	"haina.com/share/logging"
)

// UpdateMonthLineToFile ...
func UpdateMonthLineToFile(sid int32, filename string, today *protocol.KInfo) error {
	if today == nil {
		return nil
	}
	var tmp, result protocol.KInfo

	size := binary.Size(&tmp)
	bs := make([]byte, size)

	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	n, err := file.Seek(0, 2) //添加新股
	if n == 0 && err == nil {
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, today)
		if _, err = file.Write(buf.Bytes()); err != nil {
			return err
		}
		return nil
	}

	_, err = file.Seek(int64(-size), 2)
	if err != nil {
		return err
	}
	if _, err = file.Read(bs); err != nil && err != io.EOF {
		return err
	}

	buffer := bytes.NewBuffer(bs)

	if err = binary.Read(buffer, binary.LittleEndian, &tmp); err != nil && err != io.EOF {
		return err
	}

	// 判断当天是否追加了多次
	if today.NTime <= tmp.NTime {
		return nil
	}

	if tmp.NTime/100 == today.NTime/100 { //同月
		_, err = file.Seek(int64(-size), 2)
		if err != nil {
			return err
		}
		result = compareKInfo(&tmp, today)
	} else {
		result = *today
		result.NPreCPx = tmp.NLastPx //昨收价
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, &result)
	if _, err = file.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

// UpdateYearLineToFile ...
func UpdateYearLineToFile(sid int32, filename string, today *protocol.KInfo) error {
	if today == nil {
		return nil
	}
	var tmp, result protocol.KInfo

	size := binary.Size(&tmp)
	bs := make([]byte, size)

	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	n, err := file.Seek(0, 2)
	if n == 0 && err == nil { //添加新股
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, today)
		if _, err = file.Write(buf.Bytes()); err != nil {
			return err
		}
		return nil
	}

	_, err = file.Seek(int64(-size), 2)
	if err != nil {
		return err
	}
	if _, err = file.Read(bs); err != nil && err != io.EOF {
		return err
	}

	buffer := bytes.NewBuffer(bs)

	if err = binary.Read(buffer, binary.LittleEndian, &tmp); err != nil && err != io.EOF {
		return err
	}
	// 判断当天是否追加了多次
	if today.NTime <= tmp.NTime {
		return nil
	}

	if tmp.NTime/10000 == today.NTime/10000 { //同年
		_, err = file.Seek(int64(-size), 2)
		if err != nil {
			return err
		}
		result = compareKInfo(&tmp, today)
	} else {
		result = *today
		result.NPreCPx = tmp.NLastPx //昨收价
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, &result)
	if _, err = file.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

// UpdateWeekLineToFile  更新周线
func UpdateWeekLineToFile(sid int32, filename string, today *protocol.KInfo) error {
	if today == nil {
		return nil
	}
	var tmp, result protocol.KInfo

	size := binary.Size(&tmp)
	bs := make([]byte, size)

	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	n, err := file.Seek(0, 2) //添加新股
	if n == 0 && err == nil {
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, today)
		if _, err = file.Write(buf.Bytes()); err != nil {
			return err
		}
		return nil
	}

	_, err = file.Seek(int64(-size), 2)
	if err != nil {
		return err
	}
	if _, err = file.Read(bs); err != nil && err != io.EOF {
		return err
	}

	buffer := bytes.NewBuffer(bs)

	if err = binary.Read(buffer, binary.LittleEndian, &tmp); err != nil && err != io.EOF {
		return err
	}

	// 判断当天是否追加了多次
	if today.NTime <= tmp.NTime {
		return nil
	}
	b1, _ := DateAdd(tmp.NTime) //找到该日期所在周日的那天
	b2, _ := DateAdd(today.NTime)

	if b1.Equal(b2) { //同属一周
		_, err = file.Seek(int64(-size), 2)
		if err != nil {
			return err
		}
		result = compareKInfo(&tmp, today)
	} else { //不属于同周
		result = *today
		result.NPreCPx = tmp.NLastPx //昨收价
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, &result)
	if _, err = file.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

// ReadSrcFileStore ... 读第三方数据源生成的源文件（PB形式）
func ReadSrcFileStore(filename string) (*protocol.KInfoTable, error) {
	var klist protocol.KInfoTable

	fd, err := ioutil.ReadFile(filename)
	if err != nil {
		logging.Error("%v", err.Error())
		return nil, err
	}
	//解PB
	if err = proto.Unmarshal(fd, &klist); err != nil {
		logging.Error("%v", err.Error())
		return nil, err
	}
	return &klist, nil
}

// ReadHainaFileStore  binary形式
func ReadHainaFileStore(filename string) (*protocol.KInfoTable, error) {
	if !lib.IsFileExist(filename) {
		return nil, errors.New("file not exist..")
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logging.Error("%v", err.Error())
		return nil, err
	}

	size := binary.Size(&protocol.KInfo{})
	ktable := protocol.KInfoTable{}

	for i := 0; i < len(data); i += size {
		v := data[i : i+size]
		kinfo := protocol.KInfo{}
		buffer := bytes.NewBuffer(v)
		if err = binary.Read(buffer, binary.LittleEndian, &kinfo); err != nil && err != io.EOF {
			return nil, err
		}
		ktable.List = append(ktable.List, &kinfo)
	}
	return &ktable, nil
}

// WiteHainaFileStore ... binary 形式写入海纳文件系统(hgs_file)
func WiteHainaFileStore(filepath string, ktable *protocol.KInfoTable) error {
	if len(ktable.List) < 1 {
		return errors.New("The history is nill..")
	}
	buffer := new(bytes.Buffer)

	for _, v := range ktable.List {
		if err := binary.Write(buffer, binary.LittleEndian, v); err != nil {
			return err
		}
	}
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.Write(buffer.Bytes()); err != nil {
		return err
	}
	return nil
}

// AppendFile ... 文件追加数据
func AppendFile(filepath string, today *protocol.KInfo) error {
	if today == nil {
		return nil
	}
	var tmp protocol.KInfo
	size := binary.Size(&tmp)
	bs := make([]byte, size)

	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	n, err := file.Seek(0, 2)
	if n == 0 && err == nil { // 添加新股 文件大小为0字节
		logging.Info("添加新股 %d", today.NSID)
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, today)
		if _, err = file.Write(buf.Bytes()); err != nil {
			return err
		}
		return nil
	}

	_, err = file.Seek(int64(-size), 2)
	if err != nil {
		return err
	}
	if _, err = file.Read(bs); err != nil && err != io.EOF {
		return err
	}

	buffer := bytes.NewBuffer(bs)

	if err = binary.Read(buffer, binary.LittleEndian, &tmp); err != nil && err != io.EOF {
		return err
	}

	// 如果最后一根K线的日期与today的相等(表示一天内执行了多次); 如果是大于， 表示today回滚了。 此时都不执行追加操作
	if tmp.NTime >= today.NTime {
		return nil
	}

	tobuf := new(bytes.Buffer)
	if err := binary.Write(tobuf, binary.LittleEndian, today); err != nil {
		return err
	}

	// file.Read 后， 文件游标在文件末尾
	if _, err = file.Write(tobuf.Bytes()); err != nil {
		return err
	}
	return nil
}

// CheckFileSoteDir ... 核实filestore路劲是否存在，不存在则创建
func CheckFileSoteDir(sid int32, hnpath, name string) (string, bool) {
	var path string

	market := sid / 1000000
	if market == 100 {
		path = fmt.Sprintf("%s/sh/%d", hnpath, sid)
	} else if market == 200 {
		path = fmt.Sprintf("%s/sz/%d", hnpath, sid)
	} else {
		logging.Error("Monthline write file error...Invalid file path")
	}

	filename := path + "/" + name

	if !lib.IsDirExists(path) {
		err := os.MkdirAll(path, 0777)
		if err != nil {
			logging.Error("%v", err.Error())
		}
		return filename, false
	}

	if !lib.IsFileExist(path + "/" + name) {
		return filename, false
	}
	return filename, true
}

//MaybeBelongAWeek ... 判断今天与历史最新是否同属一周（wk 线做第一次生成时使用）
func MaybeBelongAWeek(klist *protocol.KInfoTable, today *protocol.KInfo) {
	if today == nil {
		return
	}
	if len(klist.List) < 1 {
		logging.Info("%v no historical week data...", today.NSID)
		klist.List = append(klist.List, today)
		return
	}
	var kinfo = *(klist.List)[len(klist.List)-1]

	b1, _ := DateAdd(kinfo.NTime) //找到该日期所在周日的那天
	b2, _ := DateAdd(today.NTime)

	if b1.Equal(b2) { //同属一周
		result := compareKInfo(&kinfo, today)
		klist.List[len(klist.List)-1] = &result
	} else { //不属于同周
		klist.List = append(klist.List, today)
	}
	return
}

// MaybeBelongAMonth ...
func MaybeBelongAMonth(klist *protocol.KInfoTable, today *protocol.KInfo) {
	if today == nil {
		return
	}
	if len(klist.List) < 1 {
		logging.Info("%v no historical month data...", today.NSID)
		klist.List = append(klist.List, today)
		return
	}
	var kinfo = *(klist.List)[len(klist.List)-1]

	if kinfo.NTime/100 == today.NTime/100 { //同月
		result := compareKInfo(&kinfo, today)
		klist.List[len(klist.List)-1] = &result
	} else {
		klist.List = append(klist.List, today)
	}
	return
}

// MaybeBelongAYear ...
func MaybeBelongAYear(klist *protocol.KInfoTable, today *protocol.KInfo) {
	if today == nil {
		return
	}
	if len(klist.List) < 1 {
		logging.Info("%v no historical year data...", today.NSID)
		klist.List = append(klist.List, today)
		return
	}
	var kinfo = *(klist.List)[len(klist.List)-1]

	if kinfo.NTime/10000 == today.NTime/10000 { //同年
		result := compareKInfo(&kinfo, today)
		klist.List[len(klist.List)-1] = &result
	} else {
		klist.List = append(klist.List, today)
	}
	return
}
