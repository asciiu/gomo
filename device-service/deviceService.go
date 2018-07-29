package main

import (
	"context"
	"database/sql"

	responseConstants "github.com/asciiu/gomo/common/constants/response"
	deviceRepo "github.com/asciiu/gomo/device-service/db/sql"
	deviceProto "github.com/asciiu/gomo/device-service/proto/device"
)

type DeviceService struct {
	DB *sql.DB
}

// AddDevice returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *DeviceService) AddDevice(ctx context.Context, req *deviceProto.AddDeviceRequest, res *deviceProto.DeviceResponse) error {
	mreq := deviceProto.GetDeviceMatchRequest{
		UserID:           req.UserID,
		DeviceType:       req.DeviceType,
		DeviceToken:      req.DeviceToken,
		ExternalDeviceID: req.ExternalDeviceID,
	}

	device, error := deviceRepo.FindDeviceMatch(service.DB, &mreq)
	switch {
	case error == sql.ErrNoRows:
		// there were no matches found therefore insert it
		di, error := deviceRepo.InsertDevice(service.DB, req)
		switch {
		case error == nil:
			res.Status = responseConstants.Success
			res.Data = &deviceProto.UserDeviceData{
				Device: &deviceProto.Device{
					DeviceID:         di.DeviceID,
					UserID:           di.UserID,
					ExternalDeviceID: di.ExternalDeviceID,
					DeviceType:       di.DeviceType,
					DeviceToken:      di.DeviceToken,
				},
			}
		default:
			res.Status = responseConstants.Error
			res.Message = error.Error()
		}

	case device != nil:
		// device match found update the device
		ureq := deviceProto.UpdateDeviceRequest{
			DeviceID:         device.DeviceID,
			UserID:           req.UserID,
			DeviceType:       req.DeviceType,
			DeviceToken:      req.DeviceToken,
			ExternalDeviceID: req.ExternalDeviceID,
		}
		du, error := deviceRepo.UpdateDevice(service.DB, &ureq)
		switch {
		case error == nil:
			res.Status = responseConstants.Success
			res.Data = &deviceProto.UserDeviceData{
				Device: &deviceProto.Device{
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
func (service *DeviceService) GetUserDevice(ctx context.Context, req *deviceProto.GetUserDeviceRequest, res *deviceProto.DeviceResponse) error {
	device, error := deviceRepo.FindDeviceByDeviceID(service.DB, req)

	if error == nil {
		res.Status = responseConstants.Success
		res.Data = &deviceProto.UserDeviceData{
			Device: &deviceProto.Device{
				DeviceID:         device.DeviceID,
				UserID:           device.UserID,
				ExternalDeviceID: device.ExternalDeviceID,
				DeviceType:       device.DeviceType,
				DeviceToken:      device.DeviceToken,
			},
		}
	} else {
		res.Status = responseConstants.Error
		res.Message = error.Error()
	}

	return nil
}

// GetUserDevices returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *DeviceService) GetUserDevices(ctx context.Context, req *deviceProto.GetUserDevicesRequest, res *deviceProto.DeviceListResponse) error {
	dvs, err := deviceRepo.FindDevicesByUserID(service.DB, req)

	if err == nil {
		res.Status = responseConstants.Success
		res.Data = &deviceProto.UserDevicesData{
			Devices: dvs,
		}
	} else {
		res.Status = responseConstants.Error
		res.Message = err.Error()
	}

	return nil
}

// RemoveDevice returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *DeviceService) RemoveDevice(ctx context.Context, req *deviceProto.RemoveDeviceRequest, res *deviceProto.DeviceResponse) error {
	error := deviceRepo.DeleteDevice(service.DB, req.DeviceID)
	if error == nil {
		res.Status = responseConstants.Success
	} else {
		res.Status = responseConstants.Error
		res.Message = error.Error()
	}
	return nil
}

// UpdateDevice returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *DeviceService) UpdateDevice(ctx context.Context, req *deviceProto.UpdateDeviceRequest, res *deviceProto.DeviceResponse) error {
	device, error := deviceRepo.UpdateDevice(service.DB, req)
	if error == nil {
		res.Status = responseConstants.Success
		res.Data = &deviceProto.UserDeviceData{
			Device: &deviceProto.Device{
				DeviceID:         device.DeviceID,
				UserID:           device.UserID,
				ExternalDeviceID: device.ExternalDeviceID,
				DeviceType:       device.DeviceType,
				DeviceToken:      device.DeviceToken,
			},
		}
	} else {
		res.Status = responseConstants.Error
		res.Message = error.Error()
	}

	return nil
}
