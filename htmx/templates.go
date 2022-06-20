package main

import "html/template"

func mustParseTemplates() *template.Template {

	return template.Must(template.ParseGlob("templates/*.html"))

}
