package service

import (
	"context"

	"github.com/ansarctica/domashka4/internal/models"
	"github.com/ansarctica/domashka4/internal/postgres"
)

type Service struct {
	repo *postgres.Repository
}

func NewService(repo *postgres.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetStudent(ctx context.Context, id int) (*models.StudentWithGroupName, error) {
	return s.repo.GetStudentByID(ctx, id)
}

func (s *Service) GetAllSchedules(ctx context.Context) ([]models.Schedule, error) {
	return s.repo.GetAllGroupSchedules(ctx)
}

func (s *Service) GetGroupSchedule(ctx context.Context, id int) ([]models.Schedule, error) {
	return s.repo.GetGroupScheduleByID(ctx, id)
}

func (s *Service) NewAttendance(ctx context.Context, attendance *models.Attendance) (int, error) {
	return s.repo.CreateAttendance(ctx, attendance)
}

func (s *Service) AttendanceBySubject(ctx context.Context, id int) ([]models.Attendance, error) {
	return s.repo.GetAttendanceBySubjectID(ctx, id)
}

func (s *Service) AttendanceByStudent(ctx context.Context, id int) ([]models.Attendance, error) {
	return s.repo.GetAttendanceByStudentID(ctx, id)
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

func (s *Service) GetSubjectGPA(ctx context.Context, studentID, subjectID int) (float64, error) {
	return s.repo.GetSubjectGPA(ctx, studentID, subjectID)
}

func (s *Service) RankingByGroup(ctx context.Context, groupID int) ([]models.StudentGPA, error) {
	return s.repo.GetGPARankingByGroup(ctx, groupID)
}

func (s *Service) RankingBySubject(ctx context.Context, subjectID int) ([]models.StudentGPA, error) {
	return s.repo.GetGPARankingBySubject(ctx, subjectID)
}

func (s *Service) RankingByGroupSubject(ctx context.Context, groupID, subjectID int) ([]models.StudentGPA, error) {
	return s.repo.GetSubjectGPARankingByGroup(ctx, subjectID, groupID)
}
