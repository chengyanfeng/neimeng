package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
	"os"
	. "neimeng/controller"
	. "neimeng/def"
	. "neimeng/task"
	. "neimeng/util"
)

func main() {

	MODE = Trim(os.Getenv("mode"))
	beego.BConfig.Listen.HTTPPort = 96                     //端口设置
	beego.BConfig.RecoverPanic = true                        //开启异常捕获
	beego.BConfig.EnableErrorsShow = true                    //打印错误记录
	beego.InsertFilter("/*", beego.BeforeRouter, BaseFilter) //路由过滤

	// 自动匹配路由
	beego.AutoRouter(&ApiController{})
	beego.AutoRouter(&JsonController{})
	beego.AutoRouter(&NewsController{})
	beego.AutoRouter(&WebController{})
	beego.AutoRouter(&PupController{})
	Mkdir("./logs")                                        //创建日志文件夹
	beego.SetLogger("file", `{"filename":"logs/run.log"}`) //定义日志文件
	beego.BeeLogger.SetLogFuncCallDepth(4)
	//调用以下函数处理接口数据

	/*Change_old()
	New_list()
	Hot_news()
	Day30()
	Works()
	Author()
	Twitter()
	Wx()
	Web()
	App()
	SelectLists()*/
	Mediabank()
	go func() { //开启协程
		InitCache() //初始化
		crontab()   //开启定时任务
	}()
	beego.Run() //启动项目
}

func crontab() {
	toolbox.AddTask("pd", toolbox.NewTask("pd", "0 */1 * * * *", func() error { //每10分钟运行以下函数
		Dhq <- func() {
		/*	Change_old()
			New_list()
			Hot_news()
			Day30()
			Works()
			Author()
			Twitter()
			Wx()
			Web()
			App()
			SelectLists()*/
			Mediabank()
		}
		return nil
	}))
	toolbox.StartTask() //开启定时任务
}
