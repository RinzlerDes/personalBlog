package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Post struct {
	ID      uint
	Title   string
	Content string
	Created time.Time
}

type PostModel struct {
	DB *pgx.Conn
}

func (postModel *PostModel) Insert(post *Post) error {
	SQLStatement := `INSERT INTO posts (title, content, created) VALUES($1, $2, current_timestamp) RETURNING id, created`

	var id uint
	var created time.Time
	err := postModel.DB.QueryRow(context.Background(), SQLStatement, post.Title, post.Content).Scan(&id, &created)
	if err != nil {
		return err
	}
	post.ID = id
	post.Created = created

	return nil
}
