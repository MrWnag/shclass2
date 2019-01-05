package routers

import (
	"newsWeb/controllers"
	"github.com/astaxie/beego"
    "github.com/astaxie/beego/context"
)

func init() {
    beego.InsertFilter("/article/*",beego.BeforeExec,filterFunc)
    beego.Router("/", &controllers.MainController{})
    //注册业务
    beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleRegister")
    //登录业务
    beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
    //文章首页
    beego.Router("/article/index",&controllers.ArticleController{},"get:ShowIndex")
    //添加文章
    beego.Router("/article/add",&controllers.ArticleController{},"get:ShowAdd;post:HandleAdd")
    //查看文章详情
    beego.Router("/article/content",&controllers.ArticleController{},"get:ShowContent")
    //编辑业务
    beego.Router("/article/update",&controllers.ArticleController{},"get:ShowUpdate;post:HandleUpdate")
    //删除业务
    beego.Router("/article/delete",&controllers.ArticleController{},"get:HandleDelete")
    //添加类型
    beego.Router("/article/addType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
    //退出登录
    beego.Router("/article/logout",&controllers.UserController{},"get:Logout")

    //删除类型
    beego.Router("/article/deleteType",&controllers.ArticleController{},"get:DeleteType")
    //go操作redis
    beego.Router("/redis",&controllers.GoRedis{},"get:ShowGet")
}

var filterFunc = func(ctx*context.Context) {
    userName := ctx.Input.Session("userName")
    if userName == nil{
        ctx.Redirect(302,"/login")
        return
    }

}
