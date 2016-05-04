package spider

import (
	//"fmt"
	"bytes"
	"strings"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"github.com/opesun/goquery"
)

type Page struct {
	Next string `json:"next"`
}

type ZhiHuData struct {
	Paging Page `json:"paging"`
	Htmls []string `json:"htmls"`
}

func GetZhihuData(keyword string) ([]DataItem, string, error) {
	resp, err := http.Get("http://www.zhihu.com/r/search?range=&type=question&offset=0&q=" + keyword)
	if err != nil {
		return nil, "", err
	}
	return ParseZhihuHTML(resp)
}

func ParseZhihuHTML(resp *http.Response) ([]DataItem, string, error) {
	resJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	var zhihudata ZhiHuData
	err = json.Unmarshal(resJson, &zhihudata)
	if err != nil {
		return nil, "", err
	}
	
	resItems := make([]DataItem, len(zhihudata.Htmls))
	for i, htmlnode := range zhihudata.Htmls {
		itemNodes, err := goquery.ParseString(htmlnode)
		if err != nil {
			return nil, "", err
		}
		title := itemNodes.Find("div.title a")
		if len(title) > 0 {
			var b bytes.Buffer
			text(&b, title[0])
			resItems[i].Title = b.String()
			resItems[i].Link = "http://www.zhihu.com" + title.Attrs("href")[0]
		}
		abstract := itemNodes.Find(".summary")
		if len(abstract) > 0 {
			var b bytes.Buffer
			text(&b, abstract[0])
			resItems[i].Abstract = strings.Replace(strings.Trim(b.String(), " \n"), "\n", " ", -1)
		}
	}
	nextPage := "http://www.zhihu.com" + zhihudata.Paging.Next
	return resItems, nextPage, nil
}
