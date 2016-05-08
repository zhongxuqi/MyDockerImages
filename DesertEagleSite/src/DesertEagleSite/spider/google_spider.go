package spider

import (
  "strings"
  "github.com/PuerkitoBio/goquery"
)

func GetGoogleData(keyword string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument("http://i.xsou.co/search?q=" + keyword)
	if err != nil {
		return nil, "", err
	}
	return ParseGoogleHTML(resp)
}

func ParseGoogleUrl(url string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument(url)
	if err != nil {
		return nil, "", err
	}
	return ParseGoogleHTML(resp)
}

func ParseGoogleHTML(resp *goquery.Document) ([]DataItem, string, error) {
  resItems := make([]DataItem, 0)
  resp.Find(".g").Each(func(i int, s *goquery.Selection) {
    resItem := DataItem{}
    resItem.Title = s.Find(".r a").First().Text()
    resItem.Link = s.Find(".r a").First().AttrOr("href", "")
    startIndex := strings.Index(resItem.Link, "http://")
    if startIndex < 0 {
      startIndex = strings.Index(resItem.Link, "https://")
    }
    if startIndex >= 0 {
      resItem.Link = resItem.Link[startIndex:
        strings.Index(resItem.Link[startIndex:], "&") + startIndex]
    } else {
      return
    }
    resItem.Abstract = s.Find(".s .st").Text()
    resItem.Image = s.Find("img").AttrOr("src", "")
    resItems = append(resItems, resItem)
  })
	nextPage := ""
	nextHtml := resp.Find(".fl")
	if len(nextHtml.Nodes) > 0 {
    nextPage = "http://i.xsou.co" + nextHtml.Last().AttrOr("href", "")
	}
	return resItems, nextPage, nil
}
