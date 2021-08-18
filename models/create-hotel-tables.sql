DROP TABLE IF EXISTS single_room;
DROP TABLE IF EXISTS double_room;
DROP TABLE IF EXISTS triple_room;
DROP TABLE IF EXISTS twin_room;

CREATE TABLE single_room (
    id INT(3) AUTO_INCREMENT NOT NULL PRIMARY KEY,
    room_number INT(3) CHECK (room_number BETWEEN 101 AND 111),

    number_of_rooms INT NOT NULL CHECK (number_of_rooms BETWEEN 1 AND 3),

    room_subtype VARCHAR(13) NOT NULL CHECK (room_subtype IN ('standard','standard_plus','deluxe')), 
    fullname VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);
INSERT INTO single_room (room_number, number_of_rooms, room_subtype, fullname, email, start_date, end_date)
VALUES
	(101, 1, 'Standard', 'John Coltrane', 'j@y.com', '2021-07-27', '2021-07-29'),
    (101, 1, 'Standard', 'Arian Pourarian', 'a@g.com', '2021-08-10', '2021-08-14'),
    (104, 2, 'Standard_plus', 'Scott Hamshire', 's@h.com', '2021-08-13', '2021-08-19');

-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------

CREATE TABLE double_room (
    id INT(3) AUTO_INCREMENT NOT NULL PRIMARY KEY,
    room_number INT(3) CHECK (room_number BETWEEN 112 AND 121),

	number_of_rooms INT NOT NULL CHECK (number_of_rooms BETWEEN 1 AND 3),
    
    room_subtype VARCHAR(13) NOT NULL CHECK (room_subtype IN ('standard','standard_plus','deluxe')), 
    fullname VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);
INSERT INTO double_room (room_number, number_of_rooms, room_subtype, fullname, email, start_date, end_date)
VALUES
	(113, 1, 'Standard', 'Daniel Welbeck', 'j@y.com', '2021-08-05', '2021-08-13'),
    (112, 1, 'Standard', 'Parsa Jalili', 'a@g.com', '2021-08-10', '2021-08-14'),
    (115, 2, 'Standard_plus', 'Cathrine Willshere', 's@h.com', '2021-08-13', '2021-08-19');

-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------

CREATE TABLE triple_room (
    id INT(3) AUTO_INCREMENT NOT NULL PRIMARY KEY,
    room_number INT(3) CHECK (room_number BETWEEN 122 AND 131),

	number_of_rooms INT NOT NULL CHECK (number_of_rooms BETWEEN 1 AND 3),
    
    room_subtype VARCHAR(13) NOT NULL CHECK (room_subtype IN ('standard','standard_plus','deluxe')), 
    fullname VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);
INSERT INTO triple_room (room_number, number_of_rooms, room_subtype, fullname, email, start_date, end_date)
VALUES
	(123, 1, 'Standard', 'Adam Cole', 'j@y.com', '2021-08-05', '2021-08-13'),
    (124, 1, 'Deluxe', 'John Smith', 'a@g.com', '2021-08-10', '2021-08-14'),
    (129, 2, 'Standard_plus', 'Paravane Zibarooy', 's@h.com', '2021-08-13', '2021-08-19');

-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------

CREATE TABLE twin_room (
    id INT(3) AUTO_INCREMENT NOT NULL PRIMARY KEY,
    room_number INT(3) CHECK (room_number BETWEEN 132 AND 142),

	number_of_rooms INT NOT NULL CHECK (number_of_rooms BETWEEN 1 AND 3),
    
    room_subtype VARCHAR(13) NOT NULL CHECK (room_subtype IN ('standard','standard_plus','deluxe')), 
    fullname VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------

