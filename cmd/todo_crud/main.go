package main

import (
	"fmt"
	"todo_crud/internal"
	"todo_crud/internal/repository/inmemory"
	"todo_crud/internal/server"
)

func main() {
	fmt.Println("Server started...")
	cfg := internal.ReadConfig()
	fmt.Println(cfg)
	db := inmemory.NewInMemoryStorage()
	srv := server.NewServer(cfg, db)

	if err := srv.Run(); err != nil {
		panic(err)
	}
}
