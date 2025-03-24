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

INSERT INTO users 
  (email, first_name, second_name, password)
VALUES
  ('readME@test.com', 'read', 'ME', '$2a$10$sgEy3LehHNpbZ7NjqDhMiejJ8gaQTcykfv1VFJL42aPN8pZJL45EW'),
  ('testME@test.com', 'test', 'ME', '$2a$10$sgEy3LehHNpbZ7NjqDhMiejJ8gaQTcykfv1VFJL42aPN8pZJL45EW'),
  ('yetAnotherTest@test.com', 'yet', 'Another', '$2a$10$sgEy3LehHNpbZ7NjqDhMiejJ8gaQTcykfv1VFJL42aPN8pZJL45EW')
;

INSERT INTO sessions
  (user_id, session_token, csrf_token, expire_on)
VALUES
  (1,
   '2-bJbG-BU5h1fKovzqoEnwOxDsz9bm1-8vVRHYav5Z29DcaDUchc0LNufSGCEjKFsXGNtn0ZF0FdcXi9_npSGg==',
   'DpwoY8fzNfVyBnJDl9mEclJoZcWW8kxtZIo-CdMMvGnGfwzrrqwogUyVnUZknwazD_MXxEop5ewgxp2S-wTtig==',
   '2025-03-13 14:33:57'),
  (2,
   'lv1qhEGQgUn1Z3tFqKdCuvq_--W2ptuGp0oV5wRajtlC0sPN9xqsAxEZ6w2RGd-JX7nrk4_rO51tJXhoSONgmw==',
   '8NufL3TrLz3_NFgDY37cjg0LTjMPTzGc2jMsOW5GnKbjD5gP1HR4SQT6nDSCgjLtPi15FfGuvGpAQJNd2ckeUA==',
   '2025-03-13 14:45:48')
;

