package goodidea

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// We could try and use CopyFrom https://pkg.go.dev/github.com/jackc/pgx/v5#hdr-Copy_Protocol
// but the transition point to where that is quicker is around 5 rows and we don't expect more than
// 5 images on a given task
func saveTaskImages(i uint32, p []string) {
	//TODO: Pass in a context, add failed rows to the context, so we can possibly return that later
	ctx := context.Background()
	for _, pth := range p {
		if _, err := DB.Exec(ctx, "INSERT INTO task_imgs(task_id, img_path) VALUES ($1, $2)", i, pth); err != nil {
			Logr.Error("Unable to save image path record", "task", i, "imagePath", pth)
		}
	}
}

// getTaskImages List all of the paths for images associated with a given task
// Not currently in the main function to show a task because we're assuming most tasks
// will not have an image
// Parameters:
//
//	i - TaskID
func getTaskImages(i uint32) ([]string, error) {
	ctx := context.Background()
	query := "SELECT img_path FROM task_imgs WHERE task_id = $1"
	imgs := make([]string, 0, 4)
	rows, err := DB.Query(ctx, query, i)
	if err != nil {
		if err == pgx.ErrNoRows {
			return imgs, nil
		}
		return imgs, err
	}

	defer rows.Close()
	for rows.Next() {
		var s string
		err = rows.Scan(&s)
		if err != nil {
			Logr.Error("Unable to marshal image path response into string", "err", err.Error())
			return imgs, err
		}
		imgs = append(imgs, s)
	}
	return imgs, nil
}
