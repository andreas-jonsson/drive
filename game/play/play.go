// +-------------------=D=r=i=v=e=-=E=n=g=i=n=e=---------------------+
// | Copyright (C) 2016-2017 Andreas T Jonsson. All rights reserved. |
// | Contact <mail@andreasjonsson.se>                                |
// +-----------------------------------------------------------------+

package play

import (
	"image"
	"log"

	"image/png"

	"github.com/andreas-jonsson/drive/game"
	"github.com/andreas-jonsson/drive/platform"
	"github.com/andreas-jonsson/drive/rasterizer"
	"github.com/andreas-jonsson/openwar/data"
)

type playState struct {
	testImage *image.Paletted
}

func NewPlayState() *playState {
	r, err := data.FS.Open("test.png")
	if err != nil {
		log.Panicln(err)
	}
	defer r.Close()

	img, err := png.Decode(r)
	if err != nil {
		log.Panicln(err)
	}

	return &playState{img.(*image.Paletted)}
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

func (s *playState) Render(backBuffer *image.Paletted) error {

	//ps := NewDefaultUVShader(backBuffer, s.testImage)
	//rasterizer.RasterizeUV(ps, 10, 10, 500, 10, 10, 500, 0, 0, 1, 0, 0, 1)

	ps := rasterizer.NewDefaultFlatShader(backBuffer)
	rasterizer.RasterizeFlat(ps, 10, 10, 500, 10, 10, 500, 10)

	return nil
}
