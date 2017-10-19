package controller

import (
	. "../datasource"
	. "../def"
	. "../util"
	"strings"
)

type NewsController struct {
	BaseController
}

// curl "localhost:8080/news/save" -d "news_class=1&title=sa" dh_srv.txt

func (this *NewsController) Save() {
	defer func() {
		if err:=recover();err!=nil{

		}
	}()
	//添加或修改新闻
	p := this.FormToP("news_class", "tab2", "tab3", "title", "hot", "href", "content", "status", "head", "organization",
		"date", "type", "name", "keys", "subtitle1", "subtitle2", "subcontent1", "subcontent2", "orgList", "widget1", "widget2", "hot_word",
		"week", "origin", "orglist", "tp", "number", "periodical", "area", "url", "layout", "video", "aid", "sid") //接收前台传过来的新闻结构体字段，
	p["news_class"] = ToInt(p["news_class"])
	if p["news_class"] == 7 || p["news_class"] == 8 {
		//处理CRT部分作品类和作者类新闻信息
		tp := ToInt(p["tp"])
		if tp == 2 {
			//tp为2代表number为转载数，反之是评论数
			p["reprint"] = ToInt(p["number"])
		} else {
			p["comment"] = ToInt(p["number"])
		}
	}
	if p["news_class"] == 1 {
		p["aid"] = ToString(p["aid"])
	}
	if p["news_class"] == 2 {
		p["sid"] = ToString(p["sid"])
	}
	if p["news_class"] == 9 {
		p["name"] = ToString(p["title"])
	}
	if p["news_class"] == 4 {
		p["hot"] = ToInt(p["hot"])
	}
	if p["news_class"] == 5 {
		if len(ToString(p["title"]))<=1 || len(ToString(p["subtitle1"]))<=1 || len(ToString(p["subcontent1"]))<=1 || ToString(p["date"])=="NaN" {
			this.EchoJsonErr("参数不全")
			panic("错误")
		}
	}
	if len(ToString(p["week"])) > 0 {
		p["week"] = ToInt(p["week"])
	}
	if len(ToString(p["status"])) > 0 {
		p["status"] = ToInt(p["status"])
	}
	if len(ToString(p["origin"])) > 0 {
		p["origin"] = ToInt(p["origin"])
	}
	if len(ToString(p["date"])) > 0 {
		p["date"] = ToInt(p["date"]) / 1000
	}
	p["ct"] = Timestamp()
	if IsOid(this.GetString("id")) {
		//如果有id字段说明是修改，反之则是添加
		p["_id"] = this.GetOid("id")
		D(News).Save(&p)
	} else {
		if p["news_class"] == 1 || p["news_class"] == 2 || p["news_class"] == 9 {
			p["id"] = NewId()
		}
		p["dh"] = 1  //1代表人工录入数据
		p["old"] = 0 //0代表实时数据
		p["_id"] = NewId()
		D(News).Add(p)
	}
	this.EchoJsonOk(p)
}

// curl "localhost:8080/news/list?page=1&rows=1" -b dh_srv.txt
func (this *NewsController) List() {
	//获取新闻类信息列表
	p := P{}
	search := this.GetString("search") //查询字段
	if !IsEmpty(search) {
		p["$or"] = []P{P{"news_class": ToInt(search)}, P{"tab2": MgoLike(search)}, P{"tab3": MgoLike(search)}}
	}
	news_class := []int{1, 2, 4, 6, 7, 8, 9, 10, 12}
	if InArra(ToInt(search), news_class) {
		p["old"] = 0 //选择实时数据
	}
	total := D(News).Find(p).Count()
	list := *D(News).Find(p).Sort("-_id").Page(this.PageParam()).All()
	for _, v := range list {
		if v["news_class"] == 9 && v["dh"] == 0 {
			v["title"] = v["name"]
		}
		if v["news_class"] == 7 || v["news_class"] == 8 {
			if v["dh"] == 0 {
				if v["tp"] == "1" {
					//区分作品列表和作者列表的评论数和转载数
					v["number"] = v["comment"]
				} else {
					v["number"] = v["reprint"]
				}
			}
		}
	}
	r := P{}
	r["total"] = total
	r["page"], _ = this.GetInt("page", 1)
	r["list"] = list
	this.EchoJson(r)
}

// curl "localhost:8080/news/remove" -d "ids=[]" dh_srv.txt
func (this *NewsController) Remove() {
	//删除
	ids := this.FormToP("ids")
	id := strings.Split(ToString(ids["ids"]), ",")
	for _, id := range id {
		if IsOid(id) {
			D(News).Remove(P{"_id": ToOid(id)})
		}
	}
	this.EchoJsonOk()
}

// curl "localhost:8080/news/total" -d "work_count=12&media_count=12" dh_srv.txt
func (this *NewsController) Total() {
	//添加或修改作品总数和传播媒体总数
	p := this.FormToP("work_count", "media_count", "week")
	p["media_count"] = ToInt(p["media_count"])
	p["work_count"] = ToInt(p["work_count"])
	p["week"] = ToInt(p["week"])
	if IsOid(this.GetString("id")) {
		//有id则是修改，反之则是添加
		p["_id"] = this.GetOid("id")
		D(Total).Save(&p)
	} else {
		q := P{}
		q["week"] = p["week"]
		q["dh"] = 1
		D(Total).Remove(q)
		p["dh"] = 1
		p["_id"] = NewId()
		D(Total).Add(p)
	}
	this.EchoJsonOk()
}

// curl "localhost:8080/news/get_total" dh_srv.txt
func (this *NewsController) Get_total() {
	//获取作品总数和媒体总数
	total := *D(Total).Find(P{}).Sort("-_id").All()
	this.EchoJsonMsg(total)
}

func (this *NewsController) Select_id() {
	dh := this.GetString("dh")
	p := P{}
	p["news_class"] = 9
	if ToInt(dh) == 0 {
		p["old"] = 0
		p["dh"] = 0
	} else {
		p["dh"] = 1
	}
	list := *D(News).Find(p).Field("id", "name").All()
	list = append(list,map[string]interface{}{"id":"","name":"空"})
	this.EchoJsonMsg(list)
}

func (this *NewsController) Product_id() {
	dh := this.GetString("dh")
	p := P{}
	p["news_class"] = 1
	if ToInt(dh) == 0 {
		p["old"] = 0
		p["dh"] = 0
	} else {
		p["dh"] = 1
	}
	list := *D(News).Find(p).Field("id", "name").All()
	list = append(list,map[string]interface{}{"id":"","name":"空"})
	this.EchoJsonMsg(list)
}
