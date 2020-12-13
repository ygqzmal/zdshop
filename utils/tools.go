package utils

import (
	"github.com/astaxie/beego"
	"path"
	"strconv"
	"time"
)

//文件上传
func HandleFile2(this *beego.Controller, image string) (message, url string) {
	//headers, err := this.GetFiles(image)
	file, header, err := this.GetFile(image)
	defer file.Close()
	if err != nil {
		return "图片上传失败", ""
	}
	//大小
	if header.Size > 500000 {
		return "图片太大无法上传", ""
	}
	//格式
	ext := path.Ext(header.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		beego.Info("图片格式不正确")
		return "图片格式不正确", ""
	}
	//防重名
	//存储 第一个参数要和GetFile的一样
	fileName := strconv.Itoa(int(time.Now().UnixNano())) + ext
	err = this.SaveToFile(image, "./upload/img/"+fileName)
	if err != nil {
		return "文件存储错误 ", ""
	}
	return "", "./upload/img"+fileName
}