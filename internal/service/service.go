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

func (s *Service) GetStudent(ctx context.Context, id int) *models.StudentWithGroupName {
	return s.repo.GetStudentByID(ctx, id)
}

func (s *Service) GetAllSchedules(ctx context.Context) models.GroupSchedule {
	return s.repo.GetAllGroupSchedules(ctx)
}

func (s *Service) GetGroupSchedule(ctx context.Context, id int) models.GroupSchedule {
	return s.repo.GetGroupScheduleByID(ctx, id)
}

func (s *Service) NewAttendance(ctx context.Context, attendance *models.Attendance) int {
	return s.repo.CreateAttendance(ctx, attendance)
}

func (s *Service) AttendanceBySubject(ctx context.Context, id int) []models.Attendance {
	return s.repo.GetAttendanceBySubjectID(ctx, id)
}

func (s *Service) AttendanceByStudent(ctx context.Context, id int) []models.Attendance {
	return s.repo.GetAttendanceByStudentID(ctx, id)
}
