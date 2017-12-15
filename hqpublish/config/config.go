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

type AccessKeys struct {
	ID     string `xml:"id"`
	Secret string `xml:"secret"`
	AESKey string `xml:"aesKey"`
}

// Database
type Database struct {
	DriverName string `xml:"driverName"`
	DataSource string `xml:"dataSource"`
}

// MongoDB
type MongoStore struct {
	Source string `xml:"source"`
}

// Redis
type RedisStore struct {
	Addr    string `xml:"addr"`
	Auth    string `xml:"auth"`
	Db      string `xml:"db"`
	Timeout int    `xml:"timeout"`
}

// Email
type EmailSetting struct {
	Addr     string `xml:"addr"`
	Password string `xml:"password"`
	Server   string `xml:"server"`
	Port     string `xml:"port"`
}

type SessionSetting struct {
	On           bool   `xml:"on"`
	ProviderName string `xml:"providerName"`
	Config       string `xml:"config"`
}
type MnsSetting struct {
	Url             string    `xml:"url"`
	AccessKeyId     string    `xml:"accessKeyId"`
	AccessKeySecret string    `xml:"accessKeySecret"`
	Queues          QueueName `xml:"queues"`
}
type QueueName struct {
	SmartCall string `xml:"exchangeSmartCall"`
}

type LogServer struct {
	On   bool   `xml:"on"`
	Addr string `xml:"addr"`
	Port string `xml:"port"`
}

type AppConfig struct {
	File           FileStore      `xml:"fileStore"`
	TTL            CacheTTL       `xml:"cacheTTL"`
	AccessKeys     AccessKeys     `xml:"accessKeys"`
	Cors           CorsSetting    `xml:"cors"`
	Db             Database       `xml:"database"`
	DbMicroLink    Database       `xml:"dbMicroLink"` //米领后台
	DbSZ           Database       `xml:"dbSZ"`
	Email          EmailSetting   `xml:"emailSetting"`
	Mns            MnsSetting     `xml:"mns"`
	Mongo          MongoStore     `xml:"mongoStore"`
	Redis          RedisStore     `xml:"redisStore"`
	RedisCache     RedisStore     `xml:"redisCache"`
	RedisMicroLink RedisStore     `xml:"redisMicroLink"` //米领后台
	Serve          ListenAndServe `xml:"listenAndServe"`
	Session        SessionSetting `xml:"session"`
	Settings       AppSettings    `xml:"appSettings"`
	Log            LogServer      `xml:"logServer"`
	Catalog        FileCatalog    `xml:"filecatalog"`
}
type FileStore struct {
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

	Bindex string `xml:"blockindex"`
}

// 公告附件目录 url地址
type FileCatalog struct {
	Url       string `xml:"url"`
	ValidTime string `xml:"validtime"`
}

type CacheTTL struct {
	Day   int `xml:"day"`
	Week  int `xml:"week"`
	Month int `xml:"month"`
	Year  int `xml:"year"`

	Min   int `xml:"hmin"`
	Min1  int `xml:"hmin1"`
	Min5  int `xml:"hmin5"`
	Min15 int `xml:"hmin15"`
	Min30 int `xml:"hmin30"`
	Min60 int `xml:"hmin60"`
	Sort  int `xml:"sort"`
	Block int `xml:"block"`

	MinK         int `xml:"minK"`
	MarketStatus int `xml:"marketStatus"`
	MIndex       int `xml:"mindex"`
	F10HomePage  int `xml:"f10HomePage"`

	FinanceChart            int `xml:"financeChart"`
	FinanceReport           int `xml:"financeReport"`
	FinanceReportStatistics int `xml:"financeReportStatistics"`
	FinanceReportForecast   int `xml:"financeReportForecast"`
}

type CorsSetting struct {
	AllowOrigin []string `xml:"allowOrigin"`
}

type ListenAndServe struct {
	Port    string `xml:"port"`
	LogPort string `xml:"logport"`
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
