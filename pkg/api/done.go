package api

import (
	"net/http"
	"time"

	"todo_app/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		sendError(w, "Не указан идентификатор задачи")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		sendError(w, "Задача не найдена")
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			sendError(w, "Ошибка при удалении задачи")
			return
		}
	} else {
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

		nextDate, err := NextDate(today, task.Date, task.Repeat)
		if err != nil {
			sendError(w, "Ошибка при вычислении следующей даты")
			return
		}

		err = db.UpdateTaskDate(id, nextDate)
		if err != nil {
			sendError(w, "Ошибка при обновлении задачи")
			return
		}
	}

	sendJSON(w, http.StatusOK, map[string]any{})
}
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		sendError(w, "Не указан идентификатор задачи")
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		sendError(w, "Задача не найдена")
		return
	}

	sendJSON(w, http.StatusOK, map[string]any{})
}
