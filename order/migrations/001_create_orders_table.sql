-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE orders (
    uuid UUID primary key,
    user_uuid UUID not null,
    total_price numeric(10, 2) not null,
    transaction_uuid UUID,
    payment_method varchar(50),
    status varchar(50) not null,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone
);

CREATE TABLE order_parts (
    order_uuid UUID not null,
    part_uuid UUID not null,
    primary key (order_uuid, part_uuid),
    foreign key (order_uuid) references orders(uuid) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_parts;
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd