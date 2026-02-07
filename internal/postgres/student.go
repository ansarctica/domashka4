package postgres

import (
	"context"
	"fmt"

	"github.com/ansarctica/domashka4/internal/models"
)

func (r *Repository) GetAllStudents(ctx context.Context, filter models.StudentFilter) ([]models.Student, error) {
	query := `
		SELECT id, name, birth_date, gender, group_id, major, course_year
		FROM students
		WHERE 1=1
	`
	var args []interface{}
	argID := 1

	if filter.GroupID != nil {
		query += fmt.Sprintf(" AND group_id = $%d", argID)
		args = append(args, *filter.GroupID)
		argID++
	}

	if filter.Major != nil {
		query += fmt.Sprintf(" AND major ILIKE $%d", argID)
		args = append(args, "%"+*filter.Major+"%")
		argID++
	}

	if filter.CourseYear != nil {
		query += fmt.Sprintf(" AND course_year = $%d", argID)
		args = append(args, *filter.CourseYear)
		argID++
	}

	query += " ORDER BY id"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argID)
		args = append(args, filter.Limit)
		argID++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argID)
		args = append(args, filter.Offset)
		argID++
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var s models.Student
		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.BirthDate,
			&s.Gender,
			&s.GroupID,
			&s.Major,
			&s.CourseYear,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, s)
	}

	return students, rows.Err()
}

func (r *Repository) GetStudentByID(ctx context.Context, id int) (*models.Student, error) {
	query := `
		SELECT id, name, birth_date, gender, group_id, major, course_year
		FROM students
		WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)

	var s models.Student
	err := row.Scan(
		&s.ID,
		&s.Name,
		&s.BirthDate,
		&s.Gender,
		&s.GroupID,
		&s.Major,
		&s.CourseYear,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) CreateStudent(ctx context.Context, s *models.Student) (int, error) {
	query := `
		INSERT INTO students (name, birth_date, gender, group_id, major, course_year)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var id int
	err := r.db.QueryRow(ctx, query,
		s.Name, s.BirthDate, s.Gender, s.GroupID, s.Major, s.CourseYear,
	).Scan(&id)
	return id, err
}

func (r *Repository) UpdateStudent(ctx context.Context, s *models.Student) error {
	query := `
		UPDATE students 
		SET name = $1, birth_date = $2, gender = $3, group_id = $4, major = $5, course_year = $6
		WHERE id = $7
	`
	_, err := r.db.Exec(ctx, query,
		s.Name, s.BirthDate, s.Gender, s.GroupID, s.Major, s.CourseYear, s.ID,
	)
	return err
}

func (r *Repository) DeleteStudent(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM students WHERE id = $1", id)
	return err
}
