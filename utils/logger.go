package utils

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
)

// RFC5424 log message levels.
const (
	LevelEmergency = iota
	LevelAlert        //1
	LevelCritical	  //2
	LevelError		  //3
	LevelWarning	  //4
	LevelNotice		  //5
	LevelInformational//6
	LevelDebug		  //7
)

func InitLogger() (err error) {
	BConfig, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil{
		fmt.Println("config init error:", err)
		return
	}
	maxlines, lerr := BConfig.Int64("log::maxlines")
	if lerr != nil {
		maxlines = 1000
	}

	logConf := make(map[string]interface{})
	logConf["filename"] = BConfig.String("log::log_path")
	level,_ := BConfig.Int("log::log_level")
	logConf["level"] = level
	logConf["maxlines"] = maxlines

	confStr, err := json.Marshal(logConf)
	if err != nil {
		fmt.Println("marshal failed,err:", err)
		return
	}
	//设置输出到文件
	//beego.SetLogger("file", `{"filename":"logs/test.log"}`)
	beego.SetLogger(logs.AdapterFile, string(confStr))
	//输是否出文件名和行号
	beego.SetLogFuncCall(true)
	//设置级别设置日志级别
	//beego.SetLevel(beego.LevelInformational)
	return
}
//调用以下方法可以写入log中
//beego.Info("index show")
