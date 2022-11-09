package repository

import (
	"fmt"
	"strings"

	"github.com/bambi2/todo-app/pkg/domain"
	"github.com/jmoiron/sqlx"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) CreateItem(listId int, todoItem domain.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	insertItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES($1, $2) RETURNING id", todoItemsTable)
	row := tx.QueryRow(insertItemQuery, todoItem.Title, todoItem.Description)
	if err := row.Scan(&itemId); err != nil {
		tx.Rollback()
		return 0, err
	}

	insertListsItemQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES($1, $2)", listsItemsTable)
	_, err = tx.Exec(insertListsItemQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAllItems(listId int) ([]domain.TodoItem, error) {
	var items []domain.TodoItem

	query := fmt.Sprintf("SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li ON ti.id = li.item_id WHERE li.list_id = $1",
		todoItemsTable, listsItemsTable)
	err := r.db.Select(&items, query, listId)

	return items, err
}

func (r *TodoItemPostgres) GetItemById(userId int, itemId int) (domain.TodoItem, error) {
	var item domain.TodoItem

	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li ON ti.id = li.item_id
							INNER JOIN %s ul ON li.list_id = ul.list_id WHERE ul.user_id = $1 AND ti.id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)
	err := r.db.Get(&item, query, userId, itemId)

	return item, err
}

func (r *TodoItemPostgres) DeleteItem(userId int, itemId int) error {
	query := fmt.Sprintf(`DELETE FROM %s ti USING %s li, %s ul
							WHERE li.item_id = ti.id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, itemId)

	return err
}

func (r *TodoItemPostgres) UpdateItem(userId int, itemId int, input domain.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	var argId = 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, input.Description)
		argId++
	}

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ",")

	query := fmt.Sprintf(`UPDATE %s ti SET %s FROM %s li, %s ul
							WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id=$%d AND ti.id = $%d`,
		todoItemsTable, setQuery, listsItemsTable, usersListsTable, argId, argId+1)

	args = append(args, userId, itemId)

	_, err := r.db.Exec(query, args...)

	return err

}
