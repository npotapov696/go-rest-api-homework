package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task описывает задачу.
type Task struct {
	ID           string   `json:"id"`           // ID задачи
	Description  string   `json:"description"`  // Заголовок
	Note         string   `json:"note"`         // Описание задачи
	Applications []string `json:"applications"` // Используемые приложения
}

// task содержит стартовый набор задач.
var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// getTasks возвращает все задачи, содержащиеся в мапе tasks.
func getTasks(w http.ResponseWriter, r *http.Request) {

	// сериализуем данные из мапы tasks
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// в заголовок записываем тип контента - данные в формате JSON
	w.Header().Set("Content-Type", "application/json")

	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)

	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}

// postTask добавляет задачу, содержащуюся в теле запроса, в мапу tasks.
func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	// записываем данные из тела запроса в буфер, обрабатываем возможную ошибку записи данных
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// десериализуем данные из буфера и записываем их в перемнную task, обрабатываем возможную ошибку десериализации
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// проверяем на наличие в мапе tasks элемента с id, совпадающим с id из новой задачи в переменной task,
	// если существует, обрабатываем ошибку
	if _, ok := tasks[task.ID]; ok {
		http.Error(w, "Задача с указанным ID уже существует", http.StatusBadRequest)
		return
	}

	// добавляем в мапу tasks новый элемент из переменной task
	tasks[task.ID] = task

	// так как все успешно, возвращаем статус успешного создания нового элемента
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func main() {
	// создаем роутер
	r := chi.NewRouter()
	// регистрируем в роутере эндпоинт `/tasks` с методом GET, для которого используется обработчик `getTasks`
	r.Get("/tasks", getTasks)
	// регистрируем в роутере эндпоинт `/tasks` с методом POST, для которого используется обработчик `postTask`
	r.Post("/tasks", postTask)

	// запускаем сервер, обрабатываем возможную ошибку
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
