package app

import (
	"encoding/json"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

//
// Copied from https://developers.google.com/sheets/api/quickstart/go, modified to fit appengine standard

var scope = "https://www.googleapis.com/auth/spreadsheets.readonly"

// use GCP deffault credentials when running in App Engine. Access is granted because the Sheet is shared with the GCP service account (PROJECT-ID@appspot.gserviceaccount.com)
func getClient(ctx context.Context) (*http.Client, error) {

	client, err := google.DefaultClient(ctx, scope)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// getClientLocal uses the token stored in the environment varibale LOCAL_TOKEN to generate a oauth2 config nd then a http client
// Mostly used for local testing since the above function fails when running locally (gcloud bug?)
// Run getLocalDevToken.go to get this token (as json)
func getClientLocal(ctx context.Context) (*http.Client, error) {

	// get a oaut2 config with the stored credentials
	cs := os.Getenv("LOCAL_CLIENT_SECRET")
	log.Infof(ctx, "Found LOCAL_CLIENT_SECRET: %s", cs)
	config, err := google.ConfigFromJSON([]byte(cs), "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return nil, fooError{Origin: err, Msg: "Unable to parse client secret from environment variable LOCAL_CLIENT_SECRET", HTTPCode: http.StatusInternalServerError}
	}

	tok := os.Getenv("LOCAL_TOKEN")
	log.Infof(ctx, "Found LOCAL_TOKEN: %s", tok)
	t := &oauth2.Token{}
	err = json.Unmarshal([]byte(tok), t)
	if err != nil {
		return nil, fooError{Origin: err, Msg: "Unable to parse token from environment variable LOCAL_TOKEN", HTTPCode: http.StatusInternalServerError}
	}

	return config.Client(ctx, t), nil
}

func newSheetsClient(ctx context.Context) (*sheets.Service, error) {
	// Fallback to use api oauth2 3legged creds when running locally. Default Credentials should work locally as well, but there's a gcloud bug .. (shocking =)
	var err error
	var client *http.Client
	if appengine.IsDevAppServer() {
		client, err = getClientLocal(ctx)
	} else {
		client, err = getClient(ctx)
	}
	if err != nil {
		log.Errorf(ctx, "Unable to retrieve Sheets Client %v", err)
		return nil, fooError{Origin: err, Msg: "Unable to retrieve Sheets Client", HTTPCode: http.StatusInternalServerError}
	}

	srv, err := sheets.New(client)
	if err != nil {
		log.Errorf(ctx, "Unable to retrieve Sheets Service %v", err)
		return nil, fooError{Origin: err, Msg: "Unable to retrieve Sheets Service", HTTPCode: http.StatusInternalServerError}
	}

	return srv, nil
}
