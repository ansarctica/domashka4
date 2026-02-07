package postgres

import (
	"context"

	"github.com/ansarctica/domashka4/internal/models"
)

func (r *Repository) GetAllSubjects(ctx context.Context) ([]models.Subject, error) {
	query := `SELECT name FROM subjects ORDER BY name`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		var s models.Subject
		if err := rows.Scan(&s.Name); err != nil {
			return nil, err
		}
		subjects = append(subjects, s)
	}

	return subjects, rows.Err()
}
