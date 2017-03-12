package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine/log"
)

//
// Copied from https://developers.google.com/sheets/api/quickstart/go, modified to fit appengine standard

// use GCP deffault credentials when running in App Engine. Access is granted because the Sheet is shared with the GCP service account (PROJECT-ID@appspot.gserviceaccount.com)
func getClient(ctx context.Context) (*http.Client, error) {
	client, err := google.DefaultClient(ctx, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return nil, err
	}

	return client, nil
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient2(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		msg := "Unable to get path to cached credential file"
		log.Errorf(ctx, "%s: %v", msg, err)
		//return nil, nil, fooError{Origin: err, Msg: "", HTTPCode: http.StatusInternalServerError}
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(ctx, config)
		saveToken(ctx, cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		msg := "Unable to read authorization code"
		log.Errorf(ctx, "%s: %v", msg, err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		msg := "Unable to retrieve token from web"
		log.Errorf(ctx, "%s: %v", msg, err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("sheets.googleapis.com-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(ctx context.Context, file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		msg := "Unable to cache oauth token"
		log.Errorf(ctx, "%s: %v", msg, err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
