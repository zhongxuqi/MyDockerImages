package spider

import (
	"bytes"
	"strings"
	"github.com/opesun/goquery"
	"github.com/opesun/goquery/exp/html"
)

type DataItem struct {
	Title string
	Abstract string
	Link string
	Image string
}

func text(b *bytes.Buffer, node *goquery.Node)  {
	if node.Type == html.TextNode {
		b.Write([]byte(strings.Trim(node.Data, " \n")))
	}
	for _, c := range node.Child {
		text(b, &goquery.Node{c})
	}
}
