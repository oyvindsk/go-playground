package app

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/gorilla/mux"
	"google.golang.org/api/drive/v3"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const (
	oauth2State = FIXME
)

var (
	conf *oauth2.Config
)

// https://cloud.google.com/docs/authentication
// https://developers.google.com/identity/protocols/OAuth2
// https://github.com/google/google-api-go-client/blob/master/GettingStarted.md
// https://gist.github.com/jfcote87/89eca3032cd5f9705ba3
// https://play.golang.org/p/GwfnKWxUWj

// https://godoc.org/golang.org/x/oauth2/google
// https://godoc.org/golang.org/x/oauth2#Token

type oauth2Token struct {
	// AccessToken is the token that authorizes and authenticates
	// the requests.
	AccessToken string `json:"access_token"`

	// TokenType is the type of token.
	// The Type method returns either this or "Bearer", the default.
	TokenType string `json:"token_type,omitempty"`

	// RefreshToken is a token that's used by the application
	// (as opposed to the user) to refresh the access token
	// if it expires.
	RefreshToken string `json:"refresh_token,omitempty"`

	// Expiry is the optional expiration time of the access token.
	//
	// If zero, TokenSource implementations will reuse the same
	// token forever and RefreshToken or equivalent
	// mechanisms for that TokenSource will not be used.
	Expiry time.Time `json:"expiry,omitempty"`

	// raw optionally contains extra metadata from the server
	// when updating a token.
	//raw interface{}
}

// Init, run by AE
func init() {

	conf = &oauth2.Config{
		ClientID:     FIXME
		ClientSecret: FIXME
		Endpoint:     google.Endpoint,
		Scopes:       []string{drive.DriveScope},
		RedirectURL:  "http://localhost:8080/oauth2callback",
	}

	r := mux.NewRouter()
	r.Handle("/", makeHandler(handlePageMain))
	r.Handle("/oauth2callback", makeHandler(handleOauth2allback)).Queries("state", "{state}", "code", "{code}")
	// google-auth-redir
	http.Handle("/", r)
}

// Page Handler
func handlePageMain(w http.ResponseWriter, r *http.Request) error {
	ctx := appengine.NewContext(r)

	//fmt.Fprintf(w, "COnf er: %#v", conf)

	// get  google token
	log.Debugf(ctx, "Looking for Oauth2 token: foo")
	token := new(oauth2Token)
	tk := datastore.NewKey(ctx, "Oauth2Tokens", "foo", 0, nil)
	err := datastore.Get(ctx, tk, token)
	if err != nil {
		if err != datastore.ErrNoSuchEntity {
			return err
		}

		// No token, get one by redirecting the user
		url := conf.AuthCodeURL(oauth2State, oauth2.AccessTypeOffline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return nil
	}

	log.Debugf(ctx, "Found Oauth2 Token: %+v", token)

	client := conf.Client(ctx, &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry, // time.Time `json:"expiry,omitempty"`
	})

	service, err := drive.New(client)
	if err != nil {
		return err // "Unable to create Drive service: %v", err)
	}

	f := driveFields("user")
	a, e := service.About.Get().Do(f)
	log.Debugf(ctx, "About:\n%+v\n\nerr:\n%v\n\n", a.User, e)

	return nil
}

type driveFields string

func (a driveFields) Get() (string, string) { return "fields", string(a) }

func handleOauth2allback(w http.ResponseWriter, r *http.Request) error {
	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.

	ctx := appengine.NewContext(r)

	//if conf == nil {
	//	doOauth2Setup(ctx, w, r)
	//}

	uv := mux.Vars(r)
	log.Debugf(ctx, "URL VARS:\n%+v", uv["code"])

	tok, err := conf.Exchange(ctx, uv["code"])
	if err != nil {
		return err
	}

	token := &oauth2Token{
		AccessToken:  tok.AccessToken,
		TokenType:    tok.TokenType,
		RefreshToken: tok.RefreshToken,
		Expiry:       tok.Expiry,
	}

	log.Debugf(ctx, "Setting Oauth2 token: foo to:\n%#v\n\n", token)

	tk := datastore.NewKey(ctx, "Oauth2Tokens", "foo", 0, nil)
	_, err = datastore.Put(ctx, tk, token)
	if err != nil {
		return err
	}
	//client = conf.Client(ctx, tok)

	return nil
}

// Helper functions
func makeHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		err := fn(w, r)

		if err != nil {
			if ferr, ok := err.(fooError); ok {
				log.Errorf(ctx, "!! %d: %s (%v)", ferr.HTTPCode, ferr.Msg, ferr.Origin)
				http.Error(w, ferr.Msg, ferr.HTTPCode)
				// could this be moved to middleware?
			} else {
				log.Errorf(ctx, "!! Error: %s", err)
			}
		}
	}
}

// Custom errror type
type fooError struct {
	Origin   error
	Msg      string
	HTTPCode int
}

func (fe fooError) Error() string {
	return fmt.Sprintf("%d: %s (%s)", fe.HTTPCode, fe.Msg, fe.Origin.Error())
}
