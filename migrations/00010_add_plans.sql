-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- drop the previous orders table
DROP TABLE IF EXISTS orders;

CREATE TABLE plans (
  id UUID PRIMARY KEY NOT NULL,
  plan_template_id text,             -- optional frontend plan template used with this plan - (Leo wanted this) 
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  user_key_id UUID NOT NULL REFERENCES user_keys (id) ON DELETE CASCADE,
  active_order_number integer NOT NULL,
  exchange_name text NOT NULL,
  market_name text NOT NULL,
  base_balance decimal NOT NULL,
  currency_balance decimal NOT NULL,
  status text NOT NULL,               -- plan status is active, inactive, or failed
  created_on TIMESTAMP DEFAULT now(),
  updated_on TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE orders (
  id UUID PRIMARY KEY NOT NULL,
  plan_id UUID NOT NULL REFERENCES plans (id) ON DELETE CASCADE,
  order_template_id text,             -- optional frontend template used for this order 
  balance_percent decimal DEFAULT 0,  -- percent of balance to use base_balance(buy) currency_balance(sell)
  side text NOT NULL,
  order_number integer NOT NULL,      -- defines the order sequence
  order_type text NOT NULL,           -- limit, market, paper
  limit_price decimal DEFAULT 0,
  next_order_id UUID,                 -- this would be the following order after this one
  status text NOT NULL,               -- pending, active, failed, etc
  created_on TIMESTAMP DEFAULT now(),
  updated_on TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE triggers (
  id UUID PRIMARY KEY NOT NULL,
  order_id UUID NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
  trigger_number integer NOT NULL,      -- defines the order sequence
  trigger_template_id text,             -- optional frontend template used for this trigger 
  name text NOT NULL,
  code jsonb NOT NULL,
  actions text[] NOT NULL,
  triggered BOOLEAN NOT NULL DEFAULT false,
  created_on TIMESTAMP DEFAULT now(),
  updated_on TIMESTAMP DEFAULT current_timestamp
);

-- Leo wanted these to represent the frontend templating system
CREATE TABLE plan_templates (
  id text PRIMARY KEY NOT NULL,       -- can be UUID or human readable text
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  title text,
  description text,
  attributes jsonb NOT NULL,
  category text,                      -- quick, planned, custom 
  created_on TIMESTAMP DEFAULT now(),
  updated_on TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE order_templates (
  id text PRIMARY KEY NOT NULL,       -- can be UUID or human readable
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  title text,
  description text,
  attributes jsonb NOT NULL,
  category text,                      -- simple, advance, custom 
  created_on TIMESTAMP DEFAULT now(),
  updated_on TIMESTAMP DEFAULT current_timestamp
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE triggers;
DROP TABLE orders;
DROP TABLE plans;
DROP TABLE plan_templates;
DROP TABLE order_templates;

