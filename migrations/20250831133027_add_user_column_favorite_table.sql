-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS favorites (
    id SERIAL PRIMARY KEY,
    food_recipe_id INT NOT NULL REFERENCES food_recipes,
    user_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS favorites;

-- +goose StatementEnd