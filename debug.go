package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func printDebug(screen *ebiten.Image) {
	for x := 0; x < world.Width; x++ {
		dx, dy := coordTranslator.toScreenXY(x, 0)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", x), int(dx)-2, int(dy)-14)
	}

	for y := 0; y < world.Height; y++ {
		dx, dy := coordTranslator.toScreenXY(0, y)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", y), int(dx)-8, int(dy))
	}

	c := cursorPosition()
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("%v", c))
	c2 := cursorRelativePosition()
	c3 := co
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%v", c2), screenWidth-100, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%v", c2), screenWidth-100, screenHeight-20)

	ebitenutil.DrawRect(screen, c.X, c.Y, 3, 3, color.White)
}
