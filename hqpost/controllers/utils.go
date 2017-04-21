package controllers

import (
	"io/ioutil"
	"os"
	"strconv"

	"haina.com/market/hqpost/models/tb_stokcode"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

//获取某一文件夹下的所有文件名或目录名
func GetFileNameList(dir string) ([]string, error) {
	var files []string
	if !lib.IsDirExists(dir) {
		logging.Error("Invalid directory files...%s", dir)
	}

	dir_list, e := ioutil.ReadDir(dir)
	if e != nil {
		logging.Error("Traverse %s error...", dir)
	}

	for _, v := range dir_list {
		files = append(files, v.Name())
	}
	return files, e
}

//核实代码表中的代码在文件中是否有相应的存在
func CheckFilesName(sh, sz []string, codes []string) bool {
	for _, code := range codes {
		if !lib.InArray(sh, code) && !lib.InArray(sz, code) {
			logging.Error("Under the SH or SZ directory has no folder:%s", code)
			return false
		}
	}
	return true
}

//结构属性转string数组
func FieldsToArrayString(src []tb_stokcode.Code) []string {
	var codes []string
	for _, v := range src {
		codes = append(codes, strconv.Itoa(int(v.SID)))
	}
	return codes
}

//打开文件
func OpenFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		logging.Error("OpenFile error...%v", err.Error())
	}
	return file, err
}

//关闭文件
func CloseFile(file *os.File) {
	file.Close()
}

//读文件固定大小字节
func ReadFiles(file *os.File, des []byte) (int, error) {
	return file.Read(des)
}
