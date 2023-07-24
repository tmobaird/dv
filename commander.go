package main

type Commander struct {
	Add func(args [] string, r ReaderWriter) ([]Todo, error)
	List func(r ReaderWriter) ([]Todo, error)
	Delete func(uid string, r ReaderWriter) ([]Todo, error)
	Done func(uid string, r ReaderWriter) ([]Todo, error)
	Undo func(uid string, r ReaderWriter) ([]Todo, error)
	Edit func(uid, newName string, r ReaderWriter) ([]Todo, error)

	ReaderWriter ReaderWriter
}
