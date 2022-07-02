package main

import (
	"github.com/labstack/echo/v4"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

func pageHome(srv server, ec echo.Context, messages dbMessages) g.Node {

	var body []g.Node

	body = append(body, Div(
		H1(g.Text("Welcome to this example page")),
		// 		P(g.Text("I hope it will make you happy. 😄 It's using TailwindCSS for styling.")),

		form2(),

		Button(
			g.Attr("hx-get", "/"),
			g.Attr("hx-swap", "outerHTML"),
			g.Text("Click Me"),
		),
	))

	var lis []g.Node

	lis = append(lis,
		g.Attr("hx-ext", "sse"),
		g.Attr("sse-connect", "/sse-element"),
		g.Attr("sse-swap", "tick-event"),
		g.Attr("hx-swap", "afterbegin"),
	)
	for _, m := range messages {
		lis = append(lis, Li(
			g.Text(m.author),
			g.Raw(" - "),
			g.Text(m.msg),
		))
	}

	body = append(body, Ul(lis...))

	return page2(srv, ec, "Hjem", "/", body)

}

// take extra headers?
func page2(srv server, ec echo.Context, title, path string, body []g.Node) g.Node {

	// HTML5 boilerplate document
	return c.HTML5(c.HTML5Props{
		Title:    "Hjem",
		Language: "no",
		Head: []g.Node{
			// Link(Rel("stylesheet"), Href("https://unpkg.com/tailwindcss@2.1.2/dist/base.min.css")),
			// Link(Rel("stylesheet"), Href("https://unpkg.com/tailwindcss@2.1.2/dist/components.min.css")),
			// Link(Rel("stylesheet"), Href("https://unpkg.com/@tailwindcss/typography@0.4.0/dist/typography.min.css")),
			// Link(Rel("stylesheet"), Href("https://unpkg.com/tailwindcss@2.1.2/dist/utilities.min.css")),
			Script(Src("https://unpkg.com/htmx.org@1.7.0")),
			Script(Src("https://unpkg.com/htmx.org/dist/ext/sse.js"), Defer()),
		},
		Body: body,
	})
}

func form1() g.Node {
	return FormEl(
		Action(""),
		Method("post"),
		Input(
			Type("text"),
			ID("author"),
			Name("author"),
			MaxLength("500"),
			g.Attr("size", "20"),
		),
		Input(
			Type("text"),
			ID("msg"),
			Name("msg"),
			MaxLength("500"),
			g.Attr("size", "100"),
		),
		Input(
			Type("submit"),
			Value("Store"),
		),
	)
}

func form2() g.Node {
	return FormEl(
		g.Attr("hx-post", "/"),
		g.Attr("hx-swap", "none"),
		Input(
			Type("text"),
			ID("author"),
			Name("author"),
			MaxLength("500"),
			g.Attr("size", "20"),
		),
		Input(
			Type("text"),
			ID("msg"),
			Name("msg"),
			MaxLength("500"),
			g.Attr("size", "100"),
		),
		Input(
			Type("submit"),
			Value("Store"),
		),
	)
}
