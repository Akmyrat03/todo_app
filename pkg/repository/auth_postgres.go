package repository

import (
	"fmt"
	todo "todo_app"

	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user todo.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) VALUES ($1, $2, $3) RETURNING id", usersTable)
	err := r.db.QueryRow(query, user.Name, user.Username, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil

}

func (r *AuthPostgres) GetUser(username, password string) (todo.User, error) {
	var user todo.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username = $1 AND password_hash = $2", usersTable)
	err := r.db.Get(&user, query, username, password)
	if err != nil {
		return todo.User{}, err
	}

	return user, nil
}
