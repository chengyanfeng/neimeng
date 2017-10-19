package controller

import (
	"bytes"
	. "../datasource"
	. "../util"
	//"errors"
	"net/url"
	"runtime"
	"strings"
)

const MAX_UPLOAD int64 = 50 * 1024 * 1024

type ApiController struct {
	BaseController
}

// curl localhost:8080/api/upload -F "bin=@D:\\rank1.png" -F "appid=566e218bf7c6d14409000029"
// curl https://api.datahunter.cn/api/upload -F "bin=@D:/git/dh_srv/demo/line.json"
func (this *ApiController) Upload() {
	f, h, err := this.GetFile("bin")
	defer func() {
		if f != nil {
			f.Close()
		}
		if err := recover(); err != nil {
			Error("Upload", err)
		}
	}()

	if err != nil {
		Error("Upload", err)
		this.EchoJsonErr(err.Error())
	} else {
		ext := ToString(Pathinfo(h.Filename)["extension"])
		if !InArray(ext, []string{"png", "jpg", "jpeg", "bmp", "gif", "json", "csv", "xlsx", "txt", "xml"}) {
			this.EchoJsonErr("文件类型不合法")
		}
		var buff bytes.Buffer
		fileSize, _ := buff.ReadFrom(f)
		if fileSize > MAX_UPLOAD {
			this.EchoJsonErr("文件尺寸不得大于", MAX_UPLOAD)
		} else {
			md5 := Md5(buff.Bytes())
			filename := JoinStr(md5, ".", ext)
			exist, _ := FileExists("upload/" + filename)
			if !exist {
				this.SaveToFile("bin", "upload/"+filename)
			} else {
				Debug("File exists, skip")
			}
			r := P{}
			r["url"] = "http://" + GetHostname() + "/upload/" + filename
			r["ext"] = ext
			r["size"] = fileSize
			this.EchoJsonMsg(r)
		}
	}
}

// curl https://api.datahunter.cn/api/uploadb64 -d "id=123&bin=data%3Aimage/png%3Bbase64,iVBORw0KGgoAAAANSUhEUgAAABYAAAAWCAMAAADzapwJAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAMAUExURQAAADVBADdMBztRADpcEEtEAEBMGktXAENSG0ZbFkBdG1JKAFZZAFZaAFheCEtTLElUKUxYKFdHI1VRIldrDEljJk5gJktwIFNiK1liN2hcAGxdAGNgA2BkDWVlC2ZuJW17I3VyPW2EEGSLLmyKJmyJL2yCOXSKJXSBO3uOPXecPHKAQnaFQHuJVIN%2BLJx5OZl2ToiNL5GRN6GXNbKCMLKGPb%2BOMbugNYWRX4eUX5GLV5GQT5SaRJmUUqSFQqiHWqiYS6iVU6ydXLKOWrCRWqyObLOafL6jXqagbragZ7GsdN2eKtqqLt2qKt%2B0Ktu4KtmyM%2BGrJOKrN%2ButMuyxL%2Bq7J%2BS0NuG9Ne%2B2NemzOfC3LPK1OMClSsquQ8W5Tdi2SN%2B%2BTdW%2BWNmxXMy7c9%2B7feW9TO2wRem8V%2Fe2TvOyXvu%2BU%2BO9dO3AKcvFVd7CVMvLadzEbtfRb%2B%2FFX%2F%2FBTfLGW%2FvQXeDKZenHY%2BbGcerGfOfXeOjUc%2FrGY%2FzAYPXOcf%2FTbfLZdP3XcvrSe%2F%2FWff%2FYdOjnfbqkjdDCk9%2Fci%2BbclPHPhf%2FUh%2F%2Fdif%2FVlv%2FYkujigvXhnP%2Fjn%2F%2Foovv%2FovD8tP%2F%2Ftv%2F%2Fw%2F3%2FxAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAALDwhhMAAAEAdFJOU%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwBT9wclAAAACXBIWXMAAA7DAAAOwwHHb6hkAAAAGXRFWHRTb2Z0d2FyZQBQYWludC5ORVQgdjMuNS4xTuc4%2BQAAANpJREFUKFNj%2BI8VMFBRmF9S0MJSl0NCAGom1Gw2bY0p06brqImhCouqqk82nKqpzIIqzKWixNtqpKAljirMKcKeG2csx8yEKszIJ1PdkCglK48sbG3Dwz0hu6cxVprVVhEsAXKJVUd%2B%2F6SU3u760oSyTi9hmPB%2FvcKqrol9mWmRYfHJHnDV%2F%2F87FCQ112Z4R0eUOEIMh3rHrsUnMCokuNgM1SXORX45ob5BWSaowgZ1Af7hqek1pqjCrpXleeYxFU36qMJu9i5AASd3T1ThNgi3XQhVGD3qcMQlAKsSeJeqplMiAAAAAElFTkSuQmCC"
func (this *ApiController) UploadB64() {
	id := this.GetString("id")
	bin := this.GetString("bin")
	if IsEmpty(id) || IsEmpty(bin) {
		this.EchoJsonErr("Incorrect id or bin params")
	} else {
		bin = strings.Replace(bin, "data:image/png;base64,", "", -1)
		filename := JoinStr(id, ".png")
		b := Base64Decode(bin)
		Mkdir("upload")
		WriteFile("upload/"+filename, b)
		r := P{}
		r["url"] = "http://" + GetHostname() + "/upload/" + filename
		r["ext"] = "png"
		r["size"] = len(b)
		this.EchoJsonMsg(r)
	}
}

// localhost:8080/api/uuid
func (this *ApiController) Uuid() {
	this.Echo(NewId().Hex())
}

// localhost:8080/api/ok
func (this *ApiController) Ok() {
	this.EchoJsonMsg(P{"ts": Timestamp()})
}

// localhost:8080/api/panic
func (this *ApiController) Panic() {
	GenPanic()
}

// localhost:8080/api/ip
func (this *ApiController) Ip() {
	this.Ctx.Output.Header("Content-Type", "text/javascript; charset=utf-8")
	this.Echo("window.remote_ip='", this.Ctx.Input.IP(), "';")
}

// localhost:8080/api/test?a=123&b=456
func (this *ApiController) Test() {
	p := P{"head": this.Ctx.Request.Header, "param": this.FormToP(), "os": runtime.GOOS, "arch": runtime.GOARCH}
	Debug(JsonEncode(p))
	this.EchoJsonMsg(p)
}

func (this *ApiController) Tests() {
	new_list := "http://115.28.173.240/home/index/info"
	param := url.Values{}
	//param.Set()
	data, _ := Post(new_list, param)
	this.EchoJsonMsg(data)
}
