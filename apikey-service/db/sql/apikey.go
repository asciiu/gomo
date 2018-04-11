package sql

import (
	"database/sql"
	"log"

	pb "github.com/asciiu/gomo/apikey-service/proto/apikey"
	"github.com/google/uuid"
)

func DeleteApiKey(db *sql.DB, apiKeyId string) error {
	_, err := db.Exec("DELETE FROM user_keys WHERE id = $1", apiKeyId)
	return err
}

func FindApiKeyById(db *sql.DB, req *pb.GetUserApiKeyRequest) (*pb.ApiKey, error) {
	var k pb.ApiKey
	err := db.QueryRow("SELECT id, user_id, exchange_name, api_key, secret, description, status FROM user_keys WHERE id = $1", req.ApiKeyId).
		Scan(&k.ApiKeyId, &k.UserId, &k.Exchange, &k.Key, &k.Secret, &k.Description, &k.Status)

	if err != nil {
		return nil, err
	}
	return &k, nil
}

func FindApiKeysByUserId(db *sql.DB, req *pb.GetUserApiKeysRequest) ([]*pb.ApiKey, error) {
	results := make([]*pb.ApiKey, 0)

	rows, err := db.Query("SELECT id, user_id, exchange_name, api_key, secret, description, status FROM user_keys WHERE user_id = $1", req.UserId)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var k pb.ApiKey
		err := rows.Scan(&k.ApiKeyId, &k.UserId, &k.Exchange, &k.Key, &k.Secret, &k.Description, &k.Status)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		results = append(results, &k)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return results, nil
}

func InsertApiKey(db *sql.DB, req *pb.ApiKeyRequest) (*pb.ApiKey, error) {
	newId := uuid.New()

	sqlStatement := `insert into user_keys (id, user_id, exchange_name, api_key, secret, description, status) values ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(sqlStatement, newId, req.UserId, req.Exchange, req.Key, req.Secret, req.Description, "unverified")

	if err != nil {
		return nil, err
	}
	apikey := &pb.ApiKey{
		ApiKeyId:    newId.String(),
		UserId:      req.UserId,
		Exchange:    req.Exchange,
		Key:         req.Key,
		Secret:      req.Secret,
		Description: req.Description,
		Status:      "unverified",
	}
	return apikey, nil
}

func UpdateApiKey(db *sql.DB, req *pb.ApiKeyRequest) (*pb.ApiKey, error) {
	sqlStatement := `UPDATE user_keys SET description = $1, status = $2 WHERE id = $3 and user_id = $4`
	_, err := db.Exec(sqlStatement, req.Description, req.Status, req.ApiKeyId, req.UserId)

	if err != nil {
		return nil, err
	}
	apikey := &pb.ApiKey{
		ApiKeyId:    req.ApiKeyId,
		UserId:      req.UserId,
		Exchange:    req.Exchange,
		Key:         req.Key,
		Secret:      req.Secret,
		Description: req.Description,
		Status:      req.Status,
	}
	return apikey, nil
}
