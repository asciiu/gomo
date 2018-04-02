package main

import (
	"context"
	"database/sql"

	deviceRepo "github.com/asciiu/gomo/device-service/db/sql"
	"github.com/asciiu/gomo/device-service/models"
	pb "github.com/asciiu/gomo/device-service/proto/device"
)

type DeviceService struct {
	DB *sql.DB
}

func (service *DeviceService) AddDevice(ctx context.Context, req *pb.AddDeviceRequest, res *pb.DeviceResponse) error {
	device := models.NewDevice(req.UserId, req.ExternalDeviceId, req.DeviceType, req.DeviceToken)
	_, error := deviceRepo.InsertDevice(service.DB, device)

	switch {
	case error == nil:
		res.Status = "success"
		res.Data = &pb.UserDeviceData{
			Device: &pb.Device{
				DeviceId:         device.Id,
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
	return nil
}

func (service *DeviceService) GetUserDevices(ctx context.Context, req *pb.GetUserDevicesRequest, res *pb.DeviceListResponse) error {
	return nil
}

func (service *DeviceService) RemoveDevice(ctx context.Context, req *pb.RemoveDeviceRequest, res *pb.Response) error {
	return nil
}

func (service *DeviceService) UpdateDevice(ctx context.Context, req *pb.UpdateDeviceRequest, res *pb.DeviceResponse) error {
	return nil
}
