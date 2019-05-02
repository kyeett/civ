package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/kyeett/ebitendrawutil"
	"github.com/oakmound/shiny/materialdesign/colornames"
	"github.com/peterhellberg/gfx"
)

type Button struct {
	gfx.Rect
	onClick func(*Button)
	onHover func(*Button)
	clicked bool
	hovered bool
}

func (b Button) ContainsAny(vs ...gfx.Vec) bool {
	for _, v := range vs {
		if b.Contains(v) {
			return true
		}
	}
	return false
}

func (b Button) render(screen *ebiten.Image) {
	if b.hovered {
		ebitenutil.DrawRect(screen, b.Min.X, b.Min.Y, b.W(), b.H(), colornames.White)
	} else {
		ebitenutil.DrawRect(screen, b.Min.X, b.Min.Y, b.W(), b.H(), colornames.Black)
		ebitendrawutil.DrawRect(screen, b.Rect, colornames.White, 2)
		ebitenutil.DebugPrintAt(screen, "Show Map", int(b.Min.X)+4, int(b.Min.Y)+6)
	}
}
