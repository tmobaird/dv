package main

import "github.com/google/uuid"

type Todo struct {
	Id        uuid.UUID `json:"id"`
    Name      string    `json:"name"`
    CreatedAt string    `json:"createdAt"`
    Done      bool      `json:"done"`
}