package models

import "errors"

var ErrNoRecord = errors.New("models: no matching record found")
var ErrPasswordLength = errors.New("password does not meet length requirement")
var ErrBlankField = errors.New("field is blank")

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
	TextGreaterThan100
)

var FormErrorsState = map[FormError]string{
	PostInsertionError: "Post was not inserted",
	PostInserted:       "Post was inserted",
	PostNotFound:       "Post not found",
	IDBelowZero:        "ID can not be below 0",
	IDIsNotNumber:      "Id you entered is not a number",
	EmptyFields:        "Input fields can't be empty",
	TextGreaterThan100: "Text can not be more than 100 characters",
}
