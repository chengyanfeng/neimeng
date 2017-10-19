package util

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/tealeg/xlsx"
	"net/smtp"
	"regexp"
	"testing"
	"time"
	//"reflect"
	. "pc.cn/def"
)

func init() {
	beego.SetLogger("console", "")
	InitCache()
}

func TestReplace(t *testing.T) {
	str := ".abc_test;test,zxcv"
	str = Replace(str, []string{".", "_", ";", ","}, "")
	Debug(str)
}

func TestIsInt(t *testing.T) {
	Debug(IsInt("8"))
	Debug(IsInt("1471318508"))
}

func TestCsv2Json(t *testing.T) {
	csv := ReadFile("C:\\Users\\admin\\Desktop\\dr.csv")
	h, r := Csv2Json(csv, nil)
	Debug(JsonEncode(h))
	Debug(JsonEncode(r))

}

func Test_Pathinfo(t *testing.T) {
	p := Pathinfo("upload/test.png")
	Debug(p)
	p = Pathinfo("\\/test.png")
	Debug(p)
	p = Pathinfo("test.png")
	Debug(p)
}

func Test_S(t *testing.T) {
	a := 123
	S("a", a)
	Debug("1", S("a"))
	S("a", a, 2)
	time.Sleep(time.Duration(1) * time.Second)
	Debug("3", S("a"))
	time.Sleep(time.Duration(1) * time.Second)
	Debug("5", S("a"))
}

func Test_HttpGet(t *testing.T) {
	str := "select count(*) from filebeat* where city LIKE 'Bei%'"
	str, err := UrlEncoded(str)
	if err != nil {
		Error(err)
	}
	url := fmt.Sprintf("http://estest.datahunter.cn/_sql?sql=%v&format=csv", str)
	Debug(url)
	str = HttpGet(url, nil, nil)
	Debug("HttpGet", str)
}

func Test_HttpPost(t *testing.T) {
	sql := "select owner,count(table_name) from all_tables group by owner"
	body, err := HttpPost("http://10.62.22.62:4567/sql", nil,
		&P{"sql": sql,
			"db": "{\"username\":\"dw\",\"password\":\"dw_2016\",\"host\":\"10.62.22.65\",\"name\":\"orcl\"}"})
	Debug(body, err)
}

func Test_HttpPostBody(t *testing.T) {
	json := `{"a":1, "b":"test"}`
	body, err := HttpPostBody("http://www.baidu.com", nil, json)
	Debug(body, err)
}

func Test_Ip2Int(t *testing.T) {
	num := Ip2Int("207.226.142.92")
	if num != 3487731292 {
		t.Error(num)
	}
}

func Test_IsEmpty(t *testing.T) {
	Debug("Test_IsEmpty")
	str := "   "
	Debug(IsEmpty(str))
	str = "123"
	Debug(IsEmpty(str))
	p := P{}
	Debug(IsEmpty(p["test"]))
	Debug(IsEmpty(p), ToString(p))
	Debug(IsEmpty(&p), ToString(&p))
	Debug("Test_IsEmpty done")
}

func Test_P_ToInt(t *testing.T) {
	p := P{"x": "123", "y": 456, "width": "-1"}
	p.ToInt("x", "y", "width")
	x := p["x"].(int)
	if x != 123 {
		t.Error(x)
	}
	Debug(p)
}

func Test_Xls(t *testing.T) {
	excelFileName := "d:/demo/test.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		Error(err)
		return
	}
	list := make([]string, 0)
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			tmp := ""
			for _, cell := range row.Cells {
				str, _ := cell.String()
				tmp += str + ","
			}
			list = append(list, tmp + "\n")
		}
	}
	Debug(len(list), list)
}

func Test_Md5(t *testing.T) {
	Debug(Md5("1"))
	Debug(Md5("1", "2"))
}

func Test_Regex(t *testing.T) {
	r := regexp.MustCompile(IP_REGEX)
	Debug(r.Match([]byte("0.0.0.0")))
	Debug(r.Match([]byte("123.123.123.123")))
	Debug(r.Match([]byte("255.255.255.255")))
	Debug(r.Match([]byte("323.123.123.123")))
	Debug(r.Match([]byte("123.03.123.123")))
	Debug(r.Match([]byte("123.123.123.a")))
	Debug(r.Match([]byte("123.123.123")))
}

func Test_SendMail(t *testing.T) {
	SendMail("support@datahunter.cn",
		"D@tahunter8",
		"smtp.exmail.qq.com:465",
		"jindaodama@qq.com",
		JoinStr("DataHunter注册验证", Timestamp()),
		JoinStr("<a href='#'>点击链接进行身份验证</a>",
			Timestamp()),
		"html")
}

func Test_SendMailTls(t *testing.T) {
	host := "smtp.exmail.qq.com"
	port := 465
	email := "support@datahunter.cn"
	password := "D@tahunter8"
	toEmail := "jindaodama@qq.com"

	header := P{}
	header["From"] = "DataHunter" + "<" + email + ">"
	header["To"] = toEmail
	header["Subject"] = JoinStr("DataHunter注册验证", Timestamp())
	header["Content-Type"] = "text/html; charset=UTF-8"

	body := JoinStr("<a href='#'>点击链接进行身份验证</a>", Timestamp())

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	auth := smtp.PlainAuth(
		"",
		email,
		password,
		host,
	)

	err := SendMailTls(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		email,
		[]string{toEmail},
		[]byte(message),
	)

	if err != nil {
		Error(err)
	}
}

func TestGetCronStr(t *testing.T) {
	Debug(GetCronStr(60))
	Debug(GetCronStr(600))
	Debug(GetCronStr(1800))
	Debug(GetCronStr(6000))
	Debug(GetCronStr(60000))
	Debug(GetCronStr(86400))
}

func TestJsonDecode(t *testing.T) {
	s := []byte(`[{"o":"name;", "n":"姓名", "type":"string"},{"o":"birth", "n":"生日", "type":"date"}]`)
	j, _ := JsonDecodeArray(s)
	Debug(JsonEncode(j))
}

func TestJson(t *testing.T) {
	var p interface{}
	err := json.Unmarshal([]byte("{}"), &p)
	if err != nil {
		Error("JsonDecode", err)
	}
	Debug(p)
	switch p.(type) {
	case map[string]interface{}:
		Debug("right type")
	}
	err = json.Unmarshal([]byte("[{}]"), &p)
	if err != nil {
		Error("JsonDecode", err)
	}
	switch p.(type) {
	case []interface{}:
		Debug("right type array")
		for _, tmp := range p.([]interface{}) {
			Debug("inner", tmp)
		}
	}
	Debug(p)
}

func TestGbk2Utf(t *testing.T) {
	str := ReadFile("C:\\Users\\admin\\Desktop\\gbk.csv")
	Debug(Gbk2Utf(str))
}

func TestXml2Json(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8" ?>
<ProductList>
    <Product>
        <sku>ABC123</sku>
        <quantity>2</quantity>
    </Product>
    <Product>
        <sku>ABC123</sku>
        <quantity>2</quantity>
    </Product>
</ProductList>`
	str, err := Xml2Json(xmlData)
	if err != nil {
		Error(err)
	}
	Debug(str)
}

func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

func TestPointer(t *testing.T) {
	tmp := []P{}
	for i := 0; i < 10; i++ {
		tmp = append(tmp, P{"i": i})
	}
	Debug(JsonEncode(tmp))
	data := []*P{}
	for _, v := range tmp {
		t := v
		Debug(&v, &t)
		data = append(data, &t)
	}
	Debug(JsonEncode(data))
}

func TestExtractFile(t *testing.T) {
	//ExtractFile("C:\\Users\\admin\\Desktop\\业务监控平台", "d:/tmp", ".java")
	ExtractFile("C:\\Users\\admin\\.m2\\repository\\org\\springframework", "d:/tmp/spring", ".jar")
}

func TestIsJson(t *testing.T) {
	Debug(IsJson([]byte("{}")))
	Debug(IsJson([]byte("[{}]")))
	Debug(IsJson([]byte(`[{"a":1, "b":"test"}]`)))
	Debug(IsJson([]byte(`[{"a":1, "b", "test"}]`)))
	Debug(IsJson([]byte("[{]")))
	Debug(IsJson([]byte("{,}")))
}

func TestRegSplit(t *testing.T) {
	Debug(RegSplit("a1b22c333d", "[0-9]+"))
	Debug(RegSplit("a,b,c,sum(a),c(b,d)", "\\(.*?\\)|(,)"))
	Debug(RegSplit("a,b,c,sum(a),c(b,d)", ","))
}

func TestToFields(t *testing.T) {
	fields := ToFields("a,b,c,sum(a) as as,c(b,d) as bd,to_date(y,m) as ym")
	Debug(JsonEncode(fields), len(fields))
	Debug(ToFields("sum(a),c(a)"))
}

func TestToDate(t *testing.T) {
	Debug(ToDate("2016年11月18日"))
	Debug(ToDate("2016/11/18"))
	Debug(ToDate("2016-11-18"))
	Debug(ToDate("2016年11月18日 15:58:00"))
	Debug(ToDate("2016-11-18 15:58:00"))
	Debug(ToDate("15:58:00"))
	Debug(ToDate("15:58"))
	Debug(ToDate("11-18-2016"))
	Debug(ToDate("11-18-16"))
	Debug(ToDate("2016asdf"))
}
