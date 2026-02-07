package service

import (
	"context"
	"errors"

	"github.com/ansarctica/domashka4/internal/models"
	"github.com/ansarctica/domashka4/internal/postgres"
)

type Service struct {
	repo *postgres.Repository
}

func NewService(repo *postgres.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllStudents(ctx context.Context, filter models.StudentFilter) ([]models.Student, error) {
	return s.repo.GetAllStudents(ctx, filter)
}
func (s *Service) GetStudent(ctx context.Context, id int) (*models.Student, error) {
	return s.repo.GetStudentByID(ctx, id)
}

func (s *Service) CreateStudent(ctx context.Context, student *models.Student) (int, error) {
	return s.repo.CreateStudent(ctx, student)
}

func (s *Service) UpdateStudent(ctx context.Context, student *models.Student) error {
	return s.repo.UpdateStudent(ctx, student)
}

func (s *Service) DeleteStudent(ctx context.Context, id int) error {
	return s.repo.DeleteStudent(ctx, id)
}

func (s *Service) GetAllGroups(ctx context.Context) ([]models.Group, error) {
	return s.repo.GetAllGroups(ctx)
}

func (s *Service) GetSchedules(ctx context.Context, groupID *int) ([]models.Schedule, error) {
	if groupID != nil {
		return s.repo.GetGroupScheduleByID(ctx, *groupID)
	}
	return s.repo.GetAllGroupSchedules(ctx)
}

func (s *Service) CreateSchedule(ctx context.Context, schedule *models.Schedule) (int, error) {
	return s.repo.CreateSchedule(ctx, schedule)
}

func (s *Service) UpdateSchedule(ctx context.Context, schedule *models.Schedule) error {
	return s.repo.UpdateSchedule(ctx, schedule)
}

func (s *Service) DeleteSchedule(ctx context.Context, id int) error {
	return s.repo.DeleteSchedule(ctx, id)
}
func (s *Service) GetAttendance(ctx context.Context, studentID *int, subjectName *string) ([]models.Attendance, error) {
	if studentID != nil {
		return s.repo.GetAttendanceByStudentID(ctx, *studentID)
	}
	if subjectName != nil {
		return s.repo.GetAttendanceBySubjectName(ctx, *subjectName)
	}
	return nil, errors.New("must provide either student_id or subject_name")
}

func (s *Service) NewAttendance(ctx context.Context, attendance *models.Attendance) (int, error) {
	return s.repo.CreateAttendance(ctx, attendance)
}

func (s *Service) UpdateAttendance(ctx context.Context, attendance *models.Attendance) error {
	return s.repo.UpdateAttendance(ctx, attendance)
}

func (s *Service) DeleteAttendance(ctx context.Context, id int) error {
	return s.repo.DeleteAttendance(ctx, id)
}
func (s *Service) GetAssignments(ctx context.Context, subjectName *string) ([]models.Assignment, error) {
	return s.repo.GetAssignments(ctx, subjectName)
}

func (s *Service) NewAssignment(ctx context.Context, assignment *models.Assignment) (int, error) {
	return s.repo.CreateAssignment(ctx, assignment)
}

func (s *Service) NewGrade(ctx context.Context, grade *models.Grade) (int, error) {
	return s.repo.CreateGrade(ctx, grade)
}

func (s *Service) GetGPA(ctx context.Context, studentID int) (float64, error) {
	return s.repo.GetGPAByStudentID(ctx, studentID)
}
func (s *Service) GetRankings(ctx context.Context, groupID *int, subjectName *string) ([]models.StudentGPA, error) {

	if subjectName != nil {
		if *subjectName == "all" {
			subjectName = nil
		} else if groupID != nil {
			return s.repo.GetSubjectGPARankingByGroup(ctx, *subjectName, *groupID)
		}
	}

	if groupID != nil {
		return s.repo.GetGPARankingByGroup(ctx, *groupID)
	}
	if subjectName != nil {
		return s.repo.GetGPARankingBySubject(ctx, *subjectName)
	}

	return nil, errors.New("must provide group_id, subject_name, or both")
}

func (s *Service) GetAllSubjects(ctx context.Context) ([]models.Subject, error) {
	return s.repo.GetAllSubjects(ctx)
}
