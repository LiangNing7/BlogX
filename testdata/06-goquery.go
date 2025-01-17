package main

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	reader, err := os.Open("uploads/index.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		fmt.Println(err)
		return
	}
	selection := doc.Find("title")
	doc.Find("head").AppendHtml(" <meta name=\"keyword1\" content=\"枫枫知道,枫枫知道个人博客,网站,开发,程序员,golang\">")
	// fmt.Println(selection.Text())
	selection.SetText("枫枫知道")
	selection.SetAttr("", "")
	html, err := doc.Html()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(html)
}
