package utils

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"os"
	"path"
	"strconv"
	"time"
	"zdshop/models"
)


//调用
/*
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

//删除某个商品全部轮播图片
func HandleFiles(id int) (message string){
	o := orm.NewOrm()
	var goodsBanner []*models.GoodsBanner
	_, err := o.QueryTable("GoodsBanner").Filter("Goods__Id", id).All(&goodsBanner,"GoodsUrl")
	if err != nil {
		beego.Info("商品轮播图片查询失败")
		return "商品轮播图片查询失败"
	}
	for _, value := range goodsBanner {
		imgPath := value.GoodsUrl
		err := os.Remove(imgPath)
		if err != nil {
			beego.Info("图片删除失败,err", err)
			return "图片删除失败"
		}
	}
	return ""
}