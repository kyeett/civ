package main

import (
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/ebitenutil"
	"golang.org/x/image/colornames"

	hextiles "github.com/kyeett/animex/resources/hex"
	"github.com/kyeett/civ/tile"
	"github.com/peterhellberg/gfx"

	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/ebitendrawutil"
)

const (
	screenWidth, screenHeight = 240, 180
)

func getTileSprite(screen *ebiten.Image, x, y int, typ tile.Tile) *ebiten.Image {

	var ix, iy int
	switch typ {
	case tile.Desert:
		ix, iy = 1, 0
	case tile.Grass:
		ix, iy = 2, 0
	case tile.Water:
		ix, iy = 3, 2
	default:
		log.Fatal("renderTile: tile not supported")
	}

	sourcePt := image.Pt(4+18*ix, 183+19*iy)
	sourceRect := coordTranslator.tileSize.Bounds().Add(sourcePt)
	return tilesImg.SubImage(sourceRect).(*ebiten.Image)
}

func getMiniTileSprite(screen *ebiten.Image, x, y int, typ tile.Tile) *ebiten.Image {

	var ix, iy int
	switch typ {
	case tile.Desert:
		ix, iy = 1, 2
	case tile.Grass:
		ix, iy = 2, 0
	case tile.Water:
		ix, iy = 0, 0
	default:
		log.Fatal("renderTile: tile not supported")
	}

	sourcePt := image.Pt(8+10*ix, 27+7*iy)
	sourceRect := coordTranslator.tileMiniSize.Bounds().Add(sourcePt)
	return tilesImg.SubImage(sourceRect).(*ebiten.Image)
}

func update(screen *ebiten.Image) error {

	for y := 0; y < world.Height; y++ {
		for x := 0; x < world.Width; x++ {
			t := getTileSprite(screen, x, y, world.At(x, y))
			screen.DrawImage(t, coordTranslator.tileDrawOptions(x, y))
		}
	}

	printDebug(screen)

	if showMiniMap {
		c := gfx.ColorWithAlpha(colornames.Black, 210)
		ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, c)
		for y := 0; y < world.Height; y++ {
			for x := 0; x < world.Width; x++ {
				t := getMiniTileSprite(screen, x, y, world.At(x, y))
				screen.DrawImage(t, coordTranslator.tileMiniDrawOptions(x, y))
			}
		}
	}

	showMiniMap = false
	mapButton.hovered = false

	positions := []gfx.Vec{cursorPosition()}
	for _, id := range ebiten.TouchIDs() {
		x, y := ebiten.TouchPosition(id)
		positions = append(positions, gfx.IV(x, y))
		fmt.Println(positions)
	}

	if mapButton.ContainsAny(positions...) {
		fmt.Println("contains!")
		mapButton.onHover(&mapButton)
		mapButton.hovered = true

	}
	mapButton.render(screen)

	return nil
}

func cursorPosition() gfx.Vec {
	x, y := ebiten.CursorPosition()
	return gfx.V(float64(x), float64(y))
}

var (
	world           World
	tilesImg        *ebiten.Image
	coordTranslator CoordTranslator
	mapButton       Button
	showMiniMap     bool
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
		ebitenutil.DebugPrintAt(screen, "MiniMap", int(b.Min.X)+4, int(b.Min.Y)+6)
	}
}

func main() {
	world = NewWorld(20, 20)
	world.Set(0, 0, tile.Grass)
	world.Set(1, 0, tile.Grass)

	coordTranslator = CoordTranslator{
		anchor:       gfx.V(30, 30),
		tileSize:     gfx.R(0, 0, 15, 13),
		tileMiniSize: gfx.R(0, 0, 7, 5),
	}

	mapButton = Button{
		Rect:    gfx.R(0, 0, 50, 30).Moved(gfx.V(screenWidth-50, 0)),
		onClick: func(b *Button) { showMiniMap = true },
		onHover: func(b *Button) { showMiniMap = true },
	}

	// Load resources
	tileData, err := hextiles.Asset("hextilesets.png")
	if err != nil {
		log.Fatal(err)
	}
	tilesImg = ebitendrawutil.ImageFromBytes(tileData)

	if err := ebiten.Run(update, screenWidth, screenHeight, 1.5, "civ"); err != nil {
		log.Fatal(err)
	}
}
