// +-------------------=D=r=i=v=e=-=E=n=g=i=n=e=---------------------+
// | Copyright (C) 2016-2017 Andreas T Jonsson. All rights reserved. |
// | Contact <mail@andreasjonsson.se>                                |
// +-----------------------------------------------------------------+

package play

import (
	"github.com/andreas-jonsson/drive/game"
	"github.com/andreas-jonsson/drive/platform"
)

type playState struct {
}

func NewPlayState() *playState {
	return &playState{}
}

func (s *playState) Name() string {
	return "play"
}

func (s *playState) Enter(from game.GameState, args ...interface{}) error {
	return nil
}

func (s *playState) Exit(to game.GameState) error {
	return nil
}

func (s *playState) Update(gctl game.GameControl) error {
	for event := gctl.PollEvent(); event != nil; event = gctl.PollEvent() {
		switch event.(type) {
		case *platform.MouseButtonEvent, *platform.QuitEvent:
			gctl.Terminate()
		}
	}
	return nil
}

func (s *playState) Render() error {
	return nil
}
