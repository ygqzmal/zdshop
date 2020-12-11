package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["zdshop/controllers:GoodsController"] = append(beego.GlobalControllerRouter["zdshop/controllers:GoodsController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["zdshop/controllers:GoodsController"] = append(beego.GlobalControllerRouter["zdshop/controllers:GoodsController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/:cid/:pageIndex`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["zdshop/controllers:GoodsController"] = append(beego.GlobalControllerRouter["zdshop/controllers:GoodsController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:gid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["zdshop/controllers:GoodsController"] = append(beego.GlobalControllerRouter["zdshop/controllers:GoodsController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:gid`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["zdshop/controllers:GoodsController"] = append(beego.GlobalControllerRouter["zdshop/controllers:GoodsController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["zdshop/controllers:ObjectController"] = append(beego.GlobalControllerRouter["zdshop/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["zdshop/controllers:ObjectController"] = append(beego.GlobalControllerRouter["zdshop/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["zdshop/controllers:ObjectController"] = append(beego.GlobalControllerRouter["zdshop/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:objectId`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["zdshop/controllers:ObjectController"] = append(beego.GlobalControllerRouter["zdshop/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:objectId`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["zdshop/controllers:ObjectController"] = append(beego.GlobalControllerRouter["zdshop/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:objectId`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

}
