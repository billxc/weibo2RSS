package main

import (
	"regexp"
	"./wbrss"
	"net/http"
	"fmt"
	"log"
)

var UID_REGEX, _ = regexp.Compile(`/weibo/rss/(?P<uid>\d+)`)

func main() {

	http.HandleFunc("/weibo/rss/", getRss)   //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func getRss(w http.ResponseWriter, r *http.Request) {
	var uid []byte
	indexs := UID_REGEX.FindAllStringSubmatchIndex(r.RequestURI,-1)
	uid = UID_REGEX.ExpandString(uid,"$uid",r.RequestURI,indexs[0])

	if uid != nil {
		content := wbrss.GetRss(string(uid))
		fmt.Fprint(w, content)
	} else {
		fmt.Fprint(w, "Invalid UID")
	}
}
