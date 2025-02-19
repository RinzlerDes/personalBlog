package models

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type FormErrors struct {
	Errors map[string][]error
}

type UserFormErrors struct {
	Field string
	Err   error
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

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

	fe.stringGT100(text, key)
}

func (fe *FormErrors) RunChecksForContent(text string, key string) {
	fe.NotBlank(text, key)
}

func (fe *FormErrors) NotValid() bool {
	return len(fe.Errors) > 0
}

func (fe *FormErrors) AddError(key string, err error) {
	fe.Errors[key] = append(fe.Errors[key], err)
}

func (fe *FormErrors) NotBlank(text string, key string) bool {
	if strings.TrimSpace(text) == "" {
		fe.AddError(key, ErrBlankField)
		return false
	}
	return true
}

func (fe *FormErrors) CheckInputIsUInt(text string, key string) {
	_, err := strconv.ParseUint(text, 10, 64)
	if err != nil {
		fe.AddError(key, fmt.Errorf("input is not a valid number"))
	}
}

func (fe *FormErrors) stringGT100(str string, key string) bool {
	if len(str) > 100 {
		fe.AddError(key, fmt.Errorf("string can not be greater than 100 characters"))
		return true
	}

	return false
}

func (fe *FormErrors) minChars(str string, key string, length int) {
	if len(str) < length {
		fe.AddError(key, ErrPasswordLength)
	}
}

func (fe *FormErrors) ValidatePassword(pass string) {
	key := "password"
	fe.NotBlank(pass, key)
	fe.minChars(pass, key, 8)
}

func (fe *FormErrors) ValidateEmail(email string) {
	key := "email"
	fe.NotBlank(email, key)
}
