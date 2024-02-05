package main

import (
	"go.kyoto.codes/v3/component"
	"go.kyoto.codes/v3/rendering"
	"go.kyoto.codes/zen/v3/slice"
)

// CommentsComponentState is the state of the CommentsComponent.
type CommentsComponentState struct {
	component.Disposable
	rendering.Template

	IDs      []int
	Comments []component.Future
}

// CommentsComponent is a component that renders a list of comments.
func CommentsComponent(ids []int, skiphx bool) component.Component {
	return func(ctx *component.Context) component.State {
		// Initialize state
		state := &CommentsComponentState{}
		// Pass arguments
		state.IDs = ids
		// Init components
		state.Comments = slice.Map(state.IDs, func(id int) component.Future {
			return component.Use(ctx, CommentComponent(id, skiphx))
		})
		// Return
		return state
	}
}
