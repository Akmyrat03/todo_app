package repository

import (
	"fmt"
	todo "todo_app"

	"github.com/jmoiron/sqlx"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(listID int, item todo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, nil
	}

	itemQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoItemsTable)

	var itemID int
	row := tx.QueryRow(itemQuery, item.Title, item.Description)
	if err = row.Scan(&itemID); err != nil {
		tx.Rollback()
		return 0, err
	}

	listItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", listsItemsTable)
	_, err = tx.Exec(listItemsQuery, listID, itemID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemID, tx.Commit()

}

func (r *TodoItemPostgres) GetAll(userID, listID int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem
	query := `SELECT ti.id, ti.title, ti.description, ti.done 
			FROM todo_items AS ti 
			INNER JOIN lists_items AS li ON ti.id= li.item_id
			INNER JOIN users_list AS ul ON ul.list_id = li.list_id 
			WHERE li.list_id = $1 AND ul.user_id = $2`
	if err := r.db.Select(&items, query, listID, userID); err != nil {
		return nil, err
	}

	return items, nil

}

func (r *TodoItemPostgres) GetByID(userID, itemID int) (todo.TodoItem, error) {
	var item todo.TodoItem
	query := `SELECT ti.id, ti.title, ti.description 
			FROM todo_items AS ti
			INNER JOIN lists_items AS li ON ti.id = li.item_id
			INNER JOIN users_list AS ul ON ul.list_id= li.list_id
			WHERE li.item_id = $1 AND ul.user_id = $2`
	err := r.db.Get(&item, query, itemID, userID)
	if err != nil {
		return todo.TodoItem{}, err
	}
	return item, nil
}

func (r *TodoItemPostgres) Delete(userID, itemID int) error {
	query := fmt.Sprintf(`DELETE FROM %s ti USING %s li, %s ul
		WHERE ti.id = li.item_id AND ul.list_id = li.list_id AND ul.user_id = $1 AND ti.id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)
	_, err := r.db.Exec(query, userID, itemID)
	if err != nil {
		return err
	}

	return nil
}
