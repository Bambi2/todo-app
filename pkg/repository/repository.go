package repository

import (
	"github.com/bambi2/todo-app/pkg/domain"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GetUser(username string, password string) (domain.User, error)
}

type TodoList interface {
	CreateList(userId int, todoList domain.TodoList) (int, error)
	GetAllLists(userId int) ([]domain.TodoList, error)
	GetListById(userId int, listId int) (domain.TodoList, error)
	DeleteList(userId int, listId int) error
	UpdateList(userId int, listId int, input domain.UpdateListInput) error
}

type TodoItem interface {
	CreateItem(listId int, todoItem domain.TodoItem) (int, error)
	GetAllItems(listId int) ([]domain.TodoItem, error)
	GetItemById(userId int, itemId int) (domain.TodoItem, error)
	DeleteItem(userId int, itemId int) error
	UpdateItem(userId int, itemId int, input domain.UpdateItemInput) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
