package world

import (
	"errors"
	"time"
)

type InputHelloWorld struct {
	Message string
}

type HelloWorldID string

type HelloWorld struct {
	ID        HelloWorldID
	Message   string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

var (
	ErrNotFound      = errors.New("hello world not found")
	ErrCannotBeEmpty = errors.New("hello world field cannot be empty")
)
