package service

import (
	"github.com/bambi2/todo-app/pkg/domain"
	"github.com/bambi2/todo-app/pkg/repository"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GenerateToken(username string, password string) (string, error)
	ParseToken(accessToken string) (int, error)
}

type TodoList interface {
	CreateList(userId int, todoList domain.TodoList) (int, error)
	GetAllLists(userId int) ([]domain.TodoList, error)
	GetListById(userId int, listId int) (domain.TodoList, error)
	DeleteList(userId int, listId int) error
	UpdateList(userId int, listId int, input domain.UpdateListInput) error
}

type TodoItem interface {
	CreateItem(userId int, listId int, todoItem domain.TodoItem) (int, error)
	GetAllItems(userId int, listId int) ([]domain.TodoItem, error)
	GetItemById(userId int, itemId int) (domain.TodoItem, error)
	DeleteItem(userId int, itemId int) error
	UpdateItem(userId int, itemId int, input domain.UpdateItemInput) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
