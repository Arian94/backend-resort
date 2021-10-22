DROP TABLE IF EXISTS courier;
CREATE TABLE courier (
  id serial PRIMARY KEY,
  fullname varchar (50) NOT NULL,
  email varchar (25) NOT NULL,
  password VARCHAR (8) NOT NULL,
  todo json,
  delivered json
);
INSERT INTO
  courier (fullname, email, password)
VALUES
  ('Ghasem Peiki', 'gh.peyki@gm.com', '12345678'),
  ('Moji Peiki', 'm.peyki@yh.com', '12345678');