package main

import (
	"os"
	"net/http"
	"sync"
	"time"
	"strconv"
	"fmt"
	"strings"
	"encoding/json"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"io/ioutil"
	. "qiniupkg.com/api.v7/conf"
	"qiniupkg.com/api.v7/kodo"
)

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

	url := r.URL.Path[1:]
	if mhandleFunc, ok := urlFuncMap[url]; ok {
		mhandleFunc(w, r)
		return
	}
}

func writeResult(w http.ResponseWriter, r *http.Request, token string) {
	var response MyResponse
	response.Status = "200"
	response.Message = "success"
	response.Token = token
	response.AccessKey = ACCESS_KEY
	respBytes, err := json.Marshal(response)
	if err != nil {
		err.Error()
		return
	}
	w.Write(respBytes)
}

func parseParam(r *http.Request, keyname string) (string, bool) {
	for _, item := range strings.Split(r.URL.RawQuery, "&") {
		if item[0:strings.Index(item, "=")] == keyname {
			return item[strings.Index(item, "=")+1:], true
		}
	}
	return "", false
}

func getListToken(w http.ResponseWriter, r *http.Request) {
	bucket, ok := parseParam(r, "bucket")
	if !ok {
		return
	}
	host := "http://rsf.qbox.me"
	path := "/list?"
	body := "bucket=" + bucket
	accessToken := getManageToken(path + body + "\n")

	req, err := http.NewRequest("POST", host+path+body, nil)
	if err != nil {
		err.Error()
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "QBox "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err.Error()
		return
	}
	respB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err.Error()
		return
	}
	respStr := string(respB)

	writeResult(w, r, respStr)
}

func getManageToken(signingStr string) string {
	mac := hmac.New(sha1.New, []byte(SECRET_KEY))
	mac.Write([]byte(signingStr))
	encodedSign := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	accessToken := ACCESS_KEY + ":" + encodedSign
	return accessToken
}

func getUploadToken(w http.ResponseWriter, r *http.Request) {
	bucket, ok := parseParam(r, "bucket")
	if !ok {
		return
	}
	key, ok := parseParam(r, "key")
	if !ok {
		return
	}
	policy := &kodo.PutPolicy{
		Scope:      bucket + ":" + key,
		Expires:    uint32(3600 + time.Now().Unix()),
		ReturnBody: "{\"w\":$(imageInfo.width),\"h\":$(imageInfo.height)}",
	}
	b, err := json.Marshal(policy)
	if err != nil {
		err.Error()
		return
	}
	encodedPutPolicy := base64.URLEncoding.EncodeToString(b)
	uptoken := ACCESS_KEY + ":" + getEncodedSign([]byte(encodedPutPolicy)) + ":" + encodedPutPolicy
	writeResult(w, r, uptoken)
}

func getDownloadToken(w http.ResponseWriter, r *http.Request) {
	encodedURL, ok := parseParam(r, "url")
	if !ok {
		return
	}
	b, err := base64.URLEncoding.DecodeString(encodedURL)
	if err != nil {
		err.Error()
		return
	}
	realURL := string(b) + "?e=" + strconv.Itoa(int(time.Now().Unix())+3600)
	downloadtoken := ACCESS_KEY + ":" + getEncodedSign([]byte(realURL))
	downloadUrl := realURL + "&token=" + downloadtoken
	writeResult(w, r, base64.URLEncoding.EncodeToString([]byte(downloadUrl)))
}

func statFile(w http.ResponseWriter, r *http.Request) {
	bucket, ok := parseParam(r, "bucket")
	if !ok {
		return
	}
	key, ok := parseParam(r, "key")
	if !ok {
		return
	}
	entry := bucket + ":" + key
	encodedEntryURI := base64.URLEncoding.EncodeToString([]byte(entry))

	actionFile(w, r, "/stat/", encodedEntryURI)
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	bucket, ok := parseParam(r, "bucket")
	if !ok {
		return
	}
	key, ok := parseParam(r, "key")
	if !ok {
		return
	}
	entry := bucket + ":" + key
	encodedEntryURI := base64.URLEncoding.EncodeToString([]byte(entry))

	actionFile(w, r, "/delete/", encodedEntryURI)
}

func actionFile(w http.ResponseWriter, r *http.Request, action, encodedEntryURI string) {
	host := "http://rs.qiniu.com"
	path := action + encodedEntryURI
	encodedSign := signString(path + "\n")
	accessToken := ACCESS_KEY + ":" + encodedSign

	req, err := http.NewRequest("GET", host+path, nil)
	if err != nil {
		err.Error()
		return
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", "QBox "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		err.Error()
		return
	}
	respB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err.Error()
		return
	}
	writeResult(w, r, string(respB))
}

func signString(data string) string {
	return getEncodedSign([]byte(data))
}

func getEncodedSign(b []byte) string {
	mac := hmac.New(sha1.New, []byte(SECRET_KEY))
	mac.Write([]byte(b))
	sign := mac.Sum(nil)
	encodedSign := base64.URLEncoding.EncodeToString(sign)
	return encodedSign
}
