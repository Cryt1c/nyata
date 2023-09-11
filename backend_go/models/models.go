package models

type Todo struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
	Id        int    `json:"id"`
}
