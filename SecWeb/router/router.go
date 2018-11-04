package router

import (
	"Seckill/SecWeb/controller/product"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/product/list", &product.ProductController{}, "*ListProducts")

}
