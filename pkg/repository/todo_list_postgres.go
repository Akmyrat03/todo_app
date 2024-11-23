package repository

import (
	"fmt"
	todo "todo_app"

	"github.com/jmoiron/sqlx"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(userID int, list todo.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var listId int
	todoListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", toDoListsTable)
	row := tx.QueryRow(todoListQuery, list.Title, list.Description)
	if err := row.Scan(&listId); err != nil {
		tx.Rollback()
		return 0, err
	}

	usersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(usersListQuery, userID, listId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return listId, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userID int) ([]todo.TodoList, error) {
	var lists []todo.TodoList
	query := `SELECT tl.id, tl.title, tl.description FROM todo_lists AS tl INNER JOIN users_list AS ul ON tl.id = ul.list_id WHERE ul.user_id = $1`
	err := r.db.Select(&lists, query, userID)
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func (r *TodoListPostgres) GetByID(userID, listID int) (todo.TodoList, error) {
	var list todo.TodoList
	query := `SELECT tl.id, tl.title, tl.description FROM todo_lists AS tl INNER JOIN users_list AS ul ON tl.id= ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`
	err := r.db.Get(&list, query, userID, listID)
	if err != nil {
		return todo.TodoList{}, err
	}

	return list, err
}

func (r *TodoListPostgres) Delete(userID, listID int) error {
	query := `DELETE FROM todo_lists AS tl USING users_list AS ul WHERE tl.id = ul.list_id AND ul.user_id = $1 AND ul.list_id = $2`
	_, err := r.db.Exec(query, userID, listID)
	if err != nil {
		return err
	}

	return nil
}
