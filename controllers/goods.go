package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"path"
	"strconv"
	"strings"
	"time"
	"zdshop/models"
)

//operations for Goods
type GoodsController struct {
	beego.Controller
}

type Goods struct {
	Goods []*models.Goods `json:"goods"`
	Code int `json:"code"`
	PageCount int `json:"page_count"`
	PageIndex int `json:"page_index"`
}

//这个好像没什么实际作用，只是为了方便一览这个类的所有方法
func (this *GoodsController) URLMapping() {
	this.Mapping("Post", this.Post)
	this.Mapping("GetOne", this.GetOne)
	this.Mapping("GetAll", this.GetAll)
	this.Mapping("Put", this.Put)
	this.Mapping("Delete", this.Delete)
}

// @Title Post
// @Description AddGoods and AddGoodsParameter and AddGoodsBanner
// @Param	name 		formData	true	"商品名称"
// @Param	brief 		formData	true	"商品简介"
// @Param	state 		formData	true	"商品状态 0-上架 1-下架"
// @Param	explain 	formData	true	"商品说明"
// @Param	categoryId 	formData	true	"商品分类id"
// @Success  200 {string} 商品添加成功
// @Failure 400 商品添加失败
// @router / [post]
func (this *GoodsController) Post() {
	//获取商品基本属性
	name := this.GetString("name")
	brief := this.GetString("brief")
	state := this.GetString("state")
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
	if name == "" || brief == "" || explain == "" || state == "" {
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
		o.Rollback()
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

	var category models.GoodsCategory
	category.Id = categoryId
	goods.Category = &category

	//添加商品基本属性
	_, err = o.Insert(&goods)
	if err != nil {
		beego.Info(err)
		resp["code"] = 400
		resp["errMsg"] = "商品添加失败"
		this.Data["json"] = resp
		return
	}

	//添加商品参数=========
	//模拟前端传来数据 参数、市场价、实际价、是否默认、库存
	//参数不能有空格（前端限制）
	parameter1 := "[限量版内存:68G 5000 4500 1 100]"
	parameter2 := "[普通版内存:68G 4000 3500 0 100]"
	parameter3 := "[特别版内存:256G 6000 5000 0 100]"
	data := []string{parameter1, parameter2, parameter3}
	//获取商品参数(数组)
	//parameters := this.GetStrings("parameters")
	parameters := data

	//获取新商品
	var newGood models.Goods
	newGood.Name = name
	err = o.Read(&newGood, "Name")
	if err != nil {
		beego.Info("获取新商品id失败: ", err)
	}
	//newId := newGood.Id
	//beego.Info(newId)

	for _, value := range parameters {
		//获取每个具体参数
		tmp := value[1 : len(value)-1]
		para := strings.Split(tmp, " ")
		if len(para) > 5 {
			beego.Info("参数不能空或添加参数过长")
			err := o.Rollback()
			beego.Info("rollback err :", err)
			resp["code"] = 400
			resp["errMsg"] = "传入参数有问题"
			this.Data["json"] = resp
			return
		}
		content := para[0]
		truePrice, err := strconv.ParseFloat(para[1], 64) //转float64
		NowPrice, err := strconv.ParseFloat(para[2], 64)
		IsDefault := para[3]
		number, err := strconv.Atoi(para[4])

		if len(content) >= 30 || content == "" {
			beego.Info("参数不能空或添加参数过长")
			err := o.Rollback()
			beego.Info("rollback err :", err)
			resp["code"] = 400
			resp["errMsg"] = "参数不能空或添加参数过长"
			this.Data["json"] = resp
			return
		}

		if err != nil {
			o.Rollback()
			resp["code"] = 400
			resp["errMsg"] = "输入价格或库存有误"
			this.Data["json"] = resp
			return
		}
		//应当对price进行判断，正则匹配是否为浮点型

		//存入数据库
		var goodsPara models.GoodsParameter
		//goodsPara.GoodId = newId

		goodsPara.Goods = &newGood
		goodsPara.Parameter = content
		goodsPara.GoodsTruePrice = truePrice
		goodsPara.GoodsNowPrice = NowPrice
		goodsPara.IsDefault = IsDefault
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

	//添加商品图片======
	//模拟前端传来数据 排序 url

	//headers, err := this.GetFiles("picture")
	//for循环每次得到以下数据 [0 ./upload/img/01.jpg] 是否默认、路径

	goodPicture1 := "./upload/img/01.jpg"
	goodPicture2 := "./upload/img/02.jpg"
	goodPicture3 := "./upload/img/03.jpg"
	pictures := []string{goodPicture1, goodPicture2, goodPicture3}
	//beego.Info(prices)
	for _, value := range pictures {
		//获取每个具体的参数
		var goodsBanner models.GoodsBanner
		goodsBanner.Goods = &newGood
		//goodsBanner.GoodsOrder = key
		goodsBanner.GoodsUrl = value

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
	//提交事务
	o.Commit()
	resp["code"] = 200
	resp["succMsg"] = "商品添加成功"
	this.Data["json"] = resp
}

// @Title Update
// @Description UpdateGoods
// @Param name pwd 	true		"body for user content"
// @Success  200 {string} 修改商品成功
// @Failure 403 lost data
// @router /id [put]
func (this *GoodsController) Put() {
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
	err := o.Read(&goods)
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

// @Title Get One
// @Description Get Good
// @Param gid path true "good ID"
// @Success  200 {string} ok
// @Failure 403 lost data
// @router /:gid [get]
func (this *GoodsController) GetOne() {
	//id := this.Ctx.Input.Param(":gid")
	id, err := this.GetInt(":gid")
	beego.Info(id)
	if err != nil {
		beego.Error("浏览器请求错误")
		this.Redirect("/", 302)
		return
	}
	o := orm.NewOrm()
	//商品基本属性
	var good models.Goods
	good.Id = id
	//err = o.Read(&good)  查询goods表Id为id的值, 想要查关联的外键的话,要加上relatedSel("model中的外键字段")
	err = o.QueryTable("Goods").RelatedSel("Category").Filter("Id", id).One(&good)
	if err != nil {
		beego.Info("商品查询失败")
		return
	}
	//商品参数(一对多), 如果加上relatedSel反而会加强查询结果，把所有关联都查询到。
	//如果想指定查询内容, 在all后面加上查询参数即可
	var parameters []*models.GoodsParameter
	_, err = o.QueryTable("GoodsParameter").Filter("Goods__Id", id).All(&parameters, "Parameter", "GoodsTruePrice", "GoodsNowPrice")
	if err != nil {
		beego.Info("商品参数查询失败")
		return
	}
	//商品图片
	var banner []*models.GoodsBanner
	_, err = o.QueryTable("GoodsBanner").Filter("Goods__Id", id).All(&banner, "GoodsUrl")
	if err != nil {
		beego.Info("商品参数查询失败")
		return
	}
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	beego.Info(123)
	resp["code"] = 200
	resp["goods"] = good
	resp["goodsPara"] = parameters
	resp["goodsBanner"] = banner
	this.Data["json"] = resp
}

// @Title Get All
// @Description Get Goods
// @Param cid 	    path true "category id"
// @Param pageIndex path false "pageIndex"
// @Success  200 {object} controllers.Goods
// @Failure 403 lost data
// @router /:cid/:pageIndex [get]
func (this *GoodsController) GetAll() {
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
	count, err := o.QueryTable("Goods").Filter("Category__Id", id).Count()
	if err != nil {
		beego.Info("获取当前分类商品数量失败")
		return
	}
	pageSize := 1
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
	var data = Goods{}
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
func (this *GoodsController) Delete() {
	id, err := this.GetInt(":gid")
	if err != nil {
		beego.Info("获取参数失败")
		return
	}
	var goods models.Goods
	goods.Id = id
	o := orm.NewOrm()
	//删除商品基本属性
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

func (this *GoodsController) Tex() {
	data := this.GetStrings("data")
	beego.Info(data)
}

//文件上传
func HandleFile(this *beego.Controller, image string) (message, url string) {
	//headers, err := this.GetFiles(image)
	file, header, err := this.GetFile(image)
	defer file.Close()
	if err != nil {
		beego.Info("图片上传失败")
		return "图片上传失败", "nil"
	}
	//大小
	if header.Size > 500000 {
		beego.Info("图片太大无法上传")
		return "图片太大无法上传", "nil"
	}
	//格式
	ext := path.Ext(header.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		beego.Info("图片格式不正确")
		return "图片格式不正确", "nil"
	}
	//防重名
	//fileName := time.Now().Format("2006-01-02-15:04:05") + ext
	//存储 第一个参数要和GetFile的一样
	//err = this.SaveToFile("image", fileName)
	//3.防止重名
	//fileName := time.Now().Format("2006-01-02-15-04-05") + ext 好像文件名重复不会覆盖
	fileName := strconv.Itoa(int(time.Now().UnixNano())) + ext //这个不会重名
	//存储
	err = this.SaveToFile(image, "./upload/img/"+fileName)
	if err != nil {
		beego.Info("文件存储错误 ", err)
		return "文件存储错误 ", "nil"
	}
	beego.Info("图片添加成功")
	return "图片添加成功", fileName
}
