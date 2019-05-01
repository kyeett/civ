package main

import (
	"log"

	"github.com/kyeett/civ/tile"
)

type World struct {
	Width, Height int
	tiles         []tile.Tile
}

func NewWorld(width, height int) World {
	tiles := []tile.Tile{}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			tiles = append(tiles, tile.Water)
		}
	}
	return World{
		Width:  width,
		Height: height,
		tiles:  tiles,
	}
}

func (w *World) At(x, y int) tile.Tile {
	if !w.pointInBounds(x, y) {
		return tile.Null
	}

	return w.tiles[x+y*w.Width]
}

func (w *World) pointInBounds(x, y int) bool {
	if x < 0 || x >= w.Width || y < 0 || y >= w.Height {
		return false
	}
	return true
}

func (w *World) Set(x, y int, t tile.Tile) {
	if !w.pointInBounds(x, y) {
		log.Fatalf("point %d,%d out of bounds", x, y)
	}

	w.tiles[x+y*w.Width] = t
}
