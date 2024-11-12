package models

type PostTemplateData struct {
	Posts              []*Post
	Post               Post
	CurrentYear        int
	PostNotFound       bool
	IDBelowZero        bool
	IDIsNotNumber      bool
	PostInsertionError bool
	PostInserted       bool
	EmptyFields        bool
}
