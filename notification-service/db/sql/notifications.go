package sql

import (
	"database/sql"

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

// func FindNotificationsByType(db *sql.DB, req *pb.GetUserKeysRequest) ([]*pb.Key, error) {
// 	results := make([]*pb.Key, 0)

// 	rows, err := db.Query("SELECT id, user_id, exchange_name, api_key, secret, description, status FROM user_keys WHERE user_id = $1", req.UserID)
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var k pb.Key
// 		err := rows.Scan(&k.KeyID, &k.UserID, &k.Exchange, &k.Key, &k.Secret, &k.Description, &k.Status)
// 		if err != nil {
// 			log.Fatal(err)
// 			return nil, err
// 		}
// 		results = append(results, &k)
// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, err
// 	}

// 	return results, nil
// }

func InsertNotification(db *sql.DB, note *notification.Notification) (*notification.Notification, error) {
	newID := uuid.New()

	sqlStatement := `insert into notifications (id, user_id, title, subtitle, description, 
		timestamp, notification_type, object_id) values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := db.Exec(sqlStatement, newID, note.UserID, note.Title, note.Subtitle,
		note.Description, note.Timestamp, note.NotificationType, note.ObjectID)

	if err != nil {
		return nil, err
	}
	n := &notification.Notification{
		NotificationID:   newID.String(),
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
