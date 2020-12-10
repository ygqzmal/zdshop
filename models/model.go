package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//用户表
type User struct {
	Id              int
	Name            string             `orm:"size(15);unique"`               //用户名
	Password        string             `orm:"size(32)"`                      //登陆密码
	Tele            string             `orm:"size(11);unique"`               //手机号码
	State           string             `orm:"size(1);type(char);default(1)"` //用户状态 0-正常 1-注销 2-停机
	RegisterTime    time.Time          `orm:"auto_now_add"`                  //注册时间, 第一次保存才设置时间
	LoginNum        int                `orm:"default(0)"`                    //登录次数
	LastLoginIp     string             `orm:"size(20);default(\"未登录过\")"`    //最后一次登录ip
	LastLoginTime   time.Time          `orm:"auto_now;type(datetime)"`       //最后一次登录时间,每次 model 保存时都会对时间自动更新
	Logistics       *Logistics         `orm:"reverse(one)"`
	Information     []*Information     `orm:"reverse(many)"`
	Admin           []*Admin           `orm:"reverse(many)"`
	BigDistribution []*BigDistribution `orm:"reverse(many)"`
	Distribution    []*Distribution    `orm:"reverse(many)"`
}

//资金流动表
type Logistics struct {
	Id          int
	User        *User     `orm:"rel(one)"`                      //用户Id
	Time        time.Time `orm:"auto_now;type(datetime)"`       //流动时间
	InFlow      string    `orm:"size(1);type(char);default(0)"` //资金流动判断 0-进账 1-出账
	AMoney      float64   `orm:"default(0)"`                    //流动金额
	Description string    `orm:"size(50);default(\"未知\")"`      //流动描述
}

//信息表
type Information struct {
	Id      int
	User1   *User     `orm:"rel(fk)"`                       //发信人id-用户表的外键
	User2   *User     `orm:"rel(fk)"`                       //收信人id-用户表的外键
	Time    time.Time `orm:"auto_now_add"`                  //发信时间
	Content string    `orm:"size(100)"`                     //发信内容
	Kind    string    `orm:"size(1);type(char);default(0)"` //0-普通的聊天消息 1-订单状态该别消息 2-通知消息
	State   string    `orm:"size(1);type(char);default(0)"` //0-未读 1-已读 2-删除
}

//管理员表
type Admin struct {
	Id                   int          //管理员Id
	User                 *User        `orm:"rel(fk);unique"`                //用户Id-用户表外键
	Role                 string       `orm:"size(1);type(char);default(1)"` //分总管理员		0-总管理员 1-分管理员
	JobPosition          string       `orm:"size(20);default(\"普通员工\")"`    //管理员职位备注
	GoodsAuthority       string       `orm:"size(1);type(char);default(0)"` //商品权限	  	0-没有权限 1-可查看权限 2-可操作商品数据权限
	OrderAuthority       string       `orm:"size(1);type(char);default(0)"` //订单权限   	0-没有权限 1-可查看权限 2-可操作商品数据权限
	DistributorAuthority string       `orm:"size(1);type(char);default(0)"` //分销权限   	0-没有权限 1-可查看权限 2-可操作商品数据权限
	CapitalAuthority     string       `orm:"size(1);type(char);default(0)"` //资金权限		0-没有权限 1-可查看权限
	MessageAuthority     string       `orm:"size(1);type(char);default(0)"` //消息权限		0-没有权限 1-可查看权限
	AccordAuthority      string       `orm:"size(1);type(char);default(0)"` //协议管理		0-没有权限 1-可查看权限 2-可操作商品数据权限
	Operation            []*Operation `orm:"reverse(many)"`
}

//管理员操作表
type Operation struct {
	Id      int
	Admin   *Admin    `orm:"rel(fk)"`   //管理员-管理员表外键
	Time    time.Time `orm:"auto_now"`  //操作时间
	Content string    `orm:"size(100)"` //操作内容
}

//领域表
type Territory struct {
	Id              int
	Name            string             `orm:"size(20);unique"` //领域名称
	Introduce       string             `orm:"type(text)"`      //领域介绍
	BigDistribution []*BigDistribution `orm:"reverse(many)"`
	Distribution    []*Distribution    `orm:"reverse(many)"`
}

//总经销商表
type BigDistribution struct {
	Id        int
	User      *User      `orm:"rel(fk);unique"`     //用户Id-用户表外键
	AdTer     *Territory `orm:"rel(fk)"`            //管辖领域-领域表外键
	AdGrade   string     `orm:"size(1);type(char)"` //管辖等级 0-国级 1-省级 2-市级 3-县级
	AsLocal   *HoldLocal `orm:"rel(fk)"`            //管辖地-管辖地表外键
	Money     float64    `orm:"default(0)"`         //资金总额
	UserMoney float64    `orm:"default(0)"`         //可提现资金
}

//管辖地表
type HoldLocal struct {
	Id              int
	Country         string             `orm:"size(30)"` //国
	Province        string             `orm:"size(30)"` //省
	Town            string             `orm:"size(30)"` //市
	Prefecture      string             `orm:"size(30)"` //县
	BigDistribution []*BigDistribution `orm:"reverse(many)"`
	Distribution    []*Distribution    `orm:"reverse(many)"`
}

//分销商表
type Distribution struct {
	Id        int
	User      *User        `orm:"rel(fk);unique"`                //用户id-用户表外键
	AdTer     *Territory   `orm:"rel(fk)"`                       //管辖领域-领域表外键
	AsLocal   *HoldLocal   `orm:"rel(fk)"`                       //管辖地-管辖地表外键
	Money     float64      `orm:"default(0)"`                    //分销商总资金
	Type      string       `orm:"size(1);type(char);default(0)"` //分销商类型 0-普通分销员 1-学生等等分销员
	Address   []*Address   `orm:"reverse(many)"`
	ShopCart  []*ShopCart  `orm:"reverse(many)"`
	OrderInfo []*OrderInfo `orm:"reverse(many)"`
}

//收货地址表
type Address struct {
	Id           int
	Distribution *Distribution `orm:"rel(fk)"`                       //分销商id
	Content      string        `orm:"size(255)"`                     //地址内容
	Name         string        `orm:"size(10)"`                      //收货人
	Tele         string        `orm:"size(11);type(char)"`           //收货电话
	IsDefault    string        `orm:"size(1);type(char);default(0)"` //默认地址 0-不是默认地址 1-默认地址
	OrderInfo    []*OrderInfo  `orm:"reverse(many)"`
}

//商品表
type Goods struct {
	Id          int
	Name        string         `orm:"size(50)"` //商品名称
	GoodsBrief  string         //商品简介 数据库类型应该为blob
	GoodsState  string         `orm:"size(1);type(char);default(0)"` //商品状态 0-上架 1-下架
	Explain     string         `orm:"type(30)"`                      //说明
	CreateTime  time.Time      `orm:"auto_now_add"`                  //商品录入时间
	UpdateTime  time.Time      `orm:"auto_now"`                      //最后修改时间
	SalesValue  int            `orm:"default(0)"`                    //销量
	Category    *GoodsCategory `orm:"rel(fk)"`                       //分类id
	ShopCart    []*ShopCart    `orm:"reverse(many)"`
	GoodsBanner []*GoodsBanner `orm:"reverse(many)"`
	OrderGoods  []*OrderGoods  `orm:"reverse(many)"`
	GoodsParameter []*GoodsParameter `orm:"reverse(many)"`
}

//商品参数表
type GoodsParameter struct {
	Id             int
	Goods		   *Goods	`orm:"rel(fk)"` //
	Parameter      string  `orm:"size(30)"` //商品参数
	Parameter2     string  //商品参数Json形式(暂时不用) 数据库类型应该为blob
	GoodsTruePrice float64 //商品市场价
	GoodsNowPrice  float64 //商品批发价
	IsDefault      string  `orm:"size(1);type(char);default(0)"` //默认参数 0-不默认 1-默认
	GoodsNumber    int     //产品库存
}

//商品图片表
type GoodsBanner struct {
	Id         int
	Goods      *Goods `orm:"rel(fk)"`    //商品id
	GoodsOrder int     //图片排序 0-封面，1-第一张...
	GoodsUrl   string `orm:"size(255)"`  //图片路径
}

//商品分类表
type GoodsCategory struct {
	Id        int
	Name      string   `orm:"size(30);unique"` //分类名称
	ShowIndex int      `orm:"default(0)"`      //分类排序
	ParentId  int      `orm:"default(0)"`      //父分类 0-为一级分类
	Goods     []*Goods `orm:"reverse(many)"`
}

//购物车表
type ShopCart struct {
	Id           int
	Distribution *Distribution `orm:"rel(fk)"`      //分销商id
	Good         *Goods        `orm:"rel(fk)"`      //商品id
	Number       int           `orm:"default(1)"`   //数量
	Time         time.Time     `orm:"auto_now_add"` //添加时间
}

//订单表
type OrderInfo struct {
	Id           int
	Distribution *Distribution `orm:"rel(fk)"`                       //分销商id
	Address      *Address      `orm:"rel(fk)"`                       //地址id
	TotalCount   int           `orm:"default(0)"`                    //商品总数
	TotalPrice   float64       `orm:"default(0)"`                    //商品总价
	OrderStatus  string        `orm:"size(1);type(char);default(0)"` //订单状态 0-未支付 1-已支付 2-已发货 3-已签收 4-退款
	CreateTime   time.Time     `orm:"auto_now_add"`                  //下单时间
	OrderGoods   []*OrderGoods `orm:"reverse(many)"`
}

//订单商品表
type OrderGoods struct {
	Id    int
	Order *OrderInfo `orm:"rel(fk)"`    //订单id
	Goods *Goods     `orm:"rel(fk)"`    //商品id
	Count int        `orm:"default(0)"` //小件数量
	Price float64    `orm:"default(0)"` //小件价格
}

func init() {
	// set default database
	dns := "root:123456@tcp(127.0.0.1:3306)/zdshop?charset=utf8mb4&parseTime=True"
	_ = orm.RegisterDataBase("default", "mysql", dns)

	// register model
	orm.RegisterModel(new(User), new(Logistics), new(Information), new(Admin), new(Operation), new(Territory))
	orm.RegisterModel(new(BigDistribution), new(HoldLocal), new(Distribution), new(Address), new(Goods))
	orm.RegisterModel(new(GoodsCategory), new(GoodsParameter), new(ShopCart), new(GoodsBanner), new(OrderInfo), new(OrderGoods))

	// create table
	//第二个参数改为true,则每次加载都会重置数据库(原有数据会丢失)
	_ = orm.RunSyncdb("default", false, true)
}
