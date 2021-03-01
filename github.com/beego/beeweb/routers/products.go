package routers

import (
	"github.com/beego/beeweb/models"
)

type ProductsRouter struct {
	baseRouter
}

func (this *ProductsRouter) Get() {
	this.TplNames = "products.html"
	this.Data["IsProducts"] = true
	this.Data["Products"] = models.Products
}
