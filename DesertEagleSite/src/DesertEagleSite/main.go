package main

import (
	// "fmt"
	"time"
	"net"
	"net/http"
	_ "DesertEagleSite/handlefunc"
	"DesertEagleSite/handlefunc"
)

func main() {
	// init tcp port
	localAddr, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:37077")
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
	http.HandleFunc("/", handlefunc.HandleMain)

	http.ListenAndServe("0.0.0.0:8089", nil)
}
