package main

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/awesome-gocui/gocui"
	"github.com/grgurc/data-based/config"
	"github.com/grgurc/data-based/query"
	"github.com/jmoiron/sqlx"
)

const (
	tableView   = "table"
	commandView = "command"
	listView    = "list"
)

type TableManager struct {
	view  *gocui.View
	query query.Query
}

// this one displays the result of the query
func (m *TableManager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if m.view != nil && m.query != nil {
		m.view.Clear()
		// m.view.SetOrigin() -> use this for scrolling the table
		log.Println(m.view.ViewBuffer())
		err := m.query.Write(m.view)
		if err != nil {
			return err
		}
	}

	if v, err := g.SetView(tableView, int(0.25*float32(maxX))+1, 0, maxX-1, maxY-6, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Clear()

		v.Title = "Query Result"
		// v.Subtitle = "test subtitel"
		v.Wrap = false
		v.Autoscroll = false

		m.view = v
	}

	// set keybindings for this one here
	// up down left right - scroll through table
	return nil
}

type CommandManager struct {
	view *gocui.View
}

// this one will display table names and search them or something
func (m *CommandManager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(commandView, 0, maxY-5, maxX-1, maxY-1, 0); err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = "Insert Query"
		v.Editable = true
		v.Overwrite = false
		v.Wrap = false
		v.Autoscroll = false

		m.view = v
	}

	// set keybindings here
	return nil
}

type ListManager struct {
	view *gocui.View
}

// this one displays a list of tables and some other stuff maybe dunno
func (m *ListManager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(listView, 0, 0, int(0.25*float32(maxX)), maxY-6, 0); err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = "List"
		v.Wrap = false
		v.Autoscroll = false

		m.view = v
	}

	// set keybindings here
	return nil
}

// TODO: put in separate file/package
type App struct {
	g  *gocui.Gui
	db *sqlx.DB

	// managers for each view/window thing
	tableManager   *TableManager
	listManager    *ListManager
	commandManager *CommandManager
}

func NewApp(g *gocui.Gui, db *sqlx.DB) *App {
	a := &App{
		g:  g,
		db: db,
	}

	g.SupportOverlaps = false
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen

	a.tableManager = &TableManager{}
	a.commandManager = &CommandManager{}
	a.listManager = &ListManager{}

	g.SetManager(a.tableManager, a.listManager, a.commandManager)
	_, err := g.SetCurrentView(tableView)
	if err != nil {
		log.Println(err) // so this stuff panics every time i guess
	}

	return a
}

func (a *App) BaseKeybindings() error {
	// quit
	if err := a.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}

	// Ctrl+Q activates Query Insert View
	if err := a.g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err := a.g.SetCurrentView(commandView)
		return err
	}); err != nil {
		return err
	}

	// Ctrl+T activates Table View
	if err := a.g.SetKeybinding("", gocui.KeyCtrlT, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err := a.g.SetCurrentView(tableView)
		return err
	}); err != nil {
		return err
	}

	// Ctrl+L
	if err := a.g.SetKeybinding("", gocui.KeyCtrlL, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err := a.g.SetCurrentView(listView)
		return err
	}); err != nil {
		return err
	}

	return nil
}

func main() {
	// logging
	f, err := os.OpenFile("./logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	wrt := io.Writer(f)
	log.SetOutput(wrt)

	// db conn + query
	db := config.NewDbFromYaml("config/config_default.yml") // replace with your config
	// running queries should be done from a goroutine
	// in order not to halt the rest of the app
	// TODO: check out goroutines in
	q, err := query.New(db, "DELETE FROM workspaces WHERE id = 100;")
	if err != nil {
		log.Panicln(err)
	}

	// gocui setup
	g, err := gocui.NewGui(gocui.OutputTrue, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	app := NewApp(g, db)
	if err := app.BaseKeybindings(); err != nil {
		log.Panicln(err)
	}
	log.Println(app)
	log.Println(app.commandManager)
	log.Println(app.listManager)
	log.Println(app.tableManager)

	// just for now
	app.tableManager.query = q

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
}
