package main

import (
	"encoding/json"
	"net/http"
)

type Todo struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
	Id        int    `json:"id"`
}

func main() {
	http.HandleFunc("/todos", todosHandler)
	http.ListenAndServe(":8080", nil)
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	todo := Todo{Id: 0, Completed: false, Name: "Do dishes"}
	todos := []Todo{todo, todo}
	json, _ := json.Marshal(todos)
	w.Write([]byte(json))
}
