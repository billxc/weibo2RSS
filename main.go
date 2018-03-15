package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"github.com/PuerkitoBio/goquery"

	"log"
)

func main() {
	fmt.Println(getData("3306934123"))
}

const BASEURL = `http://service.weibo.com/widget/widget_blog.php?uid=%s`

func getData(uid string) string {
	resp, err := http.Get(fmt.Sprintf(BASEURL, uid))
	var body []byte
	if err != nil {
		goto fail
	}

	if resp.StatusCode != 200 {
		goto fail
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		goto fail
	}

	return string(body)

fail:
	return "null"
}

func get(uid string) {
	doc, err := goquery.NewDocument(fmt.Sprintf(BASEURL, uid))
	if err != nil {
		log.Fatal(err)
	}

	var wbs []weibo;
	items := doc.Find(".wgtCell")
	items.Each(func(_ int, item *goquery.Selection) {
		titleEle := item.Find(".wgtCell_txt")
		//wb.title = titleEle.text().replace(/^\s+|\s+$/g, '');
		wb := weibo{title: titleEle.Text()}
		if len(wb.title) > 24 {
			wb.title = wb.title[0: 24] + "..."
		}
		wbs = append(wbs, wb)
	})
}

type weibo struct {
	title, description, pubDate, link string
}
