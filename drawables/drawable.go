package drawables

import "github.com/gdamore/tcell/v2"

type dimensions struct {
	x, y int
	w, h int
}

type Drawable interface {
	Draw() [][]rune
	Update(e tcell.Event)
}
