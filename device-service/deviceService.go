package main

import (
	"context"
	"database/sql"

	constRes "github.com/asciiu/gomo/common/constants/response"
	repoDevice "github.com/asciiu/gomo/device-service/db/sql"
	protoDevice "github.com/asciiu/gomo/device-service/proto/device"
)

type DeviceService struct {
	DB *sql.DB
}

// AddDevice returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *DeviceService) AddDevice(ctx context.Context, req *protoDevice.AddDeviceRequest, res *protoDevice.DeviceResponse) error {
	mreq := protoDevice.GetDeviceMatchRequest{
		UserID:           req.UserID,
		DeviceType:       req.DeviceType,
		DeviceToken:      req.DeviceToken,
		ExternalDeviceID: req.ExternalDeviceID,
	}

	device, error := repoDevice.FindDeviceMatch(service.DB, &mreq)
	switch {
	case error == sql.ErrNoRows:
		// there were no matches found therefore insert it
		di, error := repoDevice.InsertDevice(service.DB, req)
		switch {
		case error == nil:
			res.Status = constRes.Success
			res.Data = &protoDevice.UserDeviceData{
				Device: &protoDevice.Device{
					DeviceID:         di.DeviceID,
					UserID:           di.UserID,
					ExternalDeviceID: di.ExternalDeviceID,
					DeviceType:       di.DeviceType,
					DeviceToken:      di.DeviceToken,
				},
			}
		default:
			res.Status = constRes.Error
			res.Message = error.Error()
		}

	case device != nil:
		// device match found update the device
		ureq := protoDevice.UpdateDeviceRequest{
			DeviceID:         device.DeviceID,
			UserID:           req.UserID,
			DeviceType:       req.DeviceType,
			DeviceToken:      req.DeviceToken,
			ExternalDeviceID: req.ExternalDeviceID,
		}
		du, error := repoDevice.UpdateDevice(service.DB, &ureq)
		switch {
		case error == nil:
			res.Status = constRes.Success
			res.Data = &protoDevice.UserDeviceData{
				Device: &protoDevice.Device{
					DeviceID:         du.DeviceID,
					UserID:           du.UserID,
					ExternalDeviceID: du.ExternalDeviceID,
					DeviceType:       du.DeviceType,
					DeviceToken:      du.DeviceToken,
				},
			}
		default:
			res.Status = "error"
			res.Message = error.Error()
		}
	}

	return nil
}

// GetUserDevice returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *DeviceService) GetUserDevice(ctx context.Context, req *protoDevice.GetUserDeviceRequest, res *protoDevice.DeviceResponse) error {
	device, error := repoDevice.FindDeviceByDeviceID(service.DB, req)

	if error == nil {
		res.Status = constRes.Success
		res.Data = &protoDevice.UserDeviceData{
			Device: &protoDevice.Device{
				DeviceID:         device.DeviceID,
				UserID:           device.UserID,
				ExternalDeviceID: device.ExternalDeviceID,
				DeviceType:       device.DeviceType,
				DeviceToken:      device.DeviceToken,
			},
		}
	} else {
		res.Status = constRes.Error
		res.Message = error.Error()
	}

	return nil
}

// GetUserDevices returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *DeviceService) GetUserDevices(ctx context.Context, req *protoDevice.GetUserDevicesRequest, res *protoDevice.DeviceListResponse) error {
	dvs, err := repoDevice.FindDevicesByUserID(service.DB, req)

	if err == nil {
		res.Status = constRes.Success
		res.Data = &protoDevice.UserDevicesData{
			Devices: dvs,
		}
	} else {
		res.Status = constRes.Error
		res.Message = err.Error()
	}

	return nil
}

// RemoveDevice returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *DeviceService) RemoveDevice(ctx context.Context, req *protoDevice.RemoveDeviceRequest, res *protoDevice.DeviceResponse) error {
	error := repoDevice.DeleteDevice(service.DB, req.DeviceID)
	if error == nil {
		res.Status = constRes.Success
	} else {
		res.Status = constRes.Error
		res.Message = error.Error()
	}
	return nil
}

// UpdateDevice returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *DeviceService) UpdateDevice(ctx context.Context, req *protoDevice.UpdateDeviceRequest, res *protoDevice.DeviceResponse) error {
	device, error := repoDevice.UpdateDevice(service.DB, req)
	if error == nil {
		res.Status = constRes.Success
		res.Data = &protoDevice.UserDeviceData{
			Device: &protoDevice.Device{
				DeviceID:         device.DeviceID,
				UserID:           device.UserID,
				ExternalDeviceID: device.ExternalDeviceID,
				DeviceType:       device.DeviceType,
				DeviceToken:      device.DeviceToken,
			},
		}
	} else {
		res.Status = constRes.Error
		res.Message = error.Error()
	}

	return nil
}
