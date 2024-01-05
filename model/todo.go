package model

type Todo struct {
	Id   int
	Text string
	Done bool
}

type Collection struct {
	Id       int
	Name     string
	HexColor string
}

type CollectionWithTodoList struct {
	Collection
	List []Todo
}
