package models

import "errors"

var ErrNoRecord = errors.New("models: no matching record found")

type FormError int

const (
	PostNotFound FormError = iota
	IDBelowZero
	IDIsNotNumber
	PostInsertionError
	PostInserted
	EmptyFields
	EmptyTitle
	EmptyContent
)

var FormErrorsState = map[FormError]string{
	PostInsertionError: "Post was not inserted",
	PostInserted:       "Post was inserted",
	PostNotFound:       "Post not found",
	IDBelowZero:        "ID can not be below 0",
	IDIsNotNumber:      "Id you entered is not a number",
	EmptyFields:        "Input fields can't be empty",
}
