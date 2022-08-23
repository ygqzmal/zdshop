package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/smartwalle/alipay"
)

var AppId = "2021000119671615"

var PrivateKey = "MIIEpgIBAAKCAQEAwVzyS9U2ZojqFH6fya1KRQMg6TRvAg9bvtGgOQF/0QyuAlFBGbAv7fGJwOFAqdwI6UwO2oKKSeR88wr7h0ZBbAIcRuJbB5EiQDQxuRqLwWr0DUIQcLwUchFdYM+isPzxvxrha/2ZK8PeaYpYqIDS82oTYBWOeQl6B/kKmiYo3LUvnQWLHZ9ASMR33AUcXB7UrEBKcFYSWlwz+Yb54QbCtDOn3aaiMxADuDrJRyKIV1weNt9zT/6By4tNFH43nTE8g0V5da/eU6nxaH9fZAcQbSJI6N2ptmdwpvG/PrDaXw8N0S+d88OKFFB2olEOVlBqprpd9nI0wWzoWM+eA8B6HQIDAQABAoIBAQCC9tUVD5/86pbAstK+4iP/ghL0YJMfLg/RumFuutk7Wf8xip8xKugLnSDUBrypT2KpwN3+mZPYYG1QoAukn60H3pYQXJeTFHXCTpeu64/kvO/3WtuPv5OJfsFkJL3oainCag5r+iOYRN2LViyeVEKMozfkSDVrPUPtynun1FiwwskjVtVhhXxIkCPAJ3/v2iN+kDMcU6cM6qDk/wFAiznItKurNRYfMbjIIZxailXi4ENRcaeY5oNaR0io//zNmnpNLxz7cTGDIbFJTXq8aqvAiJoCaJ5j4EhzpFE4e5MNl3UPAGfkDUtpZm0LEvNyIeivSXnE2swBhxdkhFqxE80BAoGBAO3MFjD2mpovks/Ocq3+gywXIyAWndmydGjVueJvbqR3HxNS9xE2zl8oePQqkVnrRD8vo3sDIob9uKwgrcisnl7YbU9AQsdfP/D14ONgedFT/vcHR5s0OGE+oaEEmNBQqt+2vlL5XeeL4xA8+wgHdHfLlBiltReZIl+XZQ4Zv4LxAoGBANAqHvxxYIX/SOKrn8O/tbUgzN/nHlRhw906dK1T2unTVvLNuIhbGF86RYhEUrUMNHEAPRgk6gcGskqr7TC8WVTkNab5zXHSpSwN86fzL//YxchBZp8gigVB5ezUuQrnLp73mXaB2mRGwvMt2rpPwEW2zYpHdiu8PWEf7sD5ZlHtAoGBALPtXE0oCsrnfDpohzVIApB14UoCUuXJtXMtZD0E+77Ns2G3wOHxii1OPlbhbqGO9lCpBxWoxZNGn1j+UQAqPJqfP/ZbNSwN0h/Mq6Df+sx8tcrMY034MUDDVyCyjb3xi5lCeLfnnzn4CpLa0Ua9/U43Z5NOrrtwTyXtM7V7ngDBAoGBAJKA0oYAlVo5LOa6uxpdVlk/2HDMjD/+/oY4md1S4wMlxk/kETeGRUTgEgexbjQVfuL4tAbGFB8Vy21aSvi91nE0m74EmV6+TZkPyKgvM1zxB2HFBaCAmiLRGizwGbtesSUYRV1uTnG8i3/yiboOXtexrD7hxH1LYjd07efKnwSRAoGBAIe53cahkKavpuluyaFfgUafGw7HKgDlcDJMVd815w0evmFzLsJ0M4OQZg8OC6xF5ZnMn1n/Xqiio5ntpGKuzaU8Q1YSXij2DCIoWwAnqdOhIcPlOAxfFJ3Y4LFKDoTAvEb5XvRXOEhNTeKtsQU53+U9wZ5d2N48qZHK5HZc7ZE5"

var AliPublicKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwVzyS9U2ZojqFH6fya1KRQMg6TRvAg9bvtGgOQF/0QyuAlFBGbAv7fGJwOFAqdwI6UwO2oKKSeR88wr7h0ZBbAIcRuJbB5EiQDQxuRqLwWr0DUIQcLwUchFdYM+isPzxvxrha/2ZK8PeaYpYqIDS82oTYBWOeQl6B/kKmiYo3LUvnQWLHZ9ASMR33AUcXB7UrEBKcFYSWlwz+Yb54QbCtDOn3aaiMxADuDrJRyKIV1weNt9zT/6By4tNFH43nTE8g0V5da/eU6nxaH9fZAcQbSJI6N2ptmdwpvG/PrDaXw8N0S+d88OKFFB2olEOVlBqprpd9nI0wWzoWM+eA8B6HQIDAQAB"

// @Title handle pay
// @Description pay
// @Param uid path false "userId"
// @Success  200 {string} pay success
// @Failure 400 pay false
// @router /pay [post]
//处理支付
func (this *OrderController) HandlePay() {
	fmt.Println("支付...")
	//var client, _ = alipay.New(appId, aliPublicKey, privateKey, false)
	var client  = alipay.New(AppId, AliPublicKey, PrivateKey,false)
	//获取数据
	orderId := this.GetString("orderId")
	totalPrice := this.GetString("totalPrice")

	//orderId := "2022051315431"
	//totalPrice := "1000"

	var p = alipay.AliPayTradePagePay{}
	//alipay.NewRequest()
	p.NotifyURL = "http://xxx"
	//p.ReturnURL = "http://127.0.0.1:8080/v1/payok"
	p.Subject = "电子购物平台"
	p.OutTradeNo = orderId
	p.TotalAmount = totalPrice
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	var url, err = client.TradePagePay(p)
	if err != nil {
		fmt.Println(err)
	}

	var payURL = url.String()
	fmt.Println("payUrl: ", payURL)

	resp := make(map[string]interface{})
	defer this.ServeJSON()
	resp["code"] = 200
	resp["msg"] = ""
	resp["payUrl"] = payURL
	this.Data["json"] = resp
	return
	//this.Redirect(payURL, 302)
}

//支付成功
func (this *OrderController) PayOk() {
	//获取数据
	//out_trade_no=999998888777
	orderId := this.GetString("out_trade_no")

	//校验数据
	if orderId == "" {
		beego.Info("支付返回数据错误")
		this.Redirect("/user/userCenterOrder", 302)
		return
	}

	//操作数据

	o := orm.NewOrm()
	count, _ := o.QueryTable("OrderInfo").Filter("OrderId", orderId).Update(orm.Params{"Orderstatus": 2})
	if count == 0 {
		beego.Info("更新数据失败")
		this.Redirect("/user/userCenterOrder", 302)
		return
	}

	//返回视图
	this.Redirect("/user/userCenterOrder", 302)

}