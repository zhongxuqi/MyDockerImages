package main

import (
	"fmt"
	"os"
	"webspider/spider"
	"net/http"
	"io/ioutil"
)

func main() {
	fmt.Println("it is spidertest.")

	// resp, err := http.Get("http://www.baidu.com/s?ie=uft-8&word=android")
	// b, _ := ioutil.ReadAll(resp.Body)
	// file, err := os.Create("baidu.html")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// _, err = file.Write(b)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// defer file.Close()

	resItems, nextPage, err := spider.GetBaiduData("android")
	CheckError(err)
	for _, item := range resItems {
		fmt.Println(item.Title)
		fmt.Println(item.Link)
		fmt.Println(item.Abstract)
		fmt.Println(item.Image)
	}
	fmt.Println(nextPage)
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
