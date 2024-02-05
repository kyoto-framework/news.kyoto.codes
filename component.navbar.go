package main

import "go.kyoto.codes/v3/component"

type NavbarComponentState struct {
	component.Disposable
}

// NavbarComponent is a navbar component.
// It doesn't have any state, it's just a static component now.
func NavbarComponent(ctx *component.Context) component.State {
	state := &NavbarComponentState{}
	return state
}
