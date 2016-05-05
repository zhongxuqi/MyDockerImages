package main

import (
	"time"
	"net"
	"net/http"
	. "qiniupkg.com/api.v7/conf"
	"qiniupkg.com/api.v7/kodo"
)

type MyResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
	AccessKey string `json:"access_key"`
	Token string `json:"token"`
}

// var staticHandler http.Handler = http.FileServer(http.Dir("templates/"))
// var urlMap map[string] http.Handler
var urlFuncMap map[string] func(w http.ResponseWriter, r *http.Request)

func initEnv() {
	// urlMap = make(map[string] http.Handler)
	// urlMap["media/css"] = staticHandler
	// urlMap["media/js"] = staticHandler
	// urlMap["media/image"] = staticHandler

	urlFuncMap = make(map[string] func(w http.ResponseWriter, r *http.Request))
	urlFuncMap["list"] = getListToken
	urlFuncMap["upload"] = getUploadToken
	urlFuncMap["download"] = getDownloadToken
	urlFuncMap["stat"] = statFile
	urlFuncMap["delete"] = deleteFile

	kodo.SetMac(ACCESS_KEY, SECRET_KEY)
}

func main() {
	// init tcp port
	localAddr, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:37078")
	http.DefaultTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			LocalAddr: localAddr,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	
	ACCESS_KEY = "4oK4Tp4L4zVVwZl6Vk_d2C5O1wC08hfXRi9bAu-Q"
	SECRET_KEY = "sEhV0aPeFD57hNc4MzJIMmE39VEtxTL2K87TTOOB"
	initEnv()
	http.HandleFunc("/", HandleMain)

	http.ListenAndServe("0.0.0.0:8090", nil)
}
