-- +migrate Up
CREATE TABLE book_list (
    book_id uuid NOT NULL DEFAULT uuid_generate_v4 (),
    list_id uuid NOT NULL DEFAULT uuid_generate_v4 (),
    CONSTRAINT book_list_pkey PRIMARY KEY (book_id, list_id),
    CONSTRAINT books_fkey FOREIGN KEY (book_id) REFERENCES books (book_id),
    CONSTRAINT lists_fkey FOREIGN KEY (list_id) REFERENCES lists (list_id)
);

CREATE INDEX "idx_book_list_fk_list_id" ON "book_list" ("list_id");

-- +migrate Down
DROP TABLE IF EXISTS book_list;

