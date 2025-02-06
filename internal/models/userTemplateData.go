package models

import "time"

type UserTemplateData struct {
	User        User
	Users       []*User
	FormErrors  FormErrors
	CurrentYear int
}

// Marker func for TemplateData interface
func (utd *UserTemplateData) isTemplateData() {}

func NewUserTemplateData() UserTemplateData {
	return UserTemplateData{
		FormErrors:  FormErrors{Errors: make(map[string]string)},
		CurrentYear: time.Now().Year(),
	}
}
