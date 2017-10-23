package def

var PUNCTUATION []string = []string{".", ";", ",", "(", ")"}

var MODE string = ""

const (
	DbPos string = "dbpos"
	User string = "user"
	Json string = "json"
	News string = "news"
	Total string = "total"
	Cat string = "cat"
	Media string = "media"
)

const (
	IP_REGEX string = "((?:(?:25[0-5]|2[0-4]\\d|((1\\d{2})|([1-9]?\\d)))\\.){3}(?:25[0-5]|2[0-4]\\d|((1\\d{2})|([1-9]?\\d))))"
)

const (
	GENERAL_ERR int = 400
)

