CREATE TABLE sellers (
  seller_id SERIAL PRIMARY KEY,
  seller_name TEXT NOT NULL
);

CREATE TABLE categories (
  category_id SERIAL PRIMARY KEY,
  category_name TEXT NOT NULL
);

CREATE TABLE products (
  product_id SERIAL PRIMARY KEY,
  product_name TEXT NOT NULL,
  seller_id SERIAL REFERENCES sellers ON DELETE CASCADE, 
  category_id SERIAL REFERENCES categories ON DELETE CASCADE,
  product_price NUMERIC(10, 2) NOT NULL,
  CHECK (product_price > 0)
);

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  first_name TEXT NOT NULL,
  second_name TEXT NOT NULL,
  password TEXT NOT NULL
);

CREATE TABLE sessions (
  session_token TEXT PRIMARY KEY,
  csrf_token TEXT UNIQUE NOT NULL,
  expire_on TIMESTAMP WITH TIME ZONE NOT NULL
);
