package main

import (
	"net/url"

	"github.com/kyoto-framework/news.kyoto.codes/hn"
	"go.kyoto.codes/v3/component"
	"go.kyoto.codes/v3/rendering"
	"go.kyoto.codes/zen/v3/conv"
	"go.kyoto.codes/zen/v3/logic"
	"go.kyoto.codes/zen/v3/slice"
)

var (
	FEED_LIMIT = 60
)

// FeedPageState is the state for the FeedPage.
type FeedPageState struct {
	component.Disposable
	rendering.Template

	// Page state
	Page     int
	Query    string
	Category string

	// Next page url
	Next string

	// Components
	Navbar component.Future
	Stories []component.Future
}

// FeedPage is the component for the feed page.
func FeedPage(ctx *component.Context) component.State {
	// Initialize state
	state := &FeedPageState{}

	// Resolve query and category from the URL
	state.Page = conv.Int(logic.Or(ctx.Request.URL.Query().Get("page"), "1"))
	state.Query = ctx.Request.URL.Query().Get("query")
	state.Category = ctx.Request.URL.Query().Get("category")

	// Defaults
	if state.Category == "" {
		state.Category = "top"
	}

	// Fetch story ids
	storyids, err := hn.FetchStoryIds(
		state.Category,
		state.Query,
		[2]int{
			(state.Page - 1) * FEED_LIMIT,
			state.Page * FEED_LIMIT,
		},
	)
	if err != nil {
		panic(err)
	}

	// Prepare next page url
	nexturl, _ := url.Parse(ctx.Request.URL.String())
	q := nexturl.Query()
	q.Set("page", conv.String(state.Page+1))
	nexturl.RawQuery = q.Encode()
	state.Next = nexturl.String()

	// Init components
	state.Navbar = component.Use(ctx, NavbarComponent)
	state.Stories = slice.MapIndexed(storyids, func(i int, id int) component.Future {
		return component.Use(ctx, StoryComponent(id, (state.Page - 1) * FEED_LIMIT + (i+1)))
	})

	// Return
	return state
}
