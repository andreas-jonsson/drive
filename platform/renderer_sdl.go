// +build !mobile

// +-------------------=D=r=i=v=e=-=E=n=g=i=n=e=---------------------+
// | Copyright (C) 2016-2017 Andreas T Jonsson. All rights reserved. |
// | Contact <mail@andreasjonsson.se>                                |
// +-----------------------------------------------------------------+

package platform

import (
	"image"
	"image/color"
	"image/color/palette"
	"log"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const fullscreenFlag = sdl.WINDOW_FULLSCREEN //sdl.WINDOW_FULLSCREEN_DESKTOP

type Config func(*sdlRenderer) error

func ConfigWithSize(w, h int) Config {
	return func(rnd *sdlRenderer) error {
		rnd.config.windowSize = image.Point{w, h}
		return nil
	}
}

func ConfigWithTitle(title string) Config {
	return func(rnd *sdlRenderer) error {
		rnd.config.windowTitle = title
		return nil
	}
}

func ConfigWithDiv(n int) Config {
	return func(rnd *sdlRenderer) error {
		rnd.config.resolutionDiv = n
		return nil
	}
}

func ConfigWithFullscreen(rnd *sdlRenderer) error {
	rnd.config.fullscreen = true
	return nil
}

func ConfigWithDebug(rnd *sdlRenderer) error {
	rnd.config.debug = true
	return nil
}

func ConfigWithNoVSync(rnd *sdlRenderer) error {
	rnd.config.novsync = true
	return nil
}

type sdlRenderer struct {
	window           *sdl.Window
	backBuffer       *image.Paletted
	hwBuffer         *sdl.Texture
	internalRenderer *sdl.Renderer
	paletteLookup    [256]uint32

	config struct {
		windowTitle   string
		windowSize    image.Point
		resolutionDiv int
		debug, novsync,
		fullscreen bool
	}
}

func NewRenderer(configs ...Config) (*sdlRenderer, error) {
	var (
		err error
		r   sdlRenderer
		dm  sdl.DisplayMode

		sdlFlags uint32 = sdl.WINDOW_SHOWN
	)

	for _, cfg := range configs {
		if err = cfg(&r); err != nil {
			return nil, err
		}
	}

	cfg := &r.config
	if cfg.fullscreen {
		sdlFlags |= fullscreenFlag
	}

	if err = sdl.GetDesktopDisplayMode(0, &dm); err != nil {
		return nil, err
	}

	if cfg.windowSize.X <= 0 {
		cfg.windowSize.X = int(dm.W)
	}
	if cfg.windowSize.Y <= 0 {
		cfg.windowSize.Y = int(dm.H)
	}

	if cfg.resolutionDiv > 0 {
		cfg.windowSize.X /= cfg.resolutionDiv
		cfg.windowSize.Y /= cfg.resolutionDiv
	}

	r.window, err = sdl.CreateWindow(cfg.windowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, cfg.windowSize.X, cfg.windowSize.Y, sdlFlags)
	if err != nil {
		return nil, err
	}

	const width, height = 320, 200
	r.backBuffer = image.NewPaletted(image.Rect(0, 0, width, height), palette.Plan9)
	r.SetPalette(palette.Plan9)

	renderer, err := sdl.CreateRenderer(r.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err
	}

	renderer.SetLogicalSize(width, height)
	r.internalRenderer = renderer

	r.hwBuffer, err = renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		return nil, err
	}

	sdl.ShowCursor(0)
	return &r, nil
}

func (r *sdlRenderer) ToggleFullscreen() {
	isFullscreen := (r.window.GetFlags() & fullscreenFlag) != 0
	if isFullscreen {
		r.window.SetFullscreen(0)
	} else {
		r.window.SetFullscreen(fullscreenFlag)
	}
}

func (r *sdlRenderer) BackBuffer() *image.Paletted {
	return r.backBuffer
}

func (r *sdlRenderer) SetPalette(pal color.Palette) {
	r.backBuffer.Palette = pal
	for i, c := range pal {
		cr, cg, cb, _ := c.RGBA()
		r.paletteLookup[i] = cr<<24 | cg<<16 | cb<<8
	}
}
func (r *sdlRenderer) Clear() {
	pix := r.backBuffer.Pix
	black := r.backBuffer.Palette.Index(color.RGBA{0, 0, 0, 255})

	for i := range pix {
		pix[i] = uint8(black)
	}
}

func (r *sdlRenderer) Present() {
	var (
		p     unsafe.Pointer
		pitch int
	)

	if err := r.hwBuffer.Lock(nil, &p, &pitch); err != nil {
		log.Panicln(err)
	}

	for _, idx := range r.backBuffer.Pix {
		*(*uint32)(p) = r.paletteLookup[idx]
		p = unsafe.Pointer(uintptr(p) + 4)
	}

	r.hwBuffer.Unlock()
	r.internalRenderer.Copy(r.hwBuffer, nil, nil)
	r.internalRenderer.Present()
}

func (r *sdlRenderer) Shutdown() {
	r.window.Destroy()
	r.hwBuffer.Destroy()
	r.internalRenderer.Destroy()
}

func (r *sdlRenderer) SetWindowTitle(title string) {
	r.window.SetTitle(title)
}
