package kline

import (
	"fmt"
	"os"

	"strings"

	"haina.com/market/hqpost/config"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

const (
	// SH ... 上海交易市场
	SH = "/sh/"
	// SZ ... 深圳交易市场
	SZ = "/sz/"
)

// GetExchange ... 根据SID 获取交易市场
func GetExchange(sid int32) (string, error) {
	var exchange string
	if sid/100000000 == 1 {
		exchange = SH
	} else if sid/100000000 == 2 {
		exchange = SZ
	} else {
		logging.Error("Invalid sid...")
		return "", fmt.Errorf("%s", "Invalid sid...")
	}
	return exchange, nil
}

// IsExistdirInHGSFileStore ... 判断HGS_FILE是否存在该SID的源数据
func IsExistFileInHGSFileStore(cfg *config.AppConfig, kind string, sid int32) (string, bool) {
	exchange, err := GetExchange(sid)
	if err != nil {
		return "", false
	}
	dpath := fmt.Sprintf("%s%s%s/%d.dat", cfg.File.Path, exchange, kind, sid)
	if !lib.IsFileExist(dpath) {
		return dpath, false
	}
	return dpath, true
}

// IsExistdirInThirdDB ... 判断ThirdDB是否存在该SID的源数据
func IsExistdirInThirdDB(cfg *config.AppConfig, sid int32) (string, bool) {
	exchange, err := GetExchange(sid)
	if err != nil {
		return "", false
	}
	spath := fmt.Sprintf("%s%s%d", cfg.File.Finpath, exchange, sid)
	if !lib.IsDirExists(spath) {
		logging.Debug("%s", fmt.Sprintf("non-existent data sources | %d, finchina DB", sid))
		return spath, false
	}
	return spath, true
}

// CreateDir ... 创建目录
func CreateDir(path string) error {
	return os.Mkdir(path, 0755)
}

// HGSFilepath ... 得到hgs_file文件路劲
func HGSFilepath(file string) string {
	if !lib.IsFileExist(file) {
		ss := strings.Split(file, "/")
		var dir string
		for i, v := range ss {
			if i == len(ss)-1 {
				break
			}
			dir += v
			dir += "/"
		}
		if err := os.MkdirAll(dir, 0666); err != nil {
			logging.Error("Mkdir :%v", err)
		}
	}
	return file
}
