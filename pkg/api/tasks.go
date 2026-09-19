package api

import (
	"net/http"

	"todo_app/pkg/db"
)

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	search := r.FormValue("search")

	tasks, err := db.Tasks(search, 50)
	if err != nil {
		sendError(w, "Ошибка при получении списка задач")
		return
	}
	sendJSON(w, http.StatusOK, map[string]any{
		"tasks": tasks,
	})
}
