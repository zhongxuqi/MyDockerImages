package spider

import (
	"bytes"
	"strings"
	"github.com/opesun/goquery"
)

func GetWikipediaData(keyword string) ([]DataItem, string, error) {
	resp, err := goquery.ParseUrl("https://wuu.wikipedia.org/w/index.php?title=Special:搜索&profile=default&fulltext=Search&search=" + keyword)
	if err != nil {
		return nil, "", err
	}
	return ParseWikipediaHTML(&resp)
}

func ParseWikipediaHTML(resp *goquery.Nodes) ([]DataItem, string, error) {
	htmlItems := resp.Find("ul.mw-search-results li")
	resItems := make([]DataItem, len(htmlItems))
	for i, htmlnode := range htmlItems.HtmlAll() {
		itemNodes, err := goquery.ParseString(htmlnode)
		if err != nil {
			return nil, "", err
		}
		title := itemNodes.Find("div.mw-search-result-heading a")
		if len(title) > 0 {
			var b bytes.Buffer
			text(&b, title[0])
			resItems[i].Title = b.String()
			resItems[i].Link = "http://wc.yooooo.us" + title.Attrs("href")[0]
		}
		abstract := itemNodes.Find("div.searchresult")
		if len(abstract) > 0 {
			var b bytes.Buffer
			text(&b, abstract[0])
			resItems[i].Abstract = strings.Replace(strings.Trim(b.String(), " \n"), "\n", " ", -1)
		}
	}
	nextPage := ""
	nextHtml := resp.Find("p.mw-search-pager-bottom a.mw-nextlink")
	if len(nextHtml) > 0 {
		nextPage = "http://wc.yooooo.us" + nextHtml.Attrs("href")[0]
	}
	return resItems, nextPage, nil
}
