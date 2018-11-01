package controller

import (
	"Seckill/SecProxy/service"
	"github.com/astaxie/beego"
)


type SkillController struct {
	beego.Controller
}


func(p *SkillController) SecKill(){
	p.Data["json"] = "secKill"
	p.ServeJSON()
}

func(p *SkillController) SecInfo(){
	//p.Data["json"] = "secInfo"
	//p.ServeJSON()
	productId, e := p.GetInt("product_id")
	result := make(map[string]interface{})

	result["code"] = 0
	result["msg"] = "success"

	defer func() {
		p.Data["json"] = result
		p.ServeJSON()
	}()

	if e!=nil {
		result["code"] = 1001
		result["msg"] = "invalid argument"
		return
	}

	data,err,code := service.SecInfo(productId)
	if err!=nil {
		result["code"] = code
		result["msg"] = err.Error()
		return
	}

	result["code"]=code
	result["data"]=data


}