package db

import (
	"database/sql"
	"errors"
	"strconv"
)

const defaultTasksLimit = 50



type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)
	`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Tasks(limit int) ([]*Task, error) {
	if limit <= 0 {
		limit = defaultTasksLimit
	}

	rows, err := DB.Query(`
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date ASC, id ASC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task

	for rows.Next() {
		var (
			id      int64
			date    string
			title   string
			comment string
			repeat  string
		)

		if err := rows.Scan(&id, &date, &title, &comment, &repeat); err != nil {
			return nil, err
		}

		tasks = append(tasks, &Task{
			ID:      strconv.FormatInt(id, 10),
			Date:    date,
			Title:   title,
			Comment: comment,
			Repeat:  repeat,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = make([]*Task, 0)
	}

	return tasks, nil
}

var ErrTaskNotFound = errors.New("task not found")

func GetTask(id string) (*Task, error) {
	row := DB.QueryRow(`
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?
	`, id)

	var (
		dbID    int64
		date    string
		title   string
		comment string
		repeat  string
	)

	err := row.Scan(&dbID, &date, &title, &comment, &repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	return &Task{
		ID:      strconv.FormatInt(dbID, 10),
		Date:    date,
		Title:   title,
		Comment: comment,
		Repeat:  repeat,
	}, nil
}

func UpdateTask(task *Task) error {
	res, err := DB.Exec(`
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func UpdateDate(id string, next string) error {
	res, err := DB.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return ErrTaskNotFound
	}
	return nil
}
