-- +goose Up
-- +goose StatementBegin
ALTER TABLE couriers ADD COLUMN transport_type TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE couriers DROP COLUMN transport_type
-- +goose StatementEnd
