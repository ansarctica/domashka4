CREATE TABLE subjects (
    name VARCHAR(50) PRIMARY KEY
);

CREATE TABLE groups (
    id SERIAL PRIMARY KEY
);

CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    birth_date DATE,
    gender VARCHAR(10) CHECK (gender IN ('M', 'F')),
    
    group_id INT REFERENCES groups(id),
    
    major VARCHAR(100),
    course_year INT
);

CREATE TABLE attendance (
    id SERIAL PRIMARY KEY,
    student_id INT REFERENCES students(id) ON DELETE CASCADE,
    
    subject_name VARCHAR(50) REFERENCES subjects(name),
    
    visit_day DATE,
    visited BOOLEAN
);

CREATE TABLE schedule (
    id SERIAL PRIMARY KEY,
    group_id INT REFERENCES groups(id),
    
    subject_name VARCHAR(50) REFERENCES subjects(name),
    
    start_time TIME,
    end_time TIME
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255)
);

CREATE TABLE assignments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    
    subject_name VARCHAR(50) REFERENCES subjects(name),
    
    weight INT CHECK (weight <= 100),
    date DATE
);

CREATE TABLE grades (
    id SERIAL PRIMARY KEY,
    student_id INT REFERENCES students(id),
    assignment_id INT REFERENCES assignments(id),
    mark INT CHECK (mark >= 0 AND mark <= 100)
);

INSERT INTO subjects (name) VALUES
('Physics'), ('Calculus'), ('Drawing'), ('Chemistry'), 
('Physical Education'), ('Programming'), ('History'), 
('English'), ('Mathematics'), ('Psychology'), 
('Sociology'), ('Literature');

INSERT INTO groups (id) VALUES (DEFAULT), (DEFAULT), (DEFAULT), (DEFAULT);

INSERT INTO students (name, birth_date, gender, group_id, major, course_year) VALUES
('Damir', '2005-05-15', 'M', 1, 'Computer Science', 1),
('Dinara', '2005-08-22', 'F', 1, 'Computer Science', 1),
('Alexey', '2004-11-03', 'M', 1, 'Computer Science', 1),
('Ilyas', '2005-01-10', 'M', 2, 'Mechanical Engineering', 2),
('Sarah', '2005-12-30', 'F', 2, 'Mechanical Engineering', 2),
('Mukhtar', '2004-09-14', 'M', 2, 'Mechanical Engineering', 2),
('Polina', '2006-03-25', 'F', 3, 'Psychology', 1),
('Abdulla', '2005-07-07', 'M', 3, 'Psychology', 1),
('Aisulu', '2005-04-12', 'F', 3, 'Psychology', 1),
('Kaisar', '2004-06-18', 'M', 4, 'Sociology', 3),
('Irina', '2005-10-05', 'F', 4, 'Sociology', 3),
('Marat', '2005-02-28', 'M', 4, 'Sociology', 3);

INSERT INTO schedule (group_id, subject_name, start_time, end_time) VALUES
(1, 'Physics', '09:00', '10:30'),
(1, 'Calculus', '10:45', '12:15'),
(2, 'Chemistry', '09:00', '10:30'),
(3, 'History', '09:00', '10:30'),
(4, 'Psychology', '10:45', '12:15');