package main

import (
	"github.com/labstack/echo/v4"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

func pageHome(srv server, ec echo.Context, messages dbMessages) g.Node {

	div1 :=
		Div(
			H1(g.Text("Velkommen til opplagstavla")),

			form2(),

			// Button(
			// 	g.Attr("hx-get", "/"),
			// 	g.Attr("hx-swap", "outerHTML"),
			// 	g.Text("Click Me"),
			// ),
		)

	// var lis []g.Node

	// lis = append(lis,
	// 	g.Attr("hx-ext", "sse"),
	// 	g.Attr("sse-connect", "/sse-element"),
	// 	g.Attr("sse-swap", "tick-event"),
	// 	g.Attr("hx-swap", "afterbegin"),
	// )
	// for _, m := range messages {
	// 	lis = append(lis, Li(
	// 		g.Text(m.author),
	// 		g.Raw(" - "),
	// 		g.Text(m.msg),
	// 	))
	// }

	// Table with messages
	var trs = []g.Node{
		g.Attr("hx-ext", "sse"),
		g.Attr("sse-connect", "/sse-element"),
		g.Attr("sse-swap", "tick-event"),
		g.Attr("hx-swap", "afterbegin"),
	}

	for _, m := range messages {
		trs = append(trs,
			Tr(
				Td(g.Text(m.author)),
				Td(g.Text(m.msg)),
			),
		)
	}

	table := Table(
		THead(
			Tr(
				Th(g.Text("Fra")),
				Th(g.Text("Beskjed")),
			),
		),
		TBody(
			trs...,
		),
	)

	body := []g.Node{
		Main(
			div1,
			// Ul(lis...),
			H2(g.Text("Beskjeder:")),
			table,
		),
	}

	return page2(srv, ec, "Hjem", "/", body)

}

func pagePassword(srv server, ec echo.Context) g.Node {

	body := []g.Node{
		Main(
			// Ul(lis...),
			H2(g.Text("Passord:")),

			FormEl(
				// g.Attr("hx-post", "/"),
				Method("post"),
				g.Attr("hx-boost", "true"),
				g.Attr("hx-swap", "none"),
				Input(
					Type("password"),
					ID("password"),
					Name("password"),
					Placeholder("Passord.."),
					MaxLength("500"),
					g.Attr("size", "20"),
				),
				Input(
					Type("submit"),
					Value("Logg inn"),
				),
			),
		),
	}

	return page2(srv, ec, "Hjem", "/", body)

}

// take extra headers?
func page2(srv server, ec echo.Context, title, path string, body []g.Node) g.Node {

	// HTML5 boilerplate document
	return c.HTML5(c.HTML5Props{
		Title:    title,
		Language: "no",
		Head: []g.Node{
			// Link(Rel("stylesheet"), Href("https://unpkg.com/tailwindcss@2.1.2/dist/base.min.css")),
			// Link(Rel("stylesheet"), Href("https://unpkg.com/tailwindcss@2.1.2/dist/components.min.css")),
			// Link(Rel("stylesheet"), Href("https://unpkg.com/@tailwindcss/typography@0.4.0/dist/typography.min.css")),
			// Link(Rel("stylesheet"), Href("https://unpkg.com/tailwindcss@2.1.2/dist/utilities.min.css")),
			Link(Rel("stylesheet"), Href("https://unpkg.com/@picocss/pico@latest/css/pico.classless.min.css"), Type("text/css")),
			Script(Src("https://unpkg.com/htmx.org@1.7.0")),
			Script(Src("https://unpkg.com/htmx.org/dist/ext/sse.js"), Defer()),
		},
		Body: body,
	})
}

// func form1() g.Node {
// 	return FormEl(
// 		Action(""),
// 		Method("post"),
// 		Input(
// 			Type("text"),
// 			ID("author"),
// 			Name("author"),
// 			MaxLength("500"),
// 			g.Attr("size", "20"),
// 		),
// 		Input(
// 			Type("text"),
// 			ID("msg"),
// 			Name("msg"),
// 			MaxLength("500"),
// 			g.Attr("size", "100"),
// 		),
// 		Input(
// 			Type("submit"),
// 			Value("Store"),
// 		),
// 	)
// }

func form2() g.Node {
	return FormEl(
		// g.Attr("hx-post", "/"),
		Method("post"),
		g.Attr("hx-boost", "true"),
		g.Attr("hx-swap", "none"),
		Input(
			Type("text"),
			ID("author"),
			Name("author"),
			Placeholder("Fra.."),
			MaxLength("500"),
			g.Attr("size", "20"),
		),
		Input(
			Type("text"),
			ID("msg"),
			Name("msg"),
			Placeholder("Beskjed.."),
			MaxLength("500"),
			g.Attr("size", "100"),
		),
		Input(
			Type("submit"),
			Value("Legg til"),
		),

		Div(
			g.Attr("hx-get", "/htmx:sseError"),
			g.Attr("hx-trigger", "htmx:sseError"),
			g.Attr("hx-swap", "afterend"),
			g.Raw("htmx:sseError"),
		),

		Div(
			g.Attr("hx-get", "/htmx:onLoadError"),
			g.Attr("hx-swap", "afterend"),
			g.Attr("hx-trigger", "htmx:onLoadError"),
			g.Raw("htmx:onLoadError"),
		),

		Div(
			g.Attr("hx-get", "/htmx:timeout"),
			g.Attr("hx-swap", "afterend"),
			g.Attr("hx-trigger", "htmx:timeout"),
			g.Raw("htmx:timeout"),
		),

		Div(
			g.Attr("hx-get", "/htmx:responseError"),
			g.Attr("hx-swap", "afterend"),
			g.Attr("hx-trigger", "htmx:responseError"),
			g.Raw("htmx:responseError"),
		),
	)
}
