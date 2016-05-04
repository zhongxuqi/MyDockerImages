package spider

import (
	"fmt"
	"bytes"
	"strings"
	"github.com/opesun/goquery"
)

func GetDoubanData(keyword string) ([]DataItem, string, error) {
	resp, err := goquery.ParseUrl("http://www.douban.com/search?q=" + keyword)
	if err != nil {
		return nil, "", err
	}
	htmlItems := resp.Find("div.result")
	resItems := make([]DataItem, len(htmlItems))
	fmt.Println(len(resItems))
	for i, htmlnode := range htmlItems.HtmlAll() {
		itemNodes, err := goquery.ParseString(htmlnode)
		if err != nil {
			return nil, "", err
		}
		title := itemNodes.Find("div.title h3 a")
		if len(title) > 0 {
			var b bytes.Buffer
			text(&b, title[0])
			resItems[i].Title = b.String()
			resItems[i].Link = title.Attrs("href")[0]
		}
		abstract := itemNodes.Find("div.content p")
		if len(abstract) > 0 {
			var b bytes.Buffer
			text(&b, abstract[0])
			resItems[i].Abstract = strings.Trim(b.String(), " \n")
		} else {
			abstract = itemNodes.Find(".c-row div p")
			if len(abstract) > 0 {
				var b bytes.Buffer
				text(&b, abstract[0])
				resItems[i].Abstract = strings.Trim(b.String(), " \n")
			}
		}
	}
	nextPage := ""
	return resItems, nextPage, nil
}
