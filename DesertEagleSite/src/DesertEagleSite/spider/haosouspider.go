package spider

import (
	"strings"
  "github.com/PuerkitoBio/goquery"
)

func GetHaosouData(keyword string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument("https://www.so.com/s?ie=utf-8&q=" + keyword)
	if err != nil {
		return nil, "", err
	}
	return ParseHaosouHTML(resp)
}

func ParseHaosouUrl(url string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument(url)
	if err != nil {
		return nil, "", err
	}
	return ParseHaosouHTML(resp)
}

func ParseHaosouHTML(resp *goquery.Document) ([]DataItem, string, error) {
	resItems := make([]DataItem, 0)
  resp.Find("li.res-list").Each(func(i int, s *goquery.Selection) {
    resItem := DataItem{}
    resItem.Title = strings.Replace(strings.Trim(
			s.Find("h3.res-title a").First().Text(), " \n"), "\n", " ", -1)
    resItem.Link = s.Find("h3.res-title a").First().AttrOr("href", "")
		if len(s.Find("p.res-desc").Nodes) > 0 {
    	resItem.Abstract = strings.Replace(strings.Trim(
				s.Find("p.res-desc").Text(), " \n"), "\n", " ", -1)
		}
    resItem.Image = s.Find("img").AttrOr("src", "")
    resItems = append(resItems, resItem)
  })
	nextPage := ""
	nextHtml := resp.Find("a#snext")
	if len(nextHtml.Nodes) > 0 {
    nextPage = "https://www.so.com" + nextHtml.Last().AttrOr("href", "")
	}
	return resItems, nextPage, nil
}
