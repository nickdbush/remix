package main

type Option[T any] struct{}

type UserParams struct {
	Name  string
	Email Option[string]
	Age   int
}

type User struct {
	ID    string
	Name  string
	Email Option[string]
	Age   int
}
