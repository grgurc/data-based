package query

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/jmoiron/sqlx"
)

type Query interface {
	Write(w io.Writer) error
}

type SelectQuery struct {
	query    string
	colNames []string
	colTypes []string
	rows     [][]string
}

func (q *SelectQuery) Write(w io.Writer) error {
	joined := [][]string{
		q.colNames,
		q.colTypes,
	}
	joined = append(joined, q.rows...)

	b := new(bytes.Buffer)
	tw := tabwriter.NewWriter(b, 0, 0, 2, ' ', tabwriter.AlignRight|tabwriter.Debug)
	for _, row := range joined {
		tw.Write([]byte(strings.Join(row, "\t") + "\t\n"))
	}
	tw.Flush()

	// now b contains the whole table
	// we need to read the first 2 lines
	colNames, err := b.ReadString('\n')
	if err != nil {
		return err
	}
	colTypes, err := b.ReadString('\n')
	if err != nil {
		return err
	}

	// doesn't seem very efficient but oh well
	sep := strings.Repeat("-", len(colNames)-1) + "\n"

	// write back out w including separation line
	w.Write([]byte(colNames))
	w.Write([]byte(colTypes))
	w.Write([]byte(sep))
	w.Write(b.Bytes())

	return nil
}

func NewSelectQuery(db *sqlx.DB, query string) (*SelectQuery, error) {
	q := &SelectQuery{
		query: query,
	}

	rows, err := db.Queryx(query)
	if err != nil {
		return nil, err
	}

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
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
			return nil, err
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

	return q, nil
}

type ExecQuery struct {
	query        string
	rowsAffected int64
	lastInsertId int64
}

func (q *ExecQuery) Write(w io.Writer) error {
	fmt.Fprintln(w, "Query:")
	fmt.Fprintln(w, q.query)
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Rows affected:")
	fmt.Fprintln(w, q.rowsAffected)
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Last inserted id:")
	fmt.Fprintln(w, q.lastInsertId)
	fmt.Fprintln(w)

	return nil
}

func NewExecQuery(db *sqlx.DB, query string) (*ExecQuery, error) {
	res, err := db.Exec(query)
	if err != nil {
		return nil, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	last, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &ExecQuery{
		query:        query,
		rowsAffected: rows,
		lastInsertId: last,
	}, nil
}

func New(db *sqlx.DB, queryString string) (Query, error) {
	if db == nil {
		return nil, errors.New("no database connection supplied")
	}

	// this is not necessary, should be done by the app during init
	err := db.Ping()
	if err != nil {
		return nil, err
	}

	qType, _, found := strings.Cut(queryString, " ")
	if !found {
		return nil, err
	}

	switch strings.ToUpper(qType) {
	case "SELECT":
		return NewSelectQuery(db, queryString)
	default:
		return NewExecQuery(db, queryString)
	}
}

/*
// this should be reworked into something like
type Schema struct {
	Tables []Table
}

type Table struct {
	ColNames []string
	ColTypes []string
}
// and it should be used for the left widget - tables info / list
*/
/*
func getTableNames(db *sqlx.DB) []string {
	tables := []string{}

	if err := db.Select(&tables, "SHOW TABLES"); err != nil {
		panic(err)
	}

	return tables
}
*/
