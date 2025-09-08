package server

import (
	"fmt"
	"net/http"
	"time"
	"todo_crud/internal"
	todoDomain "todo_crud/internal/domain/todo/models"

	"github.com/gin-gonic/gin"
)

type TodoStorage interface {
	SaveTodo(todo todoDomain.Todo) error
	GetTodo(todoReq todoDomain.TodoRequest) (todoDomain.Todo, error)
	GetAllTodos() []todoDomain.Todo
	DeleteTodo(uid string) error
}

type ToDoAPI struct {
	httpServer *http.Server
	storage    TodoStorage
}

func NewServer(cfg internal.Config, db TodoStorage) *ToDoAPI {
	httpSrv := http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}

	api := ToDoAPI{
		httpServer: &httpSrv,
		storage:    db,
	}

	api.configRouters()

	return &api
}

func (api *ToDoAPI) Run() error {
	return api.httpServer.ListenAndServe()
}

func (api *ToDoAPI) Shutdown() error {
	return nil
}

func (api *ToDoAPI) configRouters() {
	router := gin.Default()

	tasks := router.Group("/tasks")
	{
		tasks.GET("/", api.getAllTodos)
		tasks.GET("/:id", api.getTodoByID)
		tasks.POST("/", api.createTodo)
		tasks.PUT("/:id", api.updateTodo)
		tasks.DELETE("/:id", api.deleteTodo)
	}

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello world")
	})

	api.httpServer.Handler = router
}

func (api *ToDoAPI) getAllTodos(ctx *gin.Context) {
	todos := api.storage.GetAllTodos()
	for _, todo := range todos {
		ctx.JSON(http.StatusOK, todo)
	}
}

func (api *ToDoAPI) getTodoByID(ctx *gin.Context) {
	id := ctx.Param("id")
	todo, err := api.storage.GetTodo(todoDomain.TodoRequest{UID: id})
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, todo)
}

func (api *ToDoAPI) createTodo(ctx *gin.Context) {
	var todo todoDomain.Todo
	if err := ctx.ShouldBindJSON(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := api.storage.SaveTodo(todo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, todo)
}

func (api *ToDoAPI) updateTodo(ctx *gin.Context) {
	id := ctx.Param("id")

	var input todoDomain.Todo
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	existing, err := api.storage.GetTodo(todoDomain.TodoRequest{UID: id})
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}
	if input.Title != "" {
		existing.Title = input.Title
	}
	existing.IsDone = input.IsDone
	existing.UpdatedAt = time.Now()

	if err := api.storage.SaveTodo(existing); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, existing)
}

func (api *ToDoAPI) deleteTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := api.storage.DeleteTodo(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
