package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"./wbrss"
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
	// var uid []byte
	// indexs := UID_REGEX.FindStringSubmatchIndex(r.RequestURI)
	// if indexs == nil {
	// 	fmt.Fprint(w, "Invalid UID")
	// 	return
	// }
	// uid = UID_REGEX.ExpandString(uid,"$uid",r.RequestURI,indexs)
	// the above usage is very heavy, just use the

	// if uid != nil {
	// 	content := wbrss.GetRss(string(uid))
	// 	fmt.Fprint(w, content)
	// 	return;
	// }
	strs := UID_REGEX.FindStringSubmatch(r.RequestURI)
	if len(strs) == 2 {
		println(strs[1])
		content := wbrss.GetRss(string(strs[1]))
		fmt.Fprint(w, content)
		return
	}
	fmt.Fprint(w, "Invalid UID")
}
