package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Runme to get a local token using 3-legged oauth2
// This token can be used to grant access to the Sheet, mainly when testing locally.
// Currently use the default service account when running on Appengine / GCP. This does not work locally right now, seems like gcloud bug
// !! You need a 'client_secret.json' file, see: https://developers.google.com/sheets/api/quickstart/go
func main() {
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	tok := getTokenFromWeb(config)

	printToken(tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		msg := "Unable to read authorization code"
		log.Fatalf("%s: %v", msg, err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		msg := "Unable to retrieve token from web"
		log.Fatalf("%s: %v", msg, err)
	}
	return tok
}

func printToken(token *oauth2.Token) {
	t, err := json.Marshal(token)
	if err != nil {
		log.Fatalf("Could not json encode token: %s", err)
	}
	fmt.Printf("\n\nToken - Put this in the environment varibale 'LOCAL_TOKEN' (for example by including it in app.yaml):\n%s", t)
}
