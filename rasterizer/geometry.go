// +-------------------=D=r=i=v=e=-=E=n=g=i=n=e=---------------------+
// | Copyright (C) 2016-2017 Andreas T Jonsson. All rights reserved. |
// | Contact <mail@andreasjonsson.se>                                |
// +-----------------------------------------------------------------+

package rasterizer

import (
	"image"
	"sync"

	"github.com/andreas-jonsson/warp/platform"
	"github.com/ungerik/go3d/mat4"
	"github.com/ungerik/go3d/vec2"
	"github.com/ungerik/go3d/vec3"
)

const (
	drawCallBufferSize = 255
	triangleBufferSize = 256
)

type triangle struct {
	a, b, c       vec3.T
	uva, uvb, uvc vec2.T
	texture       *image.Paletted
	color         uint8
}

type drawCall struct {
	id      uint64
	vert    []vec3.T
	uvs     []vec2.T
	colors  []uint8
	texture *image.Paletted
	mvp     mat4.T
}

type Rasterizer struct {
	drawCallChan chan drawCall
	triangleChan chan triangle
	workerWG     sync.WaitGroup
	target       *image.Paletted
}

func NewRasterizer(backBuffer *image.Paletted) *Rasterizer {
	r := &Rasterizer{
		target: backBuffer,
		make(chan drawCall, triangleBufferSize),
		make(chan triangle, triangleBufferSize),
	}

	workerWG.Add(2)

	go func() {
		var tri triangle

		for dc := range r.drawCallChan {
			numVert := len(dc.vert)

			for i := 0; i < numVert; i += 3 {
				tri.a = dc.mvp.MulVec3(&dc.vert[i])
				tri.b = dc.mvp.MulVec3(&dc.vert[i+1])
				tri.c = dc.mvp.MulVec3(&dc.vert[i+2])

				if dc.texture == nil {
					tri.color = dc.colors[i/3]
					multiSwap(&tri.a[0], &tri.a[1], &tri.b[0], &tri.b[1], &tri.c[0], &tri.c[1])
				} else {
					tri.uva = dc.uvs[i]
					tri.uvb = dc.uvs[i+1]
					tri.uvc = dc.uvs[i+2]
					tri.texture = dc.texture

					multiSwapUV(
						&tri.a[0], &tri.a[1], &tri.b[0], &tri.b[1], &tri.c[0], &tri.c[1],
						&tri.uva[0], &tri.uva[1], &tri.uvb[0], &tri.uvb[1], &tri.uvc[0], &tri.uvc[1],
					)
				}

				r.triangleChan <- tri
			}
		}

		close(r.triangleChan)
		r.workerWG.Done()
	}()

	go func() {
		for tri := range r.triangleChan {
			if tri.texture == nil {
				r.rasterizeFlat(
					int(tri.a[0]), int(tri.a[1]),
					int(tri.b[0]), int(tri.b[1]),
					int(tri.c[0]), int(tri.c[1]),
					tri.color,
				)
			} else {
				r.rasterizeTextured(
					int(tri.a[0]), int(tri.a[1]),
					int(tri.b[0]), int(tri.b[1]),
					int(tri.c[0]), int(tri.c[1]),
					tri.uva[0], tri.uva[1],
					tri.uvb[0], tri.uvb[1],
					tri.uvc[0], tri.uvc[1],
					tri.texture,
				)
			}
		}
		r.workerWG.Done()
	}()

	return r
}

func (r *Rasterizer) Wait(id uint64) {
	//TODO Implement this!
	panic("not implemented")
}

func (r *Rasterizer) Sync() {
	r.Wait(r.DrawFlat(mat4.Ident, nil, nil))
}

func (r *Rasterizer) Destroy() {
	close(r.drawCallChan)
	r.workerWG.Wait()
}

func (r *Rasterizer) DrawTextured(mvp *mat4.T, vert []vec3.T, uvs []vec2.T, texture *image.Paletted) uint64 {
	id := platform.NewId64()
	r.drawCallChan <- drawCall{id: id, mvp: mvp, vert: vert, uvs: uvs, texture: texture}
	return id
}

func (r *Rasterizer) DrawFlat(mvp *mat4.T, vert []vec3.T, colors []uint8) uint64 {
	id := platform.NewId64()
	r.drawCallChan <- drawCall{id: id, mvp: mvp, vert: vert, colors: colors}
	return id
}

func swapInt(a, b *int) {
	tmp := *a
	*a = *b
	*b = tmp
}

func swapFloat32(a, b *float32) {
	tmp := *a
	*a = *b
	*b = tmp
}

func multiSwap(x0, y0, x1, y1, x2, y2 *int) {
	if *y1 < *y0 {
		swapInt(y1, y0)
		swapInt(x1, x0)
	}

	if *y2 < *y0 {
		swapInt(y2, y0)
		swapInt(x2, x0)
	}

	if *y1 < *y2 {
		swapInt(y2, y1)
		swapInt(x2, x1)
	}
}

func multiSwapUV(x0, y0, x1, y1, x2, y2 *int, u0, v0, u1, v1, u2, v2 *float32) {
	if *y1 < *y0 {
		swapInt(y1, y0)
		swapInt(x1, x0)
		swapFloat32(u1, u0)
		swapFloat32(v1, v0)
	}

	if *y2 < *y0 {
		swapInt(y2, y0)
		swapInt(x2, x0)
		swapFloat32(u2, u0)
		swapFloat32(v2, v0)
	}

	if *y1 < *y2 {
		swapInt(y2, y1)
		swapInt(x2, x1)
		swapFloat32(u2, u1)
		swapFloat32(v2, v1)
	}
}
