package main

import (
	"net/http"
	_ "DesertEagleSite/handlefunc"
	"DesertEagleSite/handlefunc"
)

func main() {
	http.HandleFunc("/", handlefunc.HandleMain)

	http.ListenAndServe("0.0.0.0:8089", nil)
}
