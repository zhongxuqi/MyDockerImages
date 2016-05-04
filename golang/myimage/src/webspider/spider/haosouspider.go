package spider

import (
	//"fmt"
	"bytes"
	"strings"
	"github.com/opesun/goquery"
)

func GetHaosouData(keyword string) ([]DataItem, string, error) {
	resp, err := goquery.ParseUrl("http://www.haosou.com/s?ie=utf-8&shb=1&src=360sou_newhome&q=" + keyword)
	if err != nil {
		return nil, "", err
	}
	return ParseHaosouHTML(&resp)
}

func ParseHaosouHTML(resp *goquery.Nodes) ([]DataItem, string, error) {
	htmlItems := resp.Find("li.res-list")
	resItems := make([]DataItem, len(htmlItems))
	for i, htmlnode := range htmlItems.HtmlAll() {
		itemNodes, err := goquery.ParseString(htmlnode)
		if err != nil {
			return nil, "", err
		}
		title := itemNodes.Find("h3.res-title a")
		if len(title) > 0 {
			var b bytes.Buffer
			text(&b, title[0])
			resItems[i].Title = b.String()
			resItems[i].Link = title.Attrs("href")[0]
		}
		abstract := itemNodes.Find("p.res-desc")
		if len(abstract) > 0 {
			var b bytes.Buffer
			text(&b, abstract[0])
			resItems[i].Abstract = strings.Replace(strings.Trim(b.String(), " \n"), "\n", " ", -1)
		}
	}
	nextPage := ""
	nextHtml := resp.Find("a#snext")
	if len(nextHtml) > 0 {
		nextPage = "http://www.haosou.com" + nextHtml.Attrs("href")[0]
	}
	return resItems, nextPage, nil
}
