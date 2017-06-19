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

func UpdateMonthLineToFile(filename string, today *protocol.KInfo) error {
	var tmp, result protocol.KInfo

	size := binary.Size(&tmp)
	bs := make([]byte, size)

	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

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

func UpdateYearLineToFile(filename string, today *protocol.KInfo) error {
	var tmp, result protocol.KInfo

	size := binary.Size(&tmp)
	bs := make([]byte, size)

	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Seek(int64(-size), 2)
	if err != nil {
		//	logging.Error("fist seek error ...%v", err.Error())
		return err
	}
	if _, err = file.Read(bs); err != nil && err != io.EOF {
		return err
	}

	buffer := bytes.NewBuffer(bs)

	if err = binary.Read(buffer, binary.LittleEndian, &tmp); err != nil && err != io.EOF {
		return err
	}

	if tmp.NTime/10000 == today.NTime/10000 { //同年
		_, err = file.Seek(int64(-size), 2)
		if err != nil {
			//	logging.Error("fist seek error ...%v", err.Error())
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

//更新周线
func UpdateWeekLineToFile(filename string, today *protocol.KInfo) error {
	var tmp, result protocol.KInfo

	size := binary.Size(&tmp)
	bs := make([]byte, size)

	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

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

//读PB形式
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

//binary形式
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

//binary 形式
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

//日K线数据追加相应文件的操作
func AppendFile(filepath string, today *protocol.KInfo) error {
	buffer := new(bytes.Buffer)

	if err := binary.Write(buffer, binary.LittleEndian, today); err != nil {
		return err
	}

	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.Write(buffer.Bytes()); err != nil {
		return err
	}
	return nil
}

//核实filestore路劲是否存在，不存在则创建
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
