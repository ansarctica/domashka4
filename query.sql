CREATE TABLE subjects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50)
);

CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    department VARCHAR(50),
    CONSTRAINT check_department CHECK (department IN ('Инженеры', 'Гуманитарии'))
);

CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    birth_date DATE,
    gender CHAR(1),
    group_id INT,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    CONSTRAINT check_gender CHECK (gender IN ('М', 'Ж'))
);

CREATE TABLE attendance (
    id SERIAL PRIMARY KEY,
    subject_id INT,
    visit_day DATE,
    student_id INT,
    visited BOOLEAN,
    FOREIGN KEY (subject_id) REFERENCES subjects(id),
    FOREIGN KEY (student_id) REFERENCES students(id)
);

CREATE TABLE schedule (
    id SERIAL PRIMARY KEY,
    group_id INT,
    subject_id INT,
    start_time TIME,
    end_time TIME,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (subject_id) REFERENCES subjects(id)
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255)
);

CREATE TABLE assignments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    subject_id INT REFERENCES subjects(id),
    weight INT CHECK (weight <= 100),
    date DATE,
);

CREATE TABLE grades (
    id SERIAL PRIMARY KEY,
    student_id INT REFERENCES students(id),
    assignment_id INT REFERENCES assignments(id),
    mark INT CHECK (mark >= 0 AND mark <= 100)
);

INSERT INTO subjects (name) VALUES
('Физика'),
('Калкулус'),
('Черчение'),
('Химия'),
('Физ. Культура'),
('Программирование'),
('История'),
('Английский'),
('Математика'),
('Психология'),
('Социология'),
('Литература');

INSERT INTO groups (name, department) VALUES
('И-101', 'Инженеры'),
('И-102', 'Инженеры'),
('Г-201', 'Гуманитарии'),
('Г-202', 'Гуманитарии');

INSERT INTO students (name, birth_date, gender, group_id) VALUES
('Дамир', '2005-05-15', 'М', 1),
('Динара', '2005-08-22', 'Ж', 1),
('Алексей', '2004-11-03', 'М', 1),
('Ильяс', '2005-01-10', 'М', 2),
('Сара', '2005-12-30', 'Ж', 2),
('Мухтар', '2004-09-14', 'М', 2),
('Полина', '2006-03-25', 'Ж', 3),
('Абдулла', '2005-07-07', 'М', 3),
('Айсулу', '2005-04-12', 'Ж', 3),
('Кайсар', '2004-06-18', 'М', 4),
('Ирина', '2005-10-05', 'Ж', 4),
('Марат', '2005-02-28', 'М', 4);

INSERT INTO schedule (group_id, subject_id, start_time, end_time) VALUES
(1, 1, '09:00', '10:30'),
(1, 2, '10:45', '12:15'),
(1, 3, '13:00', '14:30'),
(2, 4, '09:00', '10:30'),
(2, 5, '10:45', '12:15'),
(2, 6, '13:00', '14:30'),
(3, 7, '09:00', '10:30'),
(3, 8, '10:45', '12:15'),
(3, 9, '13:00', '14:30'),
(4, 10, '10:45', '12:15'),
(4, 11, '13:00', '14:30'),
(4, 12, '14:45', '16:15');