-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE plans (
 id UUID PRIMARY KEY NOT NULL,
 user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
 user_key_id UUID NOT NULL REFERENCES user_keys (id) ON DELETE CASCADE,
 exchange_name text NOT NULL,
 market_name text NOT NULL,
 plan_order_ids uuid[] NOT NULL,
 base_balance decimal NOT NULL,
 currency_balance decimal NOT NULL,
 status text NOT NULL,               -- plan status is active, inactive, or failed
 created_on TIMESTAMP DEFAULT now(),
 updated_on TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE plan_orders (
 id UUID PRIMARY KEY NOT NULL,
 plan_id UUID NOT NULL REFERENCES plans (id) ON DELETE CASCADE,
 base_percent decimal DEFAULT 0,
 currency_percent decimal DEFAULT 0,
 side text NOT NULL,
 order_type text NOT NULL,           -- limit, market, paper
 conditions jsonb NOT NULL,
 price decimal DEFAULT 0,
 status text NOT NULL,               -- pending, active, failed, etc
 next_plan_order_id UUID,
 created_on TIMESTAMP DEFAULT now(),
 updated_on TIMESTAMP DEFAULT current_timestamp
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE plan_orders;
DROP TABLE plans;
