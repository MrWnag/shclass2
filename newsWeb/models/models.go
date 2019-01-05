package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id int
	UserName string `orm:"unique"`
	Pwd string
	Articles []*Article `orm:"reverse(many)"`//设置多对多反向关系(可互换)
}

//rel(fk) reverse(many)    rel(m2m)   reverse(many)    rel      reverse

type Article struct {
	Id int `orm:"pk;auto"`
	Title string `orm:"size(100)"`
	Time time.Time `orm:"type(datetime);auto_now"`
	Count int `orm:"default(0)"`
	Img string `orm:"null"`
	Content string
	ArticleType *ArticleType `orm:"rel(fk);set_null;null"`//设置一对多正向关系 看一对多那个表需要外键就设置rel(fk)
	Users []*User `orm:"rel(m2m)"`//设置多对多正向关系(可互换)
	price float64 `orm:"digits(10);decimals(2)"`
}

type ArticleType struct {
	Id int
	TypeName string `orm:"size(20);unique"`
	Articles []*Article `orm:"reverse(many)"`//设置一对多反向关系
}
//需要创建一对多的类型表以及多对多的用户与文章关系表

func init(){
	//建表的三步骤
	//注册数据库
	//第一个,为什么要用别名
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/newsWeb?charset=utf8")
	//注册表
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	//跑起来
	orm.RunSyncdb("default",false,true)
}
