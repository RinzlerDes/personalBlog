package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID       uint
	Name     string
	Email    string
	Password string
	Created  time.Time
}

type UserModel struct {
	DBPool *pgxpool.Pool
}

func (user *User) String() string {
	return fmt.Sprintf(`
Id			%d
Name		%s
Email		%s
Password	%s
Created		%v`,
		user.ID, user.Name, user.Email, user.Password, user.Created,
	)
}

func (userModel *UserModel) Insert(user User) (uint, time.Time, error) {
	SQLStatement := `insert into users (name, email, password, created) 
	values($1, $2, $3, current_timestamp)
	returning id, created`

	row := userModel.DBPool.QueryRow(
		context.Background(),
		SQLStatement,
		user.Name,
		user.Email,
		user.Password)

	var id uint
	var created time.Time

	err := row.Scan(&id, &created)
	if err != nil {
		logErr.Println("error scanning")
		return 0, time.Time{}, err
	}

	return id, created, nil
}

func (userModel *UserModel) Get(id uint) (User, error) {
	SQLStatement := "select * from users where id = $1"

	row := userModel.DBPool.QueryRow(
		context.Background(),
		SQLStatement,
		id)

	user := User{}

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Created,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
