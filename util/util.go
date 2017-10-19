package util

import (
	"bytes"
	"code.google.com/p/mahonia"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/httplib"
	"github.com/clbanning/mxj"
	"gopkg.in/mgo.v2/bson"
	"hash"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"path/filepath"
	. "../def"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var localCache cache.Cache

func InitCache() {
	c, err := cache.NewCache("memory", `{"interval":60}`)
	//c, err := cache.NewCache("file", `{"CachePath":"./dhcache","FileSuffix":".cache","DirectoryLevel":2,"EmbedExpiry":120}`)
	if err != nil {
		Error(err)
	} else {
		localCache = c
	}
}

type P map[string]interface{}

func (p *P) ToInt(s ...string) {
	for _, k := range s {
		v := ToString((*p)[k])
		(*p)[k] = ToInt(v)
	}
}

func (p *P) Like(s ...string) {
	for _, k := range s {
		v := ToString((*p)[k])
		if v != "" {
			(*p)[k] = &bson.RegEx{Pattern: v, Options: "i"}
		}
	}
}

func (p *P) ToP(s ...string) (r P) {
	for _, k := range s {
		v := ToString((*p)[k])
		r = *JsonDecode([]byte(v))
		(*p)[k] = r
		Debug("ToP", k, (*p)[k])
	}
	return
}

func (p *P) Get(k string, def interface{}) interface{} {
	r := (*p)[k]
	if r == nil {
		r = def
	}
	return r
}

func ToInt(s interface{}, default_v ...int) int {
	i, e := strconv.Atoi(ToString(s))
	if e != nil && len(default_v) > 0 {
		return default_v[0]
	}
	return i
}

func ToInt64(s interface{}, default_v ...int64) int64 {
	switch s.(type) {
	case int64:
		return s.(int64)
	case int:
		return int64(s.(int))
	case float64:
		return int64(s.(float64))
	}
	i64, e := strconv.ParseInt(ToString(s), 10, 64)
	if e != nil && len(default_v) > 0 {
		return default_v[0]
	}
	return i64
}

func ToFloat(s interface{}, default_v ...float64) float64 {
	f64, e := strconv.ParseFloat(ToString(s), 64)
	if e != nil && len(default_v) > 0 {
		return default_v[0]
	}
	return f64
}

func IsInt(s interface{}) bool {
	_, e := strconv.ParseInt(ToString(s), 10, 64)
	if e != nil {
		return false
	}
	return true
}

func IsFloat(s interface{}) bool {
	_, e := strconv.ParseFloat(ToString(s), 64)
	if e != nil {
		return false
	}
	return true
}

func Md5(s ...interface{}) (r string) {
	return Hash("md5", s...)
}

func Hash(algorithm string, s ...interface{}) (r string) {
	var h hash.Hash
	switch algorithm {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	case "sha2", "sha256":
		h = sha256.New()
	}
	for _, value := range s {
		switch value.(type) {
		case []byte:
			h.Write(value.([]byte))
		default:
			h.Write([]byte(ToString(value)))
		}
	}
	r = hex.EncodeToString(h.Sum(nil))
	return
}

func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func Base64Decode(s string) []byte {
	r, e := base64.StdEncoding.DecodeString(s)
	if e != nil {
		Error(e)
	}
	return r
}

func Timestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func DateTimeStr() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

func ToDate(s string) (t time.Time, e error) {
	fmt := []string{"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006/01/02 15:04:05",
		"15:04:05",
		"15:04",
		"2006/01/02",
		"2006-01-02",
		"01-02-2006",
		"01-02-06",
		"2006年01月02日 15:04:05",
		"2006年01月02日"}
	for _, f := range fmt {
		t, e = time.Parse(f, s)
		if e == nil {
			return t, e
		}
	}
	return time.Now(), e
}

func InArray(s string, a []string) bool {
	for _, x := range a {
		if x == s {
			return true
		}
	}
	return false
}

func InArra(s int, a []int) bool {
	for _, x := range a {
		if x == s {
			return true
		}
	}
	return false
}

func StartsWith(s string, a ...string) bool {
	for _, x := range a {
		if strings.HasPrefix(s, x) {
			return true
		}
	}
	return false
}

func Unset(p P, keys ...string) {
	for _, x := range keys {
		delete(p, x)
	}
}

func ReadFile(path string) string {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		Error(err)
	}
	return string(c)
}

func ReadFileBytes(path string) []byte {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		Error(err)
	}
	return c
}

func WriteFile(path string, body []byte) {
	err := ioutil.WriteFile(path, body, os.ModeAppend)
	if err != nil {
		Error(err)
	}
}

func ReadLine(path string) []string {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		Error(err)
	}
	if len(c) > 0 {
		return strings.Split(string(c), "\n")
	} else {
		return nil
	}
}

func Rand(start int, end int) int {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(end)
	if r < start {
		r = start + rand.Intn(end-start)
	}
	//time.Sleep(1 * time.Nanosecond)
	return r
}

func JsonDecode(b []byte) (p *P) {
	p = &P{}
	err := json.Unmarshal(b, p)
	if err != nil {
		Error("JsonDecode", string(b), err)
	}
	return
}

func JsonEncode(v interface{}) (r string) {
	b, err := json.Marshal(v)
	if err != nil {
		Error(err)
	}
	r = string(b)
	return
}

func IsJson(b []byte) bool {
	var j json.RawMessage
	return json.Unmarshal(b, &j) == nil
}

func JsonDecodeArray(b []byte) (p []*P, e error) {
	p = []*P{}
	e = json.Unmarshal(b, &p)
	if e != nil {
		Error("JsonDecodeArray", e)
	}
	return
}

func JsonDecodeArrays(b []byte) (p *[]P) {
	p = &[]P{}
	json.Unmarshal(b, p)
	return
}

func JoinStr(val ...interface{}) (r string) {
	for _, v := range val {
		r += ToString(v)
	}
	return
}

func Replace(src string, find []string, r string) string {
	for _, v := range find {
		src = strings.Replace(src, v, r, -1)
	}
	return src
}

func Pathinfo(url string) P {
	p := P{}
	url = strings.Replace(url, "\\", "/", -1)
	if strings.Index(url, "/") < 0 {
		url = JoinStr("./", url)
	}
	re := regexp.MustCompile("(.*)/([^/]*)\\.([^.]*)")
	match := re.FindAllStringSubmatch(url, -1)
	if len(match) > 0 {
		m0 := match[0]
		fmt.Println(m0)
		if len(m0) == 4 {
			p["basename"] = m0[0]
			p["dirname"] = m0[1]
			p["filename"] = m0[2]
			p["extension"] = strings.ToLower(m0[3])
		}
	}
	return p
}

func HttpGet(url string, header ...*P) (body string) {
	body = string(HttpGetBytes(url, header...))
	return
}

func HttpGetBytes(url string, header ...*P) (body []byte) {
	req := httplib.Get(url)
	if len(header) > 0 {
		if header[0] != nil {
			for k, v := range *header[0] {
				req.Header(ToString(k), ToString(v))
			}
		}
		if header[1] != nil {
			for k, v := range *header[1] {
				req.Param(ToString(k), ToString(v))
			}
		}
	}
	body, err := req.Bytes()
	if err != nil {
		Error("HttpGetBytes", err)
	}
	return
}

func HttpPost(url string, header *P, param *P) (body string, err error) {
	req := httplib.Post(url)
	req.SetTimeout(time.Duration(30*time.Second), time.Duration(30*time.Second))
	if header != nil {
		for k, v := range *header {
			req.Header(ToString(k), ToString(v))
		}
	}
	if param != nil {
		for k, v := range *param {
			req.Param(ToString(k), ToString(v))
		}
	}
	body, err = req.String()
	if err != nil {
		Error(err)
	}
	return
}

func HttpPostBody(url string, header *P, body string) (r string, err error) {
	req := httplib.Post(url)
	req.SetTimeout(time.Duration(30*time.Second), time.Duration(30*time.Second))
	if header != nil {
		for k, v := range *header {
			req.Header(ToString(k), ToString(v))
		}
	}
	req.Body(body)
	r, err = req.String()
	if err != nil {
		Error(err)
	}
	return
}

func HttpDelete(url string, header ...*P) (body []byte) {
	req := httplib.Delete(url)
	if len(header) > 0 {
		for k, v := range *header[0] {
			fmt.Println(ToString(k), ToString(v))
			req.Header(ToString(k), ToString(v))
		}
	}
	body, err := req.Bytes()
	if err != nil {
		Error(err)
	}
	return
}

func ToString(v interface{}) string {
	if v != nil {
		switch v.(type) {
		case bson.ObjectId:
			return v.(bson.ObjectId).Hex()
		case []byte:
			return string(v.([]byte))
		case *P, P:
			var p P
			switch v.(type) {
			case *P:
				if v.(*P) != nil {
					p = *v.(*P)
				}
			case P:
				p = v.(P)
			}
			var keys []string
			for k := range p {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			r := "P{"
			for _, k := range keys {
				r = JoinStr(r, k, ":", p[k], " ")
			}
			r = JoinStr(r, "}")
			return r
		case int64:
			return strconv.FormatInt(v.(int64), 10)
		default:
			return fmt.Sprintf("%v", v)
		}
	}
	return ""
}

func ToP(v interface{}) P {
	if v != nil {
		switch v.(type) {
		case P:
			return v.(P)
		case map[string]interface{}:
			return v.(map[string]interface{})
		}
	}
	return P{}
}

func ToStrings(v interface{}) []string {
	strs := []string{}
	if v != nil {
		switch v.(type) {
		case []interface{}:
			for _, i := range v.([]interface{}) {
				strs = append(strs, ToString(i))
			}
		}
	}
	return strs
}

func AsOids(v interface{}) []bson.ObjectId {
	oids := []bson.ObjectId{}
	if v != nil {
		switch v.(type) {
		case []interface{}:
			for _, i := range v.([]interface{}) {
				oids = append(oids, i.(bson.ObjectId))
			}
		}
	}
	return oids
}

// 记录debug信息
func Debug(v ...interface{}) {
	beego.Debug(v)
}

// 记录err信息
func Error(v ...interface{}) {
	beego.Error(v)
}

func IsEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	switch v.(type) {
	case P:
		return len(v.(P)) == 0
	}
	return ToString(v) == ""
}

func Trim(str string) string {
	return strings.TrimSpace(str)
}

func Ip2Int(ip string) int64 {
	sec := strings.Split(ip, ".")
	if len(sec) == 4 {
		return int64(ToInt(sec[0]))<<24 + int64(ToInt(sec[1]))<<16 + int64(ToInt(sec[2]))<<8 + int64(ToInt(sec[3]))
	}
	return 0
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func Csv2Json(csv string, filter []string, th ...int) (head []*P, r []*P) {
	lines := strings.Split(csv, "\n")
	head = []*P{}
	if len(lines) > 0 {
		r = []*P{}
		cols := ToFields(Trim(lines[0]))
		for _, col := range cols {
			if filter != nil {
				col = Replace(col, PUNCTUATION, "")
			}
			p := P{"o": col}
			head = append(head, &p)
		}
		first := 1
		if len(th) > 0 && th[0] < 1 {
			first = 0
			for i, v := range head {
				(*v)["o"] = JoinStr("C", i)
			}
		}
		for _, line := range lines[first:] {
			p := P{}
			row := strings.Split(Trim(line), ",")
			if len(row) >= len(head) {
				for i, v := range head {
					row[i] = Trim(row[i])
					k := ToString((*v)["o"])
					if IsInt(row[i]) {
						p[k] = ToInt64(row[i])
						if (*v)["type"] == nil {
							(*v)["type"] = "long"
						}
					} else if IsFloat(row[i]) {
						p[k] = ToFloat(row[i])
						if (*v)["type"] == nil {
							(*v)["type"] = "float"
						}
					} else {
						p[k] = row[i]
						if (*v)["type"] == nil {
							(*v)["type"] = "string"
						}
					}
				}
				r = append(r, &p)
			}
		}
	}
	return
}

func Xml2Json(src string) (s string, err error) {
	m, err := mxj.NewMapXml([]byte(src))
	return JsonEncode(m), err
}

func SendMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	Error(err)
	return err
}

func SendMailTls(addr string, auth smtp.Auth, from string, to []string, msg []byte) (err error) {

	c, err := func(addr string) (*smtp.Client, error) {
		conn, err := tls.Dial("tcp", addr, nil)
		if err != nil {
			Error("SendMail", err)
			return nil, err
		}
		//分解主机端口字符串
		host, _, _ := net.SplitHostPort(addr)
		return smtp.NewClient(conn, host)
	}(addr)
	//create smtp client
	//c, err := dial(addr)
	if err != nil {
		Error("SendMail", err)
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				Error("SendMail", err)
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}

func Mail(to string, subject string, body string) {
	if IsEmpty(to) || IsEmpty(subject) || IsEmpty(body) {
		Error("SendMail", to, subject, body)
		return
	}
	host := "smtp.exmail.qq.com"
	port := 465
	email := "support@datahunter.cn"
	password := "D@tahunter8"

	header := P{}
	header["From"] = "DataHunter" + "<" + email + ">"
	header["To"] = to
	header["Subject"] = subject
	header["Content-Type"] = "text/html; charset=UTF-8"

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
		[]string{to},
		[]byte(message),
	)

	if err != nil {
		Error(err)
	}
}

func UrlEncoded(str string) (string, error) {
	str = strings.Replace(str, "%", "%25", -1)
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func GetHostname() (host string) {
	host = os.Getenv("host")
	if IsEmpty(host) {
		host = "www.datahunter.cn"
	}
	Debug("host", host)
	return host
}

func GetCronStr(sec int) string {
	ss := sec % 60
	ii := sec / 60
	hh := sec / 3600
	if ii == 0 && hh == 0 {
		return fmt.Sprintf("0/%v * * * * *", sec)
	} else if ii > 0 && hh == 0 {
		return fmt.Sprintf("%v */%v * * * *", ss, ii)
	} else if hh > 0 {
		return fmt.Sprintf("%v %v */%v * * *", ss, ii%60, hh)
	}
	return "0/60 * * * * *"
}

func Gbk2Utf(str string) string {
	enc := mahonia.NewDecoder("gbk")
	return enc.ConvertString(str)
}

func RenderTpl(tpl string, data interface{}) string {
	var bb bytes.Buffer
	//t, err := template.ParseFiles(tpl)
	t, err := template.New(Md5(tpl)).Parse(tpl)
	if err != nil {
		Error(err)
	}
	t.Execute(&bb, data)
	return bb.String()
}

func Mkdir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func AddInOid(oids *[]bson.ObjectId, nid bson.ObjectId) {
	for _, oid := range *oids {
		if oid.Hex() == nid.Hex() {
			return
		}
	}
	*oids = append(*oids, nid)
	return
}

// 缓存接口，存 S("key", value)，取 S("key")
func S(key string, p ...interface{}) (v interface{}) {
	if len(p) == 0 {
		return localCache.Get(key)
	} else {
		if len(p) == 2 {
			var ttl int64
			switch p[1].(type) {
			case int:
				ttl = int64(p[1].(int))
			case int64:
				ttl = p[1].(int64)
			}
			localCache.Put(key, p[0], time.Duration(ttl)*time.Second)
		} else if len(p) == 1 {
			localCache.Put(key, p[0], time.Duration(1e9)*time.Second)
		}
		return p[0]
	}
}

func ExtractFile(path string, target string, ext string) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		Debug(path)
		//if !f.IsDir() {
		if strings.HasSuffix(f.Name(), ext) {
			Copy(path, target+"/"+f.Name())
		}
		//}
		return nil
	})
	Debug("filepath.Walk() %v\n", err)
}

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}
	return cerr
}

func RegSplit(text string, delimeter string) []string {
	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	laststart := 0
	result := make([]string, len(indexes)+1)
	for i, element := range indexes {
		result[i] = text[laststart:element[0]]
		laststart = element[1]
	}
	result[len(indexes)] = text[laststart:len(text)]
	return result
}

func ToFields(s string) (r []string) {
	tmp := strings.Split(s, ",")
	r = []string{}
	state := ""
	seg := ""
	for i, v := range tmp {
		if strings.Index(v, "(") > -1 && strings.Index(v, ")") < 0 {
			state = "("
			seg = v
		} else if strings.Index(v, "(") < 0 && strings.Index(v, ")") > -1 {
			state = ")"
		}
		if state == "(" {
			seg = JoinStr(seg, ",", tmp[i+1])
		} else if state == ")" {
			r = append(r, seg)
			seg = ""
			state = ""
		} else {
			r = append(r, v)
		}
	}
	return
}

func GenPanic() {
	data := []string{"1"}
	//panic(errors.New("Panic"))
	Debug(data[1])
}

// post 网络请求 ,params 是url.Values类型
func Post(apiURL string, params url.Values) (rs map[string]P, err error) {
	resp, err := http.PostForm(apiURL, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &rs)
	return
}
