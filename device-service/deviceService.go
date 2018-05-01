package main

import (
	"context"
	"database/sql"

	deviceRepo "github.com/asciiu/gomo/device-service/db/sql"
	devices "github.com/asciiu/gomo/device-service/proto/device"
)

type DeviceService struct {
	DB *sql.DB
}

// AddDevice returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *DeviceService) AddDevice(ctx context.Context, req *devices.AddDeviceRequest, res *devices.DeviceResponse) error {
	device, error := deviceRepo.InsertDevice(service.DB, req)

	switch {
	case error == nil:
		res.Status = "success"
		res.Data = &devices.UserDeviceData{
			Device: &devices.Device{
				DeviceID:         device.DeviceID,
				UserID:           device.UserID,
				ExternalDeviceID: device.ExternalDeviceID,
				DeviceType:       device.DeviceType,
				DeviceToken:      device.DeviceToken,
			},
		}
	default:
		res.Status = "error"
		res.Message = error.Error()
	}
	return nil
}

// GetUserDevice returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *DeviceService) GetUserDevice(ctx context.Context, req *devices.GetUserDeviceRequest, res *devices.DeviceResponse) error {
	device, error := deviceRepo.FindDeviceByDeviceID(service.DB, req)

	if error == nil {
		res.Status = "success"
		res.Data = &devices.UserDeviceData{
			Device: &devices.Device{
				DeviceID:         device.DeviceID,
				UserID:           device.UserID,
				ExternalDeviceID: device.ExternalDeviceID,
				DeviceType:       device.DeviceType,
				DeviceToken:      device.DeviceToken,
			},
		}
	} else {
		res.Status = "error"
		res.Message = error.Error()
	}

	return nil
}

// GetUserDevices returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *DeviceService) GetUserDevices(ctx context.Context, req *devices.GetUserDevicesRequest, res *devices.DeviceListResponse) error {
	dvs, error := deviceRepo.FindDevicesByUserID(service.DB, req)

	if error == nil {
		res.Status = "success"
		res.Data = &devices.UserDevicesData{
			Devices: dvs,
		}
	} else {
		res.Status = "error"
		res.Message = error.Error()
	}

	return nil
}

// RemoveDevice returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *DeviceService) RemoveDevice(ctx context.Context, req *devices.RemoveDeviceRequest, res *devices.DeviceResponse) error {
	error := deviceRepo.DeleteDevice(service.DB, req.DeviceID)
	if error == nil {
		res.Status = "success"
	} else {
		res.Status = "error"
		res.Message = error.Error()
	}
	return nil
}

// UpdateDevice returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *DeviceService) UpdateDevice(ctx context.Context, req *devices.UpdateDeviceRequest, res *devices.DeviceResponse) error {
	device, error := deviceRepo.UpdateDevice(service.DB, req)
	if error == nil {
		res.Status = "success"
		res.Data = &devices.UserDeviceData{
			Device: &devices.Device{
				DeviceID:         device.DeviceID,
				UserID:           device.UserID,
				ExternalDeviceID: device.ExternalDeviceID,
				DeviceType:       device.DeviceType,
				DeviceToken:      device.DeviceToken,
			},
		}
	} else {
		res.Status = "error"
		res.Message = error.Error()
	}

	return nil
}
