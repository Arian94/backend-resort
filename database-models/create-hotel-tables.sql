DROP TABLE IF EXISTS single_room;
DROP TABLE IF EXISTS double_room;
DROP TABLE IF EXISTS triple_room;
DROP TABLE IF EXISTS twin_room;

CREATE TABLE single_room (
    id INT(3) AUTO_INCREMENT NOT NULL PRIMARY KEY,
    -- room_mark INT(3) CHECK (room_mark BETWEEN 101 AND 111),
    room_mark VARCHAR(16),

    number_of_rooms INT NOT NULL CHECK (number_of_rooms BETWEEN 1 AND 3),

    room_subtype VARCHAR(13) NOT NULL CHECK (room_subtype IN ('standard','standard_plus','deluxe')), 
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);
INSERT INTO single_room (room_mark, number_of_rooms, room_subtype, full_name, email, start_date, end_date)
VALUES
	('101', 1, 'Standard', 'John Coltrane', 'j.coltrone@yh.com', '2021-07-27', '2021-07-29'),
    ('101', 1, 'Standard', 'Arian Pourarian', 'a.pourarian@gm.com', '2021-08-10', '2021-08-14'),
    ('104', 2, 'Standard_plus', 'Scott Hamshire', 's.hamshire@gm.com', '2021-08-13', '2021-08-19');

-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------

CREATE TABLE double_room (
    id INT(3) AUTO_INCREMENT NOT NULL PRIMARY KEY,
    -- room_mark INT(3) CHECK (room_mark BETWEEN 112 AND 121),
    room_mark VARCHAR(16),

	number_of_rooms INT NOT NULL CHECK (number_of_rooms BETWEEN 1 AND 3),
    
    room_subtype VARCHAR(13) NOT NULL CHECK (room_subtype IN ('standard','standard_plus','deluxe')), 
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);
INSERT INTO double_room (room_mark, number_of_rooms, room_subtype, full_name, email, start_date, end_date)
VALUES
	('113', 1, 'Standard', 'Daniel Welbeck', 'd.welbeck@yh.com', '2021-08-05', '2021-08-13'),
    ('112', 1, 'Standard', 'Arian Pourarian', 'a.pourarian@gm.com', '2021-08-10', '2021-08-14'),
    ('115', 2, 'Standard_plus', 'Cathrine Willshere', 'c.willshere@ht.com', '2021-08-13', '2021-08-19');

-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------

CREATE TABLE triple_room (
    id INT(3) AUTO_INCREMENT NOT NULL PRIMARY KEY,
    -- room_mark INT(3) CHECK (room_mark BETWEEN 122 AND 131),
    room_mark VARCHAR(16),

	number_of_rooms INT NOT NULL CHECK (number_of_rooms BETWEEN 1 AND 3),
    
    room_subtype VARCHAR(13) NOT NULL CHECK (room_subtype IN ('standard','standard_plus','deluxe')), 
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);
INSERT INTO triple_room (room_mark, number_of_rooms, room_subtype, full_name, email, start_date, end_date)
VALUES
	('123', 1, 'Standard', 'Adam Cole', 'a.cole@yh.com', '2021-08-05', '2021-08-13'),
    ('124', 1, 'Deluxe', 'John Smith', 'j.smith@gm.com', '2021-08-10', '2021-08-14'),
    ('129', 2, 'Standard_plus', 'Paravane Zibarooy', 'p.ziba@gm.com', '2021-08-13', '2021-08-19');

-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------

CREATE TABLE twin_room (
    id INT(3) AUTO_INCREMENT NOT NULL PRIMARY KEY,
    -- room_mark INT(3) CHECK (room_mark BETWEEN 132 AND 142),
    room_mark VARCHAR(16),

	number_of_rooms INT NOT NULL CHECK (number_of_rooms BETWEEN 1 AND 3),
    
    room_subtype VARCHAR(13) NOT NULL CHECK (room_subtype IN ('standard','standard_plus','deluxe')), 
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);
INSERT INTO twin_room (room_mark, number_of_rooms, room_subtype, full_name, email, start_date, end_date)
VALUES
	('135', 1, 'Standard', 'Michael Sumatra', 'm.sumatra@gm.com', '2021-10-05', '2021-10-15'),
    ('137', 1, 'Deluxe', 'Manuel Sandler', 'm.sandler@ht.com', '2021-09-17', '2021-09-21'),
    ('140', 2, 'Deluxe', 'Peter Zurich', 'p.zurich@yh.com', '2021-09-03', '2021-09-10'),
    ('141', 2, 'Standard_plus', 'Shawn Dudley', 'sh.dudley@yh.com', '2021-09-27', '2021-09-29');
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------
-- -------------------------------------------------------------------------------------------------------

