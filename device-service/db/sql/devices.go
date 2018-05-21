package sql

import (
	"database/sql"
	"log"

	pb "github.com/asciiu/gomo/device-service/proto/device"
	"github.com/google/uuid"
)

func DeleteDevice(db *sql.DB, deviceID string) error {
	_, err := db.Exec("DELETE FROM user_devices WHERE id = $1", deviceID)
	return err
}

func FindDeviceByDeviceID(db *sql.DB, req *pb.GetUserDeviceRequest) (*pb.Device, error) {
	var d pb.Device
	err := db.QueryRow("SELECT id, user_id, device_id, device_type, device_token FROM user_devices WHERE id = $1", req.DeviceID).
		Scan(&d.DeviceID, &d.UserID, &d.ExternalDeviceID, &d.DeviceType, &d.DeviceToken)

	if err != nil {
		return nil, err
	}
	return &d, nil
}

func FindDeviceMatch(db *sql.DB, req *pb.GetDeviceMatchRequest) (*pb.Device, error) {
	var d pb.Device
	query := `SELECT id, user_id, device_id, device_type, device_token 
		FROM user_devices WHERE user_id = $1 and device_type = $2 and device_id = $3`
	err := db.QueryRow(query, req.UserID, req.DeviceType, req.ExternalDeviceID).
		Scan(&d.DeviceID, &d.UserID, &d.ExternalDeviceID, &d.DeviceType, &d.DeviceToken)

	if err != nil {
		return nil, err
	}
	return &d, nil
}

func FindDevicesByUserID(db *sql.DB, req *pb.GetUserDevicesRequest) ([]*pb.Device, error) {
	results := make([]*pb.Device, 0)

	rows, err := db.Query("SELECT id, user_id, device_id, device_type, device_token FROM user_devices WHERE user_id = $1", req.UserID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var d pb.Device
		err := rows.Scan(&d.DeviceID, &d.UserID, &d.ExternalDeviceID, &d.DeviceType, &d.DeviceToken)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		results = append(results, &d)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return results, nil
}

func InsertDevice(db *sql.DB, req *pb.AddDeviceRequest) (*pb.Device, error) {
	newID := uuid.New()

	sqlStatement := `insert into user_devices (id, user_id, device_id, device_type, device_token) values ($1, $2, $3, $4, $5)`
	_, err := db.Exec(sqlStatement, newID, req.UserID, req.ExternalDeviceID, req.DeviceType, req.DeviceToken)

	if err != nil {
		return nil, err
	}
	device := &pb.Device{
		DeviceID:         newID.String(),
		UserID:           req.UserID,
		ExternalDeviceID: req.ExternalDeviceID,
		DeviceType:       req.DeviceType,
		DeviceToken:      req.DeviceToken,
	}
	return device, nil
}

func UpdateDevice(db *sql.DB, req *pb.UpdateDeviceRequest) (*pb.UpdateDeviceRequest, error) {
	sqlStatement := `UPDATE user_devices SET device_id = $1, device_type = $2, device_token = $3 WHERE id = $4`
	_, err := db.Exec(sqlStatement, req.ExternalDeviceID, req.DeviceType, req.DeviceToken, req.DeviceID)

	if err != nil {
		return nil, err
	}
	return req, nil
}
