-- +migrate Up
CREATE TABLE lists (
    list_id uuid NOT NULL DEFAULT uuid_generate_v4 (),
    name text NOT NULL,
    user_id uuid NOT NULL,
    CONSTRAINT lists_pkey PRIMARY KEY (list_id),
    CONSTRAINT users_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX "idx_lists_fk_user_id" ON "lists" ("user_id");

-- +migrate Down
DROP TABLE IF EXISTS lists;

