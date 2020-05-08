// Functions for visiting xxl and look for bikes
// not really a scraper, it just visitsa a few pages, and a json api, to look for specific info
package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// nameRE => size as an int, must be found manually for each bike?
type bikesWanted map[string][]int

type bikeIgnored struct {
	nameRE regexp.Regexp
}

type bikeResults map[string]*bikeResult // id => bikeResult

type bikeResult struct {
	name, url string
	price     int
	availible []struct {
		size     int
		store    string
		lowStock bool
	}
	unavailible []struct {
		size  int
		store string
	}
}

func findBikes() error {

	// "global" variables, to keep things easy (nothing should run in paralell here)
	uniqueBikes := make(map[string]bool) // We can see the same bike many times, unique it by relative url

	// input
	maxPrice := 15000
	wantedBikes := make(map[string][]int)
	wantedBikes[`(?i)white\s+GX\s+Lite\s+20`] = []int{1, 3, 4}       // S(52) L XL
	wantedBikes[`(?i)white\s+gx\s+pro\s+20.*herre`] = []int{1, 3, 4} // S(52) L XL

	wantedBikes[`(?i)merida\s+speeder\s+sl\s+20`] = []int{3} // Ikke gravel!

	wantedBikes[`(?i)merida\s+silex\s+200`] = []int{1, 3} // S(47) L, ikke i butikken nå
	wantedBikes[`(?i)merida\s+silex\s+300`] = []int{1, 3} // S(47) L

	wantedBikes[`(?i)scott\s+gravel\s+expert`] = []int{1, 3, 4} // S L XL

	unwantedBikes := []*regexp.Regexp{
		regexp.MustCompile(`(?i)white\s+rr\s*pro\s+20`),
		regexp.MustCompile(`(?i)white\s+rr\s+pro\s+ane\s+20`), // Selv S (size 1) er for stor
		regexp.MustCompile(`(?i)white\s+gx\s+pro\s+ane\s+20`), // Selv S (size 1) er for stor
	}

	linkREs := []string{`(?i)cyclocross`, `(?i)landeveis\s*sykkel`}

	// results
	var lastErr error
	resWanted := make(bikeResults)
	resUnwanted := make(bikeResults)
	resUnknown := make(bikeResults)

	// error handler from all of them
	// we assume no errors, so report all errors
	onErr := func(r *colly.Response, err error) {
		lastErr = fmt.Errorf("Request URL: %q failed with response: %q. Error: %s", r.Request.URL, r.Body, err)
	}

	// Initialize 3 collector:
	// 1) Find pages with bike listings
	// 2) Find interesting bikes on thos pages (wanted, unwanted and unknown)
	// 3) Find out if the interesting ones (wanted) are in stock, and where

	// 1) this one finds interesting bike page(s)
	bikePagesCollector := colly.NewCollector(
		colly.AllowedDomains("xxl.no", "www.xxl.no"),
		colly.MaxDepth(10),
	)
	bikePagesCollector.OnError(onErr)

	// 2) his one finds interesting bikes on those pages(s)
	bikesCollector := colly.NewCollector(
		colly.AllowedDomains("xxl.no", "www.xxl.no"),
		colly.MaxDepth(10),
	)
	bikesCollector.OnError(onErr)

	// 3) this one finds in-stock for a specifick bike
	bikeStockCollector := colly.NewCollector(
		colly.AllowedDomains("xxl.no", "www.xxl.no"),
		colly.MaxDepth(10),
	)
	bikeStockCollector.OnError(onErr)

	//
	//
	// 1

	bikePagesCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("bikePagesCollector: Visiting", r.URL.String())
	})

	// Find interesting bike pages
	bikePagesCollector.OnHTML("a[href]", func(e *colly.HTMLElement) {

		// Stop if we have seen any error anywhere
		if lastErr != nil {
			return
		}

		// What pages are we looking for?

		// do this link match any of them?
		visit := false
		for _, linkRE := range linkREs {

			found, err := regexp.MatchString(linkRE, e.Text)
			if err != nil {
				lastErr = fmt.Errorf("bikePagesCollector: Reg exp failesd: %s", err)
				return
			}

			if !found {
				continue
			}

			// direct link to a bike or a page with a list?
			if !strings.Contains(e.Attr("href"), "/sykkel/") {
				continue
			}

			visit = true
			break
		}

		if !visit {
			return
		}

		bikesCollector.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})

	//
	//
	// 2

	bikesCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("bikesCollector: Visiting", r.URL.String())
	})

	// Find interesting bikes
	bikesCollector.OnHTML("a[href]", func(e *colly.HTMLElement) {

		// Local helper vars
		id := strings.Replace(e.Attr("data-id"), "_Style", "", 1)
		name := e.Attr("data-brand") + " " + e.Attr("data-name")

		// Stop if we have seen any error anywhere
		if lastErr != nil {
			return
		}

		// Have we seen this bike already? (e.g. linked to from another page)
		// FIXME use id and the other map?
		if uniqueBikes[e.Attr("href")] {
			return
		}

		// has price attr
		if e.Attr("data-price") == "" {
			return
		}

		// Filter on price
		pf, err := strconv.ParseFloat(e.Attr("data-price"), 32)
		if err != nil {
			lastErr = fmt.Errorf("bikesCollector: error: could not convert price %q: err: %s", e.Attr("data-price"), err)
			return
		}
		price := int(pf)

		if price > maxPrice {
			return
		}

		// Are we activly ignoring this?
		for _, ignore := range unwantedBikes {
			if ignore.MatchString(name) {
				resUnwanted[id] = &bikeResult{
					name:  name,
					url:   e.Request.AbsoluteURL(e.Attr("href")),
					price: price,
				}
				return
			}
		}

		// fmt.Printf("bikesCollector: Link found: %q -> %s\n", name, e.Attr("href"))

		// Are we activly looking for this?
		wanted := false
		for nameRE, sizes := range wantedBikes {
			found, err := regexp.MatchString(nameRE, e.Text)
			if err != nil {
				lastErr = fmt.Errorf("bikesCollector: Reg exp failes: %s", err)
				return
			}
			if !found {
				continue
			}

			wanted = true // yes we want it (found at least once)

			// save some info in our global found map (instead of using the context to pass it on)
			// should work fine since evrything is sequential
			resWanted[id] = &bikeResult{
				name:  name,
				url:   e.Request.AbsoluteURL(e.Attr("href")),
				price: price,
			}

			// Call the next step (collector) for each size of we are iterested in
			if len(sizes) == 0 {
				lastErr = fmt.Errorf("bikesCollector: Wanted bike %q (%q) has 0 sizes", nameRE, id)
				return
			}

			for _, size := range sizes {
				ctx := colly.NewContext()
				ctx.Put("id", id)
				ctx.Put("size", size)
				err = bikeStockCollector.Request("GET", fmt.Sprintf("https://www.xxl.no/p/%s_Size_%d/stores/stock", id, size), nil, ctx, nil)
				if err != nil {
					lastErr = err
					return
				}
			}

		}

		if !wanted {
			resUnknown[id] = &bikeResult{
				name:  name,
				url:   e.Request.AbsoluteURL(e.Attr("href")),
				price: price,
			}
		}

		uniqueBikes[e.Attr("href")] = true
	})

	//
	//
	// 3

	bikeStockCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("bikeStockCollector: Visiting", r.URL.String())
	})

	bikeStockCollector.OnResponse(func(r *colly.Response) {

		// Stop if we have seen any error anywhere
		if lastErr != nil {
			return
		}

		// get som values out of the context
		id := r.Ctx.Get("id")
		size := r.Ctx.GetAny("size").(int) // will panic if not int, assume we don't fuck up

		fmt.Println("bikeStockCollector.OnResponse")

		if r.StatusCode != 200 {
			lastErr = fmt.Errorf("bikeStockCollector: Response was not 200: %d  body: %s", r.StatusCode, r.Body)
			return
		}

		sr := &StockReply{}
		err := json.Unmarshal(r.Body, sr)
		if err != nil || len(sr.StoreAvailabilities) == 0 {
			lastErr = fmt.Errorf("bikeStockCollector: Response was unexpected. JSON err: %v", err)
			return
		}

		for _, s := range sr.StoreAvailabilities {
			re := regexp.MustCompile(`(?i)oslo`)
			if re.MatchString(s.StoreData.Address.Town) {
				fmt.Print(s.StoreData.Address.Town, "   ", s.StoreData.DisplayName, "   ", s.StoreData.Name, "   ")
				fmt.Println(s.StockStatus)

				switch s.StockStatus {
				case "OUTOFSTOCK":
					resWanted[id].unavailible = append(resWanted[id].unavailible, struct {
						size  int
						store string
					}{
						size:  size,
						store: s.StoreData.DisplayName,
					})
				case "LOWSTOCK":
					resWanted[id].availible = append(resWanted[id].availible, struct {
						size     int
						store    string
						lowStock bool
					}{
						size:     size,
						store:    s.StoreData.DisplayName,
						lowStock: true,
					})
				default: // be possitive :D
					resWanted[id].availible = append(resWanted[id].availible, struct {
						size     int
						store    string
						lowStock bool
					}{
						size:     size,
						store:    s.StoreData.DisplayName,
						lowStock: false,
					})
				}

			}
		}

	})

	err := bikePagesCollector.Visit("https://www.xxl.no")
	if err != nil {
		return fmt.Errorf("Initial Visit failed with error: %s (lastErr: %v)", err, lastErr)
	}

	fmt.Printf("\nWanted:\n")
	for k, v := range resWanted {
		fmt.Printf("%s (%s)\t%d\n%s\n", v.name, k, v.price, v.url)
		fmt.Println("\tAvailible:")
		for _, vv := range v.availible {
			if vv.lowStock {
				fmt.Printf("\t\t%d\t%s\n", vv.size, vv.store)
			} else {
				fmt.Printf("!!!!\t%d\t%s\n", vv.size, vv.store)
			}
		}

		fmt.Println("\tUnavailible:")
		for _, vv := range v.unavailible {
			fmt.Printf("\t\t%d\t%s\n", vv.size, vv.store)
		}
	}

	fmt.Printf("\nUnwanted:\n")
	for k, v := range resUnwanted {
		fmt.Printf("%s (%s)\t%d\n%s\n", v.name, k, v.price, v.url)
	}

	fmt.Printf("\nUnknown:\n")
	for k, v := range resUnknown {
		fmt.Printf("%s (%s)\t%d\n%s\n", v.name, k, v.price, v.url)
	}

	return lastErr

}
