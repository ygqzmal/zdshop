@[TOC](beego+orm 实战简易小商城项目)
# 前言

<font color=#999AAA >
最近用beego做了个非常简单的商城小程序，在这里详细总结一下，主要记录model层和controller层的增删改查，以及用到一些golang的基础知识。有错的地方还望指正。
go版本1.13/ beego版本1.10 
</font>

<hr style=" border:solid; width:100px; height:1px;" color=#000000 size=1">




# 一、model层

beego一般会配合着orm这类的数据库关系映射一起使用，主要有gorm和xorm，也可以使用在其他框架内如gin框架。这里使用的orm，其实差不多触类旁通吧。


<font color=#999AAA >首先model层可以创建一个model.go文件也可以根据数据表来创建多个go文件，因为项目不算大，所以这里就只用一个model.go文件包含所有结构体。

## 1. init函数连接数据库
在model文件中需要init函数来初始化需要连接的数据库以及创建表，项目会自动创建表，我们只需要先创建好数据库就行了。这里使用的是mysql
```c
func init() {
	// set default database 首先需要注册一个默认的数据库，这里只用到一个数据库。
	dns := "root:123456@tcp(127.0.0.1:3306)/zdshop?charset=utf8mb4&parseTime=True"
	_ = orm.RegisterDataBase("default", "mysql", dns)

	// register model 这里仅仅是因为一行太长分开几行而已，这里需要将结构体转为数据表
	orm.RegisterModel(new(User), new(Logistics), new(Information), new(Admin), new(Operation), new(Territory))
	orm.RegisterModel(new(BigDistribution), new(HoldLocal), new(Distribution), new(Address), new(Goods))
	orm.RegisterModel(new(GoodsCategory), new(GoodsParameter), new(ShopCart), new(GoodsBanner), new(OrderInfo), new(OrderGoods))

	// create table
	//当想修改表结构例如添加字段，bee run 启动项目的时候不会进行修改。
	//需要第二个参数改为true,则每次加载都会重置数据库(原有数据会丢失！！)。
	_ = orm.RunSyncdb("default", false, true)
}
```
需要注意charset后面的类型不能错，否则会报错，mysql默认的数据库类型是utf8。
同时不要忘记import
```c
import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)
```

## 2.  beego在model中用结构体来映射数据表
<font color=#999AAA >代码如下（示例）：
```c
//商品表
type Goods struct {
	Id             int
	Name           string            `orm:"size(50);unique"` //商品名称
	GoodsBrief     string            `orm:"type(text)"`      //商品简介 数据库类型应该为blob
	GoodsState     int               `orm:"default(0)"`      //商品状态 0-上架 1-下架
	Explain        string            `orm:"type(30)"`        //说明
	CreateTime     time.Time         `orm:"auto_now_add"`    //商品录入时间
	UpdateTime     time.Time         `orm:"auto_now"`        //最后修改时间
	SalesValue     int               `orm:"default(0)"`      //销量
	Category       *GoodsCategory    `orm:"rel(fk)"`         //分类id 级联删除
	ShopCart       []*ShopCart       `orm:"reverse(many)"`
	GoodsBanner    []*GoodsBanner    `orm:"reverse(many)"`
	GoodsParameter []*GoodsParameter `orm:"reverse(many)"`
}
```
通过自定义struct，使用Tag标签来确定数据类型，下面列举一些常用的 
orm     | sql
-------- | -----
size | varchar
unique  | 唯一
type()  | 类型(text,char)
default() | 默认
rel(fk) | 外键
reverse(many) | 外键多对多反向关系
auto_now_add | 创建时间
auto_now | 更新时间
1. 时间类型为time.time，意思为一个时间点，在插入数据时这两个字段会自动添加
2. 我在实践中发现这两个其实没有什么区别，更新时还是需要指定字段进行更新
3. 关于外键，商品有个字段是分类ID，使用rel(fk)，指向分类struct。beego会自动给商品表建一个分类ID字段
4. 指定了外键关系，删除时会自动关联删除
5. **需要注意的是struct中不允许出现__双下划线来命名**

## 3. 设置表引擎，因为使用到事务所以需要用到InnoDB
beego创建表类型默认为MyISAM，但在实际编程中需要使用到事务操作，myisam不支持事务操作，所以需要修改表引擎为InnoDB

```c
//这个可能beego版本不同，建表默认引擎不同，没有考究过
func (g *Goods) TableEngine() string {
	return "INNODB"
}
```
再次查看，表引擎已经更改
![](https://img-blog.csdnimg.cn/2021012816152371.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3dlaXhpbl80NTE2NjUxMQ==,size_16,color_FFFFFF,t_70)
至此，model层就算完成了
# 二、controller层
## 1. 获取各种数据
<font color=#999AAA >代码如下（示例）：
```c
//获取formData数据 (表单传参)
name := this.GetString("name")
id, err := this.GetInt("ID") //返回两个参数，如果接收不到ID则会返回err
this.GetFile("Img") //接收文件
this.GetStrings("id") //接收数组
//获取path数据 (url传参)
id, err := this.Int(":id") //如果没有接收到id，并不会err报错
```

## 2. 查询

<font color=#999AAA >一般查询代码如下（示例）：
```c
o := orm.NewOrm()
var goods models.Goods
goods.Id = Id
//根据Id查询所有的数据
o.read(&goods) 	
//如果要根据其他字段查询，例如根据name字段查询所有字段
goods.Name = name
o.read(&goods, "Name")
```
但是这种方法无法指定字段进行查询，也无法进行排序、过滤等操作。beego提供高级查询(增删查改)

<font color=#999AAA >高级查询代码如下（示例）：
```c
o := orm.NewOrm()
var good models.Goods //一条数据
var goods []*models.Goods //多条数据
//一条数据
o.QueryTalbe("Goods").需要的条件.One(&good, "需要查询指定的字段1"，"需要查询指定的字段2")
//多条数据
o.QueryTalbe("Goods").需要的条件.All(&goods, "需要查询指定的字段1"，"需要查询指定的字段2")
```
下面列举几个常用的条件
.条件     | sql
-------- | -----
Filter() | where
orderBy()  |  排序
Limit () |  Limit
RelatedSel() |  关联查询
1. 如果想根据分类id查询所有该分类下的商品，（一对多），Filter("表名__Id")
```c
o.QueryTable("GoodsParameter").Filter("Goods__Id", id).All(&parameters)
```
2. 查询goods表Id为id的值, 想要查关联的外键的话,要加上relatedSel("model中的外键字段"), 会把外键所有数据查询出来而且无法指定

##  3. 增（添加）
<font color=#999AAA >代码如下（示例）：
```c
o := orm.NewOrm()
var goods models.goods
goods.Id = 1
goods.Name = "huawei"
//添加外键关联商品分类id, 不能直接goods.Category = 1, 会报错
var category models.GoodsCategory
category.Id = categoryId
goods.Category = &category 
o.Insert(&goods)
```

##  4. 改（更新）
方法一： <font color=#999AAA >代码如下（示例）：
```c
o := orm.NewOrm()
var goods models.Goods
goods.Id = id
goods.Name = name
goods.GoodsBrief = brief
//这里虽然在model中time的tag不一样，但是没办法做到自动更新时间，需要指定updateTime
_, err = o.Update(&goods, "Name", "GoodsBrief", "UpdateTime")
//下面这个更新如果没有指定字段的话，会设置为null，但是model中建表时默认不为空，所以无法这样更新
_, err = o.Update(&goods)
```
方法二：(推荐这个）<font color=#999AAA >代码如下（示例）：
```c
o.QueryTalbe("表名").Filter("id", 1).Update(orm.Params{
    "name": "",
    "GoodsBrief": "",
})
```

##  5. 删（删除）
<font color=#999AAA >代码如下（示例）：
```c
_, err := o.QueryTable("表名").Filter("id", 1).Delete()
```

##  6.  补充
**下面这个我找到一个queryTable接口博客，比我的详细，**
链接: [queryTable接口](https://www.cnblogs.com/hei-ma/articles/13716916.html).

1. queryTalbe接口有很多方法，例如统计条数Count(), 判断是否存在Exist()
2. **在相对时间进行排序时，正序：orderBy("createTime") , 倒序：orderBy("-createTime"),**

#  三. 项目完善
## 1. 使用swagger,
swagger是一个go语言和框架结合使用的自动化API文档。通过注释来完成api文档的编辑，相当好用。
还实现了restful api 接口规范。
<font color=#999AAA >代码如下（示例）：
```c
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
func (this *GoodsController) Post() {...}
```
启动项目，访问http://127.0.0.1:8080/swagger/  效果如下：
![在这里插入图片描述](https://img-blog.csdnimg.cn/20210128180447190.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3dlaXhpbl80NTE2NjUxMQ==,size_16,color_FFFFFF,t_70)

参考链接: [swagger文档](https://www.bookstack.cn/read/beego/advantage-docs.md)


## 2. 日志收集
这里在网上找到了一个beego的日志收集，出处找不到了，所以直接将源代码放出来。
```c
package utils

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
)

//说明：日志记录
//在conf下app.conf中进行配置, 只保存最近10个文件，会不停覆盖。最新记录在app.log中
//[log]
//log_level = 7
//log_path = logs/app.log
//maxlines = 100 表示每个文件最大行数
//maxsize = 1024 表示每个文件最大大小

/*
在main函数中进行调用
func init() {
	//日志写入
	_ = utils.InitLogger()
}
 */
// RFC5424 log message levels.
const (
	LevelEmergency = iota
	LevelAlert        //1
	LevelCritical	  //2
	LevelError		  //3
	LevelWarning	  //4
	LevelNotice		  //5
	LevelInformational//6
	LevelDebug		  //7
)

func InitLogger() (err error) {
	BConfig, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil{
		fmt.Println("config init error:", err)
		return
	}
	maxlines, lerr := BConfig.Int64("log::maxlines")
	if lerr != nil {
		maxlines = 1000
	}

	logConf := make(map[string]interface{})
	logConf["filename"] = BConfig.String("log::log_path")
	level,_ := BConfig.Int("log::log_level")
	logConf["level"] = level
	logConf["maxlines"] = maxlines

	confStr, err := json.Marshal(logConf)
	if err != nil {
		fmt.Println("marshal failed,err:", err)
		return
	}
	//设置输出到文件
	//beego.SetLogger("file", `{"filename":"logs/test.log"}`)
	beego.SetLogger(logs.AdapterFile, string(confStr))
	//输是否出文件名和行号
	beego.SetLogFuncCall(true)
	//设置级别设置日志级别
	//beego.SetLevel(beego.LevelInformational)
	return
}
```
创建一个log文件夹，当运行时，所有的请求和beego.Info都会写到log文件夹下的日志当中

## 3. 文件上传
自己写的一个公共方法
```c
//调用
/*
	//第二个参数为formData的key，第三个参数为限制图片上传的最大大小， 200000为2M
	errMsg, url := utils.HandleFile(&this.Controller, "img", 200000)
	if errMsg != "" {
		beego.Info(errMsg)
		return
	}
 */
//文件上传
func HandleFile(this *beego.Controller, image string, size int64) (error, url string) {
	//headers, err := this.GetFiles(image)
	file, header, err := this.GetFile(image)
	defer file.Close()
	if err != nil {
		return "图片上传失败", ""
	}
	//大小
	if header.Size > size {
		return "图片太大无法上传", ""
	}
	//格式
	ext := path.Ext(header.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		beego.Info("图片格式不正确")
		return "图片格式不正确", ""
	}
	//防重名
	fileName := strconv.Itoa(int(time.Now().UnixNano())) + ext
	//存储 第一个参数要和GetFile的一样
	err = this.SaveToFile(image, "./upload/img/"+fileName)
	if err != nil {
		return "文件存储错误 ", ""
	}
	return "", "./upload/img/"+fileName
}
```

#  四. 其他
下面说说我在写项目时遇到一些问题以及解决方法。
## 1.  把元素追加到slice头部
```c
var carts []Cart
for key, value := range goodsParaMap {
	var cart model.Cart
	...
	//应该是后添加的商品在前面显示
	//下面这个可以实现在slice头部追加元素
	carts = append([]Cart{cart}, carts...)
}
```
当想对一个slice追加一个slice的时候，第二个参数即要追加的slice，用...来展开追加

## 2. 接收和返回json数据
接收：
```c
data := this.GetString("parameters")
beego.Info("json: ", data)
var parameters []JsonPara
err = json.Unmarshal([]byte(data), &parameters)
if err != nil {
	beego.Info("json.Unmarshal err: ", err)
	beego.Info("传入数据有误") //不能为空且不能类型不对，否则无法unmarshal
	return
}
```
返回：

```c
resp := make(map[string]interface{})
defer this.ServeJSON()

resp["key1"] = value1
resp["key2"] = value2
this.Data["json"] = resp
```

## 3. orm开启事务
```c
o := orm.NewOrm()
//事务开始
o.Begin()
//事务回滚
o.Rollback()
//事务提交
o.Commit()

```

## 4. 时间类型格式化输出
如果不进行格式化输出，是没办法直接获取到下面这种时间类型的
```c
//返回数据库格式化时间：2020-12-16 21:24:35
createTime := goods.CreateTime.Format("2006-01-02 15:04:05")
```
## 5. 看到一些好的博客链接
beego其实还有一些功能但我没有用上，例如可以在controller中定义表单对应的struct和tag来进行自动校验。
在写项目的时候也遇上很多问题，查了很多资料，下面给几个我经常看的链接，都挺好的。

参考链接:
[beego脱坑](https://blog.csdn.net/yang731227/category_8058880.html)
[queryTable接口](https://www.cnblogs.com/hei-ma/articles/13716916.html).
[git使用](https://blog.csdn.net/weixin_42490398/article/details/90212418)   这是我看过git教学最好的一篇博客
[golang基础教程](https://www.liwenzhou.com/posts/Go/go_menu/)  这是李文周老师的博客 我觉得很好， 推荐给大家



## 6. 补充
beego的路由很简单，没有什么好说的。如果使用swagger创建beego项目还有自带的示例，相当友好
启动项目：
```c
bee run
```
有时候会报错，可以用这个再次生成文档
```c
bee run -gendoc=true -downdoc=true
```
# 总结
**这是第一次使用beego写项目，因此想对此总结一下，beego该有的都有，github上也有beego的支付宝微信支付源码，一样可以连接redis；不要看现在做go语言开发的少，go在服务器开发也很有优势，给我最大的感觉是效率高。之前一直使用php的TP5来写服务端，用完beego感觉会更加简单，代码更加清晰，相对代码行数也少了很多。个人的话更加倾向于beego。beego其实也可以做更大的项目，业务逻辑复杂还可以自己添加逻辑层以及dao层。beego当然也有缺陷，由于技术相对较新，一些问题网上搜了很久也没有答案；希望golang越来越棒吧。**

