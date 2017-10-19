package datasource

import (
	. "pc.cn/util"
	"testing"
)

func TestMongoD(t *testing.T) {
	D("test").Remove(P{"test": "ok"})
	D("test").Remove(P{"test": "ok2"})
	D("test").Remove(P{"test": "1"})
	D("test").Remove(P{"test": "2"})
	test := P{"test": "1", "mt": Timestamp()}
	D("test").Add(test)
	t2 := P{"test": "2", "mt": Timestamp() + 1}
	D("test").Add(t2)
	test = *D("test").Find(P{"test": "ok"}).One()
	Debug(test)
	test["mt"] = Timestamp()
	D("test").Save(&test)
	all := D("test").Find(nil).Sort("_id").All()
	Debug(JsonEncode(all))
	all = D("test").Find(nil).Sort("-_id").All()
	Debug(JsonEncode(all))
}

func TestMongoSql(t *testing.T) {
	r, _ := D("", P{
		"username": "",
		"password": "",
		"host":     "127.0.0.1",
		"name":     "practice",
	}).Sql("select * from form limit 1")
	Debug(JsonEncode(r))
}

func TestMongo_Import(t *testing.T) {
	page := 2
	m := D("", P{
		"username": "",
		"password": "",
		"host":     "127.0.0.1",
		"name":     "dh",
	})
	err := m.Import("user", func(r []*P) {
		Debug(JsonEncode(r), page)
	}, page)
	Debug(err)

}
