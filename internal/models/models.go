package models

import (
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Assignment struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	SubjectName string    `json:"subject_name"`
	Weight      int       `json:"weight"`
	Date        time.Time `json:"date"`
}

type Grade struct {
	ID           int `json:"id"`
	StudentID    int `json:"student_id"`
	AssignmentID int `json:"assignment_id"`
	Mark         int `json:"mark"`
}

type StudentGPA struct {
	StudentID int     `json:"student_id"`
	GPA       float64 `json:"gpa"`
}

type Schedule struct {
	ID        int       `json:"id"`
	GroupID   int       `json:"group_id"`
	Subject   string    `json:"subject"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type Group struct {
	ID int `json:"id"`
}

type Student struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	BirthDate  time.Time `json:"birth_date"`
	Gender     string    `json:"gender"`
	GroupID    int       `json:"group_id"`
	Major      string    `json:"major"`
	CourseYear int       `json:"course_year"`
}

type Attendance struct {
	ID          int       `json:"id"`
	SubjectName string    `json:"subject_name"`
	VisitDay    time.Time `json:"visit_day"`
	Visited     bool      `json:"visited"`
	StudentID   int       `json:"student_id"`
}

type StudentFilter struct {
	GroupID    *int
	Major      *string
	CourseYear *int
	Limit      int
	Offset     int
}

type Subject struct {
	Name string `json:"name"`
}
