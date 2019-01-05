package controllers

import
("github.com/astaxie/beego"
"github.com/gomodule/redigo/redis"
	)

type GoRedis struct {
	beego.Controller
}

func(this*GoRedis)ShowGet(){
	//链接数据库
	conn,err :=redis.Dial("tcp",":6379")
	defer 	conn.Close()
	if err != nil {
		beego.Error("redis数据库链接失败",err)
		return
	}

	//操作数据库
	resp,err :=conn.Do("mget","class3","aa")

	re,err :=redis.Values(resp,err)

	//获取值相同类型和不同类型的操作
	var string1 string
	var int1 int

	redis.Scan(re,&string1,&int1)
	beego.Info("回复值=",string1,int1)

	//关闭数据库
}