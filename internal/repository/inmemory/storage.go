package inmemory

import (
	"errors"
	"sync"
	"time"
	"todo_crud/internal/domain/todo/models"

	"github.com/google/uuid"
)

type InMemoryStorage struct {
	data map[string]models.Todo
	mu   sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]models.Todo),
	}
}

func (s *InMemoryStorage) SaveTodo(todo models.Todo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if todo.UID == "" {
		todo.UID = uuid.New().String()
		todo.CreatedAt = time.Now()
	} else {
		existing, exists := s.data[todo.UID]
		if !exists {
			return errors.New("todo not found")
		}
		todo.CreatedAt = existing.CreatedAt
	}

	todo.UpdatedAt = time.Now()
	s.data[todo.UID] = todo
	return nil
}

func (s *InMemoryStorage) GetTodo(req models.TodoRequest) (models.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todo, ok := s.data[req.UID]
	if !ok {
		return models.Todo{}, errors.New("todo not found")
	}
	return todo, nil
}

func (s *InMemoryStorage) GetAllTodos() []models.Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todos := make([]models.Todo, 0, len(s.data))
	for _, v := range s.data {
		todos = append(todos, v)
	}
	return todos
}

func (s *InMemoryStorage) DeleteTodo(uid string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[uid]; !ok {
		return errors.New("todo not found")
	}
	delete(s.data, uid)
	return nil
}
