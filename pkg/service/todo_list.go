package service

import (
	"github.com/bambi2/todo-app/pkg/domain"
	"github.com/bambi2/todo-app/pkg/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) CreateList(userId int, todoList domain.TodoList) (int, error) {
	return s.repo.CreateList(userId, todoList)
}

func (s *TodoListService) GetAllLists(userId int) ([]domain.TodoList, error) {
	return s.repo.GetAllLists(userId)
}

func (s *TodoListService) GetListById(userId int, listId int) (domain.TodoList, error) {
	return s.repo.GetListById(userId, listId)
}

func (s *TodoListService) DeleteList(userId int, listId int) error {
	return s.repo.DeleteList(userId, listId)
}

func (s *TodoListService) UpdateList(userId int, listId int, input domain.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.UpdateList(userId, listId, input)
}
