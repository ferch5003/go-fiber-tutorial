-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
   id INT PRIMARY KEY AUTO_INCREMENT,
   first_name VARCHAR(255) NOT NULL,
   last_name VARCHAR(255) NOT NULL,
   email VARCHAR(255) NOT NULL UNIQUE,
   password VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS todos (
   id INT PRIMARY KEY AUTO_INCREMENT,
   title VARCHAR(255) NOT NULL,
   description TEXT,
   completed BOOLEAN NOT NULL DEFAULT FALSE,
   user_id INT NOT NULL,
   KEY user_id_idx (user_id)
);
-- +goose StatementEnd
