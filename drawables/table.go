package drawables

import (
	"github.com/gdamore/tcell/v2"
)

type column struct {
	maxW        int
	header      []string
	body        []string
	first, last bool
}

// table has header and main body
type scrollTable struct {
	x, y    int // position of top left corner relative to drawn table data
	w, h    int // total table dimensions
	columns []*column
}

func leftPadded(str string, size int) []rune {
	res := []rune{}
	for i := 0; i < size-len(str); i++ {
		res = append(res, ' ')
	}
	res = append(res, []rune(str)...)
	return res
}

func rightPadded(str string, size int) []rune {
	res := []rune{}
	res = append(res, []rune(str)...)
	for i := 0; i < size-len(str); i++ {
		res = append(res, ' ')
	}
	return res
}

func (d *scrollTable) drawHeader() [][]rune {
	var res [][]rune
	for _, r := range d.header {
		var row []rune
		for colIdx, colVal := range r {
			row = append(row, rune(' '))
			row = append(row, leftPadded(colVal, d.maxW[colIdx])...)
			row = append(row, rune(' '))
			row = append(row, tcell.RuneVLine)
		}
		res = append(res, row)
	}

	return res
}

func (d *scrollTable) drawBody() [][]rune {
	var res [][]rune
	for _, r := range d.body {
		var row []rune
		for colIdx, colVal := range r {
			row = append(row, rune(' '))
			row = append(row, leftPadded(colVal, d.maxW[colIdx])...)
			row = append(row, rune(' '))
			row = append(row, tcell.RuneVLine)
		}
		res = append(res, row)
	}

	return res
}

func (d *scrollTable) topLine() []rune {
	var row []rune
	row = append(row, tcell.RuneULCorner)
	for i := 1; i < d.w-1; i++ {
		row = append(row, tcell.RuneHLine)
	}
	row = append(row, tcell.RuneURCorner)

	return row
}

func (d *scrollTable) botLine() []rune {
	var row []rune
	row = append(row, tcell.RuneLTee)
	for i := 1; i < d.w-1; i++ {
		row = append(row, tcell.RuneHLine)
	}
	row = append(row, tcell.RuneRTee)

	return row
}

func (d *scrollTable) Draw() [][]rune {
	var res [][]rune

	// draw header

	// top line
	res = append(res, d.topLine())

	// content - header content is always the same
	res = append(res, d.drawHeader()...)

	// bottom line
	res = append(res, d.botLine())

	// top line
	res = append(res, d.topLine())

	// content - body needs to change depending on the position of the 'cursor'
	res = append(res, d.drawBody()...)

	// bottom line
	res = append(res, d.botLine())

	return res
}

func (d *scrollTable) Update(e tcell.Event) {

}

func (d *scrollTable) MoveLeft() {

}

func NewScrollTable(w, h int, header [][]string, body [][]string) *scrollTable {
	// just in case, not going to bother with this for now
	if len(header) == 0 || len(body) == 0 {
		return &scrollTable{
			w: w,
			h: h,
			x: 0,
			y: 0,
		}
	}

	// okay, so header and body are made as [row][column]
	// need to make columns out of it
	// bruh my brain is not braining at the moment

	columns := make([]*column, len(header[0]), len(header[0]))
	for i, row := range header {
		for j, col := range row {
			// might be nil, will see if program dies :)
			columns[i].header = append(columns[i].header)
		}
	}

	var columns []*column
	for i, col := range header[0] {

	}
	for i, row := range header {
		for j, column := range row {

		}
	}
	maxW := make([]int, len(header[0]))
	for _, r := range header {
		for i, c := range r {
			maxW[i] = max(maxW[i], len(c))
		}
	}
	for _, r := range body {
		for i, c := range r {
			maxW[i] = max(maxW[i], len(c))
		}
	}

	return &scrollTable{
		w:      w,
		h:      h,
		maxW:   maxW,
		header: header,
		body:   body,
	}
}
