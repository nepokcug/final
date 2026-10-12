package api

import (
	"final/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

const q int = 50

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJson(w, map[string]string{"error": "метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}
	search := r.FormValue("search")
	tasks, err := db.Tasks(search, q)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	writeJson(w, TasksResp{Tasks: tasks}, http.StatusOK)
}
