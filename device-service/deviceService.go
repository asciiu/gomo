package main

import (
	"context"
	"database/sql"

	deviceRepo "github.com/asciiu/gomo/device-service/db/sql"
	pb "github.com/asciiu/gomo/device-service/proto/device"
)

type DeviceService struct {
	DB *sql.DB
}

func (service *DeviceService) AddDevice(ctx context.Context, req *pb.AddDeviceRequest, res *pb.DeviceResponse) error {
	device, error := deviceRepo.InsertDevice(service.DB, req)

	switch {
	case error == nil:
		res.Status = "success"
		res.Data = &pb.UserDeviceData{
			Device: &pb.Device{
				DeviceID:         device.DeviceID,
				UserID:           device.UserID,
				ExternalDeviceID: device.ExternalDeviceID,
				DeviceType:       device.DeviceType,
				DeviceToken:      device.DeviceToken,
			},
		}
		return nil

	default:
		res.Status = "error"
		res.Message = error.Error()
		return error
	}
}

func (service *DeviceService) GetUserDevice(ctx context.Context, req *pb.GetUserDeviceRequest, res *pb.DeviceResponse) error {
	device, error := deviceRepo.FindDeviceByDeviceID(service.DB, req)

	if error == nil {
		res.Status = "success"
		res.Data = &pb.UserDeviceData{
			Device: &pb.Device{
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

	return error
}

func (service *DeviceService) GetUserDevices(ctx context.Context, req *pb.GetUserDevicesRequest, res *pb.DeviceListResponse) error {
	devices, error := deviceRepo.FindDevicesByUserID(service.DB, req)

	if error == nil {
		res.Status = "success"
		res.Data = &pb.UserDevicesData{
			Device: devices,
		}
	} else {
		res.Status = "error"
		res.Message = error.Error()
	}

	return error
}

func (service *DeviceService) RemoveDevice(ctx context.Context, req *pb.RemoveDeviceRequest, res *pb.DeviceResponse) error {
	error := deviceRepo.DeleteDevice(service.DB, req.DeviceID)
	if error == nil {
		res.Status = "success"
	} else {
		res.Status = "error"
		res.Message = error.Error()
	}
	return error
}

func (service *DeviceService) UpdateDevice(ctx context.Context, req *pb.UpdateDeviceRequest, res *pb.DeviceResponse) error {
	device, error := deviceRepo.UpdateDevice(service.DB, req)
	if error == nil {
		res.Status = "success"
		res.Data = &pb.UserDeviceData{
			Device: &pb.Device{
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

	return error
}
