package spider

import (
  "bytes"
  "net/http"
  "net/url"
  "DesertEagleSite/util"
	"strings"
  "github.com/PuerkitoBio/goquery"
)

type Jar struct {
    cookies []*http.Cookie
}
func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
    jar.cookies = cookies
}
func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
    return jar.cookies
}

func GetGithubData(keyword string) ([]DataItem, string, error) {
  req, err := http.NewRequest("GET", "https://github.com/search?utf8=✓&q=" + keyword, nil)
	if err != nil {
		return nil, "", err
	}
  req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
  req.Header.Add("Upgrade-Insecure-Requests", "1")
  req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/49.0.2623.108 Chrome/49.0.2623.108 Safari/537.36")
  req.Header.Add("Accept-Encoding", "gzip, deflate, sdch")
  req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8,en;q=0.6")
  req.Header.Add("Connection", "keep-alive")
  req.Header.Add("Host", "github.com")
  jar := new(Jar)
  http.DefaultClient.Jar = jar
  resp, err := http.DefaultClient.Do(req)
  if err != nil {
		return nil, "", err
	}
  buf := bytes.NewBuffer([]byte(""))
  resp.Write(buf)
  util.Write2File("temp.html", buf.Bytes())
  doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(buf.Bytes()))
	return ParseGithubHTML(doc)
}

func ParseGithubUrl(url string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument(url)
	if err != nil {
		return nil, "", err
	}
	return ParseGithubHTML(resp)
}

func ParseGithubHTML(resp *goquery.Document) ([]DataItem, string, error) {
	resItems := make([]DataItem, 0)
  resp.Find(".repo-list .repo-list-item").Each(func(i int, s *goquery.Selection) {
    resItem := DataItem{}
    resItem.Title = s.Find("h3 a").First().Text()
    resItem.Link = "https://github.com" +
    s.Find("h3 a").First().AttrOr("href", "")
		if len(s.Find(".repo-list-description").Nodes) > 0 {
    	resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".repo-list-description").Text(), " \n"), "\n", " ", -1)
		}
    resItem.Image = s.Find("img").AttrOr("src", "")
    resItems = append(resItems, resItem)
  })
	nextPage := ""
	nextHtml := resp.Find(".pagination a.next_page")
	if len(nextHtml.Nodes) > 0 {
    nextPage = "https://github.com" + nextHtml.Last().AttrOr("href", "")
	}
	return resItems, nextPage, nil
}
