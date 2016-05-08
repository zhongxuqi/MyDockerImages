package spider

import (
	"strings"
  "github.com/PuerkitoBio/goquery"
)

func GetCSDNData(keyword string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument("http://so.csdn.net/so/search/s.do?t=blog&o=&s=&q=" + keyword )
	if err != nil {
		return nil, "", err
	}
	return ParseCSDNHTML(resp)
}

func ParseCSDNUrl(url string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument(url)
	if err != nil {
		return nil, "", err
	}
	return ParseCSDNHTML(resp)
}

func ParseCSDNHTML(resp *goquery.Document) ([]DataItem, string, error) {
	resItems := make([]DataItem, 0)
  resp.Find("dl.search-list").Each(func(i int, s *goquery.Selection) {
    resItem := DataItem{}
    resItem.Title = s.Find("dt a").First().Text()
    resItem.Link = s.Find("dt a").First().AttrOr("href", "")
		if len(s.Find("dd.search-detail").Nodes) > 0 {
    	resItem.Abstract = strings.Replace(strings.Trim(
				s.Find("dd.search-detail").Text(), " \n"), "\n", " ", -1)
		}
    resItem.Image = s.Find("img").AttrOr("src", "")
    resItems = append(resItems, resItem)
  })
	nextPage := ""
	nextHtml := resp.Find("a.btn-next")
	if len(nextHtml.Nodes) > 0 {
    nextPage = "http://so.csdn.net/so/search/s.do" + nextHtml.Last().AttrOr("href", "")
	}
	return resItems, nextPage, nil
}
