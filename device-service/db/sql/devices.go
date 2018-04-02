package sql

import (
	"database/sql"

	"github.com/asciiu/gomo/device-service/models"
)

func DeleteDevice(db *sql.DB, deviceId string) error {
	_, err := db.Exec("DELETE FROM user_devices WHERE id = $1", deviceId)
	return err
}

func FindDevice(db *sql.DB, deviceId string) (*models.UserDevice, error) {
	var d models.UserDevice
	err := db.QueryRow("SELECT id, user_id, device_id, device_type, device_token, created_on, updated_on FROM user_devices WHERE id = $1", deviceId).
		Scan(&d.Id, &d.UserId, &d.DeviceId, &d.DeviceType, &d.DeviceToken, &d.CreatedOn, &d.UpdatedOn)

	if err != nil {
		return nil, err
	}
	return &d, nil
}

//func FindUserById(db *sql.DB, userId string) (*models.User, error) {
//	var u models.User
//	err := db.QueryRow("SELECT id, first_name, last_name, email, email_verified, password_hash FROM users WHERE id = $1", userId).
//		Scan(&u.Id, &u.First, &u.Last, &u.Email, &u.EmailVerified, &u.PasswordHash)
//
//	if err != nil {
//		return nil, err
//	}
//	return &u, nil
//}

func InsertDevice(db *sql.DB, device *models.UserDevice) (*models.UserDevice, error) {
	sqlStatement := `insert into user_devices (id, user_id, device_id, device_type, device_token) values ($1, $2, $3, $4, $5)`
	_, err := db.Exec(sqlStatement, device.Id, device.UserId, device.DeviceId, device.DeviceType, device.DeviceToken)

	if err != nil {
		return nil, err
	}
	return device, nil
}

func UpdateDevice(db *sql.DB, device *models.UserDevice) (*models.UserDevice, error) {
	sqlStatement := `UPDATE user_devices SET device_id = $1, device_type = $2, device_token = $3 WHERE id = $4`
	_, err := db.Exec(sqlStatement, device.DeviceId, device.DeviceType, device.DeviceToken, device.Id)

	if err != nil {
		return nil, err
	}
	return device, nil
}
