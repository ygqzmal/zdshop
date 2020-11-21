package controllers

import "github.com/astaxie/beego"

type ShopControllers struct {
	beego.Controller
}

func Get() {
	beego.Info("Hello World")
}
