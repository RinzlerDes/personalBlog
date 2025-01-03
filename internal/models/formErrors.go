package models

import (
	"fmt"
	"strconv"
	"strings"
)

type FormErrors struct {
	Errors map[string]string
}

func (fe *FormErrors) RunChecksForId(text string, key string) {
	ok := fe.NotBlank(text, key)
	if !ok {
		return
	}
	fe.CheckInputIsUInt(text, key)
}

func (fe *FormErrors) RunChecksForTitle(text string, key string) {
	ok := fe.NotBlank(text, key)
	if !ok {
		return
	}

	tooBig := fe.stringGT100(text, key)
	if tooBig {
		return
	}
}

func (fe *FormErrors) RunChecksForContent(text string, key string) {
	ok := fe.NotBlank(text, key)
	if !ok {
		return
	}
}

func (fe *FormErrors) NotValid() bool {
	return len(fe.Errors) > 0
}

func (fe *FormErrors) AddError(key string, errorMessage string) {
	fe.Errors[key] = errorMessage
}

func (fe *FormErrors) NotBlank(text string, key string) bool {
	if strings.TrimSpace(text) == "" {
		fe.AddError(key, FormErrorsState[EmptyFields])
		return false
	}
	return true
}

func (fe *FormErrors) CheckInputIsUInt(text string, key string) {
	_, err := strconv.ParseUint(text, 10, 64)
	if err != nil {
		str := fmt.Sprintf("%s or %s\n", FormErrorsState[IDBelowZero], FormErrorsState[IDIsNotNumber])
		fe.AddError(key, str)
	}
}

func (fe *FormErrors) stringGT100(str string, key string) bool {
	if len(str) > 100 {
		fe.AddError(key, FormErrorsState[TextGreaterThan100])
		return true
	}

	return false
}
