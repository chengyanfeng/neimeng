package controller

import (
	. "../datasource"
	. "../def"
	. "../util"
	"strings"
)

type JsonController struct {
	BaseController
}

// curl "localhost:8080/json/save" -d "json={json}&type=1" dh_srv.txt

func (this *JsonController) Save() {
	//修改图形类json数据
	p := this.FormToP("json", "type", "week") //获取图形类字段信息
	p["type"] = ToInt(p["type"])
	if ToInt(p["week"]) > 0 {
		if p["type"] == 12 || p["type"] == 13 {
			p["week"] = ToInt(p["week"]) - 1
		} else {
			p["week"] = ToInt(p["week"])
		}
	}
	if p["type"] == 3 {
		json := *JsonDecode([]byte(ToString(p["json"])))
		p["sort"] = ToInt(json["val"])
	}
	if p["type"] == 12 || p["type"] == 13 {
		json := *JsonDecode([]byte(ToString(p["json"])))
		p["number"] = ToInt(json["number"])
	}
	if p["type"] == 5 || p["type"] == 8 || p["type"] == 11 {
		json := *JsonDecode([]byte(ToString(p["json"])))
		p["number"] = ToInt(json["val"])
	}
	p["ct"] = Timestamp()
	if IsOid(this.GetString("id")) {
		//有id修改，反之添加
		p["_id"] = this.GetOid("id")
		D(Json).Save(&p)
	} else {
		p["dh"] = 1
		p["old"] = 0
		p["_id"] = NewId()
		D(Json).Add(p)
	}
	this.EchoJsonOk()
}

// curl "localhost:8080/json/list?page=1&rows=1" -b dh_srv.txt
func (this *JsonController) List() {
	//图形json列表
	p := P{}
	search, _ := this.GetInt("search")
	if search != 0 {
		//按type字段查询
		p["type"] = search
	}
	tp := []int{1, 3, 4, 5, 6, 8, 11, 12, 13, 14}
	if InArra(search, tp) {
		//只显示实时数据
		p["old"] = 0
	}
	total := D(Json).Find(p).Count()
	list := *D(Json).Find(p).Page(this.PageParam()).All()
	for _, v := range list {
		if v["type"] == 12 || v["type"] == 13 {
			v["week"] = ToInt(v["week"]) + 1
		}
	}
	r := P{}
	r["total"] = total
	r["page"], _ = this.GetInt("page", 1)
	r["list"] = list
	this.EchoJson(r)
}

// curl "localhost:8080/json/remove" -d "ids=[]" dh_srv.txt
func (this *JsonController) Remove() {
	//删除
	ids := this.FormToP("ids")
	id := strings.Split(ToString(ids["ids"]), ",")
	for _, id := range id {
		D(Json).Remove(P{"_id": ToOid(id)})
	}
	this.EchoJsonOk()
}
