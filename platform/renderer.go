// +-------------------=D=r=i=v=e=-=E=n=g=i=n=e=---------------------+
// | Copyright (C) 2016-2017 Andreas T Jonsson. All rights reserved. |
// | Contact <mail@andreasjonsson.se>                                |
// +-----------------------------------------------------------------+

package platform

import (
	"image"
	"image/color"
)

type Renderer interface {
	Clear()
	Present()
	Shutdown()
	BackBuffer() *image.Paletted
	SetPalette(pal color.Palette)
	ToggleFullscreen()
	SetWindowTitle(title string)
}
