package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"zdshop/models"
	"zdshop/utils"
)

//operations for Category
type CategoryController struct {
	beego.Controller
}

// @Title Post
// @Description add category
// @Param name formData true "分类名称"
// @Param categoryId formData false "分类名称"
// @Param image formData false "分类名称"
// @Success  200 {string} 分类添加成功
// @Failure 400 分类添加失败
// @router / [post]
func (this *CategoryController) PostCategory() {
	//分类只做二级分类
	name := this.GetString("categoryName")
	if name == "" {
		beego.Info("请添加商品分类名称")
		return
	}
	resp := make(map[string]interface{})
	defer this.ServeJSON()

	o := orm.NewOrm()
	var category models.GoodsCategory
	category.Name = name
	exist := o.QueryTable("GoodsCategory").Filter("name", name).Exist()
	if exist {
		beego.Info("该分类名称已经存在, 请勿重复添加")
		resp["code"] = 400
		resp["errMsg"] = "该分类名称已经存在, 请勿重复添加"
		this.Data["json"] = resp
		return
	}

	var newCategory models.GoodsCategory
	id, err := this.GetInt("categoryId")
	if err != nil {
		//如果是一级分类, 只添加名字且pId = 0
		newCategory.Name = name
		newCategory.ParentId = 0
		_, err := o.Insert(&newCategory)
		if err != nil {
			beego.Info("插入商品分类失败")
			return
		}
		resp["code"] = 200
		resp["succMsg"] = "分类添加成功"
		this.Data["json"] = resp
		return
	}
	//否则是二级分类, 添加名字+图片且pId = id
	message, url := utils.HandleFile2(&this.Controller, "image")
	if message != "" {
		beego.Info("图片添加失败")
		return
	}
	newCategory.Name = name
	newCategory.ParentId = id
	newCategory.Image = url
	_, err = o.Insert(&newCategory)
	if err != nil {
		beego.Info("插入商品分类失败")
		return
	}
	resp["code"] = 200
	resp["succMsg"] = "分类添加成功"
	this.Data["json"] = resp
	return
}

// @Title Get First
// @Description 获取一级分类
// @Success  200 {string} 成功
// @Failure 400 失败
// @router / [get]
func (this *CategoryController) GetFirstCategory() {
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	var categorys []*models.GoodsCategory
	o := orm.NewOrm()
	_, err := o.QueryTable("GoodsCategory").Filter("ParentId", 0).OrderBy("-CreateTime").All(&categorys, "Id", "Name", "ParentId")
	if err != nil {
		beego.Info("分类查询失败：err", err)
		return
	}
	resp["code"] = 200
	resp["data"] = categorys
	this.Data["json"] = resp
	return
}

// @Title Get Second
// @Description 根据一级分类id获取二级级分类
// @Param name path true "分类id"
// @Success  200 {string} 成功
// @Failure 400 失败
// @router /:id [get]
func (this *CategoryController) GetSecondCategory() {
	id, err := this.GetInt(":id")
	if err != nil {
		beego.Info("传递参数有误")
		return
	}
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	var categorys []*models.GoodsCategory
	o := orm.NewOrm()
	_, err = o.QueryTable("GoodsCategory").Filter("ParentId", id).OrderBy("-CreateTime").All(&categorys, "Id", "Name", "Image")
	if err != nil {
		beego.Info("商品类型查询失败：err", err)
		return
	}
	resp["code"] = 200
	resp["data"] = categorys
	this.Data["json"] = resp
	return
}

// @Title Delete Category
// @Description 删除某个分类
// @Param cid path true "分类id"
// @Success  200 {string} 成功
// @Failure 400 失败
// @router /:cid [delete]
func (this *CategoryController) DeleteCategory() {
	id, err := this.GetInt(":cid")
	if err != nil {
		beego.Info("删除商品传递参数有误")
		return
	}
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	o := orm.NewOrm()
	var category models.GoodsCategory
	category.Id = id
	//可以判断这个分类在不在
	exist := o.QueryTable("GoodsCategory").Filter("Id", id).Exist()
	if !exist {
		beego.Info("该分类名称不存在")
		resp["code"] = 400
		resp["errMsg"] = "该分类名称不存在"
		this.Data["json"] = resp
		return
	}
	err = o.Read(&category)
	pId := category.ParentId
	if pId == 0 {
		//删除一级分类,以及该分类下所有分类
		o.Begin()
		_, err = o.QueryTable("GoodsCategory").Filter("ParentId__in", category.Id).Delete()
		//_, err = o.Delete(secondCategorys)
		if err != nil {
			o.Rollback()
			beego.Info(err)
		}
		_, err = o.Delete(&category)
		if err != nil {
			o.Rollback()
			beego.Info("商品删除失败")
			return
		}
		o.Commit()
		resp["code"] = 200
		resp["succMsg"] = "分类删除成功"
		this.Data["json"] = resp
		return
	}
	//删除二级分类
	_, err = o.Delete(&category)
	if err != nil {
		beego.Info("商品删除失败")
		return
	}
	resp["code"] = 200
	resp["succMsg"] = "分类删除成功"
	this.Data["json"] = resp
	return
}

// @Title Put Category
// @Description 修改某个分类名字
// @Param cid formDate true "分类id"
// @Param name formDate true "分类名称"
// @Success  200 {string} 成功
// @Failure 400 失败
// @router / [put]
func (this *GoodsController) UpdateCategory() {
	id, err := this.GetInt("cid")
	if err != nil {
		beego.Info("参数传递失败")
		return
	}
	name := this.GetString("name")
	o := orm.NewOrm()
	var category models.GoodsCategory
	category.Id = id
	category.Name = name
	_, err = o.Update(&category)
	if err != nil {
		beego.Info("更新失败")
		return
	}

}
