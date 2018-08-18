package sql

import (
	"database/sql"

	protoHistory "github.com/asciiu/gomo/history-service/proto"
	"github.com/google/uuid"
)

// Sql functions here:
// FindUserHistory
// FindObjectHistory
// FindRecentObjectHistory
// FindObjectHistoryCount
// InsertHistory
// UpdateHistoryClickedAt
// UpdateHistorySeenAt

func FindUserHistory(db *sql.DB, userID string, page, pageSize uint32) (*protoHistory.UserHistoryPage, error) {
	history := make([]*protoHistory.History, 0)

	var total uint32
	queryTotal := `SELECT count(*) FROM history WHERE user_id = $1`
	if err := db.QueryRow(queryTotal, userID).Scan(&total); err != nil {
		return nil, err
	}

	query := `SELECT id, 
		user_id, 
		title, 
		subtitle, 
		description, 
		timestamp,
		type,
		object_id,
		clicked_at,
		seen_at
		FROM history WHERE user_id = $1 
		ORDER BY timestamp OFFSET $2 LIMIT $3`

	rows, err := db.Query(query, userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var h protoHistory.History
		var clickedAt sql.NullString
		var seenAt sql.NullString
		err := rows.Scan(&h.HistoryID,
			&h.UserID,
			&h.Title,
			&h.Subtitle,
			&h.Description,
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

	result := protoHistory.UserHistoryPage{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		History:  history,
	}

	return &result, nil
}

func FindObjectHistory(db *sql.DB, req *protoHistory.HistoryRequest) (*protoHistory.UserHistoryPage, error) {
	history := make([]*protoHistory.History, 0)

	var total uint32
	queryTotal := `SELECT count(*) FROM history WHERE object_id = $1`
	if err := db.QueryRow(queryTotal, req.ObjectID).Scan(&total); err != nil {
		return nil, err
	}

	query := `SELECT id, 
		user_id, 
		title, 
		subtitle, 
		description, 
		timestamp,
		type,
		object_id,
		clicked_at,
		seen_at
		FROM history WHERE object_id = $1 
		ORDER BY timestamp OFFSET $2 LIMIT $3`

	rows, err := db.Query(query, req.ObjectID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var h protoHistory.History
		var clickedAt sql.NullString
		var seenAt sql.NullString
		err := rows.Scan(&h.HistoryID,
			&h.UserID,
			&h.Title,
			&h.Subtitle,
			&h.Description,
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

	result := protoHistory.UserHistoryPage{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		History:  history,
	}

	return &result, nil
}

func FindRecentObjectHistory(db *sql.DB, req *protoHistory.RecentHistoryRequest) ([]*protoHistory.History, error) {
	history := make([]*protoHistory.History, 0)

	query := `SELECT id, 
		user_id, 
		title, 
		subtitle, 
		description, 
		timestamp,
		type,
		object_id,
		clicked_at,
		seen_at
		FROM history WHERE object_id = $1 
		ORDER BY timestamp DESC LIMIT $2`

	rows, err := db.Query(query, req.ObjectID, req.Count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var h protoHistory.History
		var clickedAt sql.NullString
		var seenAt sql.NullString
		err := rows.Scan(&h.HistoryID,
			&h.UserID,
			&h.Title,
			&h.Subtitle,
			&h.Description,
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

func FindObjectHistoryCount(db *sql.DB, objectID string) uint32 {
	var count uint32
	queryCount := `SELECT count(*) FROM history WHERE object_id = $1`
	err := db.QueryRow(queryCount, objectID).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

func InsertHistory(db *sql.DB, history *protoHistory.History) (*protoHistory.History, error) {
	newID := uuid.New().String()

	sqlStatement := `insert into history (
		id, 
		user_id, 
		title, 
		subtitle, 
		description, 
		timestamp, 
		type, 
		object_id) 
		values ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := db.Exec(sqlStatement,
		newID,
		history.UserID,
		history.Title,
		history.Subtitle,
		history.Description,
		history.Timestamp,
		history.Type,
		history.ObjectID)

	if err != nil {
		return nil, err
	}
	n := &protoHistory.History{
		HistoryID:   newID,
		Type:        history.Type,
		UserID:      history.UserID,
		ObjectID:    history.ObjectID,
		Title:       history.Title,
		Subtitle:    history.Subtitle,
		Description: history.Description,
		Timestamp:   history.Timestamp,
	}
	return n, nil
}

func UpdateHistoryClickedAt(db *sql.DB, historyID, timestamp string) (*protoHistory.History, error) {
	stmt := `
		UPDATE history 
		SET 
			clicked_at = $1
		WHERE
			id = $2
		RETURNING id, 
		user_id, 
		title, 
		subtitle, 
		description, 
		timestamp, 
		type, 
		object_id,
		clicked_at,
		seen_at
		`

	var h protoHistory.History
	var clickedAt sql.NullString
	var seenAt sql.NullString

	err := db.QueryRow(stmt, timestamp, historyID).Scan(
		&h.HistoryID,
		&h.UserID,
		&h.Title,
		&h.Subtitle,
		&h.Description,
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

func UpdateHistorySeenAt(db *sql.DB, historyID, timestamp string) (*protoHistory.History, error) {
	stmt := `
		UPDATE history 
		SET 
			seen_at = $1
		WHERE
			id = $2
		RETURNING id, 
		user_id, 
		title, 
		subtitle, 
		description, 
		timestamp, 
		type, 
		object_id,
		clicked_at,
		seen_at
			`

	var h protoHistory.History
	var clickedAt sql.NullString
	var seenAt sql.NullString

	err := db.QueryRow(stmt, timestamp, historyID).Scan(
		&h.HistoryID,
		&h.UserID,
		&h.Title,
		&h.Subtitle,
		&h.Description,
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
