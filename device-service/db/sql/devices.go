package sql

import (
	"database/sql"
	"log"

	pb "github.com/asciiu/gomo/device-service/proto/device"
	"github.com/google/uuid"
)

func DeleteDevice(db *sql.DB, deviceId string) error {
	_, err := db.Exec("DELETE FROM user_devices WHERE id = $1", deviceId)
	return err
}

func FindDeviceByDeviceId(db *sql.DB, req *pb.GetUserDeviceRequest) (*pb.Device, error) {
	var d pb.Device
	err := db.QueryRow("SELECT id, user_id, device_id, device_type, device_token FROM user_devices WHERE id = $1", req.DeviceId).
		Scan(&d.DeviceId, &d.UserId, &d.ExternalDeviceId, &d.DeviceType, &d.DeviceToken)

	if err != nil {
		return nil, err
	}
	return &d, nil
}

func FindDevicesByUserId(db *sql.DB, req *pb.GetUserDevicesRequest) ([]*pb.Device, error) {
	results := make([]*pb.Device, 0)

	rows, err := db.Query("SELECT id, user_id, device_id, device_type, device_token FROM user_devices WHERE user_id = $1", req.UserId)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var d pb.Device
		err := rows.Scan(&d.DeviceId, &d.UserId, &d.ExternalDeviceId, &d.DeviceType, &d.DeviceToken)
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
	newId := uuid.New()

	sqlStatement := `insert into user_devices (id, user_id, device_id, device_type, device_token) values ($1, $2, $3, $4, $5)`
	_, err := db.Exec(sqlStatement, newId, req.UserId, req.ExternalDeviceId, req.DeviceType, req.DeviceToken)

	if err != nil {
		return nil, err
	}
	device := &pb.Device{
		DeviceId:         newId.String(),
		UserId:           req.UserId,
		ExternalDeviceId: req.ExternalDeviceId,
		DeviceType:       req.DeviceType,
		DeviceToken:      req.DeviceToken,
	}
	return device, nil
}

func UpdateDevice(db *sql.DB, req *pb.UpdateDeviceRequest) (*pb.UpdateDeviceRequest, error) {
	sqlStatement := `UPDATE user_devices SET device_id = $1, device_type = $2, device_token = $3 WHERE id = $4`
	_, err := db.Exec(sqlStatement, req.ExternalDeviceId, req.DeviceType, req.DeviceToken, req.DeviceId)

	if err != nil {
		return nil, err
	}
	return req, nil
}
