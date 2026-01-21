package postgres

import (
	"context"

	"sort"

	"github.com/ansarctica/domashka4/internal/models"
	"github.com/jackc/pgx/v5"
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

func (r *Repository) GetAllGroupSchedules(ctx context.Context) ([]models.Schedule, error) {
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

	result := make([]models.Schedule, 0)

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

func (r *Repository) GetGroupScheduleByID(ctx context.Context, groupID int) ([]models.Schedule, error) {
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

	result := make([]models.Schedule, 0)

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

func (r *Repository) CreateAssignment(ctx context.Context, assignment *models.Assignment) (int, error) {
	query := `
        INSERT INTO assignments (name, subject_id, weight, date)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	var id int
	err := r.db.QueryRow(ctx, query,
		assignment.Name,
		assignment.SubjectID,
		assignment.Weight,
		assignment.Date,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) CreateGrade(ctx context.Context, grade *models.Grade) (int, error) {
	query := `
        INSERT INTO grades (student_id, assignment_id, mark)
        VALUES ($1, $2, $3)
        RETURNING id
    `

	var id int
	err := r.db.QueryRow(ctx, query,
		grade.StudentID,
		grade.AssignmentID,
		grade.Mark,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) GetGPAByStudentID(ctx context.Context, studentID int) (float64, error) {
	query := `
		SELECT 
			COALESCE(SUM(g.mark * a.weight), 0),
			COALESCE(SUM(a.weight), 0)
		FROM grades g
		JOIN assignments a ON g.assignment_id = a.id
		WHERE g.student_id = $1
	`

	var weightedPoints float64
	var maxWeightedPoints float64

	err := r.db.QueryRow(ctx, query, studentID).Scan(&weightedPoints, &maxWeightedPoints)
	if err != nil {
		return 0, err
	}

	if maxWeightedPoints == 0 {
		return 0, nil
	}

	gpa := weightedPoints / maxWeightedPoints
	return gpa, nil
}

func (r *Repository) GetSubjectGPA(ctx context.Context, studentID int, subjectID int) (float64, error) {
	query := `
        SELECT 
            COALESCE(SUM(g.mark * a.weight), 0),
            COALESCE(SUM(a.weight), 0)
        FROM grades g
        JOIN assignments a ON g.assignment_id = a.id
        WHERE g.student_id = $1 AND a.subject_id = $2
    `

	var weightedPoints float64
	var maxWeightedPoints float64

	err := r.db.QueryRow(ctx, query, studentID, subjectID).Scan(&weightedPoints, &maxWeightedPoints)
	if err != nil {
		return 0, err
	}

	if maxWeightedPoints == 0 {
		return 0, nil
	}

	gpa := weightedPoints / maxWeightedPoints
	return gpa, nil
}

func (r *Repository) GetGPARankingByGroup(ctx context.Context, groupID int) ([]models.StudentGPA, error) {
	query := `
        SELECT 
            g.student_id,
            COALESCE(SUM(g.mark * a.weight), 0),
            COALESCE(SUM(a.weight), 0)
        FROM grades g
        JOIN assignments a ON g.assignment_id = a.id
        JOIN students s ON g.student_id = s.id 
        WHERE s.group_id = $1
        GROUP BY g.student_id
    `

	rows, err := r.db.Query(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanAndSortStudentGPAs(rows)
}

func (r *Repository) GetGPARankingBySubject(ctx context.Context, subjectID int) ([]models.StudentGPA, error) {
	query := `
        SELECT 
            g.student_id,
            COALESCE(SUM(g.mark * a.weight), 0),
            COALESCE(SUM(a.weight), 0)
        FROM grades g
        JOIN assignments a ON g.assignment_id = a.id
        WHERE a.subject_id = $1
        GROUP BY g.student_id
    `

	rows, err := r.db.Query(ctx, query, subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanAndSortStudentGPAs(rows)
}

func (r *Repository) GetSubjectGPARankingByGroup(ctx context.Context, subjectID int, groupID int) ([]models.StudentGPA, error) {
	query := `
        SELECT 
            g.student_id,
            COALESCE(SUM(g.mark * a.weight), 0),
            COALESCE(SUM(a.weight), 0)
        FROM grades g
        JOIN assignments a ON g.assignment_id = a.id
        JOIN students s ON g.student_id = s.id
        WHERE a.subject_id = $1 AND s.group_id = $2
        GROUP BY g.student_id
    `

	rows, err := r.db.Query(ctx, query, subjectID, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanAndSortStudentGPAs(rows)
}

func scanAndSortStudentGPAs(rows pgx.Rows) ([]models.StudentGPA, error) {
	var results []models.StudentGPA

	for rows.Next() {
		var studentID int
		var weightedSum float64
		var totalWeight float64

		if err := rows.Scan(&studentID, &weightedSum, &totalWeight); err != nil {
			return nil, err
		}

		var gpa float64
		if totalWeight > 0 {
			gpa = weightedSum / totalWeight
		}

		results = append(results, models.StudentGPA{
			StudentID: studentID,
			GPA:       gpa,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].GPA > results[j].GPA
	})

	return results, nil
}
