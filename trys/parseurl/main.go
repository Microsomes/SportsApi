package main

import (
	"fmt"
	"net/url"
)

func main() {
	url2 := "https://talkfs.com/langauges?lang=tayyab"

	u, _ := url.Parse(url2)

	q, _ := url.ParseQuery(u.RawQuery)

	fmt.Println(q["lang"][0])

}
