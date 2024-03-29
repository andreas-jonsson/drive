// +-------------------=D=r=i=v=e=-=E=n=g=i=n=e=---------------------+
// | Copyright (C) 2016-2017 Andreas T Jonsson. All rights reserved. |
// | Contact <mail@andreasjonsson.se>                                |
// +-----------------------------------------------------------------+

package rasterizer

import "image"

func (r *Rasterizer) pixelShader(x, y int, u, v float32, texture *image.Paletted) {
	textureSize := texture.Bounds().Max
	maxX := textureSize.X - 1
	maxY := textureSize.Y - 1

	tx := int(u * float32(maxX))
	ty := int(v * float32(maxY))

	if tx > maxX {
		tx = maxX
	} else if tx < 0 {
		tx = 0
	}

	if ty > maxY {
		ty = maxY
	} else if ty < 0 {
		ty = 0
	}

	r.target.SetColorIndex(x, y, texture.ColorIndexAt(tx, ty))
}

func (r *Rasterizer) rasterizeTextured(x0, y0, x1, y1, x2, y2 int, u0, v0, u1, v1, u2, v2 float32, texture *image.Paletted) {
	// Declare some variables that we'll use and where starting from y0 at the
	// top of the triangle
	dxdy1 := float32(x2 - x0)
	dxdu1 := u2 - u0
	dxdv1 := v2 - v0

	dxdy2 := float32(x1 - x0)
	dxdu2 := u1 - v0
	dxdv2 := v1 - v0

	var (
		sdx, sdu, sdv,
		edx, edu, edv,
		pu, pv float32
	)

	dy1 := float32(y2 - y0)
	dy2 := float32(y1 - y0)

	// Check for divide by zero's
	if y2-y0 != 0 {
		dxdy1 /= dy1
		dxdu1 /= dy1
		dxdv1 /= dy1
	}

	if y1-y0 != 0 {
		dxdy2 /= dy2
		dxdu2 /= dy2
		dxdv2 /= dy2
	}

	var (
		dxldy, dxrdy,
		dxldu, dxrdu,
		dxldv, dxrdv float32
	)

	// Determine our left and right points for our x value gradient..
	// e.g. the starting and ending line for our x inner loop
	if dxdy1 < dxdy2 {
		dxldy = dxdy1
		dxrdy = dxdy2
		dxldu = dxdu1
		dxrdu = dxdu2
		dxldv = dxdv1
		dxrdv = dxdv2
	} else {
		dxldy = dxdy2
		dxrdy = dxdy1
		dxldu = dxdu2
		dxrdu = dxdu1
		dxldv = dxdv2
		dxrdv = dxdv1
	}

	// Initial starting x and ending x is sdx and edx which are x0,y0...the
	// top of our triangle
	sdx = float32(x0)
	sdu = u0
	sdv = v0

	edx = float32(x0)
	edu = u0
	edv = v0

	var (
		pDeltaU,
		pDeltaV float32
	)

	for y := y0; y <= y2; y++ {
		pDeltaU = edu - sdu
		pDeltaV = edv - sdv

		if edx-sdx != 0 {
			pDeltaU /= edx - sdx
			pDeltaV /= edx - sdx
		}

		pu = sdu
		pv = sdv

		for x := int(sdx); x <= int(edx); x++ {
			pixelShader(x, y, pu, pv, texture)
			pu += pDeltaU
			pv += pDeltaV
		}

		sdx += dxldy
		sdu += dxldu
		sdv += dxldv
		edx += dxrdy
		edu += dxrdu
		edv += dxrdv
	}

	// Render bottom part of triangle.

	if dxdy1 < dxdy2 {
		dxldy = float32(x1 - x2)
		dxldu = u1 - u2
		dxldv = v1 - v2

		if y1-y2 != 0 {
			dxldy /= float32(y1 - y2)
			dxldu /= float32(y1 - y2)
			dxldv /= float32(y1 - y2)
		}

		sdx = float32(x2)
		sdu = u2
		sdv = v2
	} else {
		dxrdy = float32(x1 - x2)
		dxrdu = u1 - u2
		dxrdv = v1 - v2

		if y1-y2 != 0 {
			dxrdy /= float32(y1 - y2)
			dxrdu /= float32(y1 - y2)
			dxrdv /= float32(y1 - y2)
		}

		edx = float32(x2)
		edu = u2
		edv = v2
	}

	for y := y2; y <= y1; y++ {
		pDeltaU = edu - sdu
		pDeltaV = edv - sdv

		if edx-sdx != 0 {
			pDeltaU /= edx - sdx
			pDeltaV /= edx - sdx
		}

		pu = sdu
		pv = sdv

		for x := int(sdx); x <= int(edx); x++ {
			pixelShader(x, y, pu, pv, texture)
			pu += pDeltaU
			pv += pDeltaV
		}

		sdx += dxldy
		sdu += dxldu
		sdv += dxldv
		edx += dxrdy
		edu += dxrdu
		edv += dxrdv
	}
}

func (r *Rasterizer) rasterizeFlat(x0, y0, x1, y1, x2, y2 int, color uint8) {
	dxdy1 := float32(x2 - x0)
	dxdy2 := float32(x1 - x0)

	var (
		sdx, edx float32
	)

	dy1 := float32(y2 - y0)
	dy2 := float32(y1 - y0)

	// Check for divide by zero's
	if y2-y0 != 0 {
		dxdy1 /= dy1
	}

	if y1-y0 != 0 {
		dxdy2 /= dy2
	}

	var (
		dxldy, dxrdy float32
	)

	if dxdy1 < dxdy2 {
		dxldy = dxdy1
		dxrdy = dxdy2
	} else {
		dxldy = dxdy2
		dxrdy = dxdy1
	}

	sdx = float32(x0)
	edx = float32(x0)

	for y := y0; y <= y2; y++ {
		for x := int(sdx); x <= int(edx); x++ {
			r.target.SetColorIndex(x, y, color)
		}

		sdx += dxldy
		edx += dxrdy
	}

	// Render bottom part of triangle.

	if dxdy1 < dxdy2 {
		dxldy = float32(x1 - x2)
		if y1-y2 != 0 {
			dxldy /= float32(y1 - y2)
		}
		sdx = float32(x2)
	} else {
		dxrdy = float32(x1 - x2)
		if y1-y2 != 0 {
			dxrdy /= float32(y1 - y2)
		}
		edx = float32(x2)
	}

	for y := y2; y <= y1; y++ {
		for x := int(sdx); x <= int(edx); x++ {
			r.target.SetColorIndex(x, y, color)
		}

		sdx += dxldy
		edx += dxrdy
	}
}
