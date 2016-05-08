package handlefunc

import (
	"net/http"
	"encoding/json"
	"encoding/hex"
	"DesertEagleSite/spider"
	// "github.com/opesun/goquery"
  // "DesertEagleSite/config"
)

type SpiderObject struct {
  Url string
  Name string
  ParserName string
  Logo string
	GetDataFunc func(string) ([]spider.DataItem, string, error) `json:"-"`
	ParseFunc func(string) ([]spider.DataItem, string, error) `json:"-"`
}

type ListResponse struct {
  BaseResponse
  Webs []SpiderObject
}

type SpiderResponse struct {
	BaseResponse
	ResultData []spider.DataItem `json:",omitempty"`
	NextPage string `json:",omitempty"`
}

const (
  BAIDU = "baidu"
  ZHIHU = "zhihu"
  HAOSOU = "haosou"
  WIKIPEDIA = "wikipedia"
  BAIDUXUESHU = "baiduxuehsu"
  DOUBAN = "douban"
  JIANSHU = "jianshu"
  CSDN = "csdn"
  BING = "bing"
	GOOGLE = "google"
	STACKOVERFLOW = "Stack Overflow"
	GITHUB = "Github"
)

var SpiderMap = map[string]SpiderObject{
  BAIDU: SpiderObject{
    Url: "app/search_baidu",
    Name: "百度",
    ParserName: "BaiduData",
    Logo: "/data/baidu_logo.png",
		GetDataFunc: spider.GetBaiduData,
		ParseFunc: spider.ParseBaiduUrl,
  },
  ZHIHU: SpiderObject{
    Url: "app/search_zhihu",
    Name: "知乎",
    ParserName: "ZhihuData",
    Logo: "/data/zhihu_logo.png",
		GetDataFunc: spider.GetZhihuData,
		ParseFunc: spider.ParseZhihuUrl,
  },
  HAOSOU: SpiderObject{
    Url: "app/search_haosou",
    Name: "好搜",
    ParserName: "HaosouData",
    Logo: "/data/haosou_logo.png",
		GetDataFunc: spider.GetHaosouData,
		ParseFunc: spider.ParseHaosouUrl,
  },
  WIKIPEDIA: SpiderObject{
    Url: "app/search_wikipedia",
    Name: "维基百科",
    ParserName: "WikipediaData",
    Logo: "/data/wikipedia_logo.png",
		GetDataFunc: spider.GetWikipediaData,
		ParseFunc: spider.ParseWikipediaUrl,
  },
  BAIDUXUESHU: SpiderObject{
    Url: "app/search_baiduxueshu",
    Name: "百度学术",
    ParserName: "BaiduXueShuData",
    Logo: "/data/baidu_logo.png",
		GetDataFunc: spider.GetBaiduXueShuData,
		ParseFunc: spider.ParseBaiduXueShuUrl,
  },
  // DOUBAN: {
  //   Url: "app/search_douban",
  //   Name: "豆瓣",
  //   ParserName: "",
  // },
  // JIANSHU: SpiderObject{
  //   Url: "app/search_jianshu",
  //   Name: "简书",
  //   ParserName: "JianShuData",
  //   Logo: "/data/jianshu_logo.png",
	// 	GetDataFunc: spider.GetJianshuData,
	// 	ParseFunc: spider.ParseJianShuUrl,
  // },
  CSDN: SpiderObject{
    Url: "app/search_csdn",
    Name: "CSDN",
    ParserName: "CSDNData",
    Logo: "/data/csdn_logo.png",
		GetDataFunc: spider.GetCSDNData,
		ParseFunc: spider.ParseCSDNUrl,
  },
	BING: SpiderObject{
    Url: "app/search_bing",
    Name: "Bing",
    ParserName: "BingData",
    Logo: "/data/bing_logo.png",
		GetDataFunc: spider.GetBingData,
		ParseFunc: spider.ParseBingUrl,
  },
	GOOGLE: SpiderObject{
		Url: "app/search_google",
    Name: "Google",
    ParserName: "GoogleData",
    Logo: "/data/google_logo.png",
		GetDataFunc: spider.GetGoogleData,
		ParseFunc: spider.ParseGoogleUrl,
	},
	STACKOVERFLOW: SpiderObject{
		Url: "app/search_stackoverflow",
    Name: "StackOverflow",
    ParserName: "StackOverflowData",
    Logo: "/data/stackoverflow_logo.png",
		GetDataFunc: spider.GetStackOverflowData,
		ParseFunc: spider.ParseStackOverflowUrl,
	},
	// GITHUB: SpiderObject{
	// 	Url: "app/search_github",
  //   Name: "Github",
  //   ParserName: "GithubData",
  //   Logo: "/data/github_logo.png",
	// 	GetDataFunc: spider.GetGithubData,
	// 	ParseFunc: spider.ParseGithubUrl,
	// },
}

func initSpider() {
  urlFuncMap["app/list"] = ListSpiders

  for _, spider := range SpiderMap {
    urlFuncMap[spider.Url] = SearchData
  }

	urlFuncMap["app/custom_search"] = CustomSearch
}

func writeSpiderResult(w http.ResponseWriter, r *http.Request, resItems []spider.DataItem, nextPage string, err error) {
	var response SpiderResponse
	if err != nil {
		response.Status = "500"
		response.Message = err.Error()
	} else {
		response.Status = "200"
		response.Message = "search success"
		response.ResultData = resItems
		response.NextPage = nextPage
	}
	respBytes, _ := json.Marshal(response)
	w.Write(respBytes)
}

func ListSpiders(w http.ResponseWriter, r *http.Request) {
  webs := make([]SpiderObject, 0)
  for _, spider := range SpiderMap {
    webs = append(webs, spider)
  }
  var response ListResponse
  response.Status = "200"
  response.Message = "success"
  response.Webs = webs
  respBytes, _ := json.Marshal(response)
  w.Write(respBytes)
}

func SearchData(w http.ResponseWriter, r *http.Request) {
	if keyword, ok := parseKeyword(r, "keyword"); ok {
		for _, spider := range SpiderMap {
			if spider.Url == r.URL.Path[1:] {
				resItems, nextPage, err := spider.GetDataFunc(keyword)
				writeSpiderResult(w, r, resItems, hex.EncodeToString([]byte(nextPage)), err)
				break
			}
		}
	}
}

func CustomSearch(w http.ResponseWriter, r *http.Request) {
	PaserName, _ := parseKeyword(r, "parser_name")
	url, _ := parseKeyword(r, "url")
	if len(PaserName) == 0 || len(url) == 0 {
		return
	}
	decodeUrl, err := hex.DecodeString(url)
	if err != nil {
		return
	}
	url = string(decodeUrl)
	for _, spider := range SpiderMap {
		if spider.ParserName == PaserName {
			resItems, nextPage, err := spider.ParseFunc(url)
			writeSpiderResult(w, r, resItems, hex.EncodeToString([]byte(nextPage)), err)
			break
		}
	}
}
