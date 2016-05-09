package spider

import (
	"strings"
  "github.com/PuerkitoBio/goquery"
)

func GetBaiduData(keyword string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument("http://www.baidu.com/s?ie=uft-8&word=" + keyword)
	if err != nil {
		return nil, "", err
	}
	return ParseBaiduHTML(resp)
}

func ParseBaiduUrl(url string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument(url)
	if err != nil {
		return nil, "", err
	}
	return ParseBaiduHTML(resp)
}

func ParseBaiduHTML(resp *goquery.Document) ([]DataItem, string, error) {
	resItems := make([]DataItem, 0)
  resp.Find(".c-container").Each(func(i int, s *goquery.Selection) {
    resItem := DataItem{}
    resItem.Title = strings.Replace(strings.Trim(
			s.Find("h3 a").First().Text(), " \n"), "\n", " ", -1)
    resItem.Link = s.Find("h3 a").First().AttrOr("href", "")
		if len(s.Find(".c-abstract").Nodes) > 0 {
    	resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".c-abstract").Text(), " \n"), "\n", " ", -1)
		} else if len(s.Find(".c-row div p").Nodes) > 0 {
			resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".c-row div p").Text(), " \n"), "\n", " ", -1)
		}
    resItem.Image = s.Find("img").AttrOr("src", "")
    resItems = append(resItems, resItem)
  })
	nextPage := ""
	nextHtml := resp.Find("a.n")
	if len(nextHtml.Nodes) > 0 {
    nextPage = "http://www.baidu.com" + nextHtml.Last().AttrOr("href", "")
	}
	return resItems, nextPage, nil
}

func GetBaiduXueShuData(keyword string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument("http://xueshu.baidu.com/s?ie=uft-8&wd="+ keyword)
	if err != nil {
		return nil, "", err
	}
	return ParseBaiduXueShuHTML(resp)
}

func ParseBaiduXueShuUrl(url string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument(url)
	if err != nil {
		return nil, "", err
	}
	return ParseBaiduXueShuHTML(resp)
}

func ParseBaiduXueShuHTML(resp *goquery.Document) ([]DataItem, string, error) {
	resItems := make([]DataItem, 0)
  resp.Find(".result").Each(func(i int, s *goquery.Selection) {
    resItem := DataItem{}
    resItem.Title = s.Find("div.sc_content h3.t a").First().Text()
    resItem.Link = "http://xueshu.baidu.com" + s.Find("div.sc_content h3.t a").First().AttrOr("href", "")
		if len(s.Find(".c_abstract").Nodes) > 0 {
			resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".c_abstract").Text(), " \n"), "\n", " ", -1)
		}
    resItem.Image = s.Find("img").AttrOr("src", "")
    resItems = append(resItems, resItem)
  })
	nextPage := ""
	nextHtml := resp.Find("a.n")
	if len(nextHtml.Nodes) > 0 {
    nextPage = "http://xueshu.baidu.com" + nextHtml.Last().AttrOr("href", "")
	}
	return resItems, nextPage, nil
}
