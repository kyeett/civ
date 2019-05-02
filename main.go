package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	hextiles "github.com/kyeett/animex/resources/hex"
	"github.com/kyeett/civ/tile"
	"github.com/oakmound/shiny/materialdesign/colornames"
	"github.com/peterhellberg/gfx"

	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/ebitendrawutil"
)

const (
	screenWidth, screenHeight = 240, 180
	// screenWidth, screenHeight = 480, 360
)

func getTileSprite(screen *ebiten.Image, typ tile.Tile) *ebiten.Image {

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

func getShadowTileSprite(screen *ebiten.Image) *ebiten.Image {
	sourcePt := image.Pt(4+18*0-2, 183+19*3-2)
	sourceRect := image.Rect(0, 0, 15+4, 13+4).Add(sourcePt)
	return tilesImg.SubImage(sourceRect).(*ebiten.Image)
}

func getMiniTileSprite(screen *ebiten.Image, typ tile.Tile) *ebiten.Image {

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
	// screen.Fill(color.White)

	mouseController.Update()

	for y := 0; y < world.Height; y++ {
		for x := 0; x < world.Width; x++ {
			t := getTileSprite(screen, world.At(x, y))
			screen.DrawImage(t, coordTranslator.tileDrawOptions(x, y))
		}
	}

	for y := 0; y < world.Height; y++ {
		for x := 0; x < world.Width; x++ {

			if world.At(x, y) == tile.Grass {
				op := coordTranslator.tileDrawOptions(x, y)
				op.GeoM.Translate(-2, -2)
				screen.DrawImage(getShadowTileSprite(screen), op)
			}
		}
	}

	for y := 0; y < world.Height; y++ {
		for x := 0; x < world.Width; x++ {
			t := getTileSprite(screen, world.At(x, y))

			if world.At(x, y) == tile.Grass {
				screen.DrawImage(t, coordTranslator.tileDrawOptions(x, y))
			}
		}
	}

	if selected.onGrid {
		t := getTileSprite(screen, world.At(selected.x, selected.y))
		op := coordTranslator.tileDrawOptions(selected.x, selected.y)
		op.ColorM.Scale(1.5, 1.5, 1.5, 1)
		screen.DrawImage(t, op)
	}

	printDebug(screen)
	ui.Render(screen)

	if true {

		// if showMiniMap {
		// 	c := gfx.ColorWithAlpha(colornames.Black, 210)
		// 	ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, c)
		// 	for y := 0; y < world.Height; y++ {
		// 		for x := 0; x < world.Width; x++ {
		// 			t := getMiniTileSprite(screen, world.At(x, y))
		// 			screen.DrawImage(t, coordTranslator.tileMiniDrawOptions(x, y))
		// 		}
		// 	}
		// }

		// showMiniMap = false
		// mapButton.hovered = false

		// positions := []gfx.Vec{cursorPosition()}
		// for _, id := range ebiten.TouchIDs() {
		// 	x, y := ebiten.TouchPosition(id)
		// 	positions = append(positions, gfx.IV(x, y))
		// 	fmt.Println(positions)
		// }

		// if mapButton.ContainsAny(positions...) {
		// 	fmt.Println("contains!")
		// 	mapButton.onHover(&mapButton)
		// 	mapButton.hovered = true
		// }
		// mapButton.render(screen)

	}
	return nil
}

var (
	world           World
	tilesImg        *ebiten.Image
	coordTranslator CoordTranslator
	mapButton       Button
	showMiniMap     bool
	hoverables      []Hoverable
	mouseController MouseController
	selected        struct {
		onGrid bool
		x, y   int
	}
	ui UI
)

type Hoverable interface {
	Hovered(gfx.Vec) bool
	SetHovered(bool)
}

var (
	base = gfx.V(30, 30)
)

func OnDrag(diff gfx.Vec) {
	coordTranslator.anchor = diff
}

type UI struct {
	menu gfx.Rect
}

func (u *UI) Render(screen *ebiten.Image) {

	if selected.onGrid {
		ebitenutil.DrawRect(screen, u.menu.Min.X, u.menu.Min.Y, u.menu.W(), u.menu.H(), gfx.ColorWithAlpha(color.Black, 200))
		ebitendrawutil.DrawRect(screen, u.menu, colornames.White, 2)

		tiles := []tile.Tile{tile.Water, tile.Grass, tile.Desert}
		tw := coordTranslator.tileSize.W()
		th := coordTranslator.tileSize.H()
		padd := (u.menu.W() - float64(len(tiles))*tw) / float64(len(tiles)+1)
		for i, tile := range tiles {
			t := getTileSprite(screen, tile)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(i)*(tw+padd)+padd, u.menu.Min.Y+(u.menu.H()-th)/2)
			screen.DrawImage(t, op)
		}
	}
}

func (u *UI) Click(v gfx.Vec) bool {
	if u.menu.Contains(v) {

		switch int(3 * v.X / u.menu.W()) {
		case 0:
			world.Set(selected.x, selected.y, tile.Water)
		case 1:
			world.Set(selected.x, selected.y, tile.Grass)
		case 2:
			world.Set(selected.x, selected.y, tile.Desert)
		}

		return true
	}
	return false
}

func OnClick(v gfx.Vec) {

	if clicked := ui.Click(v); clicked {
		return
	}

	grid := struct {
		Click func() bool
	}{
		Click: func() bool {
			cx, cy := coordTranslator.fromScreen(cursorRelativePosition())
			if cx >= 0 && cy >= 0 && cx < 10 && cy < 10 {
				selected.onGrid = true
				selected.x = cx
				selected.y = cy
				return true
			}
			selected.onGrid = false
			selected.x = -1
			selected.y = -1
			return false
		},
	}

	if clicked := grid.Click(); clicked {
		return
	}

}

func main() {
	ui = UI{
		menu: gfx.R(0, 0, screenWidth-1, 50).Moved(gfx.V(0, screenHeight-50)),
	}
	world = NewWorld(10, 10)

	// Island 1
	world.Set(1, 1, tile.Grass)
	world.Set(2, 1, tile.Grass)
	world.Set(3, 1, tile.Grass)
	world.Set(1, 2, tile.Grass)
	world.Set(2, 2, tile.Grass)
	world.Set(3, 2, tile.Grass)

	// Island 2
	world.Set(1+4, 1, tile.Desert)
	world.Set(2+4, 1, tile.Desert)
	world.Set(3+4, 1, tile.Desert)
	world.Set(1+4, 2, tile.Desert)
	world.Set(2+4, 2, tile.Desert)
	world.Set(3+4, 2, tile.Desert)

	// world.Set(5, 3, tile.Grass)
	// world.Set(4, 5, tile.Desert)

	// hoverables = append(hoverables, world)

	mouseController = NewMouseController(OnClick, OnDrag)

	coordTranslator = CoordTranslator{
		anchor:       base,
		tileSize:     gfx.R(0, 0, 15, 13),
		tileMiniSize: gfx.R(0, 0, 7, 5),
	}

	mapButton = Button{
		Rect:    gfx.R(0, 0, 70, 30).Moved(gfx.V(screenWidth-70, 0)),
		onClick: func(b *Button) { showMiniMap = true },
		onHover: func(b *Button) { showMiniMap = true },
	}

	// Load resources
	tileData, err := hextiles.Asset("hextilesets.png")
	if err != nil {
		log.Fatal(err)
	}
	tilesImg = ebitendrawutil.ImageFromBytes(tileData)

	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "civ"); err != nil {
		log.Fatal(err)
	}
}
