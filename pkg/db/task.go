package db

import (
	"fmt"
	"strconv"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
func Tasks(search string, limit int) ([]Task, error) {
	var query string
	var args []any

	if search == "" {
		query = `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
		args = append(args, limit)
	} else {

		parsedDate, err := time.Parse("02.01.2006", search)
		if err == nil {
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
			args = append(args, parsedDate.Format("20060102"), limit)
		} else {

			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
			pattern := "%" + search + "%"
			args = append(args, pattern, pattern, limit)
		}
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
func GetTask(id string) (Task, error) {
	var t Task

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return t, fmt.Errorf("неверный формат идентификатора")
	}

	query := `SELECT date, title, comment, repeat FROM scheduler WHERE id = ?`
	err = DB.QueryRow(query, idInt).Scan(&t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return t, err
	}

	t.ID = id
	return t, nil
}

func UpdateTask(t Task) error {

	idInt, err := strconv.ParseInt(t.ID, 10, 64)
	if err != nil {
		return fmt.Errorf("неверный формат идентификатора")
	}
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, t.Date, t.Title, t.Comment, t.Repeat, idInt)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}
func DeleteTask(id string) error {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("неверный формат идентификатора")
	}

	res, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, idInt)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

func UpdateTaskDate(id string, date string) error {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("неверный формат идентификатора")
	}

	res, err := DB.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, date, idInt)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
