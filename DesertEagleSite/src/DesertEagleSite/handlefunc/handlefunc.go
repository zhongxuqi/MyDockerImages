package handlefunc

import (
	"os"
	"time"
	"sync"
	"strconv"
	"net/http"
	"fmt"
	"encoding/json"
	"html/template"
	"encoding/hex"
	"strings"
	"webspider/spider"
	"github.com/opesun/goquery"
)

type MyResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
	ResultData []spider.DataItem `json:",omitempty"`
	NextPage string `json:",omitempty"`
}

var iconHandler http.Handler = http.FileServer(http.Dir("html/image"))
var urlFuncMap map[string] func(w http.ResponseWriter, r *http.Request)
func init() {
	urlFuncMap = make(map[string] func(w http.ResponseWriter, r *http.Request))
	urlFuncMap["app/search_baidu"] = SearchBaidu
	urlFuncMap["app/search_zhihu"] = SearchZhihu
	urlFuncMap["app/search_haosou"] = SearchHaosou
	urlFuncMap["app/search_wikipedia"] = SearchWikipedia
	urlFuncMap["app/search_baiduxueshu"] = SearchBaiduXueShu
	urlFuncMap["app/search_douban"] = SearchDouBan
	urlFuncMap["app/search_jianshu"] = SearchJianShu
	urlFuncMap["app/search_csdn"] = SearchCSDN
	urlFuncMap["app/custom_search"] = CustomSearch

	initFile();
}

func writeResult(w http.ResponseWriter, r *http.Request, resItems []spider.DataItem, nextPage string, err error) {
	var response MyResponse
	if err != nil {
		response.Status = "400"
		response.Message = err.Error()
	} else {
		response.Status = "200"
		response.Message = "search success"
		response.ResultData = resItems
		response.NextPage = nextPage
	}
	respBytes, err := json.Marshal(response)
	w.Write(respBytes)
}

func parseKeyword(r *http.Request, keyname string) (string, bool) {
	for _, item := range strings.Split(r.URL.RawQuery, "&") {
		if item[0:strings.Index(item, "=")] == keyname {
			return item[strings.Index(item, "=")+1:], true
		}
	}
	return "", false
}

var mux sync.Mutex

func writeLog(r *http.Request) {
	mux.Lock()
	defer mux.Unlock()
	t := time.Now()
	year, month, day := t.Date()
	filename := strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-"+ strconv.Itoa(day) + ".log"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer file.Close()
	file.Write([]byte(t.Format(time.UnixDate)+"  "))
	file.Write([]byte(r.Method + "  " + r.RemoteAddr + "  " + r.URL.Path + "  " + r.URL.RawQuery + "\n"))
}

func HandleMain(w http.ResponseWriter, r *http.Request) {
	writeLog(r)
	
	// go to file server
	if r.URL.Path == "/" {
		t, err := template.ParseFiles("html/index.html")
		if err != nil {
			fmt.Println(err)
		}
		t.Execute(w, nil)
		return
	}
	if r.URL.Path == "/favicon.ico" {
		iconHandler.ServeHTTP(w, r);
	}
	if r.URL.Path[0:5] == "/html" {
		handleFileServer(w, r);
		return;
	}

	// go json server
	if strings.LastIndex(r.URL.Path, "/") <= 0 {
		return
	}
	url := r.URL.Path[1:]
	if mhandleFunc, ok := urlFuncMap[url]; ok {
		mhandleFunc(w, r)
		return
	}
}

func SearchBaidu(w http.ResponseWriter, r *http.Request) {
	if keyword, ok := parseKeyword(r, "keyword"); ok {
		resItems, nextPage, err := spider.GetBaiduData(keyword)
		writeResult(w, r, resItems, hex.EncodeToString([]byte(nextPage)), err)
	}
}

func SearchZhihu(w http.ResponseWriter, r *http.Request) {
	if keyword, ok := parseKeyword(r, "keyword"); ok {
		resItems, nextPage, err := spider.GetZhihuData(keyword)
		writeResult(w, r, resItems, hex.EncodeToString([]byte(nextPage)), err)
	}
}

func SearchHaosou(w http.ResponseWriter, r *http.Request) {
	if keyword, ok := parseKeyword(r, "keyword"); ok {
		resItems, nextPage, err := spider.GetHaosouData(keyword)
		writeResult(w, r, resItems, hex.EncodeToString([]byte(nextPage)), err)
	}
}

func SearchWikipedia(w http.ResponseWriter, r *http.Request) {
	if keyword, ok := parseKeyword(r, "keyword"); ok {
		resItems, nextPage, err := spider.GetWikipediaData(keyword)
		writeResult(w, r, resItems, hex.EncodeToString([]byte(nextPage)), err)
	}
}

func SearchBaiduXueShu(w http.ResponseWriter, r *http.Request) {
	if keyword, ok := parseKeyword(r, "keyword"); ok {
		resItems, nextPage, err := spider.GetBaiduXueShuData(keyword)
		writeResult(w, r, resItems, hex.EncodeToString([]byte(nextPage)), err)
	}
}

func SearchDouBan(w http.ResponseWriter, r *http.Request) {
	if keyword, ok := parseKeyword(r, "keyword"); ok {
		resItems, nextPage, err := spider.GetDoubanData(keyword)
		writeResult(w, r, resItems, hex.EncodeToString([]byte(nextPage)), err)
	}
}

func SearchJianShu(w http.ResponseWriter, r *http.Request) {
	if keyword, ok := parseKeyword(r, "keyword"); ok {
		resItems, nextPage, err := spider.GetJianshuData(keyword)
		writeResult(w, r, resItems, hex.EncodeToString([]byte(nextPage)), err)
	}
}

func SearchCSDN(w http.ResponseWriter, r *http.Request) {
	if keyword, ok := parseKeyword(r, "keyword"); ok {
		resItems, nextPage, err := spider.GetCSDNData(keyword)
		writeResult(w, r, resItems, hex.EncodeToString([]byte(nextPage)), err)
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
	switch PaserName {
	case "BaiduData":
		resp, err := goquery.ParseUrl(url)
		if err != nil {
			return
		}
		resItems, nextPage, err := spider.ParseBaiduHTML(&resp)
		writeResult(w, r, resItems, hex.EncodeToString([]byte(nextPage)), err)
	case "ZhihuData":
		resp, err := http.Get(url)
		if err != nil {
			return
		}
		resItems, nextPage, err := spider.ParseZhihuHTML(resp)
		writeResult(w, r, resItems, hex.EncodeToString([]byte(nextPage)), err)
	case "HaosouData":
		resp, err := goquery.ParseUrl(url)
		if err != nil {
			return
		}
		resItems, nextPage, err := spider.ParseHaosouHTML(&resp)
		writeResult(w, r, resItems, hex.EncodeToString([]byte(nextPage)), err)
	case "BaiduXueShuData":
		resp, err := goquery.ParseUrl(url)
		if err != nil {
			return
		}
		resItems, nextPage, err := spider.ParseBaiduXueShuHTML(&resp)
		writeResult(w, r, resItems, hex.EncodeToString([]byte(nextPage)), err)
	case "WikipediaData":
		resp, err := goquery.ParseUrl(url)
		if err != nil {
			return
		}
		resItems, nextPage, err := spider.ParseWikipediaHTML(&resp)
		writeResult(w, r, resItems, hex.EncodeToString([]byte(nextPage)), err)
	case "JianShuData":
		resp, err := goquery.ParseUrl(url)
		if err != nil {
			return
		}
		resItems, nextPage, err := spider.ParseJianShuHTML(&resp)
		writeResult(w, r, resItems, hex.EncodeToString([]byte(nextPage)), err)
	case "CSDNData":
		resp, err := goquery.ParseUrl(url)
		if err != nil {
			return
		}
		resItems, nextPage, err := spider.ParseCSDNHTML(&resp)
		writeResult(w, r, resItems, hex.EncodeToString([]byte(nextPage)), err)
	default:

	}
}
