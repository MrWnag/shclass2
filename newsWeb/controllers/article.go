package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"encoding/gob"
	"newsWeb/models"
	"time"
	"path"
	"math"
	"bytes"
)

type ArticleController struct {
	beego.Controller
}

//展示首页
func (this*ArticleController)ShowIndex(){
	userName := this.GetSession("userName")
	if userName == nil{
		this.Redirect("/login",302)
		return
	}
	this.Data["userName"] = userName.(string)


	//获取数据
	//select * from article;
	//创建orm对象
	o := orm.NewOrm()
	//指定表
	qs :=o.QueryTable("Article") //queryseter  查询对象
	//定义一个容器
	var articles []models.Article
	//查询
	//_,err :=qs.All(&articles)
	//if err != nil {
	//	beego.Error("查询所有文章错误")
	//	this.TplName = "index.html"
	//	return
	//}
	//分页处理
	//1.获取总记录数和总也数据
	//count,err :=qs.RelatedSel("ArticleType").Count()//获取总记录数

	//定义一个每页显示条目数
	page := 2

	//获取总页数
	//pageCount := math.Ceil(float64(count) / float64(page))
	//this.Data["count"] = count
	//this.Data["pageCount"] = int(pageCount)


	//beego.Info("count=",count,"pageIndex = ",page,"pageCount=",pageCount)


	//2.首页和末页实现
	pageIndex,err :=this.GetInt("pageIndex")
	if err != nil{
		pageIndex = 1
	}

	this.Data["pageIndex"] = pageIndex
	//获取数据
	start := page * ( pageIndex  -1)

	//获取选中的类型
	typeName :=this.GetString("select")
	var count int64
	var err1 error
	if typeName == "" {
		qs.Limit(page,start).RelatedSel("ArticleType").All(&articles)
		count,err1 =qs.RelatedSel("ArticleType").Count()

	}else{
		//select * from article where articleType.typeName = "娱乐新闻"
		//orm中一对多的查询是惰性查询
		count,err1 =qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()
		qs.Limit(page,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
	}
	if err1 != nil {
		beego.Error("查询数据条目数错误")
		this.TplName = "index.html"
		return
	}
	pageCount := math.Ceil(float64(count) / float64(page))
	this.Data["count"] = count
	this.Data["pageCount"] = int(pageCount)


	this.Data["typeName"] = typeName



	//3.上一页和下一页




	//传递数据
	this.Data["articles"] = articles





	//获取所有类型数据
	var articleTypes []models.ArticleType
	conn,err :=redis.Dial("tcp",":6379")
	//连接数据库
	if err != nil{
		beego.Error("redis数据库链接失败")
	}

	resp ,err:= redis.Bytes(conn.Do("get","articleTypes"))
	//定义一个解码器
	dec :=gob.NewDecoder(bytes.NewReader(resp))
	//解码
	dec.Decode(&articleTypes)

	if len(articleTypes) == 0{
		//获取所有类型
		o.QueryTable("ArticleType").All(&articleTypes)
		//序列化存储
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		enc.Encode(&articleTypes)

		conn.Do("set","articleTypes",buffer.Bytes())

		beego.Info("从mysql中获取数据")
	}






	////序列化和反序列化
	////1.要有一个容器，用来接受编码之后的字节流
	//var buffer bytes.Buffer
	////2.要有一个编码器
	//enc :=gob.NewEncoder(&buffer)
	////3.编码
	//enc.Encode(&articleTypes)
	//conn.Do("set","articleTypes",buffer.Bytes())
	//
	//resp,err :=conn.Do("get","articleTypes")
	////先获取字节流数据
	//types,err :=redis.Bytes(resp,err)
	////2.获取解码器
	//dec :=gob.NewDecoder(bytes.NewReader(types))
	////3.解码
	//var testTyps []models.ArticleType
	//dec.Decode(&testTyps)
	//beego.Info(testTyps)


	//把数据传递给前端
	this.Data["articleTypes"] = articleTypes

	this.Layout = "layout.html"
	this.TplName = "index.html"
}

//展示添加文章页面
func(this*ArticleController)ShowAdd(){
	//获取所有的类型数据
	//获取orm对象
	o := orm.NewOrm()
	//定义一个容器
	var articleTypes []models.ArticleType
	//查询所有
	o.QueryTable("ArticleType").All(&articleTypes)

	this.Data["articleTypes"] = articleTypes

	this.TplName = "add.html"
}

//处理添加文章业务
func(this*ArticleController)HandleAdd(){
	//获取数据
	typeName := this.GetString("select")
	title :=this.GetString("articleName")
	content :=this.GetString("content")
	file,head,err :=this.GetFile("uploadname")
	defer file.Close()

	//校验数据
	if title == "" || content == "" || err != nil{
		this.Data["errmsg"] = "添加文章失败，请重新添加！"
		this.TplName = "add.html"
		return
	}
	//beego.Info(file,head)

	//1.文件存在覆盖的问题
	//加密算法

	//当前时间
	fileName := time.Now().Format("2006-01-02-15-04-05")
	ext := path.Ext(head.Filename)
	beego.Info(head.Filename,ext)
	//2.文件类型也需要校验
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
		this.Data["errmsg"] = "上传图片格式不正确，请重新上传"
		this.TplName = "add.html"
		return
	}
	//3.文件大小校验
	if head.Size > 5000000 {
		this.Data["errmsg"] = "上传图片过大，请重新上传"
		this.TplName = "add.html"
		return
	}

	//把图片存起来
	this.SaveToFile("uploadname","./static/img/"+fileName+ext)

	//处理数据
	//数据库的插入操作
	//获取orm对象
	o := orm.NewOrm()
	//获取一个插入对象
	var article models.Article
	//给插入对象赋值
	article.Title = title
	article.Content = content
	article.Img = "/static/img/"+fileName+ext
	//插入一个带类型的文章


	//外键的类型是文章类型对象指针，所以要插入带类型文章的时候需要把相对应的类型对象插入到文章中，可以获取到类型名称，可以根据类型名称获取类型对象
	//插入文章类型g  查询操作
	var articleType models.ArticleType
	//给查询条件赋值
	articleType.TypeName = typeName
	//查询所对应
	o.Read(&articleType,"TypeName")
	//把类型对象赋值给文章
	article.ArticleType = &articleType


	//插入到数据库
	o.Insert(&article)

	//返回数据
	this.Redirect("/article/index",302)
}

//展示文章详情
func(this*ArticleController)ShowContent(){
	//获取数据
	articleId,err :=this.GetInt("articleId")
	//数据校验
	if err != nil{
		beego.Error("请求链接错误")
		this.TplName = "index.html"
		return
	}

	//数据处理
	//获取orm对象
	o := orm.NewOrm()
	//获取查询对象
	var article models.Article
	//给查询条件赋值
	article.Id = articleId
	//查询
	err = o.Read(&article)
	if err != nil{
		beego.Error("查询文章不存在")
		//this.TplName = "index.html"
		this.Redirect("/article/index",302)
		return
	}
	article.Count += 1
	//更新数据
	o.Update(&article)

	//需要添加用户浏览记录
	//多对多的数据添加

	m2m :=o.QueryM2M(&article,"Users")
	//插入的是用户的对象
	var user models.User
	userName := this.GetSession("userName")
	user.UserName = userName.(string)
	o.Read(&user,"UserName")
	//插入
	m2m.Add(user)

	//获取浏览记录
	//o.LoadRelated(&article,"Users")

	//第二种多对多查询方法
	qs := o.QueryTable("User")
	var users []models.User
	qs.Filter("Articles__Article__Id",article.Id).Distinct().All(&users)
	//把数据传递给前端
	this.Data["users"] = users


	//返回数据
	this.Data["article"] = article
	this.Layout = "layout.html"
	this.TplName = "content.html"
}

//展示编辑页面
func(this*ArticleController)ShowUpdate(){
	//获取数据
	articleId,err := this.GetInt("articleId")
	//校验数据
	if err != nil {
		beego.Error("请求链接错误")
		this.Redirect("/article/index",302)
		return
	}
	//处理数据
	//更新操作
	//获取orm
	o := orm.NewOrm()
	//获取更新对象
	var article models.Article
	//给更新条件赋值
	article.Id = articleId
	//读一下
	err = o.Read(&article)
	if err != nil{
		beego.Error("要更新的文章不存在")
		this.Redirect("/article/index",302)
		return
	}


	//返回数据
	this.Data["article"] = article
	this.TplName = "update.html"
}

//封装函数
func UploadFunc(this*beego.Controller,filePath string)string{
	file,head,err :=this.GetFile(filePath)
	defer file.Close()

	if err != nil{
		return ""
	}

	//1.文件存在覆盖的问题
	//加密算法

	//当前时间
	fileName := time.Now().Format("2006-01-02-15-04-05")
	ext := path.Ext(head.Filename)
	beego.Info(head.Filename,ext)
	//2.文件类型也需要校验
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
		this.Data["errmsg"] = "上传图片格式不正确，请重新上传"
		this.TplName = "add.html"
		return ""
	}
	//3.文件大小校验
	if head.Size > 5000000 {

		this.Data["errmsg"] = "上传图片过大，请重新上传"
		this.TplName = "add.html"
		return ""
	}

	//把图片存起来
	this.SaveToFile(filePath,"./static/img/"+fileName+ext)
	return "/static/img/"+fileName+ext
}

//处理更新文章数据
func(this*ArticleController)HandleUpdate(){
	//获取数据
	articleName :=this.GetString("articleName")
	content :=this.GetString("content")
	fileAddr := UploadFunc(&this.Controller,"uploadname")
	articleId,err := this.GetInt("articleId")
	//校验数据
	if fileAddr == ""||articleName == "" || content == "" || err != nil{
		this.Data["errmsg"] = "上传数据失败"
		this.TplName = "update.html"
		return
	}

	//处理数据
	//数据库的更新操作
	o := orm.NewOrm()
	var article models.Article
	//查询操作
	article.Id = articleId
	err = o.Read(&article)
	if err != nil{
		this.Data["errmsg"] = "更新的文章id错误"
		this.TplName = "update.html"
		return
	}
	article.Title = articleName
	article.Content = content
	article.Img = fileAddr
	o.Update(&article)

	//返回数据
	this.Redirect("/article/index",302)
}

//删除文章
func(this*ArticleController)HandleDelete(){
	//获取数据
	articleId,err :=this.GetInt("articleId")
	//校验数据
	if err !=nil{
		beego.Error("删除文章链接错误")
		this.Redirect("/article/index",302)
		return
	}
	//处理数据
	//获取orm对象
	o := orm.NewOrm()
	//获取删除对象
	var article models.Article
	//给删除条件赋值
	article.Id = articleId
	//删除
	_,err =o.Delete(&article)
	if err != nil{
		beego.Error("删除失败")
		this.Redirect("/article/index",302)
		return
	}
	//返回数据
	this.Redirect("/article/index",302)
}

//展示添加类型页面
func(this*ArticleController)ShowAddType(){

	//获取数据并展示  获取所有类型
	//获取orm对象
	o := orm.NewOrm()
	//定义一个容器
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)

	//传递给前段
	this.Data["articleTypes"] = articleTypes

	this.TplName = "addType.html"
}

//处理添加类型数据
func(this*ArticleController)HandleAddType(){
	//获取数据
	typeName :=this.GetString("typeName")
	//校验数据
	if typeName == ""{
		beego.Error("类型名称为空，请重新输入")
		this.Redirect("/article/addType",302)
		return
	}
	//处理数据
	//数据库插入操作
		//获取orm对象
		o := orm.NewOrm()
		//获取插入对象
		var articleType models.ArticleType
		//给插入对象赋值
		articleType.TypeName = typeName
		//执行插入
		_,err :=o.Insert(&articleType)
		if err != nil{
			beego.Error("类型名重复，请重新插入")
			this.TplName = "addType.html"
			return
		}

	//返回数据
	this.Redirect("/article/addType",302)
}

//删除文章类型
func(this*ArticleController)DeleteType(){
	//获取数据
	typeId,err :=this.GetInt("typeId")
	//校验数据
	if err != nil{
		beego.Error("删除类型链接错误")
		this.Redirect("/article/addType",302)
		return
	}
	//处理数据
	//获取orm对象
	o := orm.NewOrm()
	//获取删除对象
	var articleType models.ArticleType
	//给删除对象赋值
	articleType.Id = typeId
	//删除
	_,err = o.Delete(&articleType)
	if err != nil{
		beego.Error("删除失败")
		this.Redirect("/article/addType",302)
		return
	}
	//返回数据
	this.Redirect("/article/addType",302)
}
