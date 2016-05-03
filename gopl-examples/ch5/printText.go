// printText prints the text in textnodes in the html given on stdin
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	z := html.NewTokenizer(os.Stdin)

	var skip bool
LOOP:
	for {
		tt := z.Next()
		tagname, _ := z.TagName()
		tagnameStr := string(tagname)
		if tagnameStr == "script" || tagnameStr == "style" {
			if skip {
				skip = false // asume this is the closing tag, should check depth etc, but..
			} else {
				skip = true
			}
		}

		if skip {
			//	fmt.Println("skipping")
			continue
		}

		//fmt.Println("tagname:", string(tagname))
		switch tt {
		case html.ErrorToken:
			fmt.Println("err:") //, tt.String())
			break LOOP
		case html.TextToken:
			//fmt.Printf("tt: %s\n\"%s\"", tt.String(), z.Text())
			txt := z.Text()
			if len(txt) > 5 {
				fmt.Printf("\"%s\"\n", txt)
			}
		default:
			//	fmt.Printf("tt: %s\n", tt.String())
		}
	}
	os.Exit(0)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "elements: %v\n", err)
	//	os.Exit(1)
	//}
	//for _, link := range visit(nil, doc) {
	//	fmt.Println(link)
	//}
}

// visit travers the documnet recursivly and prints the content of text nodes. <script> and <style> nodes are skiped
func visit(z *html.Token) {
	if z == nil {
		return
	}

	//if n.Type == html.TextNode {
	//fmt.Printf(`"%s"`, z.String())
	//visit(z.Next())
	//}

	//visit(n.NextSibling)
	//visit(n.FirstChild)

}
