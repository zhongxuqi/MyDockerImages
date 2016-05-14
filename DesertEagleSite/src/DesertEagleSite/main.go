package main

import (
	"net/http"
	"DesertEagleSite/handlefunc"
	"net/http/cookiejar"
)

func main() {
	jar, _ := cookiejar.New(nil)
	http.DefaultClient.Jar = jar
	http.HandleFunc("/", handlefunc.HandleMain)

	http.ListenAndServe("0.0.0.0:8089", nil)
}
