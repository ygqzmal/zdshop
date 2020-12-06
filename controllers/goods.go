package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"zdshop/models"
)

type GoodsController struct {
	beego.Controller
}

// @Title AddGoods
// @Description AddGoods and AddGoodsParameter and AddGoodsBanner
// @Param	name 		formData	true	"商品名称"
// @Param	brief 		formData	true	"商品简介"
// @Param	state 		formData	true	"商品状态 0-上架 1-下架"
// @Param	explain 	formData	true	"商品说明"
// @Param	categoryId 	formData	true	"商品分类id"
// @Success  200 {string} 商品添加成功
// @Failure 400 商品添加失败
// @router /AddGoods [post]
func (this *GoodsController) AddGoods() {
	//获取商品基本属性
	name := this.GetString("name")
	brief := this.GetString("brief")
	state := this.GetString("state")
	explain := this.GetString("explain")
	categoryId, _ := this.GetInt("categoryId")

	resp := make(map[string]interface{})
	defer this.ServeJSON()

	if name == "" || brief == "" || explain == "" || state == "" {
		resp["code"] = 1
		resp["errMsg"] = "提交内容不能为空"
		this.Data["json"] = resp
		return
	}
	if len(name) > 50 {
		resp["errMsg"] = "商品名称长度过长"
		this.Data["json"] = resp
		return
	}

	o := orm.NewOrm()
	var goods models.Goods
	goods.Name = name
	goods.GoodsBrief = brief
	goods.GoodsState = state
	goods.Explain = explain

	var category models.GoodsCategory
	category.Id = categoryId
	goods.Category = &category

	//事务
	o.Begin()

	//添加商品基本属性
	_, err := o.Insert(&goods)
	if err != nil {
		resp["errMsg"] = "商品添加失败"
		this.Data["json"] = resp
		o.Rollback()
		return
	}

	//添加商品参数

	//获取商品参数(数组)
	//parameters := this.GetStrings("parameters")
	//for _, value := range parameters {
	//
	//}
	var newGood models.Goods
	newGood.Name = name
	err = o.Read(&newGood, "Name")
	if err != nil {
		beego.Info("获取新商品id失败: ", err)
	}
	newId := newGood.Id
	beego.Info(newId)
	//提交事务
	o.Commit()
	resp["succMsg"] = "商品添加成功"
	//resp["data"] = goods
	this.Data["json"] = resp
}

// @Title UpdateGoods
// @Description UpdateGoods
// @Param name pwd 	true		"body for user content"
// @Success  200 {string} 修改商品成功
// @Failure 403 lost data
// @router /UpdateGoods [post]
func (this *GoodsController) UpdateGoods() {
	//获取商品基本属性
	id, _ := this.GetInt("id")
	name := this.GetString("name")
	brief := this.GetString("brief")
	state := this.GetString("state")
	explain := this.GetString("explain")
	categoryId, _ := this.GetInt("categoryId")

	resp := make(map[string]interface{})
	defer this.ServeJSON()

	var goods models.Goods
	goods.Id = id
	o := orm.NewOrm()
	err:= o.Read(&goods)
	if err != nil {
		resp["errMsg"] = "该商品不存在"
		this.Data["json"] = resp
		return
	}
	goods.Name = name
	goods.GoodsBrief = brief
	goods.GoodsState = state
	goods.Explain = explain
	var category models.GoodsCategory
	category.Id = categoryId
	goods.Category = &category

	_, err = o.Update(&goods)
	if err != nil {
		resp["errMsg"] = "商品更新失败"
		this.Data["json"] = resp
		return
	}
	resp["succMsg"] = "商品更新成功"
	this.Data["json"] = resp
	return
}

// @Title AddCategory
// @Description AddCategory
// @Param	username pwd 	true		"body for user content"
// @Success  200 {object} models.Category
// @Failure 403 lost data
// @router /AddCategory [post]
func (this *GoodsController) AddCategory() {

}

