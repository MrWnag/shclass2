package main

import (
	_ "newsWeb/routers"
	"github.com/astaxie/beego"
	_ "newsWeb/models"
)

func main() {
	beego.AddFuncMap("pre",getPre)
	beego.AddFuncMap("next",getNext)
	beego.Run()
}

func getPre(pageIndex int)int{
	if pageIndex - 1 <= 0{
		return pageIndex
	}
	return pageIndex - 1
}

func getNext(pageIndex int,pageCount int)int{
	if pageIndex + 1 > pageCount{
		return pageCount
	}

	return pageIndex + 1
}
