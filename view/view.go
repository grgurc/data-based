package view

// View is a menu item / screen the app can be in
// example views:
// - query result table
// - edit table (edit row/column inplace)
// - query insert text box
type View interface {
	Draw()      // for now this should be enough, although i think there should be a controller
	Next() View // switches to new view (e.g. user runs a select, switch to table view)
}
