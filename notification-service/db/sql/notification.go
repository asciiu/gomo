package sql

import (
	"database/sql"
	"encoding/json"

	protoNotification "github.com/asciiu/gomo/notification-service/proto/notification"
	"github.com/google/uuid"
)

// Sql functions here:
// FindActivity
// FindUserActivity
// FindObjectActivity
// FindRecentObjectActivity
// FindObjectActivityCount
// InsertActivity
// UpdateActivityClickedAt
// UpdateActivitySeenAt

func FindActivity(db *sql.DB, activityID string) (*protoNotification.Activity, error) {

	query := `SELECT 
		id, 
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
		FROM activity_bulletin WHERE id = $1`

	var h protoNotification.Activity
	var clickedAt sql.NullString
	var seenAt sql.NullString
	err := db.QueryRow(query, activityID).Scan(
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

	return &h, nil
}

func FindUserActivity(db *sql.DB, userID string, page, pageSize uint32) (*protoNotification.UserActivityPage, error) {
	history := make([]*protoNotification.Activity, 0)

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
		var h protoNotification.Activity
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

	result := protoNotification.UserActivityPage{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Activity: history,
	}

	return &result, nil
}

func FindUserPlansActivity(db *sql.DB, userID string, page, pageSize uint32) (*protoNotification.UserActivityPage, error) {
	history := make([]*protoNotification.Activity, 0)

	var total uint32
	queryTotal := `SELECT count(*) FROM activity_bulletin WHERE user_id = $1`
	if err := db.QueryRow(queryTotal, userID).Scan(&total); err != nil {
		return nil, err
	}

	query := `SELECT 
		ab.id, 
		ab.user_id, 
		p.title, 
		ab.subtitle, 
		ab.description, 
		ab.details,
		ab.timestamp,
		ab.type,
		ab.object_id,
		ab.clicked_at,
		ab.seen_at
		FROM activity_bulletin ab
		JOIN plans p ON p.id = ab.object_id 
		WHERE ab.user_id = $1 
		ORDER BY timestamp OFFSET $2 LIMIT $3`

	rows, err := db.Query(query, userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var h protoNotification.Activity
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

	result := protoNotification.UserActivityPage{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Activity: history,
	}

	return &result, nil
}

func FindObjectActivity(db *sql.DB, req *protoNotification.ActivityRequest) (*protoNotification.UserActivityPage, error) {
	history := make([]*protoNotification.Activity, 0)

	var total uint32
	queryTotal := `SELECT count(*) FROM activity_bulletin WHERE object_id = $1`
	if err := db.QueryRow(queryTotal, req.ObjectID).Scan(&total); err != nil {
		return nil, err
	}

	query := `SELECT 
		ab.id, 
		ab.user_id, 
		p.title, 
		ab.subtitle, 
		ab.description, 
		ab.details,
		ab.timestamp,
		ab.type,
		ab.object_id,
		ab.clicked_at,
		ab.seen_at
		FROM activity_bulletin ab
		JOIN plans p ON p.id = ab.object_id
		WHERE object_id = $1 
		ORDER BY timestamp OFFSET $2 LIMIT $3`

	rows, err := db.Query(query, req.ObjectID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var h protoNotification.Activity
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

	result := protoNotification.UserActivityPage{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		Activity: history,
	}

	return &result, nil
}

func FindRecentObjectActivity(db *sql.DB, req *protoNotification.RecentActivityRequest) ([]*protoNotification.Activity, error) {
	history := make([]*protoNotification.Activity, 0)

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
		var h protoNotification.Activity
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

func InsertActivity(db *sql.DB, history *protoNotification.Activity) (*protoNotification.Activity, error) {
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
	n := &protoNotification.Activity{
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

func UpdateActivityClickedAt(db *sql.DB, activityID, timestamp string) (*protoNotification.Activity, error) {
	stmt := `
		UPDATE activity_bulletin 
		SET 
			clicked_at = $1
		WHERE
			id = $2
		RETURNING 
		id, 
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

	var h protoNotification.Activity
	var clickedAt sql.NullString
	var seenAt sql.NullString

	err := db.QueryRow(stmt, timestamp, activityID).Scan(
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

func UpdateActivitySeenAt(db *sql.DB, activityID, timestamp string) (*protoNotification.Activity, error) {
	stmt := `
		UPDATE activity_bulletin 
		SET 
			seen_at = $1
		WHERE
			id = $2
		RETURNING 
		id, 
		user_id, 
		title, 
		subtitle, 
		description, 
		details,
		timestamp, 
		type, 
		object_id,
		clicked_at,
		seen_at`

	var h protoNotification.Activity
	var clickedAt sql.NullString
	var seenAt sql.NullString

	err := db.QueryRow(stmt, timestamp, activityID).Scan(
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
