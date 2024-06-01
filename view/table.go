package view

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/tabwriter"

	"github.com/gdamore/tcell/v2"
	"github.com/grgurc/data-based/query"
)

// this displays the result of running a select
// - it needs to handle user input (for scrolling and exiting)
// - in the future might enable some kind of insert/edit mode... will see
type Table struct {
	screen tcell.Screen // the table draws on the tcell screen

	header [][]rune // these are the rows and columns
	body   [][]rune // they are calculated
	x, y   int      // cursor

}

// going to attempt using tabwriter package
func NewTable(screen tcell.Screen, query *query.SelectQuery) *Table {
	table := &Table{
		screen: screen,
		x:      0,
		y:      0,
	}

	if query.Error() != "" {
		table.header = [][]rune{[]rune(query.Error())}
		return table
	}

	joined := [][]string{
		query.ColNames,
		query.ColTypes,
	}
	joined = append(joined, query.Rows...)

	res := new(bytes.Buffer)
	tabWriter := tabwriter.NewWriter(res, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	for i, row := range joined {
		tabWriter.Write([]byte(strings.Join(row, "\t") + "\t"))
		if i < len(joined)-1 {
			tabWriter.Write([]byte("\n"))
		}
	}
	tabWriter.Flush()

	rows := strings.Split(res.String(), "\n")
	for i, row := range rows {
		if i < 2 { // this is hardcoded, mucho no good...
			table.header = append(table.header, []rune(row))
		} else {
			table.body = append(table.body, []rune(row))
		}
	}

	fmt.Println("rows", len(rows))
	fmt.Println("header", len(table.header))
	fmt.Println("body", len(table.body))
	return table
}

func (t *Table) Render() {
	// put all content to render in one big [][]rune
	// TODO: start header columns from x
	// TODO: start body rows from y
	var content [][]rune

	for _, row := range t.header {
		var contentRow []rune
		contentRow = append(contentRow, row...)[t.x:] // from cursor and right
		content = append(content, contentRow)
	}

	lineRow := []rune(strings.Repeat("-", len(content[0])))
	content = append(content, lineRow)

	var body [][]rune
	for _, row := range t.body {
		var contentRow []rune
		contentRow = append(contentRow, row...)
		contentRow = contentRow[t.x:] // from cursor and right
		body = append(body, contentRow)
	}

	content = append(content, body[t.y:]...) // this will fail when cursor goes too far down

	w, h := t.screen.Size()

	for y, row := range content {
		if y >= h {
			break
		}
		for x, val := range row {
			if x >= w {
				break
			}
			t.screen.SetContent(x, y, val, nil, tcell.StyleDefault)
		}
	}

	t.screen.Show()
	return

	// get screen size
	w, h = t.screen.Size()
	// set padding for w, h
	const paddingX, paddingY int = 10, 3
	// dimensions of screen inside padding
	w, h = w-2*paddingX, h-2*paddingY

	xScr, yScr := 0, 0
	for y := range min(h, len(t.header)) {
		for x := range min(w, len(t.header[y])) {
			xScr, yScr = x+paddingX, y+paddingY
			xMax, yMax := w-paddingX, h-paddingY
			if xScr < xMax && yScr < yMax {
				t.screen.SetContent(xScr, yScr, t.header[y][x], nil, tcell.StyleDefault)
			}
		}
	}

	yScr++
	for x := range min(w, len(t.header[0])) {
		xScr = x + paddingX
		t.screen.SetContent(xScr, yScr, '-', nil, tcell.StyleDefault)
	}
	yScr++

	t.screen.Show()
	return

	for y, row := range t.header {
		for x, val := range row {
			t.screen.SetContent(x, y, val, nil, tcell.StyleDefault)
		}
	}

	// TODO - padding
	x1, y1 := 5, 5 // padded dimensions
	x2, y2 := w-5, h-5

	// cW, cH := x2-x1, y2-y1           // content area width and height
	for y, row := range t.header {
		if y+y1 > y2 {
			break
		}
		for x, val := range row {
			if x+x1 > x2 {
				break
			}
			t.screen.SetContent(x+x1, y+y1, val, nil, tcell.StyleDefault)
		}
	}

	t.screen.Show()
}

func (t *Table) CursorUp() {
	if t.y > 0 {
		t.y--
	}
	log.Println("UP:", t.y)
}

func (t *Table) CursorDown() {
	_, h := t.screen.Size()
	if t.y < len(t.body)-(h-3) {
		t.y++
	}
	log.Println("DOWN:", t.y)
}

func (t *Table) CursorLeft() {
	if t.x-3 > 0 {
		t.x -= 3
	} else {
		t.x = 0
	}
	log.Println("LEFT:", t.x)
}

func (t *Table) CursorRight() {
	w, _ := t.screen.Size()
	if t.x+3 < len(t.header[0])-w {
		t.x += 3
	} else {
		t.x = len(t.header[0]) - w
	}
	log.Println("RIGHT:", t.x)
}

func (t *Table) Loop() {
	for {
		e := t.screen.PollEvent()

		switch e := e.(type) {
		case *tcell.EventKey:
			if e.Key() == tcell.KeyCtrlC {
				return
			}
			// this one should not quit the app
			// instead it should return to the previous screen
			// and then run its Loop() i guess or something like that
			// for now tho, just going to quit
			if e.Key() == tcell.KeyEsc {
				return // this should not quit
			}
			// these 4 attempt to move the cursor
			if e.Key() == tcell.KeyUp {
				t.CursorUp()
			}
			if e.Key() == tcell.KeyDown {
				t.CursorDown()
			}
			if e.Key() == tcell.KeyRight {
				t.CursorRight()
			}
			if e.Key() == tcell.KeyLeft {
				t.CursorLeft()
			}
		}

		t.Render()
	}
}

// TODO: implement this
// cuts off in regards to padding
// draws box around content -> TODO

/*
type screenBuffer struct {
	screen tcell.Screen
}

func (b screenBuffer) Render(content [][]rune) {
	w, h := b.screen.Size()
	x1, x2 := 10, w-10
	y1, y2 := 4, h-4

	// TODO: implement this shit
	for y := 0; y < len(content); y++ {
		if y > h-8 {
			break
		}
		for x := 0; x < len(content[y]); x++ {
			if x > w-20 {
				break
			}
		}
	}
}
*/
