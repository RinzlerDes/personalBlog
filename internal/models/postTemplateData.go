package models

type PostTemplateData struct {
	Post               Post
	PostNotFound       bool
	IDBelowZero        bool
	IDIsNotNumber      bool
	PostInsertionError bool
	PostInserted       bool
}
