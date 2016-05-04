package spider

import (
	//"fmt"
	"github.com/opesun/goquery"
	"bytes"
	"strings"
)

func GetCSDNData(keyword string) ([]DataItem, string, error) {
	resp, err := goquery.ParseUrl("http://so.csdn.net/so/search/s.do?t=blog&o=&s=&q=" + keyword )
	if err != nil {
		return nil, "", err
	}
	return ParseCSDNHTML(&resp)
}

func ParseCSDNHTML(resp *goquery.Nodes) ([]DataItem, string, error) {
	htmlItems := resp.Find("dl.search-list")
	resItems := make([]DataItem, len(htmlItems))
	for i, htmlnode := range htmlItems.HtmlAll() {
		itemNodes, err := goquery.ParseString(htmlnode)
		if err != nil {
			return nil, "", err
		}
		title := itemNodes.Find("dt a")
		if len(title) > 0 {
			var b bytes.Buffer
			text(&b, title[0])
			resItems[i].Title = b.String()
			resItems[i].Link = title.Attrs("href")[0]
		}
		abstract := itemNodes.Find("dd.search-detail")
		if len(abstract) > 0 {
			var b bytes.Buffer
			text(&b, abstract[0])
			resItems[i].Abstract = strings.Replace(strings.Trim(b.String(), " \n"), "\n", " ", -1)
		}
	}
	nextPage := ""
	nextHtml := resp.Find("a.btn-next")
	if len(nextHtml) > 0 {
		nextPage = "http://so.csdn.net/so/search/s.do" + nextHtml.Attrs("href")[0]
	}
	return resItems, nextPage, nil
}
