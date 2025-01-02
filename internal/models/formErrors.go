package models

import (
	"fmt"
	"strconv"
	"strings"
)

type FormErrors struct {
	Errors map[string]string
}

func (fe *FormErrors) Valid() bool {
	return len(fe.Errors) == 0
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

func (fe *FormErrors) RunChecksForId(text string, key string) {
	ok := fe.NotBlank(text, key)
	if !ok {
		return
	}
	fe.CheckInputIsUInt(text, key)
}
