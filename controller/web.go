package controller

import (
	"github.com/astaxie/beego/toolbox"
	. "../datasource"
	. "../def"
	. "../util"
	"sort"
	"strings"
	"time"
	"fmt"
)

type WebController struct {
	BaseController
}

var list_dh = 0  //新闻线索部分
var topic_dh = 0 //选题部分
var hot_dh = 0   //热点部分
var track_dh = 0 //CRT部分
var trace_dh = 0 //新媒体运营
var trac_dh = 0  //新媒体追踪
var task = 0     //定时任务状态

// curl "localhost:8080/web/title"  dh_srv.txt

func (this *WebController) Title() {
	tab2 := []string{"互联网", "新华电讯"}
	tab3 := []string{"政治", "财经", "社会", "科技", "娱乐", "体育"}
	r := P{}
	r["orign_list"] = tab2
	r["trade_list"] = tab3
	this.EchoJsonMsg(r)
}

// curl "localhost:8080/web/list" -d "tab2=77&week=1&tab3=77&page=1&rows=10" dh_srv.txt

func (this *WebController) List() {
	//新闻线索模块
	week, _ := this.GetInt("week") //获取时间周期
	//以下为算出每个周期的开始时间，人工录入的新闻都是按时间戳筛选的，接口的数据可直接调用week字段
	var st int64
	local := time.Now().Local().Format("2006-01-02")
	today, _ := time.Parse("2006-01-02 15:04:05", local + " 00:00:00")
	now := time.Now().Unix()
	p := this.FormToP("tab2", "tab3")
	if week == 3 {
		st = today.Unix() - 30 * 3600 * 24 - 8 * 3600
	} else if week == 2 {
		st = today.Unix() - 7 * 3600 * 24 - 8 * 3600
	} else {
		st = today.Unix() - 8 * 3600
	}
	if p["tab2"] == "新华电讯" {
		//新华电讯下的数据是按时间戳筛选的
		if week == 1 {
			p["date"] = P{"$gte": st, "$lte": now}
		} else {
			et := today.Unix() - 8 * 3600 - 1
			p["date"] = P{"$gte": st, "$lte": et}
		}
	} else {
		p["week"] = 1
	}
	p["news_class"] = 4
	p["old"] = 0
	p["dh"] = list_dh
	total := D(News).Find(p).Count()
	sort := "-date"
	if p["tab2"] == "互联网" {
		sort = "-hot"
	} else {

	}
	list := *D(News).Find(p).Sort(sort).Limit(100).All()
	if len(list) == 0 {
		//如果实时数据为空则开启备份数据
		q := P{}
		q["news_class"] = 4
		if p["tab3"] != nil {
			q["tab3"] = p["tab3"]
		}
		if p["tab2"] != nil {
			q["tab2"] = p["tab2"]
		}
		q["dh"] = list_dh
		q["old"] = 1
		list = *D(News).Find(q).Sort(sort).Limit(100).All()
	}
	for _, v := range list {
		v["name"] = v["title"] //前端用name接受
		if week == 1 {
			//修改当日时间
			v["date"] = today.Unix()
		}
	}
	totals := 0
	r := P{}
	r["total"] = total + totals
	r["page"], _ = this.GetInt("page", 1) //分页机制暂时没用到
	r["list"] = list
	this.EchoJsonMsg(r)
}

//curl "localhost:8080/web/topic" -d "page=1&rows=10&organization=aaa" dh_srv.txt

func (this *WebController) Topic() {
	//选题总览
	p := P{}
	p["news_class"] = 1
	p["dh"] = topic_dh
	p["old"] = 0
	total := D(News).Find(p).Count()
	list := *D(News).Find(p).Sort("-date").Limit(100).All()
	if len(list) == 0 {
		p := P{}
		p["news_class"] = 1
		p["dh"] = topic_dh
		p["old"] = 1
		list = *D(News).Find(p).Sort("-date").Limit(100).All()
	}
	r := P{}
	r["total"] = total
	//r["page"], _ = this.GetInt("page", 1)
	index := 1
	for _, v := range list {
		//为每条信息添加序号
		v["index"] = index
		index++
	}
	r["list"] = list
	this.EchoJsonMsg(r)
}

//curl "localhost:8080/web/hot" dh_srv.txt

func (this *WebController) Hot() {
	//互联网实时热点
	var hot []interface{}
	status := 0
	area := []string{}
	news := []map[string]interface{}{}
	hots := *D(Json).Find(P{"type": 4, "dh": hot_dh}).Sort("-ct").Limit(34).All()                //获取全局热点

	cat_number := []int{}
	cat := *D(Cat).Find(nil).All()
	for _, s := range cat {
		count := D(News).Find(P{"dh": hot_dh, "news_class": 30, "old": status, "tab3": ToString(s["chname"])}).Count()
		if count != 0 {
			cat_number = append(cat_number, count)
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(cat_number)))
	other_total := 0
	if len(cat_number) > 7 {
		cat_other := cat_number[7:]
		for _, v := range cat_other {
			other_total = other_total + v
		}
		cat_number = cat_number[:7]
	} else {

	}

	new := *D(News).Find(P{"news_class": 6, "dh": hot_dh, "old": 0}).Sort("-date").Limit(30).All() //热点新闻
	if len(new) == 0 {
		new = *D(News).Find(P{"news_class": 6, "dh": hot_dh, "old": 1}).Sort("-date").Limit(30).All()
	}
	for _, v := range new {
		//获取热点地区
		if !InArray(ToString(v["area"]), area) {
			area = append(area, ToString(v["area"]))
		}
	}
	for _, o := range area {
		//将新闻按地区分类
		city := map[string]interface{}{}
		city["area"] = o
		list := *D(News).Find(P{"news_class": 6, "dh": hot_dh, "area": o, "old": 0}).Limit(4).Sort("-date").All()
		if len(list) == 0 {
			list = *D(News).Find(P{"news_class": 6, "dh": hot_dh, "old": 1}).Sort("-date").Limit(4).All()
		}
		city["list"] = list
		news = append(news, city)

	}
	for _, o := range hots {
		//为全局热点分别选取最近一条新闻
		json := *JsonDecode([]byte(ToString(o["json"])))
		hot_liist := *D(News).Find(P{"news_class": 6, "dh": hot_dh, "area": json["area"], "old": 0}).Sort("-date").One()
		if len(hot_liist) == 0 {
			hot_liist = *D(News).Find(P{"news_class": 6, "dh": hot_dh, "area": json["area"], "old": 1}).Sort("-date").One()
		}
		if len(ToString(hot_liist["title"])) != 0 {
			json["title"] = hot_liist["title"]
			json["href"] = hot_liist["href"]
			json["date"] = hot_liist["date"]
			hot = append(hot, JsonEncode(json))
		}
	}

	r := P{}
	r["hot"] = hot
	r["news"] = news
	this.EchoJsonMsg(r)
}

//curl "localhost:8080/web/operative" -d "week=1" dh_srv.txt

func (this *WebController) Operative() {
	//新媒体运营
	week, _ := this.GetInt("week")
	client := *D(News).Find(P{"news_class": 10, "dh": trace_dh, "origin": 1, "week": week}).Sort("-ct").One()  //客户端
	web := *D(News).Find(P{"news_class": 10, "dh": trace_dh, "origin": 2, "week": week}).Sort("-ct").One()     //网站
	wx := *D(News).Find(P{"news_class": 10, "dh": trace_dh, "origin": 3, "week": week}).Sort("-ct").One()      //微信
	twitter := *D(News).Find(P{"news_class": 10, "dh": trace_dh, "origin": 4, "week": week}).Sort("-ct").One() //微博
	client["hot_word"] = strings.Split(ToString(client["hot_word"]), "@")                                      //将字符串截取为数组
	web["hot_word"] = strings.Split(ToString(web["hot_word"]), "@")
	wx["hot_word"] = strings.Split(ToString(wx["hot_word"]), "@")
	twitter["hot_word"] = strings.Split(ToString(twitter["hot_word"]), "@")
	r := P{}
	r["client"] = client
	r["web"] = web
	r["wx"] = wx
	r["twitter"] = twitter
	this.EchoJsonMsg(r)
}


//curl "localhost:8080/web/track" -d "week=1" dh_srv.txt

func (this *WebController) Track() {
	//CRT新媒体追踪
	var total int
	index := 1
	var works, auths = P{}, P{}
	week, _ := this.GetInt("week")
	work_count := *D(Total).Find(P{"week": week, "dh": track_dh}).One() //获取作品总数
	total = ToInt(work_count["work_count"])
	works_reprint := *D(News).Find(P{"dh": track_dh, "news_class": 7, "week": week, "old": 0, "tp": "2"}).Sort("-reprint").Limit(10).All() //作品转载列表
	if len(works_reprint) == 0 {
		works_reprint = *D(News).Find(P{"dh": track_dh, "news_class": 7, "week": week, "old": 1, "tp": "2"}).Sort("-reprint").Limit(10).All()
	}
	works_comment := *D(News).Find(P{"dh": track_dh, "news_class": 7, "week": week, "old": 0, "tp": "1"}).Sort("-comment").Limit(10).All() //作品评论列表
	if len(works_comment) == 0 {
		works_comment = *D(News).Find(P{"dh": track_dh, "news_class": 7, "week": week, "old": 1, "tp": "1"}).Sort("-reprint").Limit(10).All()
	}
	auths_reprint := *D(News).Find(P{"dh": track_dh, "news_class": 8, "week": week, "old": 0, "tp": "2"}).Sort("-reprint").Limit(10).All() //作者转载列表
	if len(auths_reprint) == 0 {
		auths_reprint = *D(News).Find(P{"dh": track_dh, "news_class": 8, "week": week, "old": 1, "tp": "2"}).Sort("-reprint").Limit(10).All()
	}
	auths_comment := *D(News).Find(P{"dh": track_dh, "news_class": 8, "week": week, "old": 0, "tp": "1"}).Sort("-comment").Limit(10).All() //作者评论列表
	if len(auths_comment) == 0 {
		auths_comment = *D(News).Find(P{"dh": track_dh, "news_class": 8, "week": week, "old": 1, "tp": "1"}).Sort("-reprint").Limit(10).All()
	}
	for _, v := range works_reprint {
		//添加序号
		v["index"] = index
		index++
	}
	index = 1
	for _, v := range works_comment {
		v["index"] = index
		index++
	}
	index = 1
	for _, v := range auths_reprint {
		v["index"] = index
		index++
	}
	index = 1
	for _, v := range auths_comment {
		v["index"] = index
		index++
	}
	index = 1
	works["reprint"] = works_reprint
	works["comment"] = works_comment
	auths["reprint"] = auths_reprint
	auths["comment"] = auths_comment
	r := P{}
	r["total"] = total
	r["works"] = works
	r["auths"] = auths
	this.EchoJsonMsg(r)
}

//curl "localhost:8080/web/article30" dh_srv.txt

func (this *WebController) Article30() {
	//文章发布30天
	//var month []string
	list := *D(Json).Find(P{"type": 14, "dh": track_dh, "old": 0}).Sort("-_id").One()
	if track_dh == 0 {
		day30 := ToString(list["json"]) //对象转json字符串
		this.EchoJsonMsg(day30)
	} else {
		//转换人工录入数据格式
		months := strings.Replace(ToString(list["json"]), "@", ",", -1) //将@替换为，
		str := months
		/*for _, v := range months {
			month = append(month, v)
		}*/
		this.EchoJsonMsg(str)
	}
}

//curl "localhost:8080/web/change" dh_srv.txt

func (this *WebController) Change() {
	//前台数据展示切换 0：接口数据，1：人工录入
	tp := this.GetString("tp")
	val := this.GetString("val")
	switch tp {
	case "list_dh":
		list_dh = ToInt(val)
	case "topic_dh":
		topic_dh = ToInt(val)
	case "hot_dh":
		hot_dh = ToInt(val)
	case "track_dh":
		track_dh = ToInt(val)
	case "trace_dh":
		trace_dh = ToInt(val)
	case "trac_dh":
		trac_dh = ToInt(val)
	}
	this.EchoJsonOk()
}

//curl "localhost:8080/web/data_status" dh_srv.txt

func (this *WebController) Data_status() {
	//获取每个模块的数据类型
	p := P{}
	p["list_dh"] = list_dh
	p["topic_dh"] = topic_dh
	p["hot_dh"] = hot_dh
	p["track_dh"] = track_dh
	p["trace_dh"] = trace_dh
	p["trac_dh"] = trac_dh
	this.EchoJsonMsg(p)
}

func (this *WebController) Change_crontab() {
	status := this.GetString("status")
	if ToInt(status) == 1 {
		task = 1
		toolbox.StopTask()
	} else {
		task = 0
		toolbox.StartTask()
	}
	this.EchoJsonOk()
}

func (this *WebController) Task_status() {
	this.EchoJsonMsg(task)
}

func(this *WebController)Media_brank(){
	fmt.Print("1221")
	p:=P{"backups":0}

	list:=*D(Media).Find(p).Sort("序号").All()
	if len(list)==0{
		p:=P{"backups":1}
		list:=*D(Media).Find(p).Sort("序号").All()
		this.EchoJsonMsg(list)
	}else {
	this.EchoJsonMsg(list)
	}

	}


func(this *WebController)First_line(){
	fmt.Print("1221")
	tp := this.GetString("type")
	p:=P{"backups":0,"type":tp}

	list:=*D(Firstline).Find(p).Sort("序号").All()
	if len(list)==0{
		p:=P{"backups":1,"type":tp}
		list:=*D(Firstline).Find(p).Sort("序号").All()
		this.EchoJsonMsg(list)
	}else {
		this.EchoJsonMsg(list)
	}

}

