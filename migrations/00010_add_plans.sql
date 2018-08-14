-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- drop the previous orders table
DROP TABLE IF EXISTS orders;

CREATE TABLE plans (
  id UUID PRIMARY KEY NOT NULL,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  exchange_name text NOT NULL,
  market_name text NOT NULL,
  active_currency_symbol text NOT NULL,
  active_currency_balance decimal NOT NULL,
  last_executed_plan_depth integer NOT NULL,
  last_executed_order_id UUID NOT NULL,
  plan_template_id text,             -- optional frontend plan template used with this plan - (Leo wanted this) 
  close_on_complete boolean DEFAULT false, -- sets status of this plan to 'closed' when the last order of plan finishes
  status text NOT NULL,              -- plan status is active, inactive, or failed
  created_on TIMESTAMP DEFAULT now(),
  updated_on TIMESTAMP DEFAULT current_timestamp
);


-- USDT -> BTC -> ETH -> BCH -> BTC -> USDT
-- Step 1: put in order to buy BTC using $200 in tether.
-- Step 2: using 100% of my BTC balance buy ETH.
-- Step 3: using 100% of my ETH balance sell for BCH.
-- Step 4: using 100% of my BCH balance sell for USDT.
CREATE TABLE orders (
  id UUID PRIMARY KEY NOT NULL,
  user_key_id UUID NOT NULL REFERENCES user_keys (id) ON DELETE CASCADE,
  parent_order_id UUID NOT NULL,      -- parent order of 0 means no parent 
  plan_id UUID NOT NULL REFERENCES plans (id) ON DELETE CASCADE,
  plan_depth integer NOT NULL,
  exchange_name text NOT NULL,
  market_name text NOT NULL,
  initial_currency_symbol text NOT NULL,
  initial_currency_balance decimal DEFAULT 0,
  initial_currency_traded decimal DEFAULT 0,
  initial_currency_remainder decimal DEFAULT 0,
  final_currency_symbol text,         -- they may be nullable 
  final_currency_balance decimal DEFAULT 0,   
  grupo text,                         -- leo will use el grupo to pull grouped orders (e.g. "scale-in" and scale-out orders)
  order_priority integer NOT NULL,    -- when two orders are at the same depth this can be used to indicate what your order preference is in terms of what get's executed first
  order_template_id text,             -- optional frontend template used for this order 
  order_type text NOT NULL,           -- limit, market, paper
  side text NOT NULL,                 -- buy, sell
  limit_price decimal DEFAULT 0,      -- limit price of order when order type is limit
  status text NOT NULL,               -- pending, active, failed, etc
  created_on TIMESTAMP DEFAULT now(),
  updated_on TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE triggers (
  id UUID PRIMARY KEY NOT NULL,
  order_id UUID NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
  trigger_template_id text,    -- optional frontend template used for this trigger 
  index integer NOT NULL,      -- assigned by Weo 
  title text,                  -- this will be human readable and displayed on client 
  name text NOT NULL,          -- used to identify specific trigger type so we can associate with the form input on the client
  code text NOT NULL,
  actions text[] NOT NULL,
  triggered BOOLEAN NOT NULL DEFAULT false,
  triggered_price decimal,     -- only if triggered is true 
  triggered_condition text,    -- e.g. "6056.11 <= 6400"
  triggered_timestamp TIMESTAMP,  
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

