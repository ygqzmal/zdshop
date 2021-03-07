package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"zdshop/models"
)

type CartController struct {
	beego.Controller
}

type Cart struct {
	GoodsParameter models.GoodsParameter `json:"goodsParameter"`
	Count          int                   `json:"count"`
}

// @Title Post
// @Description add cart
// @Param pid formData true "参数id"
// @Param count formData true "添加数量"
// @Success  200 {string} 购物车添加成功
// @Failure 400 购物车添加失败
// @router / [post]
func (this *CartController) PostCart() {
	//获取数据
	//gid, err1 := this.GetInt("gid")
	pid, err1 := this.GetInt("pid")
	count, err2 := this.GetInt("count")
	resp := make(map[string]interface{})
	defer this.ServeJSON()

	//校验数据
	if err1 != nil || err2 != nil {
		resp["code"] = 1
		resp["msg"] = "传递的数据不正确"
		this.Data["json"] = resp
		return
	}
	//获取用户id
	//userName := this.GetSession("userName")
	//userName := "Jack"
	//o := orm.NewOrm()
	//var user models.User
	//user.Name = userName.(string)
	//user.Name = userName
	//err := o.Read(&user, "Name")
	//if err != nil {
	//	beego.Info("该用户不存在")
	//	return
	//}
	//处理数据
	//购物车数据存在redis中, 用hash
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		beego.Info("redis连接错误")
		return
	}
	defer conn.Close()
	//先获取原来的数量，然后给数量加起来
	//hgetall不会在原来的基础上进行累加,hget会在原来进行上进行累加
	//例如原来购物车里面有这个商品10个数量,再次添加购物车的时候会在10的基础上进行添加
	//应该是没有hgetall hashName key 这样的写法的,所以无法获取到数据 就为0
	//preCount, err := redis.Int(conn.Do("hgetall", "cart_"+strconv.Itoa(user.Id), gid))
	//_, err = conn.Do("hset", "cart_"+strconv.Itoa(user.Id), gid, count+preCount)
	//rep, err := conn.Do("hlen", "cart_"+strconv.Itoa(user.Id))

	//商品数量
	preCount, err := redis.Int(conn.Do("hget", "cart_1", pid))
	_, err = conn.Do("hset", "cart_1", pid, count+preCount)
	rep, err := conn.Do("hlen", "cart_1")
	//回复助手函数
	cartCount, _ := redis.Int(rep, err)
	resp["code"] = 200
	resp["msg"] = "OK"
	resp["cartCount"] = cartCount
	this.Data["json"] = resp
	return
}

// @Title Get
// @Description Get cart
// @Success  200 {object} controllers.Cart
// @Failure 400 显示失败
// @router / [get]
func (this *CartController) GetCart() {
	//用户信息
	//userName := this.GetSession("userName")

	//从redis中获取数据
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		beego.Info("redis连接错误")
		return
	}
	defer conn.Close()

	o := orm.NewOrm()
	//var user models.User
	//user.Name = userName.(string)
	//err = o.Read(&user, "Name")
	//if err != nil {
	//	beego.Info("查询用户失败")
	//	return
	//}
	//goodsMap, _ := redis.IntMap(conn.Do("hgetall", "cart_"+strconv.Itoa(user.Id))) //map[string]int
	goodsParaMap, _ := redis.IntMap(conn.Do("hgetall", "cart_1")) //map[string]int
	//beego.Info(goodsMap)
	//map[1:3 3:2]  map[key:value key:value] key:skuId(参数ID) value:count(商品数量)
	var carts []Cart
	for key, value := range goodsParaMap {
		var goodsParameter models.GoodsParameter
		pid, _ := strconv.Atoi(key)
		goodsParameter.Id = pid
		//err := o.Read(&goodsParameter)
		err := o.QueryTable("GoodsParameter").Filter("Id", pid).One(&goodsParameter, "Id", "Goods", "Parameter", "GoodsNowPrice")
		if err != nil {
			beego.Info(err)
			continue
		}
		gid := goodsParameter.Goods
		var goods models.Goods
		err = o.QueryTable("Goods").Filter("Id", gid).One(&goods, "Id", "Name")
		if err != nil {
			beego.Info(err)
			continue
		}
		goodsParameter.Goods = &goods
		var cart Cart
		cart.GoodsParameter = goodsParameter
		cart.Count = value
		//应该是后添加的商品在前面显示
		//下面这个可以实现在slice头部追加元素
		carts = append([]Cart{cart}, carts...)

	}
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	resp["code"] = 200
	resp["goods"] = carts
	this.Data["json"] = resp
	return
}

// @Title Put
// @Description Put cart
// @Param pid path true "pid for parameter id"
// @Param count path true "count for parameter"
// @Success  200 {string} ok
// @Failure 400 update失败
// @router /:pid/:count [put]
func (this *CartController) PutCartCount() {
	//更新参数和数量
	pid, err1 := this.GetInt(":pid")
	count, err2 := this.GetInt(":count")
	resp := make(map[string]interface{})
	defer this.ServeJSON()

	if err1 != nil || err2 != nil {
		resp["code"] = 1
		resp["errmsg"] = "请求数据不正确"
		this.Data["json"] = resp
		return
	}
	//userName := this.GetSession("userName")
	//o := orm.NewOrm()
	//var user models.User
	//user.Name = userName.(string)
	//o.Read(&user, "Name")

	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		resp["code"] = 2
		resp["errmsg"] = "redis数据库连接失败"
		this.Data["json"] = resp
		return
	}
	defer conn.Close()
	//_, err = conn.Do("hset", "cart_"+strconv.Itoa(user.Id), skuid, count)
	_, err = conn.Do("hset", "cart_1", pid, count)
	if err != nil {
		beego.Info(err)
		return
	}
	resp["code"] = 200
	this.Data["json"] = resp
	return
}

// @Title Put
// @Description Put cart
// @Param pid path true "pid for parameter id"
// @Success  200 {string} ok
// @Failure 400 update失败
// @router /:oldId/:newId/:count [put]
func (this *CartController) PutCartPara() {
	//更新参数和数量
	deleteId, err1 := this.GetInt(":oldId")
	count, _ := this.GetInt(":count")
	pid, err2 := this.GetInt(":newId")
	resp := make(map[string]interface{})
	defer this.ServeJSON()

	if err1 != nil || err2 != nil {
		resp["code"] = 1
		resp["errmsg"] = "请求数据不正确"
		this.Data["json"] = resp
		return
	}
	//userName := this.GetSession("userName")
	//o := orm.NewOrm()
	//var user models.User
	//user.Name = userName.(string)
	//o.Read(&user, "Name")

	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		resp["code"] = 2
		resp["errmsg"] = "redis数据库连接失败"
		this.Data["json"] = resp
		return
	}
	defer conn.Close()
	//删除原来所选的参数
	_, err = conn.Do("hdel", "cart_1", deleteId)
	if err != nil {
		beego.Info(err)
		return
	}
	//_, err = conn.Do("hset", "cart_"+strconv.Itoa(user.Id), skuid, count)
	_, err = conn.Do("hset", "cart_1", pid, count)
	if err != nil {
		beego.Info(err)
		return
	}
	resp["code"] = 200
	resp["succMsg"] = "商品参数更新成功"
	this.Data["json"] = resp
	return
}

// @Title Get
// @Description Get parameter
// @Param pid path true "pid for goodsParameter id"
// @Success  200 {string} 删除成功
// @Failure 400 lost data
// @router /:pid [delete]
func (this *CartController) DeleteCart() {
	//获取数据
	pid, err := this.GetInt(":pid")
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	//校验数据
	if err != nil {
		resp["code"] = 1
		resp["errmsg"] = "请求数据不正确"
		this.Data["json"] = resp
		return
	}
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	defer conn.Close()
	if err != nil {
		resp["code"] = 2
		resp["errmsg"] = "redis数据库连接失败"
		this.Data["json"] = resp
		return
	}
	//获取userID
	//o := orm.NewOrm()
	//var user models.User
	//userName := this.GetSession("userName")
	//user.Name = userName.(string)
	//o.Read(&user, "Name")

	//conn.Do("hdel", "cart_"+strconv.Itoa(user.Id), skuid)
	_, err = conn.Do("hdel", "cart_1", pid)
	if err != nil {
		beego.Info(err)
		return
	}
	//返回数据
	resp["code"] = 200
	resp["errmsg"] = "删除成功"
	this.Data["json"] = resp
	return
}

// @Title Get
// @Description Get parameter
// @Param gid path true "gid for goods id"
// @Success  200 {object} models.GoodsParameter
// @Failure 400 lost data
// @router /:gid [get]
func (this *CartController) GetParameters() {
	gid, err := this.GetInt(":gid")
	if err != nil {
		beego.Info("参数传递错误")
		return
	}
	var parameters []*models.GoodsParameter
	o := orm.NewOrm()
	_, err = o.QueryTable("GoodsParameter").Filter("Goods__Id",gid).All(&parameters,"Id","parameter","GoodsTruePrice","GoodsNowPrice")
	if err != nil {
		beego.Info("查询失败, err: ", err)
		return
	}
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	resp["code"] = 200
	resp["parameters"] = parameters
	this.Data["json"] = resp
	return
}
