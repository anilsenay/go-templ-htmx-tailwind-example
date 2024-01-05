package model

type Todo struct {
	Id   int
	Text string
	Done bool
}

type TodoList struct {
	Id       int
	Name     string
	HexColor string
}
