package spider

import (
	"strings"
  "github.com/PuerkitoBio/goquery"
)

func GetBingData(keyword string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument("http://cn.bing.com/search?q=" + keyword)
	if err != nil {
		return nil, "", err
	}
	return ParseBingHTML(resp)
}

func ParseBingUrl(url string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument(url)
	if err != nil {
		return nil, "", err
	}
	return ParseBingHTML(resp)
}

func ParseBingHTML(resp *goquery.Document) ([]DataItem, string, error) {
	resItems := make([]DataItem, 0)
  resp.Find("#b_results li.b_algo, #b_results li.b_ans").Each(func(i int, s *goquery.Selection) {
    resItem := DataItem{}
		if len(s.Find("h2 a").Nodes) > 0 {
	    resItem.Title = strings.Replace(strings.Trim(
				s.Find("h2 a").First().Text(), " \n"), "\n", " ", -1)
	    resItem.Link = s.Find("h2 a").First().AttrOr("href", "")
		} else if len(s.Find("h5 a").Nodes) > 0 {
			resItem.Title = strings.Replace(strings.Trim(
				s.Find("h5 a").First().Text(), " \n"), "\n", " ", -1)
	    resItem.Link = s.Find("h5 a").First().AttrOr("href", "")
		} else {
			return
		}
		if len(s.Find(".b_rich p").Nodes) > 0 {
    	resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".b_rich p").Text(), " \n"), "\n", " ", -1)
		} else if len(s.Find(".b_caption p").Nodes) > 0 {
			resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".b_caption p").Text(), " \n"), "\n", " ", -1)
		} else if len(s.Find(".b_overflow p").Nodes) > 0 {
			resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".b_overflow p").Text(), " \n"), "\n", " ", -1)
		}
    resItem.Image = s.Find("img").AttrOr("src", "")
    resItems = append(resItems, resItem)
  })
	nextPage := ""
	nextHtml := resp.Find("a.sb_pagN")
	if len(nextHtml.Nodes) > 0 {
    nextPage = "http://cn.bing.com" + nextHtml.Last().AttrOr("href", "")
	}
	return resItems, nextPage, nil
}
