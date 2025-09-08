package models

import "time"

type Todo struct {
	UID       string `json:"uid"`
	Title     string `json:"title"`
	Value     string `json:"value"`
	IsDone    bool   `json:"is_done"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TodoRequest struct {
	UID   string `json:"uid"`
	Title string `json:"title"`
}
