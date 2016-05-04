package spider

import (
	"bytes"
	"strings"
	"github.com/opesun/goquery"
)

func GetBaiduData(keyword string) ([]DataItem, string, error) {
	resp, err := goquery.ParseUrl("http://www.baidu.com/s?ie=uft-8&word=" + keyword)
	if err != nil {
		return nil, "", err
	}
	return ParseBaiduHTML(&resp)
}

func ParseBaiduHTML(resp *goquery.Nodes) ([]DataItem, string, error) {
	htmlItems := resp.Find(".c-container")
	resItems := make([]DataItem, len(htmlItems))
	for i, htmlnode := range htmlItems.HtmlAll() {
		itemNodes, err := goquery.ParseString(htmlnode)
		if err != nil {
			return nil, "", err
		}
		title := itemNodes.Find("h3 a")
		if len(title) > 0 {
			var b bytes.Buffer
			text(&b, title[0])
			resItems[i].Title = b.String()
			resItems[i].Link = title.Attrs("href")[0]
		}
		abstract := itemNodes.Find(".c-abstract")
		if len(abstract) > 0 {
			var b bytes.Buffer
			text(&b, abstract[0])
			resItems[i].Abstract = strings.Trim(b.String(), " \n")
		} else {
			abstract = itemNodes.Find(".c-row div p")
			if len(abstract) > 0 {
				var b bytes.Buffer
				text(&b, abstract[0])
				resItems[i].Abstract = strings.Replace(strings.Trim(b.String(), " \n"), "\n", " ", -1)
			}
		}
		image := itemNodes.Find("img")
		if len(image) > 0 {
			resItems[i].Image = image.Attrs("src")[0]
		}
	}
	nextPage := ""
	nextHtml := resp.Find("a.n")
	if len(nextHtml) == 1 {
		nextPage = "http://www.baidu.com" + nextHtml.Attrs("href")[0]
	} else if len(nextHtml) == 2 {
		nextPage = "http://www.baidu.com" + nextHtml.Attrs("href")[1]
	}
	return resItems, nextPage, nil
}

func GetBaiduXueShuData(keyword string) ([]DataItem, string, error) {
	resp, err := goquery.ParseUrl("http://xueshu.baidu.com/s?ie=uft-8&wd="+ keyword)
	if err != nil {
		return nil, "", err
	}
	return ParseBaiduXueShuHTML(&resp)
}

func ParseBaiduXueShuHTML(resp *goquery.Nodes) ([]DataItem, string, error) {
	htmlItems := resp.Find(".result")
	resItems := make([]DataItem, len(htmlItems))
	for i, htmlnode := range htmlItems.HtmlAll() {
		itemNodes, err := goquery.ParseString(htmlnode)
		if err != nil {
			return nil, "", err
		}
		title := itemNodes.Find("div.sc_content h3.t a")
		if len(title) > 0 {
			var b bytes.Buffer
			text(&b, title[0])
			resItems[i].Title = b.String()
			resItems[i].Link = title.Attrs("href")[0]
		}
		abstract := itemNodes.Find(".c_abstract")
		if len(abstract) > 0 {
			var b bytes.Buffer
			text(&b, abstract[0])
			resItems[i].Abstract = strings.Replace(strings.Trim(b.String(), " \n"), "\n", " ", -1)
		}
	}
	nextPage := ""
	nextHtml := resp.Find("a.n")
	if len(nextHtml) > 0 {
		nextPage = "http://www.baidu.com" + nextHtml.Attrs("href")[0]
	}
	return resItems, nextPage, nil
}
