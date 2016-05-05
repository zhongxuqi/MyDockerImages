package spider

import (
	//"fmt"
	"bytes"
	"strings"
	"github.com/opesun/goquery"
)

func GetJianshuData(keyword string) ([]DataItem, string, error) {
	resp, err := goquery.ParseUrl("http://www.jianshu.com/search?utf8=%E2%9C%93&q=" + keyword)
	if err != nil {
		return nil, "", err
	}
	return ParseJianShuHTML(&resp)
}

func ParseJianShuHTML(resp *goquery.Nodes) ([]DataItem, string, error) {
	htmlItems := resp.Find("ul.list li")
	resItems := make([]DataItem, len(htmlItems))
	for i, htmlnode := range htmlItems.HtmlAll() {
		itemNodes, err := goquery.ParseString(htmlnode)
		if err != nil {
			return nil, "", err
		}
		title := itemNodes.Find("h4.title a")
		if len(title) > 0 {
			var b bytes.Buffer
			text(&b, title[0])
			resItems[i].Title = b.String()
			resItems[i].Link = "http://www.jianshu.com" + title.Attrs("href")[0]
		}
		abstract := itemNodes.Find("p")
		if len(abstract) > 0 {
			var b bytes.Buffer
			text(&b, abstract[0])
			resItems[i].Abstract = strings.Replace(strings.Trim(b.String(), " \n"), "\n", " ", -1)
		}
	}
	nextPage := ""
	nextHtml := resp.Find("div.pagination ul li.next a")
	if len(nextHtml) > 0 {
		nextPage = "http://www.jianshu.com" + nextHtml.Attrs("href")[0]
	}
	return resItems, nextPage, nil
}

