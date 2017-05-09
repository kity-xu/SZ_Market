package controllers

import (
	"io/ioutil"
	"os"
	"strconv"

	"haina.com/market/hqinit/models/tb_security"
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
func FieldsToArrayString(src []tb_security.SecurityCode) []string {
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

//24
func StringToByte_SECURITY_CODE_LEN(str string) [SECURITY_CODE_LEN]byte {
	var ss [SECURITY_CODE_LEN]byte
	for i, v := range []byte(str) {
		if i == INDUSTRY_CODE_LEN {
			break
		}
		ss[i] = v
	}
	return ss
}

//40
func StringToByte_SECURITY_NAME_LEN(str string) [SECURITY_NAME_LEN]byte {
	var ss [SECURITY_NAME_LEN]byte
	for i, v := range []byte(str) {
		if i == INDUSTRY_CODE_LEN {
			break
		}
		ss[i] = v
	}
	return ss
}

//8
func StringToByte_SECURITY_DESC_LEN(str string) [SECURITY_DESC_LEN]byte {
	var ss [SECURITY_DESC_LEN]byte
	for i, v := range []byte(str) {
		if i == INDUSTRY_CODE_LEN {
			break
		}
		ss[i] = v
	}
	return ss
}

//8
func StringToByte_INDUSTRY_CODE_LEN(str string) [INDUSTRY_CODE_LEN]byte {
	var ss [INDUSTRY_CODE_LEN]byte
	for i, v := range []byte(str) {
		if i == INDUSTRY_CODE_LEN {
			break
		}
		ss[i] = v
	}
	return ss
}

//16
func StringToByte_SECURITY_ISIN_LEN(str string) [SECURITY_ISIN_LEN]byte {
	var ss [SECURITY_ISIN_LEN]byte
	for i, v := range []byte(str) {
		if i == INDUSTRY_CODE_LEN {
			break
		}
		ss[i] = v
	}
	return ss
}

//4
func StringToByte_4(str string) [4]byte {
	var ss [4]byte
	for i, v := range []byte(str) {
		if i == INDUSTRY_CODE_LEN {
			break
		}
		ss[i] = v
	}
	return ss
}
