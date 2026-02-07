package postgres

import (
	"context"

	"github.com/ansarctica/domashka4/internal/models"
)

func (r *Repository) CreateAttendance(ctx context.Context, a *models.Attendance) (int, error) {
	query := `
		INSERT INTO attendance (subject_name, visit_day, visited, student_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id int
	err := r.db.QueryRow(ctx, query,
		a.SubjectName, a.VisitDay, a.Visited, a.StudentID,
	).Scan(&id)
	return id, err
}
func (r *Repository) GetAttendanceBySubjectName(ctx context.Context, subjectName string) ([]models.Attendance, error) {
	query := `
		SELECT id, subject_name, visit_day, visited, student_id
		FROM attendance
		WHERE subject_name = $1
	`
	return r.scanAttendance(ctx, query, subjectName)
}

func (r *Repository) GetAttendanceByStudentID(ctx context.Context, studentID int) ([]models.Attendance, error) {
	query := `
		SELECT id, subject_name, visit_day, visited, student_id
		FROM attendance
		WHERE student_id = $1
	`
	return r.scanAttendance(ctx, query, studentID)
}

func (r *Repository) scanAttendance(ctx context.Context, query string, args ...interface{}) ([]models.Attendance, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Attendance
	for rows.Next() {
		var a models.Attendance
		if err := rows.Scan(&a.ID, &a.SubjectName, &a.VisitDay, &a.Visited, &a.StudentID); err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, rows.Err()
}

func (r *Repository) UpdateAttendance(ctx context.Context, a *models.Attendance) error {
	query := `
		UPDATE attendance 
		SET subject_name = $1, visit_day = $2, visited = $3, student_id = $4
		WHERE id = $5
	`
	_, err := r.db.Exec(ctx, query,
		a.SubjectName, a.VisitDay, a.Visited, a.StudentID, a.ID,
	)
	return err
}

func (r *Repository) DeleteAttendance(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM attendance WHERE id = $1", id)
	return err
}
