package models

import (
	"context"
	"errors"
	"fmt"
	"personalBlog/internal/loggers"
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

var logErr = loggers.LogErr
var logInfo = loggers.LogInfo

func (post *Post) String() string {
	return fmt.Sprintf(
		`Id         %d
        Title       %s
        Content     %s
        Created     %v`,
		post.ID, post.Title, post.Content, post.Created,
	)
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

func (postModel *PostModel) Get(id uint) (Post, error) {
	post := Post{
		ID: id,
	}

	SQLStatement := `SELECT * FROM posts WHERE id = $1`

	err := postModel.DB.QueryRow(
		context.Background(),
		SQLStatement,
		id,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.Created,
	)

	if err != nil {
		logErr.Println("error scanning", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return Post{}, ErrNoRecord
		}
		return Post{}, err
	}

	return post, nil
}
