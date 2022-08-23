package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["zdshop/controllers:CartController"] = append(beego.GlobalControllerRouter["zdshop/controllers:CartController"],
        beego.ControllerComments{
            Method: "PostCart",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:CartController"] = append(beego.GlobalControllerRouter["zdshop/controllers:CartController"],
        beego.ControllerComments{
            Method: "GetCart",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:CartController"] = append(beego.GlobalControllerRouter["zdshop/controllers:CartController"],
        beego.ControllerComments{
            Method: "GetParameters",
            Router: "/:gid",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:CartController"] = append(beego.GlobalControllerRouter["zdshop/controllers:CartController"],
        beego.ControllerComments{
            Method: "PutCartPara",
            Router: "/:oldId/:newId/:count",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:CartController"] = append(beego.GlobalControllerRouter["zdshop/controllers:CartController"],
        beego.ControllerComments{
            Method: "DeleteCart",
            Router: "/:pid",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:CartController"] = append(beego.GlobalControllerRouter["zdshop/controllers:CartController"],
        beego.ControllerComments{
            Method: "PutCartCount",
            Router: "/:pid/:count",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:CategoryController"] = append(beego.GlobalControllerRouter["zdshop/controllers:CategoryController"],
        beego.ControllerComments{
            Method: "PostCategory",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:CategoryController"] = append(beego.GlobalControllerRouter["zdshop/controllers:CategoryController"],
        beego.ControllerComments{
            Method: "GetFirstCategory",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:CategoryController"] = append(beego.GlobalControllerRouter["zdshop/controllers:CategoryController"],
        beego.ControllerComments{
            Method: "UpdateCategory",
            Router: "/",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:CategoryController"] = append(beego.GlobalControllerRouter["zdshop/controllers:CategoryController"],
        beego.ControllerComments{
            Method: "DeleteCategory",
            Router: "/:cid",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:CategoryController"] = append(beego.GlobalControllerRouter["zdshop/controllers:CategoryController"],
        beego.ControllerComments{
            Method: "GetSecondCategory",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:GoodsController"] = append(beego.GlobalControllerRouter["zdshop/controllers:GoodsController"],
        beego.ControllerComments{
            Method: "PostGoods",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:GoodsController"] = append(beego.GlobalControllerRouter["zdshop/controllers:GoodsController"],
        beego.ControllerComments{
            Method: "GetAllGoods",
            Router: "/:cid/:pageIndex",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:GoodsController"] = append(beego.GlobalControllerRouter["zdshop/controllers:GoodsController"],
        beego.ControllerComments{
            Method: "DeleteGoods",
            Router: "/:gid",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:GoodsController"] = append(beego.GlobalControllerRouter["zdshop/controllers:GoodsController"],
        beego.ControllerComments{
            Method: "GetOneGoods",
            Router: "/:gid",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:GoodsController"] = append(beego.GlobalControllerRouter["zdshop/controllers:GoodsController"],
        beego.ControllerComments{
            Method: "PostGoodsBannerImg",
            Router: "/img",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:GoodsController"] = append(beego.GlobalControllerRouter["zdshop/controllers:GoodsController"],
        beego.ControllerComments{
            Method: "PutGoodsBannerImg",
            Router: "/img",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:GoodsController"] = append(beego.GlobalControllerRouter["zdshop/controllers:GoodsController"],
        beego.ControllerComments{
            Method: "DelGoodsBannerImg",
            Router: "/img/:bid",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:GoodsController"] = append(beego.GlobalControllerRouter["zdshop/controllers:GoodsController"],
        beego.ControllerComments{
            Method: "DeleteParameter",
            Router: "/para/:pid",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:GoodsController"] = append(beego.GlobalControllerRouter["zdshop/controllers:GoodsController"],
        beego.ControllerComments{
            Method: "PutGoodsAttribute",
            Router: "/putAttribute",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:GoodsController"] = append(beego.GlobalControllerRouter["zdshop/controllers:GoodsController"],
        beego.ControllerComments{
            Method: "PutGoodsBanners",
            Router: "/putBanner",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:GoodsController"] = append(beego.GlobalControllerRouter["zdshop/controllers:GoodsController"],
        beego.ControllerComments{
            Method: "PutGoodsParameters",
            Router: "/putParameter",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:GoodsController"] = append(beego.GlobalControllerRouter["zdshop/controllers:GoodsController"],
        beego.ControllerComments{
            Method: "PutGoodsState",
            Router: "/putState",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:GoodsController"] = append(beego.GlobalControllerRouter["zdshop/controllers:GoodsController"],
        beego.ControllerComments{
            Method: "HandleSearch",
            Router: "/search",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:GoodsController"] = append(beego.GlobalControllerRouter["zdshop/controllers:GoodsController"],
        beego.ControllerComments{
            Method: "Tex",
            Router: "/tex",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:ObjectController"] = append(beego.GlobalControllerRouter["zdshop/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:ObjectController"] = append(beego.GlobalControllerRouter["zdshop/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:ObjectController"] = append(beego.GlobalControllerRouter["zdshop/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Get",
            Router: "/:objectId",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:ObjectController"] = append(beego.GlobalControllerRouter["zdshop/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/:objectId",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:ObjectController"] = append(beego.GlobalControllerRouter["zdshop/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/:objectId",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:OrderController"] = append(beego.GlobalControllerRouter["zdshop/controllers:OrderController"],
        beego.ControllerComments{
            Method: "PostOrder",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:OrderController"] = append(beego.GlobalControllerRouter["zdshop/controllers:OrderController"],
        beego.ControllerComments{
            Method: "GetAllOrder",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:OrderController"] = append(beego.GlobalControllerRouter["zdshop/controllers:OrderController"],
        beego.ControllerComments{
            Method: "PutOrderAddr",
            Router: "/",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:OrderController"] = append(beego.GlobalControllerRouter["zdshop/controllers:OrderController"],
        beego.ControllerComments{
            Method: "DeleteOrder",
            Router: "/:oid",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:OrderController"] = append(beego.GlobalControllerRouter["zdshop/controllers:OrderController"],
        beego.ControllerComments{
            Method: "GetOneOrderInfo",
            Router: "/:oid",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:OrderController"] = append(beego.GlobalControllerRouter["zdshop/controllers:OrderController"],
        beego.ControllerComments{
            Method: "PutOrderState",
            Router: "/:uid",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:OrderController"] = append(beego.GlobalControllerRouter["zdshop/controllers:OrderController"],
        beego.ControllerComments{
            Method: "AcceptOrder",
            Router: "/:uid",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:OrderController"] = append(beego.GlobalControllerRouter["zdshop/controllers:OrderController"],
        beego.ControllerComments{
            Method: "GetPayOrder",
            Router: "/:uid/:state",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["zdshop/controllers:OrderController"] = append(beego.GlobalControllerRouter["zdshop/controllers:OrderController"],
        beego.ControllerComments{
            Method: "HandlePay",
            Router: "/pay",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
