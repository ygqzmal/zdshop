package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"os"
	"path"
	"strconv"
	"time"
	"zdshop/models"
	"zdshop/utils"
)

//operations for Goods
type GoodsController struct {
	beego.Controller
}

type Goods struct {
	Goods     []*models.Goods `json:"goods"`
	Code      int             `json:"code"`
	PageCount int             `json:"page_count"`
	PageIndex int             `json:"page_index"`
}

type JsonPara struct {
	Parameter string
	TruePrice float64
	NowPrice  float64
	IsDefault int
	Number    int
}

// @Title Post
// @Description AddGoods and AddGoodsParameter and AddGoodsBanner
// @Param	name 		formData	true	"商品名称"
// @Param	brief 		formData	true	"商品简介"
// @Param	state 		formData	true	"商品状态 1-上架 2-下架"
// @Param	explain 	formData	true	"商品说明"
// @Param	categoryId 	formData	true	"商品分类id"
// @Param	parameters 	formData	true	"商品参数"
// @Param	img 	    formData	true	"商品默认图片"
// @Param	imgs 	    formData	false	"商品轮播图片"
// @Success  200 {string} 商品添加成功
// @Failure 400 商品添加失败
// @router / [post]
func (this *GoodsController) PostGoods() {
	//获取商品基本属性
	name := this.GetString("name")
	brief := this.GetString("brief")
	state, _ := this.GetInt("state")
	explain := this.GetString("explain")
	categoryId, _ := this.GetInt("categoryId")
	//返回map
	resp := make(map[string]interface{})
	defer this.ServeJSON()

	//事务
	o := orm.NewOrm()
	err := o.Begin()
	beego.Info("事务开启：", err)
	//检验数据
	if name == "" || brief == "" || explain == "" || categoryId == 0 {
		resp["code"] = 400
		resp["errMsg"] = "提交内容不能为空"
		this.Data["json"] = resp
		err := o.Rollback()
		beego.Info("rollback err :", err)
		return
	}
	if len(name) > 50 {
		resp["code"] = 400
		resp["errMsg"] = "商品名称长度过长"
		this.Data["json"] = resp
		err := o.Rollback()
		beego.Info("rollback err :", err)
		return
	}

	var goods models.Goods
	goods.Name = name
	exist := o.QueryTable("Goods").Filter("name", name).Exist()
	if exist {
		o.Rollback()
		beego.Info("该商品名称已经存在")
		return
	}

	goods.GoodsBrief = brief
	goods.GoodsState = state
	goods.Explain = explain
	//添加外键关联商品分类id
	var category models.GoodsCategory
	category.Id = categoryId
	goods.Category = &category

	//添加商品基本属性
	_, err = o.Insert(&goods)
	if err != nil {
		o.Rollback()
		beego.Info(err)
		resp["code"] = 400
		resp["errMsg"] = "商品添加失败,商品详情内容过大,网络连接出现问题"
		this.Data["json"] = resp
		return
	}

	//添加商品参数
	//前端传来数据 参数、市场价、实际价、是否默认、库存
/*	[
/	{
//		"parameter":"321",
//		"truePrice":222,
//		"nowPrice":111,
//		"isDefault":1,
//		"number":200
//	},
//	{
//		"parameter":"1 2 3",
//		"truePrice":2222,
//		"nowPrice":1111,
//		"isDefault":1,
//		"number":2000
//	}/
]*/
	data := this.GetString("parameters")
	beego.Info("json: ", data)
	var parameters []JsonPara
	err = json.Unmarshal([]byte(data), &parameters)
	if err != nil {
		o.Rollback()
		beego.Info("json.Unmarshal err: ", err)
		beego.Info("传入数据有误") //不能为空且不能类型不对，否则无法unmarshal
		return
	}

	beego.Info(parameters)
	for _, value := range parameters {
		//获取每个具体参数
		parameter := value.Parameter
		truePrice := value.TruePrice
		nowPrice := value.NowPrice
		isDefault := value.IsDefault
		number := value.Number
		if len(parameter) >= 30 {
			beego.Info("参数过长")
			err := o.Rollback()
			beego.Info("rollback err :", err)
			resp["code"] = 400
			resp["errMsg"] = "参数过长"
			this.Data["json"] = resp
			return
		}
		//存入数据库
		var goodsPara models.GoodsParameter
		goodsPara.Goods = &goods
		goodsPara.Parameter = parameter
		goodsPara.GoodsTruePrice = truePrice
		goodsPara.GoodsNowPrice = nowPrice
		goodsPara.IsDefault = isDefault
		goodsPara.GoodsNumber = number

		_, err = o.Insert(&goodsPara)
		if err != nil {
			o.Rollback()
			beego.Info(err)
			resp["code"] = 400
			resp["errMsg"] = "参数添加失败, 商品添加失败"
			this.Data["json"] = resp
			return
		}
	}
	//添加商品图片
	//先添加默认图片
	if _, _, err = this.GetFile("img"); err != nil {
		resp["code"] = 400
		resp["errMsg"] = "默认图片不能为空"
		this.Data["json"] = resp
		beego.Info("默认图片不能为空")
		return
	}
	errMsg, url := utils.HandleFile(&this.Controller, "img", 200000)
	if errMsg != "" {
		beego.Info("默认图片上传失败, err:", errMsg)
		return
	}
	var goodsBanner models.GoodsBanner
	goodsBanner.Goods = &goods
	goodsBanner.IsDefault = 1
	goodsBanner.GoodsUrl = url
	_, err = o.Insert(&goodsBanner)
	if err != nil {
		o.Rollback()
		beego.Info("默认图片添加失败 err: ", err)
		return
	}
	headers, err := this.GetFiles("imgs")
	//判断一下有多少张轮播图(限定6张之内)
	if len(headers) > 6 {
		beego.Info("图片轮播图片不得超过6张")
		return
	}
	//如果有数据则处理
	if err == nil {
		beego.Info("存在轮播图")
		for i := range headers {
			//获取每个具体的参数
			if headers[i].Size > 200000 {
				beego.Info("图片太大")
				return
			}
			ext := path.Ext(headers[i].Filename)
			if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
				beego.Info("图片格式不正确")
				return
			}
			fileName := strconv.Itoa(int(time.Now().UnixNano())) + ext
			err = this.SaveToFile("imgs", "./upload/img/"+fileName)
			if err != nil {
				beego.Info("文件存储错误 ")
				return
			}

			var goodsBanner models.GoodsBanner
			goodsBanner.Goods = &goods
			goodsBanner.GoodsUrl = "./upload/img/" + fileName
			goodsBanner.IsDefault = 0

			_, err = o.Insert(&goodsBanner)
			beego.Info("轮播图插入成功")
			if err != nil {
				o.Rollback()
				beego.Info(err)
				resp["code"] = 400
				resp["errMsg"] = "添加轮播图片失败, 商品添加失败"
				this.Data["json"] = resp
				return
			}
		}
	}
	//提交事务
	o.Commit()
	beego.Info("商品添加成功")
	resp["code"] = 200
	resp["succMsg"] = "商品添加成功"
	this.Data["json"] = resp
}

// @Title Get One
// @Description Get Good
// @Param gid path true "good ID"
// @Success  200 {string} ok
// @Failure 403 lost data
// @router /:gid [get]
func (this *GoodsController) GetOneGoods() {
	id, err := this.GetInt(":gid")
	if err != nil {
		beego.Error("浏览器请求错误")
		this.Redirect("/", 302)
		return
	}
	o := orm.NewOrm()
	//商品基本属性
	var good models.Goods
	//good.Id = id
	//err = o.Read(&good) //这个查询只会将外键id查询出来,而且是查询所有,无法指定查询字段
	//添加的分类id一定要存在，否则会报错
	//查询goods表Id为id的值, 想要查关联的外键的话,要加上relatedSel("model中的外键字段"), 会把外键所有数据查询出来而且无法指定
	//err = o.QueryTable("Goods").RelatedSel("Category").Filter("Id", id).One(&good,"Id","Name","GoodsBrief","Explain","SalesValue")
	//想对外键的表进行某个字段查询做不到，只能查询外键的id
	err = o.QueryTable("Goods").Filter("Id", id).One(&good, "Id", "Name", "GoodsBrief", "Explain", "SalesValue", "category_id")
	if err != nil {
		beego.Info("商品查询失败: err", err)
		return
	}
	//商品参数(一对多), 如果加上relatedSel反而会加强查询结果，把所有关联都查询到。
	//如果想指定查询内容, 在all后面加上查询参数即可
	var parameters []*models.GoodsParameter
	_, err = o.QueryTable("GoodsParameter").Filter("Goods__Id", id).All(&parameters, "Id", "Parameter", "GoodsTruePrice", "GoodsNowPrice", "IsDefault","GoodsNumber")
	if err != nil {
		beego.Info("商品参数查询失败")
		return
	}
	//商品图片
	var banner []*models.GoodsBanner
	_, err = o.QueryTable("GoodsBanner").Filter("Goods__Id", id).OrderBy("Id").All(&banner, "Id", "GoodsUrl", "IsDefault")
	if err != nil {
		beego.Info("商品参数查询失败")
		return
	}
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	resp["code"] = 200
	resp["goods"] = good
	resp["goodsPara"] = parameters
	resp["goodsBanner"] = banner
	this.Data["json"] = resp
}

// @Title Get All
// @Description Get Goods
// @Param cid 	    path true "category id"
// @Param pageIndex path true "pageIndex"
// @Success  200 {object} controllers.Goods
// @Failure 403 lost data
// @router /:cid/:pageIndex [get]
func (this *GoodsController) GetAllGoods() {
	beego.Info("GetAll")
	//展示一个分类下全部商品(名字+默认图片+默认参数)
	id, err := this.GetInt(":cid")
	if err != nil {
		beego.Info("请求参数有问题")
		return
	}
	o := orm.NewOrm()
	//分页展示。每页10条数据
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}
	beego.Info(pageIndex)
	count, err := o.QueryTable("Goods").Filter("Category__Id", id).Filter("GoodsState", 0).Count()
	if err != nil {
		beego.Info("获取当前分类商品数量失败")
		return
	}
	pageSize := 2  //每页多少条数据
	pageCount := math.Ceil(float64(count) / float64(pageSize))
	start := (pageIndex - 1) * pageSize

	//商品基本属性
	var goods []*models.Goods
	//OrderBy里面加上-表示desc
	_, err = o.QueryTable("Goods").Filter("Category__Id", id).Filter("GoodsState", 0).OrderBy("-CreateTime").Limit(pageSize, start).All(&goods, "name", "Id")
	if err != nil {
		beego.Info("商品查询失败 err: ", err)
		return
	}

	//这里虽然只是查询一条数据, 但是也要用切片, 因为para在赋值的给good中的GoodsParameter的时候需要切片类型   model中设计时决定(GoodsParameter []*GoodsParameter)
	var para []*models.GoodsParameter
	var banner []*models.GoodsBanner
	for _, good := range goods {
		//商品默认参数
		err = o.QueryTable("GoodsParameter").Filter("Goods__Id", good.Id).Filter("IsDefault", 1).One(&para, "Parameter", "GoodsTruePrice", "GoodsNowPrice")
		if err != nil {
			beego.Info("商品参数查询失败 err: ", err)
			return
		}
		good.GoodsParameter = para

		//商品默认图片
		err = o.QueryTable("GoodsBanner").Filter("Goods__Id", good.Id).Filter("IsDefault", 1).One(&banner, "GoodsUrl")
		if err != nil {
			beego.Info("商品图片查询失败 err: ", err)
			return
		}
		good.GoodsBanner = banner
	}

	//resp := make(map[string]interface{})
	defer this.ServeJSON()
	//resp["code"] = 200
	//resp["pageCount"] = pageCount
	//resp["pageIndex"] = pageIndex
	//resp["goods"] = goods
	data := Goods{}
	data.Goods = goods
	data.PageIndex = pageIndex
	data.PageCount = int(pageCount)
	data.Code = 200
	this.Data["json"] = data
	beego.Info(goods)
}

// @Title Delete
// @Description delete the good
// @Param gid path true	"body for user content"
// @Success  200 {string} ok
// @Failure 403 lost data
// @router /:gid [delete]
func (this *GoodsController) DeleteGoods() {
	id, err := this.GetInt(":gid")
	if err != nil {
		beego.Info("获取参数失败")
		return
	}
	o := orm.NewOrm()
	//删除商品轮播图片
	var goodsImg []*models.GoodsBanner
	_, err = o.QueryTable("GoodsBanner").Filter("Goods__Id", id).All(&goodsImg, "GoodsUrl")
	if err != nil {
		beego.Info("商品轮播图片查询失败")
		return
	}
	for _, value := range goodsImg {
		imgPath := value.GoodsUrl
		err := os.Remove(imgPath)
		if err != nil {
			beego.Info("图片删除失败")
			return
		}
	}
	//删除商品基本属性
	var goods models.Goods
	goods.Id = id
	_, err = o.Delete(&goods)
	if err != nil {
		beego.Info("商品删除失败")
		return
	}
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	resp["code"] = 200
	resp["succMsg"] = "商品删除成功"
	this.Data["json"] = resp
}

// @Title Update Goods Attribute
// @Description UpdateGoodsAttribute
// @Param gid formDate true	"id for goods goodsId"
// @Param cid formDate true	"cid for goods categoryId"
// @Param name formDate true	"name for goods goodsName"
// @Success  200 {string} 修改商品成功
// @Failure 403 lost data
// @router /putAttribute [put]
func (this *GoodsController) PutGoodsAttribute() {
	var goods models.Goods
	o := orm.NewOrm()
	//获取商品基本属性
	id, err := this.GetInt("id")
	if err != nil {
		beego.Info("传递id有误")
		return
	}
	name := this.GetString("name")
	brief := this.GetString("brief")
	state, err := this.GetInt("state")
	if err != nil {
		beego.Info("请选择商品状态")
		return
	}
	explain := this.GetString("explain")
	//分类id一定要传
	categoryId, err := this.GetInt("categoryId")
	if err != nil {
		beego.Info("分类id不能为空")
		return
	}
	goods.Id = id
	goods.Name = name
	goods.GoodsBrief = brief
	goods.Explain = explain
	goods.GoodsState = state
	var category models.GoodsCategory
	category.Id = categoryId
	goods.Category = &category
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	//这里虽然在model中time的tag不一样，但是没办法做到自动更新时间，需要指定updateTime
	_, err = o.Update(&goods, "Name", "GoodsBrief", "Explain", "GoodsState", "Category", "UpdateTime")
	//下面这个更新如果没有指定字段的话，会设置为null，但是model中建表时默认不为空，所以无法这样更新
	//_, err = o.Update(&goods)
	if err != nil {
		beego.Info(err)
		resp["errMsg"] = "商品更新失败"
		this.Data["json"] = resp
		return
	}
	resp["succMsg"] = "商品更新成功"
	this.Data["json"] = resp
	return
}

// @Title Update Goods Parameters
// @Description UpdateGoodsParameters
// @Param gid path true	"id for goods"
// @Success  200 {string} 修改商品成功
// @Failure 403 lost data
// @router /putParameter [put]
func (this *GoodsController) PutGoodsParameters() {
	pid, err := this.GetInt("pid")
	if err != nil {
		beego.Info("传递参数id有误")
		return
	}
	gid, err := this.GetInt("gid")
	if err != nil {
		beego.Info("传递商品id有误")
		return
	}
	parameter := this.GetString("parameter")
	truePrice, err := this.GetFloat("truePrice")
	if err != nil {
		beego.Info("请输入正确价格", err)
		return
	}
	nowPrice, err := this.GetFloat("nowPrice")
	if err != nil {
		beego.Info("请输入正确价格", err)
		return
	}
	number, err := this.GetInt("number")
	if err != nil {
		beego.Info(err)
		return
	}
	isDefault, err := this.GetInt("isDefault")
	if err != nil {
		beego.Info("请选择默认参数")
		return
	}
	o := orm.NewOrm()
	o.Begin()
	//如果改为默认参数
	if isDefault == 1 {
		//更新原来默认参数
		var goodsDefaultPara models.GoodsParameter
		err = o.QueryTable("GoodsParameter").Filter("Goods__Id", gid).Filter("IsDefault", 1).One(&goodsDefaultPara, "Id", "IsDefault")
		goodsDefaultPara.IsDefault = 0
		_, err = o.Update(&goodsDefaultPara, "IsDefault")
		if err != nil {
			o.Rollback()
			beego.Info(err)
			return
		}
	}
	var goodsParameter models.GoodsParameter
	goodsParameter.Id = pid
	goodsParameter.Parameter = parameter
	goodsParameter.GoodsTruePrice = truePrice
	goodsParameter.GoodsNowPrice = nowPrice
	goodsParameter.GoodsNumber = number
	goodsParameter.IsDefault = isDefault
	_, err = o.Update(&goodsParameter, "Parameter", "GoodsTruePrice", "GoodsNowPrice", "GoodsNumber", "IsDefault")
	if err != nil {
		o.Rollback()
		beego.Info("商品参数更新失败", err)
		return
	}
	o.Commit()
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	resp["code"] = 200
	resp["succMsg"] = "商品参数更新成功"
	this.Data["json"] = resp
	return
}

// @Title Update Goods Banner
// @Description UpdateGoodsBanner
// @Param gid path true	"id for goods"
// @Success  200 {string} 修改商品成功
// @Failure 403 lost data
// @router /putBanner [put]
func (this *GoodsController) PutGoodsBanners() {

}

// @Title Update Goods State
// @Description UpdateGoodsState
// @Param gid formDate true	"id for goods id"
// @Param state formDate true "商品上下架"
// @Success  200 {string} 修改商品成功
// @Failure 403 lost data
// @router /putState [put]
func (this *GoodsController) PutGoodsState() {
	gid, err := this.GetInt("gid")
	if err != nil {
		beego.Info("参数传递错误")
		return
	}
	state, err := this.GetInt("state")
	if err != nil {
		beego.Info("商品上下架参数传递错误")
		return
	}
	var goods models.Goods
	goods.Id = gid
	goods.GoodsState = state
	o := orm.NewOrm()
	_, err = o.Update(&goods, "GoodsState", "UpdateTime")
	if err != nil {
		beego.Info("商品上下架更新失败")
		return
	}
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	resp["code"] = 200
	if state == 0 {
		resp["succMsg"] = "该商品上架成功"
		this.Data["json"] = resp
		return
	}
	resp["succMsg"] = "该商品下架成功"
	this.Data["json"] = resp
	return
}

// @Title Delete Parameter
// @Description Delete One Parameter
// @Param pid path true	"id for parameter id"
// @Success  200 {string} 参数删除成功
// @Failure 403 lost data
// @router /para/:pid [delete]
func (this *GoodsController) DeleteParameter() {
	//删除某个参数
	pid, err := this.GetInt(":pid")
	if err != nil {
		beego.Info(err)
		return
	}
	var parameter models.GoodsParameter
	parameter.Id = pid
	o := orm.NewOrm()
	_, err = o.Delete(&parameter)
	if err != nil {
		beego.Info("商品删除失败: err", err)
		return
	}
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	resp["code"] = 200
	resp["succMsg"] = "商品参数删除成功"
	this.Data["json"] = resp
	return
}

// @Title Put BannerImage
// @Description update goods BannerImage
// @Param bid formData true	"body for user content"
// @Param img formData true "default image"
// @Success  200 {string} ok
// @Failure 403 lost data
// @router /img [put]
func (this *GoodsController) PutGoodsBannerImg() {
	//更新商品一张轮播图片, 不管是否默认
	bId, _ := this.GetInt("bid")

	o := orm.NewOrm()
	var goodsBanner models.GoodsBanner
	//删除原有图片
	err := o.QueryTable("GoodsBanner").Filter("Id", bId).One(&goodsBanner, "GoodsUrl")
	if err != nil {
		beego.Info(err)
		return
	}
	oldUrl := goodsBanner.GoodsUrl
	err = os.Remove(oldUrl)
	if err != nil {
		beego.Info(err)
		return
	}
	errMsg, url := utils.HandleFile(&this.Controller, "img", 200000)
	if errMsg != "" {
		beego.Info(errMsg)
		return
	}
	count, _ := o.QueryTable("GoodsBanner").Filter("Id", bId).Update(orm.Params{"GoodsUrl": url})
	if count == 0 {
		beego.Info("更新失败 err")
		return
	}
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	resp["code"] = 200
	resp["succMsg"] = "更新成功"
	this.Data["json"] = resp
	return
}

// @Title Delete BannerImage
// @Description update default goods image
// @Param bid path true	"body for user content"
// @Success  200 {string} ok
// @Failure 403 lost data
// @router /img/:bid [delete]
func (this *GoodsController) DelGoodsBannerImg() {
	bid, err := this.GetInt(":bid")
	if err != nil {
		beego.Info(err)
		return
	}
	o := orm.NewOrm()
	var goodsBanner models.GoodsBanner
	err = o.QueryTable("GoodsBanner").Filter("Id", bid).One(&goodsBanner, "Id", "GoodsUrl")
	if err != nil {
		beego.Info(err)
		return
	}
	//删除图片
	err = os.Remove(goodsBanner.GoodsUrl)
	//删除数据
	_, err = o.Delete(&goodsBanner)
	if err != nil {
		beego.Info(err)
		return
	}
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	resp["code"] = 200
	resp["succMsg"] = "图片删除成功"
	this.Data["json"] = resp
	return
}

// @Title Post BannerImage
// @Description post goods BannerImage
// @Param gid formData true "id for goods id"
// @Param img formData true	"body for user content"
// @Success  200 {string} ok
// @Failure 403 lost data
// @router /img [post]
func (this *GoodsController) PostGoodsBannerImg() {
	gid, err := this.GetInt("gid")
	if err != nil {
		beego.Info(err)
		return
	}
	resp := make(map[string]interface{})
	defer this.ServeJSON()

	o := orm.NewOrm()
	count, err := o.QueryTable("GoodsBanner").Filter("Goods__Id", gid).Count()
	if err != nil {
		beego.Info(err)
		return
	}
	headers, err := this.GetFiles("imgs")
	if err != nil {
		beego.Info(err)
		return
	}
	newCount := len(headers)
	beego.Info(newCount)
	if count > 6 || newCount+int(count) > 6 {
		resp["code"] = 400
		resp["errMsg"] = "轮播图片不能超过6张"
		this.Data["json"] = resp
		return
	}
	//数据处理
	o.Begin()
	var goods models.Goods
	goods.Id = gid
	for i, _ := range headers {
		//获取每个具体的参数
		if headers[i].Size > 200000 {
			beego.Info("图片太大")
			return
		}
		ext := path.Ext(headers[i].Filename)
		if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
			beego.Info("图片格式不正确")
			return
		}
		fileName := strconv.Itoa(int(time.Now().UnixNano())) + ext
		err = this.SaveToFile("imgs", "./upload/img/"+fileName)
		if err != nil {
			beego.Info("文件存储错误 ")
			return
		}

		var goodsBanner models.GoodsBanner
		goodsBanner.Goods = &goods
		goodsBanner.GoodsUrl = "./upload/img/" + fileName
		goodsBanner.IsDefault = 0

		_, err = o.Insert(&goodsBanner)
		if err != nil {
			o.Rollback()
			beego.Info(err)
			resp["code"] = 400
			resp["errMsg"] = "添加图片失败, 商品添加失败"
			this.Data["json"] = resp
			return
		}
	}
	o.Commit()
	resp["code"] = 200
	resp["succMsg"] = "商品轮播图片添加成功"
	this.Data["json"] = resp
	return
}

// @Title search
// @Description 搜索某个商品
// @Param goodsName formData true "good name"
// @Param img formData true	"body for user content"
// @Success  200 {string} ok
// @Failure 403 lost data
// @router /search [post]
func (this *GoodsController) HandleSearch() {
	//获取数据
	goodsName := this.GetString("goodsName")
	o := orm.NewOrm()
	var goods []models.Goods
	//校验数据
	if goodsName == "" {
		_, err := o.QueryTable("Goods").All(&goods)
		if err != nil {
			beego.Info(err)
		}
		return
	}
	//icontains 模糊查找 i表示忽略大小写
	count, err := o.QueryTable("Goods").Filter("Name__icontains", goodsName).All(&goods)
	if err != nil {
		beego.Info(err)
	}

	resp := make(map[string]interface{})
	defer this.ServeJSON()

	if count == 0 {
		resp["code"] = 200
		resp["msg"] = "查找不到相关商品"
		this.Data["json"] = resp
		return
	}

	resp["code"] = 200
	resp["goods"] = goods
	resp["count"] = count
	this.Data["json"] = resp
	return
}


// @Title Tex
// @Description Tex
// @Param gid path true	"body for user content"
// @Success  200 {string} ok
// @Failure 403 lost data
// @router /tex [post]
func (this *GoodsController) Tex() {
	id, _ := this.GetInt("gid")
	o := orm.NewOrm()
	var goods models.Goods
	err := o.QueryTable("Goods").Filter("Id", id).One(&goods, "Name", "CreateTime", "UpdateTime")
	if err != nil {
		beego.Info(err)
		return
	}
	//返回数据库格式化时间：2020-12-16 21:24:35
	createTime := goods.CreateTime.Format("2006-01-02 15:04:05")
	//返回当前格式化时间
	now := time.Now()
	nowTime := now.Format("2006-01-02 15:04:05")
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	resp["code"] = 200
	resp["data"] = goods
	resp["createTime"] = createTime
	resp["nowTime"] = nowTime
	this.Data["json"] = resp
	return
}
