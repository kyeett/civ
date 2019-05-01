package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/peterhellberg/gfx"
)

type CoordTranslator struct {
	anchor       gfx.Vec
	tileSize     gfx.Rect
	tileMiniSize gfx.Rect
}

func (ct CoordTranslator) toScreen(x, y int) gfx.Vec {
	dx, dy := ct.toScreenXY(x, y)
	return gfx.V(dx, dy)
}

func (ct CoordTranslator) toScreenXY(x, y int) (float64, float64) {
	return (ct.anchor.X + float64(x)*(ct.tileSize.W()-1) + float64(y)*(ct.tileSize.W()-1)/2),
		(ct.anchor.Y + float64(y)*9)
}

func (ct CoordTranslator) toMiniScreenXY(x, y int) (float64, float64) {
	return (ct.anchor.X + float64(x)*(ct.tileMiniSize.W()-1) + float64(y)*(ct.tileMiniSize.W()-1)/2),
		(ct.anchor.Y + float64(y)*3)
}

func (ct CoordTranslator) tileDrawOptions(x, y int) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	dx, dy := ct.toScreenXY(x, y)
	op.GeoM.Translate(dx, dy)
	return op
}

func (ct CoordTranslator) tileMiniDrawOptions(x, y int) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	dx, dy := ct.toMiniScreenXY(x, y)
	op.GeoM.Translate(dx, dy)
	return op
}
