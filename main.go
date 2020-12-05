package main

import (
	"github.com/astaxie/beego"
	_ "zdshop/routers"
	//初始化orm, 操作数据库
	_ "zdshop/models"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
