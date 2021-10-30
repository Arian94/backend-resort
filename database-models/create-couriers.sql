-- DROP TYPE IF EXISTS orders CASCADE;
-- CREATE TYPE orders AS (
--   name text,
--   numberOfMeals integer
-- );
-- DROP TYPE IF EXISTS food_order CASCADE;
-- CREATE TYPE food_order AS (
--   _id text,
--   address  text,
--   email text,
--   orderDate timestamp,
--   orderState text,
--   orders orders[],
--   receiver text,
--   receiverPhoneNumber text,
--   totalPrice integer
-- );
DROP TABLE IF EXISTS courier;
CREATE TABLE courier (
  id serial PRIMARY KEY,
  fullname varchar (50) NOT NULL,
  email varchar (25) NOT NULL UNIQUE,
  phoneNumber CHAR (11) NOT NULL UNIQUE,
  password VARCHAR (8) NOT NULL,
  orders_list jsonb
);
INSERT INTO
  courier (fullname, email, phoneNumber, password)
VALUES
  (
    'Ghasem Peiki',
    'gh.peyki@gm.com',
    '09127563254',
    '12345678'
  ),
  (
    'Jafar Shahi',
    'j.shahi@gm.com',
    '09127563255',
    '12345678'
  ),
  (
    'Omid Baghi',
    'o.baghi@gm.com',
    '09127563256',
    '12345678'
  ),
  (
    'Pooria Vali',
    'p.vali@gm.com',
    '09127563257',
    '12345678'
  ),
  (
    'Moji Peiro',
    'm.peiro@yh.com',
    '09127563258',
    '12345678'
  );