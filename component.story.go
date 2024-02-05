package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/kyoto-framework/news.kyoto.codes/hn"
	"go.kyoto.codes/v3/component"
)

// StoryComponentState is the state of the StoryComponent.
type StoryComponentState struct {
	component.Disposable
	hn.Story

	ID int
	Index int

	// Formatted fields
	URLFmt  string
	TimeFmt string
}

// StoryComponent is a feed single entry component.
func StoryComponent(id, index int) component.Component {
	return func(ctx *component.Context) component.State {
		// Initialize state
		state := &StoryComponentState{}

		// Pass arguments
		state.ID = id
		state.Index = index

		// Define formatters
		fmturl := func(u string) string {
			_u, _ := url.Parse(u)
			return _u.Host
		}
		fmttime := func(t int) string {
			passed := time.Since(time.Unix(int64(t), 0))
			if passed > 24*time.Hour {
				return fmt.Sprintf("%vd ago", int(passed.Hours()/24))
			} else if passed > 1*time.Hour {
				return fmt.Sprintf("%vh ago", int(passed.Hours()))
			} else if passed > 1*time.Minute {
				return fmt.Sprintf("%vm ago", int(passed.Minutes()))
			} else {
				return "just now"
			}
		}

		// Fetch story
		state.Story, _ = hn.FetchStory(id)

		// Format
		state.URLFmt = fmturl(state.URL)
		state.TimeFmt = fmttime(state.Time)

		// Return
		return state
	}
}
