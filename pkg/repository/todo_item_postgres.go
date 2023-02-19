package repository

import (
	"fmt"

	"github.com/100pecheneK/go-todo-rest-api.git/internal/models"
	"github.com/jmoiron/sqlx"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db}
}

func (r *TodoItemPostgres) Create(listId int, input models.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var itemId int
	createItemQuery := fmt.Sprintf(`INSERT INTO %s (title, description) values ($1, $2) RETURNING id`, todoItemsTable)
	row := tx.QueryRow(createItemQuery, input.Title, input.Description)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemQuery := fmt.Sprintf(`INSERT INTO %s (list_id, todo_id) values ($1, $2)`, listsItemsTable)
	_, err = tx.Exec(createListItemQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return itemId, tx.Commit()

}

func (r *TodoItemPostgres) GetAll(userId, listId int) ([]models.TodoItem, error) {
	var items []models.TodoItem

	query := fmt.Sprintf(`
	SELECT ti.id, ti.title, ti.description, ti.done
	FROM 
	%s ti
	INNER JOIN %s li ON li.todo_id=ti.id 
	INNER JOIN %s ul ON ul.list_id=li.list_id
	WHERE 
	li.list_id=$1 AND ul.user_id=$2`,
		todoItemsTable, listsItemsTable, usersListsTable)

	if err := r.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *TodoItemPostgres) GetById(userId, id int) (models.TodoItem, error) {
	var item models.TodoItem

	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done 
	FROM 
	%s ti
INNER JOIN %s li ON li.todo_id=ti.id
INNER JOIN %s ul ON ul.list_id=li.list_id
WHERE 
ti.id=$1 AND ul.user_id=$2
	`, todoItemsTable, listsItemsTable, usersListsTable)

	if err := r.db.Get(&item, query, id, userId); err != nil {
		return item, err
	}
	return item, nil
}
