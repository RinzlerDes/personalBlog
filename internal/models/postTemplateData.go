package models

import "time"

var insertionMessage = "Post inserted"

type PostTemplateData struct {
	Posts                   []*Post
	Post                    Post
	CurrentYear             int
	InsertionErrorMessage   string
	InsertionMessage        *string
	FormErrors              FormErrors
	PostInserted            bool
	InitialInsertionMessage string
}

// Marker function for TemplateData interface
func (ptd *PostTemplateData) isTemplateData() {}

func NewPostTemplateData() PostTemplateData {
	return PostTemplateData{
		CurrentYear:      time.Now().Year(),
		InsertionMessage: &insertionMessage,
		FormErrors:       FormErrors{Errors: make(map[string]string)},
	}
}
