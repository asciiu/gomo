package sql

import (
	"database/sql"
	"encoding/json"

	protoActivity "github.com/asciiu/gomo/activity-bulletin/proto"
	"github.com/google/uuid"
)

// Sql functions here:
// FindUserActivity
// FindObjectActivity
// FindRecentObjectActivity
// FindObjectActivityCount
// InsertActivity
// UpdateActivityClickedAt
// UpdateActivitySeenAt

func FindUserActivity(db *sql.DB, userID string, page, pageSize uint32) (*protoActivity.UserActivityPage, error) {
	history := make([]*protoActivity.Activity, 0)

	var total uint32
	queryTotal := `SELECT count(*) FROM activity_bulletin WHERE user_id = $1`
	if err := db.QueryRow(queryTotal, userID).Scan(&total); err != nil {
		return nil, err
	}

	query := `SELECT id, 
		user_id, 
		title, 
		subtitle, 
		description, 
		details,
		timestamp,
		type,
		object_id,
		clicked_at,
		seen_at
		FROM activity_bulletin WHERE user_id = $1 
		ORDER BY timestamp OFFSET $2 LIMIT $3`

	rows, err := db.Query(query, userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var h protoActivity.Activity
		var clickedAt sql.NullString
		var seenAt sql.NullString
		err := rows.Scan(&h.ActivityID,
			&h.UserID,
			&h.Title,
			&h.Subtitle,
			&h.Description,
			&h.Details,
			&h.Timestamp,
			&h.Type,
			&h.ObjectID,
			&clickedAt,
			&seenAt,
		)

		if err != nil {
			return nil, err
		}
		if clickedAt.Valid {
			h.ClickedAt = clickedAt.String
		}
		if seenAt.Valid {
			h.SeenAt = seenAt.String
		}

		history = append(history, &h)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	result := protoActivity.UserActivityPage{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Activity: history,
	}

	return &result, nil
}

func FindObjectActivity(db *sql.DB, req *protoActivity.ActivityRequest) (*protoActivity.UserActivityPage, error) {
	history := make([]*protoActivity.Activity, 0)

	var total uint32
	queryTotal := `SELECT count(*) FROM activity_bulletin WHERE object_id = $1`
	if err := db.QueryRow(queryTotal, req.ObjectID).Scan(&total); err != nil {
		return nil, err
	}

	query := `SELECT id, 
		user_id, 
		title, 
		subtitle, 
		description, 
		details,
		timestamp,
		type,
		object_id,
		clicked_at,
		seen_at
		FROM activity_bulletin WHERE object_id = $1 
		ORDER BY timestamp OFFSET $2 LIMIT $3`

	rows, err := db.Query(query, req.ObjectID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var h protoActivity.Activity
		var clickedAt sql.NullString
		var seenAt sql.NullString
		err := rows.Scan(&h.ActivityID,
			&h.UserID,
			&h.Title,
			&h.Subtitle,
			&h.Description,
			&h.Details,
			&h.Timestamp,
			&h.Type,
			&h.ObjectID,
			&clickedAt,
			&seenAt,
		)

		if err != nil {
			return nil, err
		}
		if clickedAt.Valid {
			h.ClickedAt = clickedAt.String
		}
		if seenAt.Valid {
			h.SeenAt = seenAt.String
		}

		history = append(history, &h)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	result := protoActivity.UserActivityPage{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		Activity: history,
	}

	return &result, nil
}

func FindRecentObjectActivity(db *sql.DB, req *protoActivity.RecentActivityRequest) ([]*protoActivity.Activity, error) {
	history := make([]*protoActivity.Activity, 0)

	query := `SELECT id, 
		user_id, 
		title, 
		subtitle, 
		description, 
		details,
		timestamp,
		type,
		object_id,
		clicked_at,
		seen_at
		FROM activity_bulletin WHERE object_id = $1 
		ORDER BY timestamp DESC LIMIT $2`

	rows, err := db.Query(query, req.ObjectID, req.Count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var h protoActivity.Activity
		var clickedAt sql.NullString
		var seenAt sql.NullString
		err := rows.Scan(&h.ActivityID,
			&h.UserID,
			&h.Title,
			&h.Subtitle,
			&h.Description,
			&h.Details,
			&h.Timestamp,
			&h.Type,
			&h.ObjectID,
			&clickedAt,
			&seenAt,
		)

		if err != nil {
			return nil, err
		}
		if clickedAt.Valid {
			h.ClickedAt = clickedAt.String
		}
		if seenAt.Valid {
			h.SeenAt = seenAt.String
		}

		history = append(history, &h)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return history, nil
}

func FindObjectActivityCount(db *sql.DB, objectID string) uint32 {
	var count uint32
	queryCount := `SELECT count(*) FROM activity_bulletin WHERE object_id = $1`
	err := db.QueryRow(queryCount, objectID).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

func InsertActivity(db *sql.DB, history *protoActivity.Activity) (*protoActivity.Activity, error) {
	newID := uuid.New().String()

	sqlStatement := `insert into activity_bulletin (
		id, 
		user_id, 
		title, 
		subtitle, 
		description, 
		details,
		timestamp, 
		type, 
		object_id) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	str, err := json.Marshal(history.Details)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(sqlStatement,
		newID,
		history.UserID,
		history.Title,
		history.Subtitle,
		history.Description,
		str,
		history.Timestamp,
		history.Type,
		history.ObjectID)

	if err != nil {
		return nil, err
	}
	n := &protoActivity.Activity{
		ActivityID:  newID,
		Type:        history.Type,
		UserID:      history.UserID,
		ObjectID:    history.ObjectID,
		Title:       history.Title,
		Subtitle:    history.Subtitle,
		Description: history.Description,
		Details:     history.Details,
		Timestamp:   history.Timestamp,
	}
	return n, nil
}

func UpdateActivityClickedAt(db *sql.DB, historyID, timestamp string) (*protoActivity.Activity, error) {
	stmt := `
		UPDATE activity_bulletin 
		SET 
			clicked_at = $1
		WHERE
			id = $2
		RETURNING id, 
		user_id, 
		title, 
		subtitle, 
		description, 
		details,
		timestamp, 
		type, 
		object_id,
		clicked_at,
		seen_at
		`

	var h protoActivity.Activity
	var clickedAt sql.NullString
	var seenAt sql.NullString

	err := db.QueryRow(stmt, timestamp, historyID).Scan(
		&h.ActivityID,
		&h.UserID,
		&h.Title,
		&h.Subtitle,
		&h.Description,
		&h.Details,
		&h.Timestamp,
		&h.Type,
		&h.ObjectID,
		&clickedAt,
		&seenAt)

	if err != nil {
		return nil, err
	}
	if clickedAt.Valid {
		h.ClickedAt = clickedAt.String
	}
	if seenAt.Valid {
		h.SeenAt = seenAt.String
	}

	return &h, nil
}

func UpdateActivitySeenAt(db *sql.DB, historyID, timestamp string) (*protoActivity.Activity, error) {
	stmt := `
		UPDATE activity_bulletin 
		SET 
			seen_at = $1
		WHERE
			id = $2
		RETURNING id, 
		user_id, 
		title, 
		subtitle, 
		description, 
		details,
		timestamp, 
		type, 
		object_id,
		clicked_at,
		seen_at
			`

	var h protoActivity.Activity
	var clickedAt sql.NullString
	var seenAt sql.NullString

	err := db.QueryRow(stmt, timestamp, historyID).Scan(
		&h.ActivityID,
		&h.UserID,
		&h.Title,
		&h.Subtitle,
		&h.Description,
		&h.Details,
		&h.Timestamp,
		&h.Type,
		&h.ObjectID,
		&clickedAt,
		&seenAt)

	if err != nil {
		return nil, err
	}
	if clickedAt.Valid {
		h.ClickedAt = clickedAt.String
	}
	if seenAt.Valid {
		h.SeenAt = seenAt.String
	}

	return &h, nil

	return nil, err
}
