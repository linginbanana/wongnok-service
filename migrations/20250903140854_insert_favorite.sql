-- +goose Up
-- +goose StatementBegin
INSERT INTO favorites (food_recipe_id, user_id, created_at, updated_at)
VALUES
	(1, '392de118-0c0c-40e4-a628-9b77b1354c42', NOW(), NOW()),
	(2, '392de118-0c0c-40e4-a628-9b77b1354c42', NOW(), NOW()),
	(3, '392de118-0c0c-40e4-a628-9b77b1354c42', NOW(), NOW()),
	(4, '392de118-0c0c-40e4-a628-9b77b1354c42', NOW(), NOW()),
	(5, '392de118-0c0c-40e4-a628-9b77b1354c42', NOW(), NOW());
-- +goose StatementEnd


