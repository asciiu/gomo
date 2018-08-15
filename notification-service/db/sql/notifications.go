package sql

import (
	"database/sql"
	"log"

	notification "github.com/asciiu/gomo/notification-service/proto"
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

func FindNotificationsByType(db *sql.DB, req *notification.GetNotifcationsByType) (*notification.UserNotificationsPage, error) {
	notifications := make([]*notification.Notification, 0)

	var total uint32
	queryTotal := `SELECT count(*) FROM notifications WHERE user_id = $1 and notification_type = $2`
	err := db.QueryRow(queryTotal, req.UserID, req.NotificationType).Scan(&total)

	query := `SELECT id, user_id, title, subtitle, description, timestamp, notification_type, 
		object_id FROM notifications WHERE user_id = $1 and notification_type = $2 ORDER BY timestamp OFFSET $3 LIMIT $4`

	rows, err := db.Query(query, req.UserID, req.NotificationType, req.Page, req.PageSize)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var n notification.Notification
		err := rows.Scan(&n.NotificationID, &n.UserID, &n.Title, &n.Subtitle, &n.Description,
			&n.Timestamp, &n.NotificationType, &n.ObjectID)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		notifications = append(notifications, &n)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	result := notification.UserNotificationsPage{
		Page:          req.Page,
		PageSize:      req.PageSize,
		Total:         total,
		Notifications: notifications,
	}

	return &result, nil
}

func FindNotifications(db *sql.DB, req *notification.GetNotifcationsByType) (*notification.UserNotificationsPage, error) {
	notifications := make([]*notification.Notification, 0)

	var total uint32
	queryTotal := `SELECT count(*) FROM notifications WHERE user_id = $1`
	err := db.QueryRow(queryTotal, req.UserID).Scan(&total)

	query := `SELECT id, user_id, title, subtitle, description, timestamp, notification_type, 
		object_id FROM notifications WHERE user_id = $1 ORDER BY timestamp OFFSET $2 LIMIT $3`

	rows, err := db.Query(query, req.UserID, req.Page, req.PageSize)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var n notification.Notification
		err := rows.Scan(&n.NotificationID, &n.UserID, &n.Title, &n.Subtitle, &n.Description,
			&n.Timestamp, &n.NotificationType, &n.ObjectID)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		notifications = append(notifications, &n)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	result := notification.UserNotificationsPage{
		Page:          req.Page,
		PageSize:      req.PageSize,
		Total:         total,
		Notifications: notifications,
	}

	return &result, nil
}

func InsertNotification(db *sql.DB, note *notification.Notification) (*notification.Notification, error) {
	newID := uuid.New().String()

	sqlStatement := `insert into notifications (
		id, 
		user_id, 
		title, 
		subtitle, 
		description, 
		timestamp, 
		notification_type, 
		object_id) 
		values ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := db.Exec(sqlStatement,
		newID,
		note.UserID,
		note.Title,
		note.Subtitle,
		note.Description,
		note.Timestamp,
		note.NotificationType,
		note.ObjectID)

	if err != nil {
		return nil, err
	}
	n := &notification.Notification{
		NotificationID:   newID,
		NotificationType: note.NotificationType,
		UserID:           note.UserID,
		ObjectID:         note.ObjectID,
		Title:            note.Title,
		Subtitle:         note.Subtitle,
		Description:      note.Description,
		Timestamp:        note.Timestamp,
	}
	return n, nil
}
