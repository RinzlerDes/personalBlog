package models

import "time"

type UserTemplateData struct {
	User        User
	Users       []*User
	FormErrors  FormErrors
	CurrentYear int
}

func (utd *UserTemplateData) isTemplateData() {}

func NewUserTemplateData() UserTemplateData {
	return UserTemplateData{
		FormErrors:  FormErrors{},
		CurrentYear: time.Now().Year(),
	}
}
