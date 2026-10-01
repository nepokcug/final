package db

import (
	"database/sql"
	"errors"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	// определите запрос
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?,?,?,?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func Tasks(search string, limit int) ([]*Task, error) {
	if search == "" {
		// 1. Если поиск пустой
		rows, err := DB.Query("SELECT id,date,title,comment,repeat FROM scheduler ORDER BY date LIMIT ?", limit)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return scanTasks(rows)
	} else if t, err := time.Parse("02.01.2006", search); err == nil { // 2. Если парсится -> время
		date := t.Format("20060102")
		rows, err := DB.Query("SELECT id,date,title,comment,repeat FROM scheduler WHERE date=? ORDER BY title LIMIT ?", date, limit)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return scanTasks(rows)
	} else { // 3. Если не пусто и не время -> запрос
		rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? ORDER BY date LIMIT ?",
			"%"+search+"%", limit)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return scanTasks(rows)
	}
}

// scanTasks читает строки из rows и возвращает слайс задач
func scanTasks(rows *sql.Rows) ([]*Task, error) {
	var tasks []*Task

	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	task := &Task{}
	err := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id=?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("задача не найдена")
		}
		return nil, err
	}
	return task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// была применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("некорректный id для обновления информации")
	}
	return nil
}
