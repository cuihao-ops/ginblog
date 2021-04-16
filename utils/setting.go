package utils

import (
	"fmt"

	"gopkg.in/ini.v1"
)

var (
	AppName  string
	AppMode  string
	HttpPort string
	JwtKey   string

	Db         string
	DbHost     string
	DbPort     string
	DbPassWord string
	DbName     string
	DbUser     string

	AccessKey  string
	SecretKey  string
	Bucket     string
	QiniuSever string
)

func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误！", err)
	}
	LoadData(file)
	LoanServer(file)
	LoadQiniu(file)
}

func LoanServer(file *ini.File) {
	AppName = file.Section("server").Key("AppName").MustString("ginblog")
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":3000")
	JwtKey = file.Section("server").Key("JwtKey").MustString("sss999wefaf")
}

func LoadData(file *ini.File) {
	Db = file.Section("database").Key("Db").MustString("mysql")
	DbHost = file.Section("database").Key("DbHost").MustString("localhost")
	DbPort = file.Section("database").Key("DbPort").MustString("3306")
	DbUser = file.Section("database").Key("DbUser").MustString("root")
	DbPassWord = file.Section("database").Key("DbPassWord").MustString("123456")
	DbName = file.Section("database").Key("DbName").MustString("ginblog")
}

func LoadQiniu(file *ini.File) {
	AccessKey = file.Section("qiniu").Key("AccessKey").MustString("")
	SecretKey = file.Section("qiniu").Key("SecretKey").MustString("")
	Bucket = file.Section("qiniu").Key("Bucket").MustString("")
	QiniuSever = file.Section("qiniu").Key("QiniuSever").MustString("")
}
