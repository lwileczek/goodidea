package goodidea

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
	Comments    []Comment  `json:"comments"`
	Images      []string   `json:"images"` //paths to get all of the images associated with task
}

func getAllTasks(limit int16) ([]Task, error) {
	ctx := context.Background()
	query := `SELECT
	id,
	title,
	-- Don't return the entire body, this is for a summary page
	SUBSTRING(body, 1, 256) AS body,
	score
FROM
	tasks
WHERE
	deleted_at IS NULL
ORDER BY
	score DESC`

    if limit > 0 {
        query += fmt.Sprintf("\nLIMIT %d", limit)
    }

	tasks := make([]Task, 0, 8)
	rows, err := DB.Query(ctx, query)
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
			&task.Title,
			&task.Body,
			&task.Score,
		)
		if err != nil {
			log.Println("Unable to marshal DB response into struct")
			return tasks, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func getSomeTasks(clue string) ([]Task, error) {
	ctx := context.Background()
	query := `SELECT
	id,
	title,
	-- Don't return the entire body, this is for a summary page
	SUBSTRING(body, 1, 256) AS body,
	score
FROM
	tasks
WHERE
	deleted_at IS NULL AND
	title ilike '%' || $1 || '%'
ORDER BY
	score DESC`

	tasks := make([]Task, 0, 8)
	rows, err := DB.Query(ctx, query, clue)
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
			&task.Title,
			&task.Body,
			&task.Score,
		)
		if err != nil {
			log.Println("Unable to marshal DB response into struct")
			return tasks, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func getTasksByID(id uint32) (Task, error) {
	ctx := context.Background()
	query := "SELECT id, status, title, body, score, completed_at, created_at FROM tasks WHERE id = $1 and deleted_at IS NULL"

	var task Task
	if err := DB.QueryRow(ctx, query, id).Scan(
		&task.ID,
		&task.Status,
		&task.Title,
		&task.Body,
		&task.Score,
		&task.CompletedAt,
		&task.CreatedAt,
	); err != nil {
		Logr.Error("Could not query task from DB", "err", err)
		return task, err
	}
	return task, nil
}

func addTask(title, details string) (uint32, error) {
	ctx := context.Background()
	query := "INSERT INTO tasks(title, body) VALUES($1, $2) RETURNING id"

	var newID uint32
	if err := DB.QueryRow(ctx, query, title, details).Scan(&newID); err != nil {
		return 0, err
	}
	return newID, nil
}

func updateTaskScore(id uint32, inc bool) (int32, error) {
	var score int32
	ctx := context.Background()
	var symbol = "-"
	if inc {
		symbol = "+"
	}
	query := fmt.Sprintf("UPDATE tasks SET score = score %s 1 WHERE id = $1", symbol)
	if _, err := DB.Exec(ctx, query, id); err != nil {
		return 0, err
	}

	if err := DB.QueryRow(ctx, "SELECT score FROM tasks WHERE id = $1", id).Scan(&score); err != nil {
		return 0, err
	}

	return score, nil
}
