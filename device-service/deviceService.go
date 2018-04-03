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
				DeviceId:         device.DeviceId,
				UserId:           device.UserId,
				ExternalDeviceId: device.ExternalDeviceId,
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
	device, error := deviceRepo.FindDeviceByDeviceId(service.DB, req)

	if error == nil {
		res.Status = "success"
		res.Data = &pb.UserDeviceData{
			Device: &pb.Device{
				DeviceId:         device.DeviceId,
				UserId:           device.UserId,
				ExternalDeviceId: device.ExternalDeviceId,
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
	return nil
}

func (service *DeviceService) RemoveDevice(ctx context.Context, req *pb.RemoveDeviceRequest, res *pb.Response) error {
	error := deviceRepo.DeleteDevice(service.DB, req.DeviceId)
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
				DeviceId:         device.DeviceId,
				UserId:           device.UserId,
				ExternalDeviceId: device.ExternalDeviceId,
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
