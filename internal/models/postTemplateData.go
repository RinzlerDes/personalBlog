package models

import "time"

// type PostTemplateData struct {
// 	Posts              []*Post
// 	Post               Post
// 	CurrentYear        int
// 	PostNotFound       bool
// 	IDBelowZero        bool
// 	IDIsNotNumber      bool
// 	PostInsertionError bool
// 	PostInserted       bool
// 	EmptyFields        bool
// }

var insertionMessage = "Post inserted"

type PostTemplateData struct {
	Posts                 []*Post
	Post                  Post
	CurrentYear           int
	InsertionErrorMessage string
	InsertionMessage      *string
	FormErrors       map[string]string
	PostInserted          bool
}

func NewPostTemplateData() PostTemplateData {
	return PostTemplateData{
		CurrentYear:      time.Now().Year(),
		InsertionMessage: &insertionMessage,
		FormErrors:  make(map[string]string),
	}
}
