package main

import (
	"github.com/astaxie/beego"
	_ "zdshop/routers"
	"zdshop/utils"
	//初始化orm, 操作数据库
	_ "zdshop/models"
	
)

func init() {
	//日志写入
	_ = utils.InitLogger()
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	//bee run
	beego.Run()
}


