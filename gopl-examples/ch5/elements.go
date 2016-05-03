// elements prints the elements and their count in the html given on stdin
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "elements: %v\n", err)
		os.Exit(1)
	}
	elements := make(map[string]int)
	elements = visit(elements, doc)
	fmt.Printf("Elements:\n%++v\n\n", elements)
	//for _, link := range visit(nil, doc) {
	//	fmt.Println(link)
	//}
}

// visit counts elements found on the page, recursivly, using 1 shared map (ok since it's not parallel..?)
func visit(elem map[string]int, n *html.Node) map[string]int {
	if n == nil {
		return elem
	}

	if n.Type == html.ElementNode {
		elem[n.Data]++
	}
	//for c := n.FirstChild; c != nil; c = c.NextSibling {
	//	links = visit(links, c)
	//}

	elem = visit(elem, n.NextSibling)
	elem = visit(elem, n.FirstChild)

	return elem
}
