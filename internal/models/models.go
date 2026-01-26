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
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	SubjectID int       `json:"subject_id"`
	Weight    int       `json:"weight"`
	Date      time.Time `json:"date"`
}

type Grade struct {
	ID           int `json:"id"`
	StudentID    int `json:"student_id"`
	AssignmentID int `json:"assignment_id"`
	Mark         int `json:"mark"`
}

type StudentGPA struct {
	StudentID int
	GPA       float64
}

type Schedule struct {
	ID        int       `json:"id"`
	GroupID   int       `json:"group_id"`
	Subject   string    `json:"subject"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type Group struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Department string `json:"department"`
}

type Student struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	BirthDate time.Time `json:"birth_date"`
	Gender    string    `json:"gender"`
	GroupID   int       `json:"group_id"`
}

type StudentWithGroupName struct {
	Student
	GroupName string `json:"group_name"`
}

type Attendance struct {
	ID        int       `json:"id"`
	SubjectID int       `json:"subject_id"`
	VisitDay  time.Time `json:"visit_day"`
	Visited   bool      `json:"visited"`
	StudentID int       `json:"student_id"`
}
