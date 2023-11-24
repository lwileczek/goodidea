CREATE TABLE task_imgs (
    -- The unique record ID
    id           SERIAL not null PRIMARY KEY,
    -- The ID of the task the image is related too
    task_id      INTEGER NOT NULL,
    -- The path the image will be served from which will be used in an <img></img>
    img_path VARCHAR(255),
    -- Status, is this complete or not
    CONSTRAINT fk_task
        FOREIGN KEY (task_id)
            REFERENCES tasks(id) ON DELETE CASCADE
);

/* Index our images table so two images cannot have the same name for a given task and
 * we can quickly find images  for a given task.*/
CREATE UNIQUE INDEX task_imgs_idx ON task_imgs (task_id, img_path);
