package main

import (
	"database/sql"
)

type todoItemModel struct{
	Id int `json:"id"`
	Description string `json:"name"`
	Completed bool `json:"completed"`
}


func (p *todoItemModel) getTask(db *sql.DB) error {
	return db.QueryRow("SELECT DESCRIPTION, COMPLETED FROM tasks WHERE id=$1",
		p.Id).Scan(&p.Description, &p.Completed)
}

func (p *todoItemModel) updateTask(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE tasks SET DESCRIPTON=$1, COMPLETED=$2 WHERE id=$3",
			p.Description, p.Completed, p.Id)

	return err
}

func (p *todoItemModel) deleteTask(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM tasks WHERE id=$1", p.Id)

	return err
}

func (p *todoItemModel) createTask(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO tasks(DESCRIPTION, COMPLETED) VALUES($1, $2) RETURNING id",
		p.Description, p.Completed).Scan(&p.Id)

	if err != nil {
		return err
	}

	return nil
}

func getTasks(db *sql.DB, start, count int) ([]todoItemModel, error) {
	rows, err := db.Query(
		"SELECT ID, DESCRIPTION, COMPLETED FROM tasks LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := []todoItemModel{}

	for rows.Next() {
		var p todoItemModel
		if err := rows.Scan(&p.Id, &p.Description, &p.Completed); err != nil {
			return nil, err
		}
		tasks = append(tasks, p)
	}

	return tasks, nil
}

