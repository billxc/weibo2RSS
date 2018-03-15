package wbrss

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"regexp"
	"strings"
	"time"
	"strconv"
)

const (
	BASEURL     = `http://service.weibo.com/widget/widget_blog.php?uid=%s`
	WB_TEMPLATE = `<item>
	<title><![CDATA[%s]]></title>
	<description><![CDATA[%s]]></description>
	<pubDate>%s</pubDate>
	<guid>%s</guid>
	<link>%s</link>
	</item>`

	RSS_TEMPLATE = `<rss version="2.0">
	<channel>
	<title>%s的微博</title>
	<link>http://weibo.com/%s</link>
	<description>%s的微博RSS，使用 Weibo2RSS(https://github.com/billxc/weibo2RSS) 构建</description>
	<language>zh-cn</language>
	<lastBuildDate>%s</lastBuildDate>
	<ttl>300</ttl>
	%s
	</channel>
	</rss>
	`
)

func GetRss(uid string) string {
	doc, err := goquery.NewDocument(fmt.Sprintf(BASEURL, uid))
	if err != nil {
		log.Fatal(err)
	}

	var wbs []weibo
	items := doc.Find(".wgtCell")
	items.Each(func(_ int, item *goquery.Selection) {
		titleEle := item.Find(".wgtCell_txt")
		//wb.title = titleEle.text().replace(/^\s+|\s+$/g, '');
		wb := weibo{}
		wb.title = getTitle(titleEle.Text())
		html, _ := titleEle.Html()
		wb.description = getDesp(html)
		html, _ = item.Find(".link_d").Html()
		wb.pubDate = getPubDate(html)
		href, exists := item.Find(".wgtCell_tm a").Attr("href")
		if exists {
			wb.link = href
		} else {
			wb.link = ""
		}

		wbs = append(wbs, wb)
	})

	name := doc.Find(".userNm").Text()
	return fmt.Sprintf(RSS_TEMPLATE, name, uid, name, time.Now().String(), formatWbList(wbs))
}


func getDesp(input string) string {
	input = strings.TrimSpace(input)
	input = strings.Replace(input, "thumbnail", "large", -1)
	return input
}

func getTitle(input string) string {
	input = strings.TrimSpace(input)
	if len([]rune(input)) > 24 {
		input = string([]rune(input)[0: 24]) + "..."
	}
	return input
}

func getPubDate(input string) string {
	now := time.Now()

	minutesAgo, _ := regexp.Compile(`(\d+)分钟前`)
	hourAndMinute, _ := regexp.Compile(`今天 (\d+):(\d+)`)
	mdhm, _ := regexp.Compile(`(\d+)月(\d+)日 (\d+):(\d+)`)

	if m := minutesAgo.FindString(input); m != "" {
		if i, err := strconv.Atoi(m); err == nil {
			return now.Add(time.Duration(int64(-1*i) * int64(time.Minute))).String()
		}
	}

	if ms := hourAndMinute.FindAllString(input, -1); ms != nil {
		if len(ms) >= 2 {
			hour, he := strconv.Atoi(ms[0])
			minute, me := strconv.Atoi(ms[1])
			if he == nil && me == nil {
				return time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, time.UTC).String()
			}
		}
	}

	if ms := mdhm.FindAllString(input, -1); ms != nil {
		if len(ms) >= 4 {
			month, moe := strconv.Atoi(ms[1])
			day, de := strconv.Atoi(ms[1])
			hour, he := strconv.Atoi(ms[2])
			minute, me := strconv.Atoi(ms[3])
			if he == nil && me == nil && moe == nil && de == nil {
				return time.Date(now.Year(), time.Month(month), day, hour, minute, 0, 0, time.UTC).String()
			}
		}
	}
	return input
}

type weibo struct {
	title, description, pubDate, link string
}

func (wb weibo) String() string {
	return fmt.Sprintf(WB_TEMPLATE, wb.title, wb.description, wb.pubDate, wb.link, wb.link)
}

func formatWbList(wbList []weibo) string {
	html := ""
	for _, wb := range wbList {
		html += wb.String()
	}
	return html
}
