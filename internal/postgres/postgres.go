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

func (r *Repository) GetStudentByID(ctx context.Context, id int) (*models.StudentWithGroupName, error) {
	query := `
		SELECT s.id, s.name, s.birth_date, s.gender, s.group_id, g.name
		FROM students s
		JOIN groups g ON s.group_id = g.id
		WHERE s.id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var s models.StudentWithGroupName
	err := row.Scan(
		&s.ID,
		&s.Name,
		&s.BirthDate,
		&s.Gender,
		&s.GroupID,
		&s.GroupName,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *Repository) GetAllGroupSchedules(ctx context.Context) (models.GroupSchedule, error) {
	query := `
		SELECT id, group_id, subject, start_time, end_time
		FROM schedule
		ORDER BY group_id, start_time
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(models.GroupSchedule, 0)

	for rows.Next() {
		var s models.Schedule
		if err := rows.Scan(&s.ID, &s.GroupID, &s.Subject, &s.StartTime, &s.EndTime); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repository) GetGroupScheduleByID(ctx context.Context, groupID int) (models.GroupSchedule, error) {
	query := `
		SELECT id, group_id, subject, start_time, end_time
		FROM schedule
		WHERE group_id = $1
		ORDER BY start_time
	`

	rows, err := r.db.Query(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(models.GroupSchedule, 0)

	for rows.Next() {
		var s models.Schedule
		if err := rows.Scan(&s.ID, &s.GroupID, &s.Subject, &s.StartTime, &s.EndTime); err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repository) CreateAttendance(ctx context.Context, attendance *models.Attendance) (int, error) {
	query := `
		INSERT INTO attendance (subject_id, visit_day, visited, student_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int
	err := r.db.QueryRow(ctx, query,
		attendance.SubjectID,
		attendance.VisitDay,
		attendance.Visited,
		attendance.StudentID,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) GetAttendanceBySubjectID(ctx context.Context, subjectID int) ([]models.Attendance, error) {
	query := `
		SELECT id, subject_id, visit_day, visited, student_id
		FROM attendance
		WHERE subject_id = $1
	`

	rows, err := r.db.Query(ctx, query, subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.Attendance, 0)

	for rows.Next() {
		var a models.Attendance
		if err := rows.Scan(&a.ID, &a.SubjectID, &a.VisitDay, &a.Visited, &a.StudentID); err != nil {
			return nil, err
		}
		result = append(result, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repository) GetAttendanceByStudentID(ctx context.Context, studentID int) ([]models.Attendance, error) {
	query := `
		SELECT id, subject_id, visit_day, visited, student_id
		FROM attendance
		WHERE student_id = $1
	`

	rows, err := r.db.Query(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.Attendance, 0)

	for rows.Next() {
		var a models.Attendance
		if err := rows.Scan(&a.ID, &a.SubjectID, &a.VisitDay, &a.Visited, &a.StudentID); err != nil {
			return nil, err
		}
		result = append(result, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
