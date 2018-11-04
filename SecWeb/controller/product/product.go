package product

import "github.com/astaxie/beego"

type ProductController struct {
	beego.Controller
}

func (p *ProductController) ListProduct() {}
