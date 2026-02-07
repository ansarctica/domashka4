package postgres

import (
	"context"
	"sort"

	"github.com/ansarctica/domashka4/internal/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetAssignments(ctx context.Context, subjectName *string) ([]models.Assignment, error) {
	query := `SELECT id, name, subject_name, weight, date FROM assignments`

	var args []interface{}
	if subjectName != nil {
		query += ` WHERE subject_name = $1`
		args = append(args, *subjectName)
	}

	query += ` ORDER BY date DESC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []models.Assignment
	for rows.Next() {
		var a models.Assignment
		if err := rows.Scan(&a.ID, &a.Name, &a.SubjectName, &a.Weight, &a.Date); err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}
	return assignments, rows.Err()
}

func (r *Repository) CreateAssignment(ctx context.Context, a *models.Assignment) (int, error) {
	query := `
        INSERT INTO assignments (name, subject_name, weight, date)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	var id int
	err := r.db.QueryRow(ctx, query, a.Name, a.SubjectName, a.Weight, a.Date).Scan(&id)
	return id, err
}

func (r *Repository) CreateGrade(ctx context.Context, g *models.Grade) (int, error) {
	query := `
        INSERT INTO grades (student_id, assignment_id, mark)
        VALUES ($1, $2, $3)
        RETURNING id
    `
	var id int
	err := r.db.QueryRow(ctx, query, g.StudentID, g.AssignmentID, g.Mark).Scan(&id)
	return id, err
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
	return r.calculateGPA(ctx, query, studentID)
}

func (r *Repository) GetSubjectGPA(ctx context.Context, studentID int, subjectName string) (float64, error) {
	query := `
        SELECT 
            COALESCE(SUM(g.mark * a.weight), 0),
            COALESCE(SUM(a.weight), 0)
        FROM grades g
        JOIN assignments a ON g.assignment_id = a.id
        WHERE g.student_id = $1 AND a.subject_name = $2
    `
	return r.calculateGPA(ctx, query, studentID, subjectName)
}

func (r *Repository) calculateGPA(ctx context.Context, query string, args ...interface{}) (float64, error) {
	var weightedPoints float64
	var maxWeightedPoints float64

	err := r.db.QueryRow(ctx, query, args...).Scan(&weightedPoints, &maxWeightedPoints)
	if err != nil {
		return 0, err
	}

	if maxWeightedPoints == 0 {
		return 0, nil
	}

	return weightedPoints / maxWeightedPoints, nil
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

func (r *Repository) GetGPARankingBySubject(ctx context.Context, subjectName string) ([]models.StudentGPA, error) {
	query := `
        SELECT 
            g.student_id,
            COALESCE(SUM(g.mark * a.weight), 0),
            COALESCE(SUM(a.weight), 0)
        FROM grades g
        JOIN assignments a ON g.assignment_id = a.id
        WHERE a.subject_name = $1
        GROUP BY g.student_id
    `
	rows, err := r.db.Query(ctx, query, subjectName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanAndSortStudentGPAs(rows)
}

func (r *Repository) GetSubjectGPARankingByGroup(ctx context.Context, subjectName string, groupID int) ([]models.StudentGPA, error) {
	query := `
        SELECT 
            g.student_id,
            COALESCE(SUM(g.mark * a.weight), 0),
            COALESCE(SUM(a.weight), 0)
        FROM grades g
        JOIN assignments a ON g.assignment_id = a.id
        JOIN students s ON g.student_id = s.id
        WHERE a.subject_name = $1 AND s.group_id = $2
        GROUP BY g.student_id
    `
	rows, err := r.db.Query(ctx, query, subjectName, groupID)
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
