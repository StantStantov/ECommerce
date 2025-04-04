INSERT INTO market.sellers
  (id, name)
VALUES
  ('f4d234ff-7aa5-4986-954c-8c2cc61ea0fc', 'Balam Industries'),
  ('7e13d4e2-408b-494f-a611-1950a3a36616', 'Arquebus Corporation')
;

INSERT INTO market.categories
  (id, name)
VALUES
  ('c735f60a-bebf-4d2f-a016-190a883eb99f', 'Head'), 
  ('7670dd24-fffd-4ede-8fad-17613ec6f2ba', 'Core'),
  ('70b0d225-f526-4c8b-aafd-cdea3f2977d2', 'Arms'),
  ('10021a86-d948-4c54-bdf2-00df93a22add', 'Legs')
;

WITH
    ins(id, name, seller_name, category_name, price) AS (
        VALUES
            (
                '02cab72f-e225-4c3d-b725-faaa5d66ca74',
                'HD-011 MELANDER',
                'Balam Industries',
                'Head',
                75000
            ),
            (
                '961dd7d0-95b3-492e-833e-c33875a64d0f',
                'HD-033M VERRILL',
                'Balam Industries',
                'Head',
                205000
            ),
            (
                '55a47ff6-d8e2-491b-9378-8d536c1f4e44',
                'VP-44S',
                'Arquebus Corporation',
                'Head',
                124000
            ),
            (
                '29f9978e-77f1-4e03-a054-9244d6bb00d0',
                'VP-44D',
                'Arquebus Corporation',
                'Head',
                231000
            )
    )
    INSERT INTO market.products
    (id, name, seller_id, category_id, price)
SELECT CAST(ins.id AS uuid), ins.name, sellers.id, categories.id, ins.price
FROM ins
JOIN market.sellers ON ins.seller_name = sellers.name
JOIN market.categories ON ins.category_name = categories.name
;

INSERT INTO market.users 
  (id, email, first_name, second_name, password)
VALUES
  ('ad43dfbf-1152-478c-a595-e3ebe5ad0085', 
    'readME@test.com', 'read', 'ME', 
    '$2a$10$sgEy3LehHNpbZ7NjqDhMiejJ8gaQTcykfv1VFJL42aPN8pZJL45EW'),
  ('eeb70fd3-47ae-40d8-b404-e1a2bd31afc5', 
    'testME@test.com', 'test', 'ME', 
    '$2a$10$sgEy3LehHNpbZ7NjqDhMiejJ8gaQTcykfv1VFJL42aPN8pZJL45EW'),
  ('02f95483-7934-41b7-af1e-40eaf67817fc', 
    'yetAnotherTest@test.com', 'yet', 'Another', 
    '$2a$10$sgEy3LehHNpbZ7NjqDhMiejJ8gaQTcykfv1VFJL42aPN8pZJL45EW')
;

INSERT INTO market.sessions
  (user_id, session_token, csrf_token, expire_on)
VALUES
  ('ad43dfbf-1152-478c-a595-e3ebe5ad0085',
   '2-bJbG-BU5h1fKovzqoEnwOxDsz9bm1-8vVRHYav5Z29DcaDUchc0LNufSGCEjKFsXGNtn0ZF0FdcXi9_npSGg==',
   'DpwoY8fzNfVyBnJDl9mEclJoZcWW8kxtZIo-CdMMvGnGfwzrrqwogUyVnUZknwazD_MXxEop5ewgxp2S-wTtig==',
   '2025-03-13 14:33:57'),
  ('eeb70fd3-47ae-40d8-b404-e1a2bd31afc5',
   'lv1qhEGQgUn1Z3tFqKdCuvq_--W2ptuGp0oV5wRajtlC0sPN9xqsAxEZ6w2RGd-JX7nrk4_rO51tJXhoSONgmw==',
   '8NufL3TrLz3_NFgDY37cjg0LTjMPTzGc2jMsOW5GnKbjD5gP1HR4SQT6nDSCgjLtPi15FfGuvGpAQJNd2ckeUA==',
   '2025-03-13 14:45:48')
;

