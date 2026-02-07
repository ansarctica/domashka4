package postgres

import (
	"context"

	"github.com/ansarctica/domashka4/internal/models"
)

func (r *Repository) GetAllGroupSchedules(ctx context.Context) ([]models.Schedule, error) {
	query := `
		SELECT id, group_id, subject_name, start_time, end_time
		FROM schedule
		ORDER BY group_id, start_time
	`
	return r.scanSchedules(ctx, query)
}

func (r *Repository) GetGroupScheduleByID(ctx context.Context, groupID int) ([]models.Schedule, error) {
	query := `
		SELECT id, group_id, subject_name, start_time, end_time
		FROM schedule
		WHERE group_id = $1
		ORDER BY start_time
	`
	return r.scanSchedules(ctx, query, groupID)
}

func (r *Repository) scanSchedules(ctx context.Context, query string, args ...interface{}) ([]models.Schedule, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Schedule
	for rows.Next() {
		var s models.Schedule
		if err := rows.Scan(&s.ID, &s.GroupID, &s.Subject, &s.StartTime, &s.EndTime); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, rows.Err()
}

func (r *Repository) CreateSchedule(ctx context.Context, s *models.Schedule) (int, error) {
	query := `
		INSERT INTO schedule (group_id, subject_name, start_time, end_time)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id int
	err := r.db.QueryRow(ctx, query, s.GroupID, s.Subject, s.StartTime, s.EndTime).Scan(&id)
	return id, err
}

func (r *Repository) UpdateSchedule(ctx context.Context, s *models.Schedule) error {
	query := `
		UPDATE schedule 
		SET group_id = $1, subject_name = $2, start_time = $3, end_time = $4
		WHERE id = $5
	`
	_, err := r.db.Exec(ctx, query, s.GroupID, s.Subject, s.StartTime, s.EndTime, s.ID)
	return err
}

func (r *Repository) DeleteSchedule(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM schedule WHERE id = $1", id)
	return err
}
