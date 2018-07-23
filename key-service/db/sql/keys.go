package sql

import (
	"database/sql"

	pb "github.com/asciiu/gomo/key-service/proto/key"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func DeleteKey(db *sql.DB, keyID string) error {
	_, err := db.Exec("DELETE FROM user_keys WHERE id = $1", keyID)
	return err
}

func FindKeys(db *sql.DB, req *pb.GetKeysRequest) ([]*pb.Key, error) {
	results := make([]*pb.Key, 0)

	rows, err := db.Query(`SELECT 
		id, 
		user_id, 
		exchange_name, 
		api_key, 
		secret, 
		description, 
		status FROM user_keys WHERE id in $1`, pq.Array(req.KeyIDs))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var k pb.Key
		err := rows.Scan(&k.KeyID,
			&k.UserID,
			&k.Exchange,
			&k.Key,
			&k.Secret,
			&k.Description,
			&k.Status)

		if err != nil {
			return nil, err
		}
		results = append(results, &k)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil

}

func FindKeyByID(db *sql.DB, req *pb.GetUserKeyRequest) (*pb.Key, error) {
	var k pb.Key
	err := db.QueryRow("SELECT id, user_id, exchange_name, api_key, secret, description, status FROM user_keys WHERE id = $1", req.KeyID).
		Scan(&k.KeyID, &k.UserID, &k.Exchange, &k.Key, &k.Secret, &k.Description, &k.Status)

	if err != nil {
		return nil, err
	}
	return &k, nil
}

func FindKeysByUserID(db *sql.DB, req *pb.GetUserKeysRequest) ([]*pb.Key, error) {
	results := make([]*pb.Key, 0)

	rows, err := db.Query("SELECT id, user_id, exchange_name, api_key, secret, description, status FROM user_keys WHERE user_id = $1", req.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var k pb.Key
		err := rows.Scan(&k.KeyID, &k.UserID, &k.Exchange, &k.Key, &k.Secret, &k.Description, &k.Status)
		if err != nil {
			return nil, err
		}
		results = append(results, &k)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}

func InsertKey(db *sql.DB, req *pb.KeyRequest) (*pb.Key, error) {
	newID := uuid.New()

	sqlStatement := `insert into user_keys (id, user_id, exchange_name, api_key, secret, description, status) values ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(sqlStatement, newID, req.UserID, req.Exchange, req.Key, req.Secret, req.Description, "unverified")

	if err != nil {
		return nil, err
	}
	apikey := &pb.Key{
		KeyID:       newID.String(),
		UserID:      req.UserID,
		Exchange:    req.Exchange,
		Key:         req.Key,
		Secret:      req.Secret,
		Description: req.Description,
		Status:      "unverified",
	}
	return apikey, nil
}

func UpdateKeyDescription(db *sql.DB, req *pb.KeyRequest) (*pb.Key, error) {
	sqlStatement := `UPDATE user_keys SET description = $1 WHERE id = $2 and user_id = $3 RETURNING exchange_name, api_key, description, status`

	var k pb.Key
	err := db.QueryRow(sqlStatement, req.Description, req.KeyID, req.UserID).
		Scan(&k.Exchange, &k.Key, &k.Description, &k.Status)

	if err != nil {
		return nil, err
	}
	apikey := &pb.Key{
		KeyID:       req.KeyID,
		UserID:      req.UserID,
		Exchange:    k.Exchange,
		Key:         k.Key,
		Description: k.Description,
		Status:      k.Status,
	}
	return apikey, nil
}

func UpdateKeyStatus(db *sql.DB, req *pb.Key) (*pb.Key, error) {
	sqlStatement := `UPDATE user_keys SET status = $1 WHERE id = $2 and user_id = $3 RETURNING exchange_name, api_key, description, status`

	var k pb.Key
	err := db.QueryRow(sqlStatement, req.Status, req.KeyID, req.UserID).
		Scan(&k.Exchange, &k.Key, &k.Description, &k.Status)

	if err != nil {
		return nil, err
	}
	apikey := &pb.Key{
		KeyID:       req.KeyID,
		UserID:      req.UserID,
		Exchange:    k.Exchange,
		Key:         k.Key,
		Description: k.Description,
		Status:      k.Status,
	}
	return apikey, nil
}
