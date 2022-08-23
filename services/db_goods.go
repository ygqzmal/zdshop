package services

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"zdshop/models"
)

func QueryGoodsInfo(id int) (models.Goods, []*models.GoodsParameter, []*models.GoodsBanner, error) {
	o := orm.NewOrm()
	//商品基本属性
	var good models.Goods
	//good.Id = id
	//err = o.Read(&good) //这个查询只会将外键id查询出来,而且是查询所有,无法指定查询字段
	//添加的分类id一定要存在，否则会报错
	//查询goods表Id为id的值, 想要查关联的外键的话,要加上relatedSel("model中的外键字段"), 会把外键所有数据查询出来而且无法指定
	//err = o.QueryTable("Goods").RelatedSel("Category").Filter("Id", id).One(&good,"Id","Name","GoodsBrief","Explain","SalesValue")
	//想对外键的表进行某个字段查询做不到，只能查询外键的id
	err := o.QueryTable("Goods").Filter("Id", id).One(&good, "Id", "Name", "GoodsBrief", "Explain", "SalesValue", "category_id")
	if err != nil {
		beego.Info("商品查询失败: err", err)
		return good, nil, nil, err
	}
	//商品参数(一对多), 如果加上relatedSel反而会加强查询结果，把所有关联都查询到。
	//如果想指定查询内容, 在all后面加上查询参数即可
	var parameters []*models.GoodsParameter
	_, err = o.QueryTable("GoodsParameter").Filter("Goods__Id", id).All(&parameters, "Id", "Parameter", "GoodsTruePrice", "GoodsNowPrice", "IsDefault", "GoodsNumber")
	if err != nil {
		beego.Info("商品参数查询失败")
		return good, nil, nil, err
	}
	//商品图片
	var banner []*models.GoodsBanner
	_, err = o.QueryTable("GoodsBanner").Filter("Goods__Id", id).OrderBy("Id").All(&banner, "Id", "GoodsUrl", "IsDefault")
	if err != nil {
		beego.Info("商品参数查询失败")
		return good, nil, nil, err
	}
	return good, parameters, banner, nil
}
