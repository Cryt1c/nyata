package com.example.nyata;

import java.util.ArrayList;
import java.util.concurrent.atomic.AtomicLong;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class TodoController {

	private final AtomicLong counter = new AtomicLong();
	private ArrayList<Todo> todos = new ArrayList<Todo>();

	@GetMapping("/todos")
	public ArrayList<Todo> getTodos() {
		return todos;
	}

	@PostMapping("/todo")
	public Todo createTodo(@RequestBody Todo todo) {
		Todo newTodo = new Todo(counter.incrementAndGet(), todo.name(), todo.completed());
		todos.add(newTodo);
		return newTodo;
	}
}
