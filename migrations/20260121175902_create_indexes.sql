-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX ux_couriers_phone ON couriers(phone);
CREATE INDEX idx_couriers_status ON couriers(status);

CREATE UNIQUE INDEX ux_delivery_order_id ON delivery(order_id);
CREATE INDEX idx_delivery_courier_deadline ON delivery(courier_id, deadline);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_delivery_courier_deadline;
DROP INDEX IF EXISTS ux_delivery_order_id;

DROP INDEX IF EXISTS idx_couriers_status;
DROP INDEX IF EXISTS ux_couriers_phone;
-- +goose StatementEnd
