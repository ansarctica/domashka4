package postgres

import (
	"context"

	"github.com/ansarctica/domashka4/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetStudentByID(ctx context.Context, id int) *models.StudentWithGroupName {
	query := `
		SELECT s.id, s.name, s.birth_date, s.gender, s.group_id, g.name
		FROM students s
		JOIN groups g ON s.group_id = g.id
		WHERE s.id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var s models.StudentWithGroupName
	_ = row.Scan(
		&s.ID,
		&s.Name,
		&s.BirthDate,
		&s.Gender,
		&s.GroupID,
		&s.GroupName,
	)

	return &s
}

func (r *Repository) GetAllGroupSchedules(ctx context.Context) models.GroupSchedule {
	query := `
		SELECT id, group_id, subject, start_time, end_time
		FROM schedule
		ORDER BY group_id, start_time
	`

	rows, _ := r.db.Query(ctx, query)
	defer rows.Close()

	result := make(models.GroupSchedule, 0)

	for rows.Next() {
		var s models.Schedule
		_ = rows.Scan(&s.ID, &s.GroupID, &s.Subject, &s.StartTime, &s.EndTime)
		result = append(result, s)
	}

	return result
}

func (r *Repository) GetGroupScheduleByID(ctx context.Context, groupID int) models.GroupSchedule {
	query := `
		SELECT id, group_id, subject, start_time, end_time
		FROM schedule
		WHERE group_id = $1
		ORDER BY start_time
	`

	rows, _ := r.db.Query(ctx, query, groupID)
	defer rows.Close()

	result := make(models.GroupSchedule, 0)

	for rows.Next() {
		var s models.Schedule
		_ = rows.Scan(&s.ID, &s.GroupID, &s.Subject, &s.StartTime, &s.EndTime)
		result = append(result, s)
	}

	return result
}

func (r *Repository) CreateAttendance(ctx context.Context, attendance *models.Attendance) int {
	query := `
		INSERT INTO attendance (subject_id, visit_day, visited, student_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int
	_ = r.db.QueryRow(ctx, query,
		attendance.SubjectID,
		attendance.VisitDay,
		attendance.Visited,
		attendance.StudentID,
	).Scan(&id)

	return id
}

func (r *Repository) GetAttendanceBySubjectID(ctx context.Context, subjectID int) []models.Attendance {
	query := `
		SELECT id, subject_id, visit_day, visited, student_id
		FROM attendance
		WHERE subject_id = $1
	`

	rows, _ := r.db.Query(ctx, query, subjectID)
	defer rows.Close()

	result := make([]models.Attendance, 0)

	for rows.Next() {
		var a models.Attendance
		_ = rows.Scan(
			&a.ID,
			&a.SubjectID,
			&a.VisitDay,
			&a.Visited,
			&a.StudentID,
		)
		result = append(result, a)
	}

	return result
}

func (r *Repository) GetAttendanceByStudentID(ctx context.Context, studentID int) []models.Attendance {
	query := `
		SELECT id, subject_id, visit_day, visited, student_id
		FROM attendance
		WHERE student_id = $1
	`

	rows, _ := r.db.Query(ctx, query, studentID)
	defer rows.Close()

	result := make([]models.Attendance, 0)

	for rows.Next() {
		var a models.Attendance
		_ = rows.Scan(
			&a.ID,
			&a.SubjectID,
			&a.VisitDay,
			&a.Visited,
			&a.StudentID,
		)
		result = append(result, a)
	}

	return result
}
