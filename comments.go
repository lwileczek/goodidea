package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type Comment struct {
	//Unique ID for this post
	ID uint64 `json:"id"`
	//The foreign key back to the tasks
	TaskID uint32 `json:"taskId"`
	//Optionally the person can leave a username
	User *string `json:"user"`
	//The comment box
	Content string `json:"comment"`
	//The time the comment was created
	CreatedAt time.Time `json:"createdAt"`
}

func getAllTaskComments(taskID uint32) ([]Comment, error) {
	ctx := context.Background()
	query := `SELECT
		username,
		content,
		created_at
	FROM 
		comments
	WHERE
		task_id = $1`

	comments := make([]Comment, 0, 4)
	rows, err := db.Query(ctx, query, taskID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return comments, nil
		}
		return comments, err
	}
	defer rows.Close()
	for rows.Next() {
		var c Comment
		err = rows.Scan(
			&c.User,
			&c.Content,
			&c.CreatedAt,
		)
		if err != nil {
			log.Println("Unable to marshal DB response into struct")
			return comments, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func addComment(t uint32, u *string, c string) error {
	ctx := context.Background()
	query := "INSERT INTO comments(task_id, username, comment) VALUES($1, $2, $3)"

	_, err := db.Exec(ctx, query, t, u, c)
	return err
}
