package main

import (
	"go.kyoto.codes/v3/component"
	"go.kyoto.codes/v3/rendering"
	"go.kyoto.codes/zen/v3/conv"
	"go.kyoto.codes/zen/v3/httpx"
)

// StoryPageState is the state for the StoryPage.
type StoryPageState struct {
	component.Disposable
	rendering.Template

	Navbar   component.Future
	Story    component.Future
	Comments component.Future
}

// StoryPage is the component for the story page.
// It includes the story itself and the comments.
func StoryPage(ctx *component.Context) component.State {
	// Initialize state
	state := &StoryPageState{}

	// Resolve story id from the URL
	storyid := conv.Int(httpx.Path(ctx.Request.URL.Path).GetAfter("story"))

	// Initialize components
	state.Navbar = component.Use(ctx, NavbarComponent)
	state.Story = component.Use(ctx, StoryComponent(storyid, 0))
	state.Comments = component.Use(ctx, CommentsComponent(state.Story().(*StoryComponentState).Children, false))

	// Return
	return state
}
