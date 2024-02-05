package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"

	"go.kyoto.codes/v3/rendering"
	"go.kyoto.codes/zen/v3/mapx"
	"go.kyoto.codes/zen/v3/templatex"
)

// setup prepares globals and application mux.
func setup() *http.ServeMux {
	// Setup kyoto
	rendering.TEMPLATE_FUNCMAP = mapx.Merge(
		rendering.TEMPLATE_FUNCMAP,
		templatex.FuncMap,
		template.FuncMap{
			"html": func (s string) template.HTML { return template.HTML(s) },
		},
	)

	// Initialize mux
	mux := http.NewServeMux()

	// Setup assets
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets-dist"))))

	// Setup pages
	mux.Handle("/", rendering.Handler(FeedPage))
	mux.Handle("/story/", rendering.Handler(StoryPage))

	// Setup components (htmx)
	mux.Handle("/htmx/comment", rendering.Handler(CommentComponent(0, false)))

	// Return mux
	return mux
}

// main is the entrypoint of the application.
func main() {
	// Parse arguments
	addr := flag.String("http", ":8000", "Serving address")
	flag.Parse()

	// Serve
	fmt.Printf("Serving on %s\n", *addr)
	if err := http.ListenAndServe(*addr, setup()); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
