package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
	"zdshop/models"

	//"github.com/astaxie/beego/orm"
	//"github.com/gomodule/redigo/redis"
	//"strconv"
	//"zdshop/models"
)

type OrderController struct {
	beego.Controller
}

type ReOrderInfo struct {
	OrderInfo  *models.OrderInfo
	CreateTime string
}

type ReOrderGood struct {
	OrderGood *models.OrderGoods
	Parameter string
	GoodName  string
}

//type ReOrderInfo struct {
//	OrderInfo models.OrderInfo
//	OrderGoods []*models.OrderGoods
//}

// @Title Get
// @Description add cart
// @Param state path true "订单状态"
// @Param uid path true "用户id"
// @Success  200 {string} 订单获取成功
// @Failure 400 订单获取失败
// @router /:uid/:state [get]
func (this *OrderController) GetPayOrder() {
	//0-未支付 1-已支付
	state, _ := this.GetInt(":state")
	uid, _ := this.GetInt(":uid")
	//redis存的是参数id,不是商品id
	o := orm.NewOrm()
	conn, _ := redis.Dial("tcp", "127.0.0.1:6379")
	defer conn.Close()
	//获取用户数据
	//var user models.User
	//userName := this.GetSession("userName")
	//user.Name = userName.(string)
	//o.Read(&user, "Name")
	var orderInfo []models.OrderInfo
	_, err := o.QueryTable("OrderInfo").Filter("Distribution__Id", uid).Filter("OrderStatus", state).OrderBy("-CreateTime").All(&orderInfo)
	if err != nil {
		beego.Info(err)
		return
	}
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	resp["orderInfos"] = orderInfo
	resp["code"] = 200
	this.Data["json"] = resp
	return
	//如果查询订单不能把一对多关联的也查出来就自己查
}

// @Title Get
// @Description add cart
// @Param uid path true "用户id"
// @Success  200 {string} 订单添加成功
// @Failure 400 订单添加失败
// @router / [get]
func (this *OrderController) GetAllOrder() {
	var orderInfos []*models.OrderInfo
	o := orm.NewOrm()
	_, err := o.QueryTable("OrderInfo").OrderBy("-CreateTime").All(&orderInfos, "Id", "OrderId", "TotalPrice", "OrderStatus", "CreateTime")
	if err != nil {
		beego.Info(err)
		return
	}
	beego.Info(orderInfos)
	var orderInfo []ReOrderInfo
	for _, value := range orderInfos {
		var temp ReOrderInfo
		temp.OrderInfo = value
		temp.CreateTime = value.CreateTime.Format("2006-01-02 15:04:05")

		orderInfo = append(orderInfo, temp)

	}

	resp := make(map[string]interface{})
	defer this.ServeJSON()

	resp["code"] = 200
	resp["orderInfos"] = orderInfo
	this.Data["json"] = resp
	return
}

// @Title Put
// @Description add cart
// @Param oid path true "orderId"
// @Success  200 {string} 订单添加成功
// @Failure 400 订单添加失败
// @router /:uid [get]
func (this *OrderController) PutOrderState() {
	//后台发货后添加订单号, 订单状态改变
}

// @Title Put
// @Description add cart
// @Param orderId   formData true "orderId"
// @Param addressId formData true "addressId"
// @Success  200 {string} 订单添加成功
// @Failure 400 订单添加失败
// @router / [put]
func (this *OrderController) PutOrderAddr() {
	//修改订单地址, 前提是还没发货
	aid, err1 := this.GetInt("addressId")
	oid, err2 := this.GetInt("orderId")
	if err1 != nil || err2 != nil {
		beego.Info("参数传递错误")
		return
	}
	o := orm.NewOrm()
	_, err := o.QueryTable("OrderInfo").Filter("Id", oid).Update(orm.Params{"Address": aid})
	if err != nil {
		beego.Info(err)
		return
	}

	resp := make(map[string]interface{})
	defer this.ServeJSON()

	resp["code"] = 200
	resp["msg"] = "订单地址修改成功"
	this.Data["json"] = resp
	return
}

// @Title Put
// @Description add cart
// @Param oid path true "orderId"
// @Success  200 {string} 订单添加成功
// @Failure 400 订单添加失败
// @router /:uid [get]
func (this *OrderController) AcceptOrder() {
	//确认收货, 修改订单状态-5-已完成
}

// @Title Post
// @Description add cart
// @Param pid formData true "参数id"
// @Param count formData true "添加数量"
// @Success  200 {string} 订单添加成功
// @Failure 400 订单添加失败
// @router / [post]
func (this *OrderController) PostOrder() {
	/*
		这里两点说明一下：
		1. 并发问题：当两个用户同时购买同一件商品时，两边同时提交订单，商品数量相加可能超过库存量
		解决方法：先查询数据库，将库存量取出来，在真正执行进行插入数据代码的时候进行过滤判断
		2. 用户多了，一直无法购买成功。
		解决方法：设置一个循环次数，将解决方法一放入循环体中，进行循环判断。
	*/

	//获取数据
	//addrId, _ := this.GetInt("addrId")
	pIds := this.GetStrings("pIds")
	//pIds会传来如：[1, 3, 4]

	totalCount, _ := this.GetInt("totalCount")   //
	totalPrice, _ := this.GetFloat("totalPrice") //总价

	resp := make(map[string]interface{})
	defer this.ServeJSON()
	//校验数据
	if len(pIds) == 0 {
		resp["code"] = 400
		resp["msg"] = "参数获取失败"
		this.Data["json"] = resp
		return
	}
	//处理数据
	//向订单表中插入数据
	o := orm.NewOrm()

	o.Begin() //标识事务的开始

	//userName := this.GetSession("userName")
	//var user models.User
	//user.Name = userName.(string)
	//o.Read(&user, "Name")

	var order models.OrderInfo
	//order.OrderId = time.Now().Format("200601021504") + strconv.Itoa(user.Id)
	order.OrderId = time.Now().Format("200601021504") + "1"
	//order.User = &user//分销商id
	order.OrderStatus = 0
	order.TotalCount = totalCount
	order.TotalPrice = totalPrice

	//查询地址
	var addr models.Address
	//addr.Id = addrId
	//o.Read(&addr)
	order.Address = &addr

	//分销商
	var temp models.Distribution
	//temp.Id = 1
	order.Distribution = &temp
	//执行插入操作
	_, err := o.Insert(&order)
	if err != nil {
		beego.Info(err)
		o.Rollback()
		return
	}

	//向订单商品表中插入数据
	conn, _ := redis.Dial("tcp", "127.0.0.1:6379")

	for _, pid := range pIds {
		id, _ := strconv.Atoi(pid)

		var goodsParameter models.GoodsParameter
		goodsParameter.Id = id
		i := 3
		for i > 0 {
			err := o.Read(&goodsParameter)
			if err != nil {
				beego.Info(err)
				return
			}
			var orderGoods models.OrderGoods
			orderGoods.GoodsParameter = &goodsParameter
			orderGoods.OrderInfo = &order

			//count, _ := redis.Int(conn.Do("hget", "cart_"+strconv.Itoa(user.Id), id))
			count, _ := redis.Int(conn.Do("hget", "cart_1", id))
			if count > goodsParameter.GoodsNumber {
				resp["code"] = 2
				resp["msg"] = "商品库存不足"
				this.Data["json"] = resp
				o.Rollback() //标识事务的回滚
				return
			}

			preCount := goodsParameter.GoodsNumber

			//time.Sleep(time.Second * 5)
			//beego.Info(preCount, user.Id)

			orderGoods.Count = count
			orderGoods.Price = float64(count) * goodsParameter.GoodsNowPrice

			_, err = o.Insert(&orderGoods)
			if err != nil {
				beego.Info(err)
				o.Rollback()
				return
			}

			goodsParameter.GoodsNumber -= count
			goodsParameter.SalesValue += count

			updateCount, _ := o.QueryTable("GoodsParameter").Filter("Id", goodsParameter.Id).Filter("GoodsNumber", preCount).Update(orm.Params{"GoodsNumber": goodsParameter.GoodsNumber, "SalesValue": goodsParameter.SalesValue})
			//如果更新条数为0, 说明有人同时在购买, 则不进行更新操作
			//因为第一个用户在下单时候已经对商品库存进行消耗, 但此时操作尚未完全结束, 因此在更新库存之前, 如果有人进行购买那么库存就会出错
			if updateCount == 0 {
				if i > 0 {
					i -= 1
					continue
				}
				resp["code"] = 3
				resp["msg"] = "商品库存改变,订单提交失败"
				this.Data["json"] = resp
				o.Rollback() //标识事务的回滚
				return
			} else {
				//conn.Do("hdel", "cart_"+strconv.Itoa(user.Id), goods.Id)
				conn.Do("hdel", "cart_4", pid)
				break
			}
		}
	}
	//返回数据
	o.Commit() //提交事务
	resp["code"] = 5
	resp["msg"] = "ok"
	this.Data["json"] = resp
}

// @Title Put
// @Description add cart
// @Param oid path true "orderId"
// @Success  200 {string} 订单添加成功
// @Failure 400 订单添加失败
// @router /:oid [delete]
func (this *OrderController) DeleteOrder() {
	//取消订单, 在订单未发货之前
	//使用:id 方法接收path参数，传不传都没有err返回
	id, _ := this.GetInt(":oid")
	o := orm.NewOrm()
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	var orderInfo models.OrderInfo
	err := o.QueryTable("orderInfo").Filter("Id", id).One(&orderInfo, "Id", "OrderStatus")
	if err != nil {
		beego.Info(err)
		return
	}
	if orderInfo.OrderStatus != 0 && orderInfo.OrderStatus != 1 {
		resp["code"] = 400
		resp["msg"] = "订单无法取消"
		this.Data["json"] = resp
		return

	}
	//_, err = o.QueryTable("orderInfo").Filter("Id", id).Delete()
	//if err != nil {
	//	beego.Info(err)
	//	return
	//}
	_, err = o.Delete(&orderInfo)
	if err != nil {
		beego.Info(err)
		return
	}
	resp["code"] = 200
	resp["msg"] = "订单取消成功"
	this.Data["json"] = resp
	return
}

// @Title Get
// @Description add cart
// @Param oid path true "orderId"
// @Success  200 {string} 订单查询成功
// @Failure 400 订单查询失败
// @router /:oid [get]
func (this *OrderController) GetOneOrderInfo() {

	oid, _ := this.GetInt(":oid")

	var orderInfo models.OrderInfo
	var orderGoods []*models.OrderGoods
	o := orm.NewOrm()
	err := o.QueryTable("OrderInfo").Filter("Id", oid).One(&orderInfo, "Id", "OrderId", "TotalPrice", "OrderStatus", "CreateTime")
	_, err = o.QueryTable("OrderGoods").Filter("OrderInfo__Id", oid).All(&orderGoods)
	if err != nil {
		beego.Info(err)
		return
	}
	var reOrderGoods []ReOrderGood
	for _, value := range orderGoods {
		var temp ReOrderGood
		var parameter models.GoodsParameter
		err := o.QueryTable("GoodsParameter").Filter("Id", value.GoodsParameter.Id).One(&parameter, "Goods", "Parameter")
		if err != nil {
			beego.Info(err)
			return
		}
		var good models.Goods
		err = o.QueryTable("Goods").Filter("Id", parameter.Goods.Id).One(&good, "Name")
		if err != nil {
			beego.Info(err)
			return
		}
		temp.OrderGood = value
		temp.Parameter = parameter.Parameter
		temp.GoodName = good.Name

		reOrderGoods = append(reOrderGoods, temp)
	}

	resp := make(map[string]interface{})
	defer this.ServeJSON()

	resp["code"] = 200
	resp["orderInfo"] = orderInfo
	resp["orderGoods"] = reOrderGoods
	resp["createTime"] = orderInfo.CreateTime.Format("2006-01-02 15:04:05")
	this.Data["json"] = resp
	return

}