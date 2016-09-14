package foo

// Middleware?

// Sessions

// Context data from middleware etc.

// Data to templates

// Create handlers that:
// - Have access to shared data without making it global
// - Avoids having to duplicate error handling (logging + http.Erro) in all handlers, just return instead (handle this in middleware?)

// Error handling, with http redirect for 404, 500 etc
//      https://golang.org/doc/articles/wiki/#tmp_9
//      https://codegangsta.gitbooks.io/building-web-apps-with-go/content/controllers/index.html
//      https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/11.1.html

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type templateData struct {
	TemplateFile  string
	Name          string
	PageLoadedCnt int
	Res           int
}

type templateMap map[string]*template.Template

type fooError struct {
	Origin   error
	Msg      string
	HTTPCode int
}

func (fe fooError) Error() string {
	return fmt.Sprintf("%d: %s (%s)", fe.HTTPCode, fe.Msg, fe.Origin.Error())
}

func init() {

	// Parse the templates, first the base and then all the other files.
	// Avoid a global variable, which is tempting but a pain when testing
	// Not convinces this way is worth it though.
	// must.. prefix is go go idiom to say that the function panics on errors, which is fine here
	templates := mustLoadTemplates("templates/", "templates/base.tmpl")

	n := negroni.New()

	// Session store. Uses the gorilla session lib, storing the data in an encrypted cookie
	store := cookiestore.New([]byte("^sHFTpot46653s4f654#^$322")) // secret used to encrypt the cookie
	n.Use(sessions.Sessions("foo", store))

	r := mux.NewRouter()
	r.Handle("/", makeHandler(handlePageMain, templates))

	n.UseHandler(r)
	http.Handle("/", n)

}

func makeHandler(fn func(http.ResponseWriter, *http.Request, templateMap) *fooError, templates templateMap) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fe := fn(w, r, templates)
		ctx := appengine.NewContext(r)
		log.Errorf(ctx, "!! %d: %s (%v)", fe.HTTPCode, fe.Msg, fe.Origin)
		http.Error(w, fe.Msg, fe.HTTPCode)
		// could this be moved to middleware?
	}
}

func handlePageMain(w http.ResponseWriter, r *http.Request, templates templateMap) *fooError {

	ctx := appengine.NewContext(r)
	log.Infof(ctx, "handlePageMain()")
	session := sessions.GetSession(r)

	// Session values
	var pageLoadedCnt int
	var ok bool
	if pageLoadedCnt, ok = session.Get("page_loaded_count").(int); ok {
		log.Infof(ctx, "Index: DID find page_loaded_count in session: %d", pageLoadedCnt)
	} else {
		log.Infof(ctx, "Index: did not find page_loaded_count in session")
	}

	pageLoadedCnt++
	session.Set("page_loaded_count", pageLoadedCnt)

	// Context data - In go 1.7 this can be added to r *http.Request, but ae does not have go 1.7 yet
	ctx2 := context.WithValue(ctx, "KEY", 44) // FIXME this is not how you're "suppose" to key :)

	// "Bussiness Logic" - aka do Stuff
	res, err := doStuff(ctx2, 2, 2)

	if err != nil {
		return err
	}

	data := templateData{
		"index.tmpl",
		"Main page",
		pageLoadedCnt,
		res,
	}

	err = renderTemplate(w, templates, data)

	if err != nil {
		return err
	}

	return nil

}

func renderTemplate(w http.ResponseWriter, templates templateMap, data templateData) *fooError {

	t, ok := templates[data.TemplateFile]
	if !ok {
		return &fooError{nil, fmt.Sprintf("User tried to load unknown template: %s", data.TemplateFile), http.StatusNotFound}
	}

	//log.Infof(ctx, "Loading template with data: %#v", data)

	err := t.ExecuteTemplate(w, "base", data) // base is the base template, defined in the template itself
	if err != nil {
		return &fooError{err, "", http.StatusInternalServerError}
	}
	return nil

}

func mustLoadTemplates(dir, baseTmpl string) templateMap {
	templates := make(map[string]*template.Template)
	tBase := template.Must(template.ParseFiles(baseTmpl))
	tFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		panic("Cant read template dir" + dir)
	}

	// Loop the rest of the template files
	for _, f := range tFiles {
		if f.Name() == baseTmpl {
			continue // already loaded above
		}
		// clone the base we already loaded
		t, err := tBase.Clone()
		if err != nil {
			panic("Cant clone template base")
		}
		// Create a new template set by parsing this template into the set already containing base
		templates[f.Name()] = template.Must(t.ParseFiles(filepath.Join(dir, f.Name())))
	}

	return templates

}
