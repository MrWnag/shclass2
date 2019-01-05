package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"newsWeb/models"
)

type UserController struct {
	beego.Controller
}

func (this*UserController)ShowRegister(){
	this.TplName = "register.html"
}
//业务四步骤: 请求-> 路由 ->  控制器   ->  页面

//注册实现 ： 获取数据->校验数据->操作数据->返回数据



//注册业务
func(this*UserController)HandleRegister(){
	//获取前端传递的数据
	userName:=this.GetString("userName")
	pwd := this.GetString("password")
	beego.Info(userName,pwd)

	//校验数据
	if userName == "" || pwd == ""{
		beego.Error("用户名或者密码不能为空")
		this.TplName = "register.html"
		return
	}
	//把数据插入到数据库
	//数据库的插入操作
	//1.获取orm对象
	o := orm.NewOrm()
	//2.获取插入对象
	var user models.User
	//3.给插入对象赋值
	user.UserName = userName
	user.Pwd = pwd
	//4.执行插入操作
	count,err :=o.Insert(&user)
	if err != nil{
		beego.Error("用户注册失败")
		this.TplName = "register.html"
		return
	}
	beego.Info("插入的数据条数=",count)
	//给前段返回一个页面
	this.Redirect("/login",302)

}

//登录业务
func(this*UserController)ShowLogin(){
	userName := this.Ctx.GetCookie("userName")
	if userName != ""{
		this.Data["userName"] = userName
		this.Data["checked"] = "checked"
	}else{
		this.Data["userName"] = userName
		this.Data["checked"] = ""
	}
	this.TplName = "login.html"
}

//处理登录数据
//命名即注释
func(this*UserController)HandleLogin(){
	//获取数据
	userName :=this.GetString("userName")
	pwd := this.GetString("password")
	//校验数据
	if userName == "" || pwd == "" {
		//beego.Error("用户名或者密码不能为空")
		this.Data["err"] = "用户名或者密码不能为空"
		this.TplName = "login.html"
		return
	}
	//处理数据
		//数据库的查询操作
		//1.获取orm对象
		o := orm.NewOrm()
		//2.获取查询对象
		var user models.User
		//3.给查询条件赋值
		user.UserName = userName
		//4.查询
		err := o.Read(&user,"UserName")
		if err != nil{
			this.Data["err"] = "用户名不存在"
			this.TplName = "login.html"
			return
		}

		if user.Pwd != pwd{
			this.Data["err"] = "密码错误"
			this.TplName = "login.html"
			return
		}
	//返回数据

	//作用设置cookie
	//cookie的key
	//cookie中存的value
	//设置生效时间
	remember := this.GetString("remember")
	//beego.Info("remember =",remember)
	if remember == "on"{
		this.Ctx.SetCookie("userName",userName,3600 * 24)
	}else{
		this.Ctx.SetCookie("userName",userName,-1)
	}

	//this.Ctx.WriteString("登录成功")
	this.SetSession("userName",userName)
	this.Redirect("/article/index",302)
}

//退出登录
func(this*UserController)Logout(){
	//删除session
	this.DelSession("userName")
	this.Redirect("/login",302)
}