package config

import (
	"os"
	"os/exec"
	"path/filepath"

	"haina.com/share/lib"
)

var config *AppConfig

type AppSettings struct {
	AllowOrigin string     `xml:"allowOrigin"`
	EncryFactor string     `xml:"encryFactor"`
	Environment string     `xml:"environment"`
	Listen      string     `xml:"listen"`
	Projects    []Projects `xml:"projects"`
}

type Projects struct {
	AppId      string `xml:"appId"`
	ConfigFile string `xml:"configFile"`
}

// log
type LogServer struct {
	On   bool   `xml:"on"`
	Addr string `xml:"addr"`
	Port string `xml:"port"`
}

// Redis
type RedisStore struct {
	Addr    string `xml:"addr"`
	Auth    string `xml:"auth"`
	Db      string `xml:"db"`
	Timeout int    `xml:"timeout"`
}

// MongoDB
type MongoStore struct {
	Source string `xml:"source"`
}

// fileSystem
type FileStore struct {
	Finpath string `xml:"finpath"`
	Finday  string `xml:"finday"`
	Findex  string `xml:"findex"`

	Path  string `xml:"path"`
	Day   string `xml:"day"`
	Index string `xml:"index"`
	Week  string `xml:"week"`
	Month string `xml:"month"`
	Year  string `xml:"year"`
	Min   string `xml:"hmin"`
	Min5  string `xml:"hmin5"`
	Min15 string `xml:"hmin15"`
	Min30 string `xml:"hmin30"`
	Min60 string `xml:"hmin60"`
}

type CorsSetting struct {
	AllowOrigin []string `xml:"allowOrigin"`
}

type ListenAndServe struct {
	Port    string `xml:"port"`
	LogPort string `xml:"logport"`
}

type AppConfig struct {
	Cors     CorsSetting    `xml:"cors"`
	Mongo    MongoStore     `xml:"mongoStore"`
	File     FileStore      `xml:"fileStore"`
	Redis    RedisStore     `xml:"redisStore"`
	Serve    ListenAndServe `xml:"listenAndServe"`
	Settings AppSettings    `xml:"appSettings"`
	Log      LogServer      `xml:"logServer"`
}

func Default(appID string) *AppConfig {
	if config == nil {
		var cfg AppConfig
		lib.LoadConfig(appID, &cfg)
		config = &cfg
	}
	return config
}

func Reload() {
	config = nil
}

// ------------------------------------------------------------------------

func getCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	return filepath.Dir(file)
}
