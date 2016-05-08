package main

import (
	"net/http"
	"DesertEagleSite/handlefunc"
	// "DesertEagleSite/wordtool"
)

func main() {
	http.HandleFunc("/", handlefunc.HandleMain)

	http.ListenAndServe("0.0.0.0:8089", nil)
}
