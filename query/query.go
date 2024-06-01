package query

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Query interface {
	Query() string // returns query string
	Run()          // actually runs the query
	Error() string
}

type SelectQuery struct {
	db       *sqlx.DB
	query    string
	ColNames []string
	ColTypes []string
	Rows     [][]string
	err      error
}

func (q *SelectQuery) Query() string {
	return q.query
}

func (q *SelectQuery) Run() {
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
		q.ColNames = append(q.ColNames, t.Name())
		q.ColTypes = append(q.ColTypes, strings.ToUpper(t.DatabaseTypeName()))
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
		q.Rows = append(q.Rows, stringValues)
	}
}

func (q *SelectQuery) Error() string {
	if q.err != nil {
		return q.err.Error()
	}

	return ""
}

type ExecQuery struct {
	db           *sqlx.DB
	query        string
	rowsAffected int64
	err          error
}

func (q *ExecQuery) Query() string {
	return q.query
}

func (q *ExecQuery) Run() {
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

func (q *ExecQuery) Error() string {
	if q.err != nil {
		return q.err.Error()
	}

	return ""
}

type FailedQuery struct {
	err error
}

func (q *FailedQuery) Query() string {
	return ""
}

func (q *FailedQuery) Run() {}

func (q *FailedQuery) Error() string {
	return q.err.Error()
}

// figure out if select or other type of query and create it
// Question for self: Should this return an error, and if so, which kind and where should it be handled?
func NewQuery(db *sqlx.DB, query string) Query {
	if db == nil {
		panic("No database connection supplied")
	}

	err := db.Ping()
	if err != nil {
		return &FailedQuery{
			err: errors.New("Couldn't connect to database"),
		}
	}

	qType, _, found := strings.Cut(query, " ")
	if !found {
		return &FailedQuery{
			err: errors.New("Malformed query"),
		}
	}

	switch strings.ToUpper(qType) {
	case "SELECT":
		return &SelectQuery{
			db:    db,
			query: query,
		}
	default:
		return &ExecQuery{
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
