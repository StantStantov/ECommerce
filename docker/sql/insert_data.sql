INSERT INTO sellers
  (seller_name)
VALUES
  ('HUAWEI'),
  ('Lenovo'),
  ('Apple')
;

INSERT INTO categories
  (category_name)
VALUES
  ('Laptops'), 
  ('Phones'),
  ('Electronics')
;

INSERT INTO products              
  (product_name, seller_id, category_id, product_price)
VALUES
  ('MateBook', 1, 1, 15000),
  ('ThinkPad', 2, 1, 10000),
  ('MacBook 1', 3, 1, 15000),
  ('MacBook 2', 3, 1, 15000),
  ('MacBook 3', 3, 1, 15000),
  ('Iphone', 3, 2, 15000)
;

