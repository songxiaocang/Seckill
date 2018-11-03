package controller

import (
	"Seckill/SecProxy/service"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strings"
	"time"
)

type SkillController struct {
	beego.Controller
}

func (p *SkillController) SecKill() {
	productId, e := p.GetInt("product_id")

	result := make(map[string]interface{})
	result["code"] = 0
	result["msg"] = "success"
	defer func() {
		p.Data["json"] = result
		p.ServeJSON()
	}()

	if e != nil {
		result["code"] = service.ErrorInvalidRequest
		result["msg"] = "invalid request"
		return
	}

	secRequest := service.NewSecRequest()
	secRequest.ProductId = productId
	secRequest.Source = p.GetString("src")
	secRequest.AuthCode = p.GetString("authCode")
	secRequest.SecTime = p.GetString("secTime")
	secRequest.Nance = p.GetString("nance")

	secRequest.UserId, e = p.GetInt("user_id")
	secRequest.AccessTime = time.Now()
	if e != nil {
		result["code"] = service.ErrorInvalidRequest
		result["msg"] = "invalid request"
		return
	}

	secRequest.UserAuthSign = p.Ctx.GetCookie("userAuthSign")
	if len(p.Ctx.Request.RemoteAddr) > 0 {
		secRequest.ClientAddr = strings.Split(p.Ctx.Request.RemoteAddr, ":")[0]
	}
	secRequest.ClientReference = p.Ctx.Request.Referer()
	secRequest.CloseNotify = p.Ctx.ResponseWriter.CloseNotify()

	logs.Debug("request params:%v", secRequest)

	data, err, code := service.SecKill(secRequest)
	if err != nil {
		result["code"] = service.ErrorInvalidRequest
		result["msg"] = err.Error()
		return
	}

	result["code"] = code
	result["data"] = data
	return
}

func (p *SkillController) SecInfo() {
	//p.Data["json"] = "secInfo"
	//p.ServeJSON()
	productId, e := p.GetInt("product_id")
	logs.Debug("receive arg:%d", productId)
	result := make(map[string]interface{})

	result["code"] = 0
	result["msg"] = "success"

	defer func() {
		p.Data["json"] = result
		p.ServeJSON()
	}()

	if e != nil {
		data, err, code := service.SecInfoList()
		if err != nil {
			result["code"] = code
			result["msg"] = err.Error()

			logs.Error("invalid request,get all product failed:%v", err)
			return
		}

		result["code"] = code
		result["data"] = data
	} else {
		data, err, code := service.SecInfo(productId)
		if err != nil {
			result["code"] = code
			result["msg"] = err.Error()

			logs.Error("invalid request,get product_id failed:%v", err)
			return
		}

		result["code"] = code
		result["data"] = data
	}

}
