package query

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/grgurc/data-based/drawables"
	"github.com/jmoiron/sqlx"
)

type Query interface {
	Query() string
	Run()
	Drawable(w, h int) drawables.Drawable
}

type selectQuery struct {
	db       *sqlx.DB
	query    string
	colNames []string
	colTypes []string
	rows     [][]string
	err      error
}

func (q *selectQuery) Query() string {
	return q.query
}

func (q *selectQuery) Run() {
	rows, err := q.db.Queryx(q.query)
	if err != nil {
		q.err = err
		return
	}

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		q.err = err
		return
	}

	for _, t := range colTypes {
		q.colNames = append(q.colNames, t.Name())
		q.colTypes = append(q.colTypes, strings.ToUpper(t.DatabaseTypeName()))
	}

	for rows.Next() {
		values := make([]interface{}, len(colTypes))
		for i := range values {
			values[i] = new(sql.NullString)
		}

		err = rows.Scan(values...)
		if err != nil {
			q.err = err
			return
		}

		stringValues := make([]string, len(values))
		for i, v := range values {
			if v.(*sql.NullString).Valid {
				stringValues[i] = v.(*sql.NullString).String
			} else {
				stringValues[i] = "NULL"
			}
		}
		q.rows = append(q.rows, stringValues)
	}

	fmt.Println(q.colNames, q.colTypes)
}

// TODO -> rework how dimensions are passed in if needed
func (q *selectQuery) Drawable(w, h int) drawables.Drawable {
	return drawables.NewTable(w, h, [][]string{q.colNames, q.colTypes}, q.rows)
}

type execQuery struct {
	db           *sqlx.DB
	query        string
	rowsAffected int64
	err          error
}

func (q *execQuery) Query() string {
	return q.query
}

func (q *execQuery) Run() {
	res, err := q.db.Exec(q.query)
	if err != nil {
		q.err = err
		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		q.err = err
		return
	}
	q.rowsAffected = rows
}

func (q *execQuery) Drawable(w, h int) drawables.Drawable {
	return nil
}

func (q *execQuery) Success() bool {
	return q.err == nil
}

type failedQuery struct {
	err error
}

func (q *failedQuery) Query() string {
	return ""
}

func (q *failedQuery) Run() {}

func (q *failedQuery) Drawable(w, h int) drawables.Drawable {
	return nil
}

// figure out if select or other type of query and create it
// Question for self: Should this return an error, and if so, which kind and where should it be handled?
func NewQuery(db *sqlx.DB, query string) Query {
	if db == nil {
		panic("No database connection supplied")
	}

	err := db.Ping()
	if err != nil {
		return &failedQuery{
			err: errors.New("Couldn't connect to database"),
		}
	}

	qType, _, found := strings.Cut(query, " ")
	if !found {
		return &failedQuery{
			err: errors.New("Malformed query"),
		}
	}

	switch strings.ToUpper(qType) {
	case "SELECT":
		return &selectQuery{
			db:    db,
			query: query,
		}
	default:
		return &execQuery{
			db:    db,
			query: query,
		}
	}
}

/*
func getTableNames(db *sqlx.DB) []string {
	tables := []string{}

	if err := db.Select(&tables, "SHOW TABLES"); err != nil {
		panic(err)
	}

	return tables
}
*/
