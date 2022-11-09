package service

import (
	"github.com/bambi2/todo-app/pkg/domain"
	"github.com/bambi2/todo-app/pkg/repository"
)

type TodoItemService struct {
	todoItemRepo repository.TodoItem
	todoListRepo repository.TodoList
}

func NewTodoItemService(todoItemRepo repository.TodoItem, todoListRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{todoItemRepo: todoItemRepo, todoListRepo: todoListRepo}
}

func (s *TodoItemService) CreateItem(userId int, listId int, todoItem domain.TodoItem) (int, error) {
	_, err := s.todoListRepo.GetListById(userId, listId)
	//the list doesn't exist and or doesn't belong to the user
	if err != nil {
		return 0, err
	}
	return s.todoItemRepo.CreateItem(listId, todoItem)
}

func (s *TodoItemService) GetAllItems(userId int, listId int) ([]domain.TodoItem, error) {
	_, err := s.todoListRepo.GetListById(userId, listId)
	//the list doesn't exist and or doesn't belong to the user
	if err != nil {
		return nil, err
	}

	return s.todoItemRepo.GetAllItems(listId)
}

func (s *TodoItemService) GetItemById(userId int, itemId int) (domain.TodoItem, error) {
	return s.todoItemRepo.GetItemById(userId, itemId)
}

func (s *TodoItemService) DeleteItem(userId int, itemId int) error {
	return s.todoItemRepo.DeleteItem(userId, itemId)
}

func (s *TodoItemService) UpdateItem(userId int, itemId int, input domain.UpdateItemInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.todoItemRepo.UpdateItem(userId, itemId, input)
}
