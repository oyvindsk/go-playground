// Package htmlstatemachinemetadataextractorthingy html-statemachine-metadata-extractor-thingy
// this was exttacted from the blogengine github.com/oyvindsk/web-oyvindsk.com/
// not testet much, but I keep it since  I like the state machine and I might need something similar in the future
package htmlstatemachinemetadataextractorthingy

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// cutMetadata looks for metadata in the input html (written for ascidoc(tor) html)
// it and extracts the magic metadata lines and removes them, along with the corresponding html (that is now useless)
// loops the html tokens and creates a simple state machine
func cutMetadata(input io.Reader) (io.Reader, blogMetadata, error) {

	z := html.NewTokenizer(input)

	/*
		States at the beginning of the big for loop:

		Current State								Input that switches state 			Next state
		----------------------------------------------------------------------------------------------------------
		(lfop) Looking for opening token			div token with openblock class 		lfml
		(lfml) Looking for magic lines				line that starts with "||"			iml
		(iml)  In magic								line that does not start with ..	lfet
		(lfet) Looking for tokens we should exclude last tokens in exclude list			done
		(done) Done looking and exluding			n/a									n/a
	*/

	// state machine variables
	state := "lfot"                 // current state
	var prevToken html.Token        // we sometimes have to look back to the previous token
	var endTokensToExlcude []string // usen when excluding tokens around the metadata lines. Typically 3 tokens before and 3 after.

	// results of the state machine loop (not the func)
	var magicLines []string
	var body strings.Builder
	var err error

MACHINE:
	for {

		// Advance to next token
		tt := z.Next()
		if tt == html.ErrorToken {
			// This includes EOF, break out and deal with it later
			err = z.Err()
			break MACHINE
		}

		thisToken := z.Token() // The token we are currenlty looking at, as opposed to prevToken

		// Switch on the 5 known states. See above.
		// this could of course be something other than a string, otoh ..
		// we do not really enforce that all transitions are valid, but that would require a bug in the code (?)
		switch state {

		case "lfot":

			// Look for opening div of metadata, with class openblock
			var found bool
			if tt == html.StartTagToken && thisToken.Data == "div" {
				if found, _ = findAttr(thisToken.Attr, "class", "openblock"); found {
					state = "lfml" // fmt.Println("\n\t==>\t Looking for magic lines")

					// add this div to the list of tokens we want to exlude after the magic lines (in lfet)
					endTokensToExlcude = append(endTokensToExlcude, thisToken.Data)
				}
			}

			// Include token unless it was the opening div we are looking for
			if !found {
				body.WriteString(thisToken.String())
			}

		case "lfml":

			if thisToken.Type.String() == "Text" && strings.HasPrefix(thisToken.Data, "||") {
				state = "iml" // fmt.Println("\n\t==>\t In magic")
				break
			}

			// Add tokens we see before the firts line of magic
			// to the list of tokens we want to exlude after the magic lines (in lfet)
			if thisToken.Type == html.StartTagToken {
				endTokensToExlcude = append(endTokensToExlcude, thisToken.Data)
			}

		case "iml":

			// Save the magic lines(s) for later
			// syntax from ascidoc(tor) puts it on 1 line with a \n, so ..
			magicLines = append(magicLines, strings.Split(prevToken.String(), "\n")...)

			if thisToken.Type.String() != "Text" || !strings.HasPrefix(thisToken.Data, "||") {
				state = "lfet" // fmt.Printf("\n\t==>\t Looking for tags we should exclude\n")
			}

		case "lfet":

			if prevToken.Type == html.EndTagToken && prevToken.Data == endTokensToExlcude[len(endTokensToExlcude)-1] {
				endTokensToExlcude = endTokensToExlcude[:len(endTokensToExlcude)-1]
			}

			if len(endTokensToExlcude) == 0 {
				state = "done" //	fmt.Println("\n\t==>\t DONE!")
			}

		case "done":
			body.WriteString(thisToken.String())

		default:
			err = fmt.Errorf("unknown state seen: %q", state)
			break MACHINE
		}

		prevToken = thisToken
	}

	// Any parse / state machine error from?
	if err != nil {
		if err != io.EOF {
			return nil, blogMetadata{}, fmt.Errorf("cutMetadata: error when running state machine: %s", err)
		}
		err = nil
	}

	// Convert the magic lines we found into blogMetadata
	metadata, err := blogMetadataFromMagicLines(magicLines)
	if err != nil {
		return nil, blogMetadata{}, fmt.Errorf("cutMetadata: %s", err)
	}

	return strings.NewReader(body.String()), metadata, nil
}

func findAttr(attrs []html.Attribute, key, val string) (bool, int) {
	for i := range attrs {
		if attrs[i].Key == key {
			if attrs[i].Val == val {
				return true, i // assume only 1 match
			}
		}
	}
	return false, 0
}

// blogMetadata is the metadata we expect to find in the magic lines
type blogMetadata struct {
	author    string
	title     string
	subtitle  string
	date      string
	servePath string
	tags      []string
}

// blogMetadataFromMagicLines takes the magi line(s) found by cutMetadata and returns a blogMetadata struct
//
// Input looks like this:
// "|| Adam Morse || Too many tools and frameworks: subTTT || 2015 || /foo/bar || Subtitle: The definitive guide to the javascript tooling landscape in 2015"
// "|| foo bar go golang javascript"
func blogMetadataFromMagicLines(magicLines []string) (blogMetadata, error) {
	if !(len(magicLines) == 1 || len(magicLines) == 2) {
		return blogMetadata{}, fmt.Errorf("blogMetadataFromMagicLines: Expect 1 or 2 magix lines, got: %d", len(magicLines))
	}

	// First line, || separated, everything but the tags
	l1 := regexp.MustCompile(`\s?\|\|\s?`).Split(magicLines[0], 100)
	l1 = l1[1:] // first is always bogus since we start out line with ||

	m := blogMetadata{
		author:    l1[0],
		title:     l1[1],
		subtitle:  l1[4],
		date:      l1[2],
		servePath: l1[3],
	}

	// add tags if any
	if len(magicLines) > 1 && len(magicLines[1]) > 4 {
		if !strings.HasPrefix(magicLines[1], "|| ") {
			return blogMetadata{}, fmt.Errorf("blogMetadataFromMagicLines: Tag line invalid, must start with '|| '")
		}

		//l2 := regexp.MustCompile(`\|?\|?\s`).Split(magicLines[1], 100)
		m.tags = strings.Fields(magicLines[1][3:]) // split on space after '|| '
	}

	return m, nil
}
