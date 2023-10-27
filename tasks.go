package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type Task struct {
	ID          uint32     `json:"id"`
	Status      bool       `json:"status"`
	Title       string     `json:"title"`
	Body        *string    `json:"body"`
	Score       int32      `json:"score"`
	CompletedAt *time.Time `json:"completedAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	DeletedAt   time.Time  `json:"deletedAt"`
}

func getAllTasks() ([]Task, error) {
	ctx := context.Background()
	query := "SELECT id, status, title, body, score, completed_at, created_at FROM tasks WHERE deleted_at IS NULL ORDER BY score DESC"

	tasks := make([]Task, 0, 8)
	rows, err := db.Query(ctx, query)
	if err != nil {
		if err == pgx.ErrNoRows {
			return tasks, nil
		}
		return tasks, err
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		err = rows.Scan(
			&task.ID,
			&task.Status,
			&task.Title,
			&task.Body,
			&task.Score,
			&task.CompletedAt,
			&task.CreatedAt,
		)
		if err != nil {
			log.Println("Unable to marshal DB response into struct")
			return tasks, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func addTask(title, details string) error {
	ctx := context.Background()
	query := "INSERT INTO tasks(title, body) VALUES($1, $2)"

	_, err := db.Exec(ctx, query, title, details)
	return err
}

func updateTaskScore(id uint32, inc bool) (int32, error) {
	var score int32
	ctx := context.Background()
	var symbol = "-"
	if inc {
		symbol = "+"
	}
	query := fmt.Sprintf("UPDATE tasks SET score = score %s 1 WHERE id = $1", symbol)
	if _, err := db.Exec(ctx, query, id); err != nil {
		return 0, err
	}

	if err := db.QueryRow(ctx, "SELECT score FROM tasks WHERE id = $1", id).Scan(&score); err != nil {
		return 0, err
	}

	return score, nil
}
