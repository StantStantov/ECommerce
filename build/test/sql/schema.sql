CREATE SCHEMA market;
SET search_path to market;

CREATE TABLE market.sellers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE market.categories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE market.products (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  seller_id UUID  REFERENCES sellers(id) ON DELETE CASCADE, 
  category_id UUID REFERENCES categories(id) ON DELETE CASCADE,
  price NUMERIC(10, 2) NOT NULL,
  CHECK (price > 0)
);

CREATE TABLE market.users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT NOT NULL UNIQUE,
  first_name TEXT NOT NULL,
  second_name TEXT NOT NULL,
  password TEXT NOT NULL
);

CREATE TABLE market.sessions (
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  session_token TEXT PRIMARY KEY,
  csrf_token TEXT UNIQUE NOT NULL,
  expire_on TIMESTAMP WITH TIME ZONE NOT NULL
);

