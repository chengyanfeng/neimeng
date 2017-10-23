package controller

import (
	"github.com/astaxie/beego"
	. "neimeng/datasource"
	. "neimeng/def"
	. "neimeng/util"
	"time"
	"strings"
	"fmt"

)

type PupController struct {
	BaseController
}

func New_list() {
	Debug("腾讯新闻接口")
	cat_url := beego.AppConfig.String("news_cat")
	cat := cat_url + "?token=2b69b7139ad5ed2a"                 //获取新闻分类
	//cat := "http://pubopinion.hubpd.com/dailyPaper/index.php/Home/Industry/getAll?token=2b69b7139ad5ed2a" //获取新闻分类
	header := P{}                                              //定义header对象
	header["Content-Type"] = "application/json; charset=utf-8" //为header对象赋值
	cats, err := HttpPost(cat, &header, nil)                   //发起请求，携带header和参数
	jd_cat := *JsonDecode([]byte(cats))                        //将返回的数据先转换成字节型再转换成对象类型
	slice_cat, err_cat := jd_cat["data"].([]interface{})       //将jd_cat["data"]转换成数组，err_cat为布尔型，为真表示转换成功反之转换失败
	//打印接口状态日志
	if err != nil {
		Debug("访问接口失败")
	} else {
		Debug("访问接口成功")
	}
	if len(slice_cat) == 0 {
		Debug("返回数据为空")
	} else {
		Debug("返回数据正常")
	}
	news_url := beego.AppConfig.String("news")
	if err_cat {
		for _, o := range slice_cat {
			//遍历数组
			oo, ok1 := o.(map[string]interface{}) //将数组中每个子集转换为对象
			if ok1 {
				url := news_url + "?token=2b69b7139ad5ed2a&strade=all&flag=3&ftrade=" + ToString(oo["enname"])
				url1 := news_url + "?token=2b69b7139ad5ed2a&strade=all&flag=2&ftrade=" + ToString(oo["enname"])
				url2 := news_url + "?token=2b69b7139ad5ed2a&strade=all&flag=1&ftrade=" + ToString(oo["enname"])
				//url := "http://pubopinion.hubpd.com/dailyPaper/index.php/Home/Industry/getTopicList?token=2b69b7139ad5ed2a&strade=all&flag=3&ftrade=" + ToString(oo["enname"])
				//url1 := "http://pubopinion.hubpd.com/dailyPaper/index.php/Home/Industry/getTopicList?token=2b69b7139ad5ed2a&strade=all&flag=2&ftrade=" + ToString(oo["enname"])
				//url2 := "http://pubopinion.hubpd.com/dailyPaper/index.php/Home/Industry/getTopicList?token=2b69b7139ad5ed2a&strade=all&flag=1&ftrade=" + ToString(oo["enname"])
				urls := []string{url2, url1, url}
				for i := 1; i < 4; i++ {
					//分别获取某一行业下的日周月新闻列表数据
					r, _ := HttpPost(urls[i - 1], &header, nil) //发起请求
					jd := *JsonDecode([]byte(r))              //转换返回的数据
					slices, err := jd["data"].([]interface{})
					if err {
						for _, vv := range slices {
							v, ok := vv.(map[string]interface{})
							if ok {
								p := P{} //定义对象
								for _, o := range slice_cat {
									oo, ok1 := o.(map[string]interface{})
									if ok1 {
										if ToString(v["category"]) == ToString(oo["enname"]) {
											p["tab3"] = ToString(oo["chname"]) //对照英文类别名称赋值中文类别名称
										}
									}
								}
								//为对象赋值，最后插入数据
								p["news_class"] = 4
								p["tab2"] = "互联网"
								p["title"] = v["name"]
								//拼接url
								url := beego.AppConfig.String("pin_url") + ToString(v["flag"]) + "&objId=" + ToString(v["topic_id"]) + "&token=2b69b7139ad5ed2a"
								p["href"] = url
								p["old"] = 0
								p["week"] = i
								p["hot"] = ToInt(v["heat"])
								date, _ := time.Parse("2006-01-02 15:04:05", ToString(v["pubtime"])) //转换日期格式
								p["date"] = date.Unix() - 8 * 3600                                     //日期转时间戳
								p["ct"] = Timestamp()
								p["dh"] = 0
								p["_id"] = NewId()
								D(News).Add(p) //D函数为实例化数据表，后续进行相应的增删改查，Add为添加操作
								//}
							}
						}
					}
				}
			}
		}
		//}
		//更新新闻分类表
		for _, s := range slice_cat {
			ss, ok1 := s.(map[string]interface{})
			if ok1 {
				count := D(Cat).Find(P{"enname": ToString(ss["enname"])}).Count()
				if count > 0 {

				} else {
					p := P{}
					p["enname"] = ToString(ss["enname"])
					p["chname"] = ToString(ss["chname"])
					D(Cat).Add(p)
				}
			}
		}
	}
}

// 热度指数&热点新闻
func Hot_news() {
	Debug("腾讯热点新闻接口")
	news_url := beego.AppConfig.String("news_hot")
	new_list := news_url + "?token=2b69b7139ad5ed2a"
	//new_list := "http://pubopinion.hubpd.com/dailyPaper/index.php/Home/search/getMainTopic?token=2b69b7139ad5ed2a"
	header := P{}
	header["Content-Type"] = "application/json; charset=utf-8"
	r, error := HttpPost(new_list, &header, nil)
	jd := *JsonDecode([]byte(r))
	slices, err := jd["data"].([]interface{})
	if error != nil {
		Debug("访问接口失败")
	} else {
		Debug("访问接口成功")
	}
	if len(slices) == 0 {
		Debug("返回数据为空")
	} else {
		Debug("返回数据正常")
	}
	if err {
		for _, vv := range slices {
			v, right := vv.(map[string]interface{})
			if right {
				tip, nice := v["tips"].([]interface{})
				if nice {
					for _, oo := range tip {
						o, ok := oo.(map[string]interface{})
						if ok {

							//添加热点新闻
							p := P{}
							p["news_class"] = 6
							p["area"] = v["name"]
							p["title"] = o["title"]
							p["old"] = 0
							//拼接url
							url := o["url"]
							if ToInt(o["flag"]) == 1 {
								url = beego.AppConfig.String("pin_url") + ToString(o["flag"]) + "&objId=" + ToString(o["topicId"]) + "&token=2b69b7139ad5ed2a"
							}
							p["href"] = url
							date, _ := time.Parse("2006-01-02 15:04:05", ToString(o["time"])) //时间格式化
							p["date"] = date.Unix() - 8 * 3600                                  //将日期转换为时间戳
							p["ct"] = Timestamp()
							p["dh"] = 0
							p["_id"] = NewId()
							D(News).Add(p)
						}
					}
					p := P{}
					p["type"] = 3
					p["dh"] = 0
					p["old"] = 0
					p["sort"] = ToInt(v["value"])
					q := P{}
					q["key"] = v["name"]
					q["val"] = v["value"]
					p["json"] = JsonEncode(q)
					p["ct"] = Timestamp()
					p["_id"] = NewId()
					D(Json).Add(p)
				}

			}
		}
	}
}

// 选题列表(全部)
func SelectLists() {
	Debug("选题列表(全部)接口")
	new_list := beego.AppConfig.String("selectlist")
	//new_list := "http://chief.hubpd.com/DYPDNewsCommandWeb/api/querySelectList"
	header := P{}
	param := P{}
	header["Content-Type"] = "application/json; charset=utf-8"
	//param["userId"] = "bfb89395ab46415d87a50dfebb6e5a90"
	r, error := HttpPostBody(new_list, &header, JsonEncode(param))
	jd := *JsonDecode([]byte(r))
	slices, err := jd["records"].([]interface{})
	if error != nil {
		Debug("访问接口失败")
	} else {
		Debug("访问接口成功")
	}
	if len(slices) == 0 {
		Debug("返回数据为空")
	} else {
		Debug("返回数据正常")
	}
	if err {
		for _, vv := range slices {
			v, ok := vv.(map[string]interface{})
			if ok {
				p := P{}
				p["news_class"] = 1
				p["dh"] = 0
				p["id"] = v["id"]
				p["old"] = 0
				p["kid"] = 1
				p["name"] = v["title"]
				p["keys"] = v["keyword"]
				p["head"] = v["charger"]
				p["organization"] = v["department"]
				switch ToInt(v["status"]) {
				case 101:
					p["status"] = 0
				case 102:
					p["status"] = 1
				case 103:
					p["status"] = 2
				default:
					p["status"] = 0
				}
				switch ToInt(v["type"]) {
				case 0:
					p["type"] = "自主"
				case 1:
					p["type"] = "上报"
				default:
					p["type"] = "自主"
				}
				date, _ := time.Parse("2006/01/02 15:04:05", ToString(v["exeStartTime"]) + " 00:00:00")
				p["date"] = date.Unix() - 8 * 3600
				p["href"] = v["detailUrl"]
				p["ct"] = Timestamp()
				p["_id"] = NewId()
				D(News).Add(p)
			}

		}
	}

}

// 作品列表
func Works() {
	i := 1
	url := []string{"http://api.crt.test.hubpd.com/crt/screen/mediatype/origin/count/yesterday",
		"http://api.crt.test.hubpd.com/crt/screen/mediatype/origin/count/week",
		"http://api.crt.test.hubpd.com/crt/screen/mediatype/origin/count/month"}
	for _, v := range url {
		r := HttpGet(v)
		jd := *JsonDecode([]byte(r))
		if i == 1 {
			Debug("当天作品总数接口")
		} else if i == 2 {
			Debug("本周作品总数接口")
		} else {
			Debug("本月作品总数接口")
		}
		if r == "" {
			Debug("访问接口失败")
		} else {
			Debug("访问接口成功")
		}
		if len(jd) == 0 {
			Debug("返回数据为空")
		} else {
			Debug("返回数据正常")
		}
		total_id := *D(Total).Find(P{"week": i, "dh": 0}).One()
		p := P{}
		p["work_count"] = jd["data"]
		if IsOid(ToString(total_id["_id"])) {
			p["_id"] = ToOid(ToString(total_id["_id"]))
			D(Total).Save(&p)
		} else {
			p["week"] = i
			p["dh"] = 0
			p["_id"] = NewId()
			D(Total).Add(p)
		}
		i++
	}
	i = 1
	url = []string{"http://api.crt.test.hubpd.com/crt/screen/reprinted/rank/yesterday",
		"http://api.crt.test.hubpd.com/crt/screen/reprinted/rank/week",
		"http://api.crt.test.hubpd.com/crt/screen/reprinted/rank/month"}
	for _, v := range url {
		r := HttpGet(v)
		jd := *JsonDecode([]byte(r))
		slice, ok := jd["data"].([]interface{})
		if i == 1 {
			Debug("当天作品转载列表接口")
		} else if i == 2 {
			Debug("本周作品转载列表接口")
		} else {
			Debug("本月作品转载列表接口")
		}
		if r == "" {
			Debug("访问接口失败")
		} else {
			Debug("访问接口成功")
		}
		if len(slice) == 0 {
			Debug("返回数据为空")
		} else {
			Debug("返回数据正常")
		}
		if ok {
			for _, vv := range slice {
				v := vv.(map[string]interface{})
				p := P{}
				p["news_class"] = 7
				p["name"] = v["name"]
				p["href"] = v["url"]
				p["reprint"] = ToInt(v["value"])
				p["week"] = i
				p["old"] = 0
				p["tp"] = "2"
				p["ct"] = Timestamp()
				p["dh"] = 0
				p["_id"] = NewId()
				D(News).Add(p)
			}
		}
		i++
	}

}

// 作者列表
func Author() {
	i := 1
	url := []string{"http://api.crt.test.hubpd.com/crt/screen/author/reprinted/rank/yesterday",
		"http://api.crt.test.hubpd.com/crt/screen/author/reprinted/rank/week",
		"http://api.crt.test.hubpd.com/crt/screen/author/reprinted/rank/month"}
	for _, v := range url {
		r := HttpGet(v)
		jd := *JsonDecode([]byte(r))
		slice, ok := jd["data"].([]interface{})
		if i == 1 {
			Debug("当天作者转载列表接口")
		} else if i == 2 {
			Debug("本周作者转载列表接口")
		} else {
			Debug("本月作者转载列表接口")
		}
		if r == "" {
			Debug("访问接口失败")
		} else {
			Debug("访问接口成功")
		}
		if len(slice) == 0 {
			Debug("返回数据为空")
		} else {
			Debug("返回数据正常")
		}
		if ok {
			for _, vv := range slice {
				v := vv.(map[string]interface{})
				p := P{}
				p["news_class"] = 8
				p["name"] = v["name"]
				p["reprint"] = ToInt(v["value"])
				p["periodical"] = v["adminName"]
				p["week"] = i
				p["old"] = 0
				p["tp"] = "2"
				p["ct"] = Timestamp()
				p["dh"] = 0
				p["_id"] = NewId()
				D(News).Add(p)
			}
		}
		i++
	}
}

// 30天
func Day30() {
	Debug("文章30天")
	var array []interface{}
	r := HttpGet("http://api.crt.test.hubpd.com/crt/screen/origin/count/d30")
	jd := *JsonDecode([]byte(r))
	slices, err := jd["data"].([]interface{})
	if r == "" {
		Debug("访问接口失败")
	} else {
		Debug("访问接口成功")
	}
	if len(slices) == 0 {
		Debug("返回数据为空")
	} else {
		Debug("返回数据正常")
	}
	if err {
		for _, vv := range slices {
			v, ok := vv.(map[string]interface{})
			if ok {
				p := P{}
				p["number"] = v["value"]
				day := beego.Substr(ToString(v["name"]), 6, 2)
				p["date"] = day
				array = append(array, p)
			}
		}
	}
	q := P{}
	q["type"] = 14
	q["dh"] = 0
	q["old"] = 0
	q["json"] = JsonEncode(array)
	q["ct"] = Timestamp()
	D(Json).Add(q)
}

// 微博
func Twitter() {
	Debug("微博")
	num := []int{1, 2, 3}
	for _, v := range num {
		//new_list := "http://10.101.67.1:8000/uar/api/getwbdata"
		//new_list := "http://10.101.67.1:8000/uar/api/getwbdata"
		header := P{}
		header["Content-Type"] = "application/json; charset=utf-8"
		param := P{}
		param["week"] = v
		r, error := HttpPostBody(beego.AppConfig.String("twitter"), &header, JsonEncode(param))
		//r, error := HttpPostBody(new_list, &header, JsonEncode(param))
		jd := *JsonDecode([]byte(r))
		slices, err := jd["data"].(map[string]interface{})
		if error != nil {
			Debug("访问接口失败")
		} else {
			Debug("访问接口成功")
		}
		if len(slices) == 0 {
			Debug("返回数据为空")
		} else {
			Debug("返回数据正常")
		}
		if err {
			str := ""
			hot_word, ok := slices["hot_word"].([]interface{})
			if ok {
				for _, vv := range hot_word {
					msg := ToString(JsonEncode(vv))
					msg = strings.Replace(msg, "@", "*", -1)
					str = str + msg + "@"
				}
			}
			p := P{}
			p["news_class"] = 10
			p["old"] = 0
			p["origin"] = 4
			p["hot_word"] = str
			p["widget1"] = ToString(JsonEncode(slices["widget1"]))
			p["widget2"] = ToString(JsonEncode(slices["widget2"]))
			p["week"] = v
			p["ct"] = Timestamp()
			p["dh"] = 0
			p["_id"] = NewId()
			D(News).Add(p)
		}
	}

}

// 客户端
func App() {
	Debug("客户端")
	num := []int{1, 2, 3}
	for _, v := range num {
		new_list := "http://uar.hubpd.com/uar/mongoscreen/getClientHotArticle"
		header := P{}
		header["Content-Type"] = "application/json; charset=utf-8"
		param := P{}
		param["week"] = v
		r, error := HttpPostBody(new_list, &header, JsonEncode(param))
		jd := *JsonDecode([]byte(r))
		slices, err := jd["data"].(map[string]interface{})
		if error != nil {
			Debug("访问接口失败")
		} else {
			Debug("访问接口成功")
		}
		if len(slices) == 0 {
			Debug("返回数据为空")
		} else {
			Debug("返回数据正常")
		}
		if err {
			str := ""
			hot_word, ok := slices["hot_word"].([]interface{})
			if ok {
				for _, vv := range hot_word {
					str = str + ToString(JsonEncode(vv)) + "@"
				}
			}
			p := P{}
			p["news_class"] = 10
			p["old"] = 0
			p["origin"] = 1
			p["hot_word"] = str
			p["widget1"] = ToString(JsonEncode(slices["widget1"]))
			p["widget2"] = ToString(JsonEncode(slices["widget2"]))
			p["week"] = v
			p["ct"] = Timestamp()
			p["dh"] = 0
			p["_id"] = NewId()
			D(News).Add(p)
		}
	}

}

// 微信
func Wx() {
	Debug("微信")
	num := []int{1, 2, 3}
	for _, v := range num {
		new_list := "http://10.101.67.1:8000/uar/api/getwcdata"
		//new_list := "http://10.101.67.1:8000/uar/api/getwcdata"
		header := P{}
		header["Content-Type"] = "application/json; charset=utf-8"
		param := P{}
		param["week"] = v
		//r, error := HttpPostBody(beego.AppConfig.String("wx"), &header, JsonEncode(param))
		r, error := HttpPostBody(new_list, &header, JsonEncode(param))
		jd := *JsonDecode([]byte(r))
		slices, err := jd["data"].(map[string]interface{})
		if error != nil {
			Debug("访问接口失败")
		} else {
			Debug("访问接口成功")
		}
		if len(slices) == 0 {
			Debug("返回数据为空")
		} else {
			Debug("返回数据正常")
		}
		if err {
			str := ""
			hot_word, ok := slices["hot_word"].([]interface{})
			if ok {
				for _, vv := range hot_word {
					str = str + ToString(JsonEncode(vv)) + "@"
				}
			}
			p := P{}
			p["news_class"] = 10
			p["old"] = 0
			p["origin"] = 3
			p["hot_word"] = str
			p["widget1"] = ToString(JsonEncode(slices["widget1"]))
			p["widget2"] = ToString(JsonEncode(slices["widget2"]))
			p["week"] = v
			p["ct"] = Timestamp()
			p["dh"] = 0
			p["_id"] = NewId()
			D(News).Add(p)
		}
	}

}

// 网站
func Web() {
	Debug("网站")
	num := []int{1, 2, 3}
	for _, v := range num {
		//new_list := "http://10.101.67.1:8000/uar/api/getwebdata"
		header := P{}
		header["Content-Type"] = "application/json; charset=utf-8"
		param := P{}
		param["week"] = v
		r, error := HttpPostBody(beego.AppConfig.String("web"), &header, JsonEncode(param))
		//r, error := HttpPostBody(new_list, &header, JsonEncode(param))
		jd := *JsonDecode([]byte(r))
		slices, err := jd["data"].(map[string]interface{})
		if error != nil {
			Debug("访问接口失败")
		} else {
			Debug("访问接口成功")
		}
		if len(slices) == 0 {
			Debug("返回数据为空")
		} else {
			Debug("返回数据正常")
		}
		if err {
			str := ""
			hot_word, ok := slices["hot_word"].([]interface{})
			if ok {
				for _, vv := range hot_word {
					str = str + ToString(JsonEncode(vv)) + "@"
				}
			}
			p := P{}
			p["news_class"] = 10
			p["old"] = 0
			p["origin"] = 2
			p["hot_word"] = str
			p["widget1"] = ToString(JsonEncode(slices["widget1"]))
			p["widget2"] = ToString(JsonEncode(slices["widget2"]))
			p["week"] = v
			p["ct"] = Timestamp()
			p["dh"] = 0
			p["_id"] = NewId()
			D(News).Add(p)
		}
	}

}

func Change_old() {
	news_class := []int{1, 2, 6, 7, 8, 9, 12, 30}
	for i := 0; i < 8; i++ {
		q := P{}
		m := P{}
		q["news_class"] = news_class[i]
		q["dh"] = 0
		q["old"] = 0
		count := D(News).Find(q).Count()
		m["news_class"] = news_class[i]
		m["dh"] = 0
		m["old"] = 1
		count1 := D(News).Find(m).Count()
		if count > 0 && count1 > 0 {
			sel := P{}
			sel["old"] = 1
			sel["dh"] = 0
			sel["news_class"] = news_class[i]
			D(News).Remove(sel)
		}
		p := P{}
		p["dh"] = 0
		p["old"] = 0
		p["news_class"] = news_class[i]
		list := *D(News).Find(p).All()
		for _, v := range list {
			q := P{}
			q["_id"] = v["_id"]
			q["old"] = 1
			D(News).Save(&q)

		}
	}

	tp := []int{1, 3, 4, 5, 6, 7, 8, 10, 11, 12, 13, 14}
	for i := 0; i < 12; i++ {
		q := P{}
		m := P{}
		q["type"] = tp[i]
		q["dh"] = 0
		q["old"] = 0
		count := D(Json).Find(q).Count()
		m["type"] = tp[i]
		m["dh"] = 0
		m["old"] = 1
		count1 := D(Json).Find(m).Count()
		if count > 0 && count1 > 0 {
			sel := P{}
			sel["old"] = 1
			sel["dh"] = 0
			sel["type"] = tp[i]
			D(Json).Remove(sel)
		}
		p := P{}
		p["dh"] = 0
		p["old"] = 0
		p["type"] = tp[i]
		list := *D(Json).Find(p).All()
		for _, v := range list {
			q := P{}
			q["_id"] = v["_id"]
			q["old"] = 1
			D(Json).Save(&q)

		}
	}
	change_old_10()
	change_old_4()
}

func change_old_10() {
	for i := 1; i < 5; i++ {
		q := P{}
		m := P{}
		q["news_class"] = 10
		q["dh"] = 0
		q["old"] = 0
		q["origin"] = i
		count := D(News).Find(q).Count()
		m["news_class"] = 10
		m["dh"] = 0
		m["old"] = 1
		m["origin"] = i
		count1 := D(News).Find(m).Count()
		if count > 0 && count1 > 0 {
			sel := P{}
			sel["old"] = 1
			sel["dh"] = 0
			sel["news_class"] = 10
			sel["origin"] = i
			D(News).Remove(sel)
		}
		p := P{}
		p["dh"] = 0
		p["old"] = 0
		p["news_class"] = 10
		p["origin"] = i
		list := *D(News).Find(p).All()
		for _, v := range list {
			q := P{}
			q["_id"] = v["_id"]
			q["old"] = 1
			D(News).Save(&q)
		}
	}
}

func change_old_4() {
	var tab = []string{"互联网", "新华电讯"}
	for _, v := range tab {
		q := P{}
		m := P{}
		q["tab2"] = v
		m["tab2"] = v
		q["news_class"] = 4
		q["dh"] = 0
		q["old"] = 0
		count := D(News).Find(q).Count()
		m["news_class"] = 4
		m["dh"] = 0
		m["old"] = 1
		count1 := D(News).Find(m).Count()
		if count > 0 && count1 > 0 {
			sel := P{}
			sel["old"] = 1
			sel["dh"] = 0
			sel["news_class"] = 4
			D(News).Remove(sel)
		}
		p := P{}
		p["dh"] = 0
		p["old"] = 0
		p["news_class"] = 4
		list := *D(News).Find(p).All()
		for _, v := range list {
			q := P{}
			q["_id"] = v["_id"]
			q["old"] = 1
			D(News).Save(&q)
		}
	}
}

//获取媒体资源库
func Mediabank(){
	Debug("媒体资源库")
	//优先修改备份字段
	priority:=P{"backups":0}//优先字段
	duplicate:=P{"backups":1}//备份字段
	prioritylist:=*D(Media).Find(priority).All()
	fmt.Println("aaaa")
	fmt.Println(len(prioritylist))
	count:=D(Media).Find(duplicate).Count()
	for i:=0;i<count;i++{
		err:=D(Media).Remove(duplicate)
		if err !=nil {
			fmt.Println("没有备份数据")
			Debug("没有备份数据")
		}
	}
	if len(prioritylist)>0{
		//先删除备份数据
		for k,_:=range prioritylist {
			//把优先数据修改为备份数据
			fmt.Println(k)
			D(Media).Upsert(priority,duplicate)
		}
	}
	var url="http://zycf.northnews.cn:8443/cre/api/json/basicsearch"
	header := P{}
	header["Content-Type"] = "text/plain;charset=UTF-8"
	header["Access-Control-Allow-Headers"]="Content-Type"
	header["Access-Control-Allow-Credentials"]="true"
	conditions:=make([]P,2)
	conditions[0]=P{
		"id":"type",
		"value": "7",
	}
	conditions[0]["operator"]=1
	conditions[1]=P{
		"id": "FOLDERID",
		"value": "PRODUCT_ROOT_FOLDERID",
	}
	dataparam:=P{}
	param := P{}
	param["conditions"]=conditions
	param["start"]=0
	param["limit"]=3
	param["userId"]="a3d9511c1bb94c9ca1e1fa52d54aa1ab"
	param["token"]="ST-53-R4FZlngYZ2eAzdofZd7J-cas"
	param["orderBy"]="lastModify desc"
	param["extendResultFields"]="作者,产品状态"
	data, error := HttpPostBody(url, &header, JsonEncode(param))
	if error != nil {
		fmt.Println("访问接口失败")
		Debug("访问接口失败")
	} else {
		fmt.Println("访问接口成功")
		Debug("访问接口成功")
		md := *JsonDecode([]byte(data))

		mddata, err :=md["itemList"].([]interface{})
		if err{
			fmt.Print("获取成功")
		}
		totalCount:=md["totalCount"]
		dataparam["totalCount"]=totalCount

		for k,v :=range mddata {

			vdata:=v.(map[string]interface{})
			dataparam["创建时间"] =vdata["created"]
			dataparam["提交时间"] =vdata["lastModify"]
			dataparam["名称"] =vdata["creatorName"]
			dataparam["资源详情"] =vdata["detailUrl"]
			dataparam["序号"]=k
			dataparam["backups"]=0
			auth:=vdata["extendAttributes"].([]interface{})
			for k,authv:=range auth{

				audata:=authv.(map[string]interface{})
				if k==0{
					dataparam["作者"]=audata["id"]


				}else {
					dataparam["状态"]=audata["value"]

				}

			}
			err:=D(Media).Add(dataparam)
			if err!=nil{
				D(Media).Remove(priority)
				for k,_:=range priority {
					//相同的循环方式把备份修改回来。
					fmt.Println(k)
					D(Media).Upsert(duplicate, priority)
				}
			}

		}

	}


}

