package postgres

import (
	"context"

	"github.com/ansarctica/domashka4/internal/models"
)

func (r *Repository) GetAllGroups(ctx context.Context) ([]models.Group, error) {
	query := `SELECT id FROM groups ORDER BY id`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var g models.Group
		if err := rows.Scan(&g.ID); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}

	return groups, rows.Err()
}
