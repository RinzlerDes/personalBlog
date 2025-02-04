package models

import "time"

type userTemplateData struct {
	User        User
	Users       []*User
	FormErrors  FormErrors
	CurrentYear int
}

func (utd *userTemplateData) isTemplateData() {}

func NewUserTemplateData() userTemplateData {
	return userTemplateData{
		FormErrors:  FormErrors{},
		CurrentYear: time.Now().Year(),
	}
}
