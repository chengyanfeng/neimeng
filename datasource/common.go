package datasource

import (
	"strings"
)


func FilterAs(o string) (n string) {
	tmp := strings.Split(o, " as ")
	if len(tmp) > 1 {
		n = tmp[0]
	} else {
		n = o
	}
	return
}
