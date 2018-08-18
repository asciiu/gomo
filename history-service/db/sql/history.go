package sql

import (
	"database/sql"

	protoHistory "github.com/asciiu/gomo/history-service/proto"
	"github.com/google/uuid"
)

// func DeleteKey(db *sql.DB, keyID string) error {
// 	_, err := db.Exec("DELETE FROM user_keys WHERE id = $1", keyID)
// 	return err
// }

// func FindKeyByID(db *sql.DB, req *pb.GetUserKeyRequest) (*pb.Key, error) {
// 	var k pb.Key
// 	err := db.QueryRow("SELECT id, user_id, exchange_name, api_key, secret, description, status FROM user_keys WHERE id = $1", req.KeyID).
// 		Scan(&k.KeyID, &k.UserID, &k.Exchange, &k.Key, &k.Secret, &k.Description, &k.Status)

// 	if err != nil {
// 		return nil, err
// 	}
// 	return &k, nil
// }

// func FindHistoryByType(db *sql.DB, req *protoHistory.GetNotifcationsByType) (*protoHistory.UserNotificationsPage, error) {
// 	notifications := make([]*notification.Notification, 0)

// 	var total uint32
// 	queryTotal := `SELECT count(*) FROM notifications WHERE user_id = $1 and notification_type = $2`
// 	err := db.QueryRow(queryTotal, req.UserID, req.NotificationType).Scan(&total)

// 	query := `SELECT id, user_id, title, subtitle, description, timestamp, notification_type,
// 		object_id FROM notifications WHERE user_id = $1 and notification_type = $2 ORDER BY timestamp OFFSET $3 LIMIT $4`

// 	rows, err := db.Query(query, req.UserID, req.NotificationType, req.Page, req.PageSize)
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var n notification.Notification
// 		err := rows.Scan(&n.NotificationID, &n.UserID, &n.Title, &n.Subtitle, &n.Description,
// 			&n.Timestamp, &n.NotificationType, &n.ObjectID)
// 		if err != nil {
// 			log.Fatal(err)
// 			return nil, err
// 		}

// 		notifications = append(notifications, &n)
// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, err
// 	}

// 	result := notification.UserNotificationsPage{
// 		Page:          req.Page,
// 		PageSize:      req.PageSize,
// 		Total:         total,
// 		Notifications: notifications,
// 	}

// 	return &result, nil
// }

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
