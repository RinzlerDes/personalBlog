package models

import (
	"context"
	"errors"
	"fmt"
	"personalBlog/internal/loggers"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	ID      uint
	Title   string
	Content string
	Created time.Time
}

type PostModel struct {
	DBPool *pgxpool.Pool
}

var (
	logErr  = loggers.LogErr
	logInfo = loggers.LogInfo
)

func (post *Post) String() string {
	return fmt.Sprintf(
		`Id         %d
        Title       %s
        Content     %s
        Created     %v`,
		post.ID, post.Title, post.Content, post.Created,
	)
}

func (postModel *PostModel) Insert(post *Post) (uint, error) {
	SQLStatement := `INSERT INTO posts (title, content, created) VALUES($1, $2, current_timestamp) RETURNING id, created`

	var id uint
	var created time.Time
	err := postModel.DBPool.QueryRow(context.Background(), SQLStatement, post.Title, post.Content).Scan(&id, &created)
	if err != nil {
		return 0, err
	}
	// post.ID = id
	// post.Created = created

	return id, nil
}

func (postModel *PostModel) Get(id uint) (Post, error) {
	post := Post{
		ID: id,
	}

	SQLStatement := `SELECT * FROM posts WHERE id = $1`

	err := postModel.DBPool.QueryRow(
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

func (postModel *PostModel) Latest(limit uint) ([]*Post, error) {
	sqlStatement := `Select * FROM posts ORDER BY created DESC LIMIT $1;`

	rows, err := postModel.DBPool.Query(context.Background(), sqlStatement, limit)
	if err != nil {
		logErr.Printf("error get lastest %d rows: %s", limit, err.Error())
		return nil, err
	}
	defer rows.Close()

	posts := make([]*Post, 0, limit)

	for rows.Next() {
		post := &Post{}
		rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Created,
		)
		posts = append(posts, post)
	}

	return posts, nil
}
