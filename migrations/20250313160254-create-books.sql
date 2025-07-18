-- +migrate Up
CREATE TABLE books (
    book_id uuid NOT NULL DEFAULT uuid_generate_v4 (),
    title text NOT NULL,
    author text NOT NULL,
    publisher text,
    book_description text,
    genre text,
    pages integer,
    rating integer,
    image_link text,
    user_notes text,
    user_id uuid NOT NULL,
    CONSTRAINT books_pkey PRIMARY KEY (book_id),
    CONSTRAINT users_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX "idx_books_fk_user_id" ON "books" ("user_id");

-- +migrate Down
DROP TABLE IF EXISTS books;

