package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// - Begynn på start (hviss error abort)
// - Let etter cyclecross linken  (hviss error abort)
// - Let etter riktig sykkel linken  (hviss error abort)
// - Velg størrelse  (hviss error abort)
// - Hent lagerstatus for alle butikker  (hviss error abort)
// - grep ut de i Oslo  (hviss error abort)
// - Mail output ??
//
//  Errors: Mail

func main() {

	// Init 3 collectors

	// this one finds interesting bike page(s)
	bikePagesCollector := colly.NewCollector(
		colly.AllowedDomains("xxl.no", "www.xxl.no"),
		colly.MaxDepth(10),
	)

	// error handler from all of them
	// we assume no errors, so report all errors
	bikePagesCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		os.Exit(1)
	})

	// this one finds interesting bikes on those pages(s)
	bikesCollector := bikePagesCollector.Clone()

	// this one extracs relevant info from the interesting bikes
	bikeInfoCollector := bikePagesCollector.Clone()

	//
	//
	// 1

	bikePagesCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("bikePagesCollector: Visiting", r.URL.String())
	})

	// Find interesting bike pages
	bikePagesCollector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// fmt.Printf("bikePagesCollector: Link found: %q -> %s\n", e.Text, link)

		found, err := regexp.MatchString(`(?is)cyclocross`, e.Text)
		if err != nil {
			fmt.Println("bikePagesCollector: Reg exp failes: ", err)
			os.Exit(1)
		}

		if !found {
			return
		}

		// direct link to a bike or a page with a list?
		if !strings.Contains(link, "/sykkel/") {
			return
		}
		bikesCollector.Visit(e.Request.AbsoluteURL(link))
	})

	//
	//
	// 2

	//	priceRE := regexp.MustCompile(`(?is)(\d+)\\u00a0(\d+)\s*,-`)
	//priceRE := regexp.MustCompile(`(?is)00a0`)

	// Find interesting bikes
	bikesCollector.OnHTML("a[href]", func(e *colly.HTMLElement) {

		// has price attr
		if e.Attr("data-price") == "" {
			return
		}

		pf, err := strconv.ParseFloat(e.Attr("data-price"), 32)
		if err != nil {
			fmt.Printf("bikesCollector: error: could not convert price %q: err: %s\n", e.Attr("data-price"), err)
			os.Exit(1)
		}
		price := int(pf)

		if price > 15000 {
			return
		}

		link := e.Attr("href")

		fmt.Printf("bikesCollector: Link found: %q -> %s\n", e.Text, link)
		fmt.Printf("%d    : %q\n", price, e.Attr("data-price"))

		// prices := priceRE.FindAllString(e.Text, 5)
		// if prices == nil {
		// 	return
		// }

		//fmt.Printf("PRICES: (%d): %#v\n\n", len(prices), prices)

		// found, err := regexp.MatchString(`(?is)cyclocross`, e.Text)
		// if err != nil {
		// 	fmt.Println("bikesCollector: Reg exp failes: ", err)
		// 	os.Exit(1)
		// }

		// if !found {
		// 	return
		// }

		// // direct link to a bike or a page with a list?
		// if !strings.Contains(link, "/sykkel/") {
		// 	return
		// }
		// bikesCollector.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	bikesCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("bikesCollector: Visiting", r.URL.String())
	})

	//
	//
	// 3

	// Before making a request print "Visiting ..."
	bikeInfoCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("bikeInfoCollector: Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	bikePagesCollector.Visit("https://www.xxl.no")
}
