package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/peterhellberg/gfx"
)

func cursorPosition() gfx.Vec {
	x, y := ebiten.CursorPosition()
	return gfx.V(float64(x), float64(y))
}

func cursorRelativePosition() gfx.Vec {
	return cursorPosition().Sub(coordTranslator.anchor)
}

func touchPosition(id int) gfx.Vec {
	x, y := ebiten.TouchPosition(id)
	return gfx.V(float64(x), float64(y))
}

func touchRelativePosition(id int) gfx.Vec {
	return touchPosition(id).Sub(coordTranslator.anchor)
}

type MouseController struct {
	strokes map[Source]struct{}
	onClick func(gfx.Vec)
	onDrag  func(gfx.Vec)
}

func NewMouseController(onClick, onDrag func(gfx.Vec)) MouseController {
	return MouseController{
		strokes: map[Source]struct{}{},
		onClick: onClick,
		onDrag:  onDrag,
	}
}

type Source interface {
	update()
	expired() bool
	dragged() bool
	diff() gfx.Vec
	position() gfx.Vec
}

type MouseSource struct {
	startPos, startOffset, currentPos gfx.Vec
}

func (ms *MouseSource) update() {
	ms.currentPos = cursorPosition()
}

func (ms *MouseSource) expired() bool {
	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
}

func (ms *MouseSource) dragged() bool {
	return ms.currentPos != ms.startPos
}

func (ms *MouseSource) diff() gfx.Vec {
	return ms.currentPos.Sub(ms.startPos).Add(ms.startOffset)
}

func (ms *MouseSource) position() gfx.Vec {
	return ms.currentPos
}

type TouchSource struct {
	id                                int
	startPos, startOffset, currentPos gfx.Vec
}

func (ts *TouchSource) update() {
	ts.currentPos = touchPosition(ts.id)
}

func (ts *TouchSource) expired() bool {
	return inpututil.IsTouchJustReleased(ts.id)
}

func (ts *TouchSource) dragged() bool {
	return ts.currentPos != ts.startPos
}

func (ts *TouchSource) diff() gfx.Vec {
	return ts.currentPos.Sub(ts.startPos).Add(ts.startOffset)
}

func (ts *TouchSource) position() gfx.Vec {
	return ts.currentPos
}

func (mc *MouseController) Update() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		fmt.Println("Mouse pressed at", cursorPosition())
		st := MouseSource{
			startOffset: coordTranslator.anchor,
			startPos:    cursorPosition(),
			currentPos:  cursorPosition(),
		}
		fmt.Println(st)
		mc.strokes[&st] = struct{}{}
	}
	for _, id := range inpututil.JustPressedTouchIDs() {
		fmt.Println("Touch pressed at ", touchPosition(id))
		st := TouchSource{
			id:          id,
			startOffset: coordTranslator.anchor,
			startPos:    touchPosition(id),
			currentPos:  touchPosition(id),
		}
		mc.strokes[&st] = struct{}{}
	}

	for s := range mc.strokes {
		if s.expired() {
			mc.handleRelease(s)
			delete(mc.strokes, s)
			continue
		}
		s.update()
	}

	for s := range mc.strokes {
		if s.dragged() {
			mc.onDrag(s.diff())
		}
	}
}

func (mc *MouseController) handleRelease(s Source) {
	if s.dragged() {
	} else {
		mc.onClick(s.position())
	}
}
