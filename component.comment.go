package main

import (
	"github.com/kyoto-framework/news.kyoto.codes/hn"
	"go.kyoto.codes/v3/component"
	"go.kyoto.codes/v3/rendering"
)

// CommentComponentState is the state of the CommentComponent.
type CommentComponentState struct {
	component.Universal
	rendering.Template
	hn.Comment

	ID       int
	Toggled  bool
	Comments component.Future `json:"-"`
}

// CommentComponent is a component that renders a comment.
// This is a recursive component that also renders children comments.
// Because of that, it have a skiphx option to skip hx-state decode form the request.
// Otherwise, the hx-state will be decoded also into children which will lead to infinite loop.
func CommentComponent(id int, skiphx bool) component.Component {
	return func(ctx *component.Context) component.State {
		// Initialize state
		state := &CommentComponentState{}
		// Pass arguments
		state.ID = id
		// Handle toggling
		if ctx.Request.Method == "POST" && !skiphx {
			// Unmarshal state
			state.Unmarshal(state, ctx.Request.FormValue("hx-state"))
			// Toggle
			state.Toggled = !state.Toggled
		}
		// Load comment if not loaded yet
		if state.Time == 0 {
			state.Comment, _ = hn.FetchComment(id)
		}
		// Load children.
		// We're passing skiphx to avoid incorrect state unmarshalling and infinite loop.
		if state.Toggled {
			state.Comments = component.Use(ctx, CommentsComponent(state.Children, true))
		}
		// Return
		return state
	}
}
