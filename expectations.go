package sqlmock

import (
	"database/sql/driver"
	"reflect"
	"regexp"
)

type expectation interface {
	fulfilled() bool
	setError(err error)
}

// common expectation

type commonExpectation struct {
	triggered bool
	err       error
}

func (e *commonExpectation) fulfilled() bool {
	return e.triggered
}

func (e *commonExpectation) setError(err error) {
	e.err = err
}

// query based expectation
type queryBasedExpectation struct {
	commonExpectation
	sqlRegex *regexp.Regexp
	args     []driver.Value
}

func (e *queryBasedExpectation) queryMatches(sql string) bool {
	return e.sqlRegex.MatchString(sql)
}

func (e *queryBasedExpectation) argsMatches(args []driver.Value) bool {
	if len(args) != len(e.args) {
		return false
	}
	for k, v := range e.args {
		vi := reflect.ValueOf(v)
		ai := reflect.ValueOf(args[k])
		switch vi.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if vi.Int() != ai.Int() {
				return false
			}
		case reflect.Float32, reflect.Float64:
			if vi.Float() != ai.Float() {
				return false
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if vi.Uint() != ai.Uint() {
				return false
			}
		case reflect.String:
			if vi.String() != ai.String() {
				return false
			}
		default:
			// compare types like time.Time based on type only
			if vi.Kind() != ai.Kind() {
				return false
			}
		}
	}
	return true
}

// begin transaction
type expectedBegin struct {
	commonExpectation
}

// tx commit
type expectedCommit struct {
	commonExpectation
}

// tx rollback
type expectedRollback struct {
	commonExpectation
}

// query expectation
type expectedQuery struct {
	queryBasedExpectation

	rows driver.Rows
}

// exec query expectation
type expectedExec struct {
	queryBasedExpectation

	result driver.Result
}
