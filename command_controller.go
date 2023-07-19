package main

type CommandController struct {
	add func([] string) ([]Todo, error)
	list func() ([]Todo, error)
	delete func(uid string) ([]Todo, error)
}
