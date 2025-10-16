-- +goose Up
-- +goose StatementBegin
ALTER TABLE food_recipes ADD COLUMN user_id VARCHAR(100);
ALTER TABLE food_recipes ADD CONSTRAINT food_recipes_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE food_recipes
DROP COLUMN user_id;

-- +goose StatementEnd