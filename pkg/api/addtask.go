package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"todoapp/pkg/db"
)

type addTaskResp struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, addTaskResp{Error: err.Error()})
		return
	}

	task.Title = strings.TrimSpace(task.Title)
	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, addTaskResp{Error: "title is required"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, addTaskResp{Error: err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, addTaskResp{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, addTaskResp{ID: strconv.FormatInt(id, 10)})
}

func checkDate(task *db.Task) error {
	now := dateOnly(time.Now())

	task.Date = strings.TrimSpace(task.Date)
	task.Repeat = strings.TrimSpace(task.Repeat)

	// если date пустая — ставим сегодня
	if task.Date == "" {
		task.Date = now.Format(dateLayout)
	}

	// date должна быть в нужном формате 
	t, err := time.Parse(dateLayout, task.Date)
	if err != nil {
		return errors.New("invalid date")
	}
	t = dateOnly(t)

	// если repeat задан — проверяем формат и заранее считаем следующую дату
	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err // в т.ч. неподдерживаемый формат
		}
	}

	// если task.Date меньше сегодняшней даты — корректируем
	if now.After(t) {
		if task.Repeat == "" {
			// без повторения — ставим сегодня
			task.Date = now.Format(dateLayout)
		} else {
			task.Date = next
		}
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
