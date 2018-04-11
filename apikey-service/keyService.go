package main

import (
	"context"
	"database/sql"

	pb "github.com/asciiu/gomo/apikey-service/proto/apikey"
)

type KeyService struct {
	DB *sql.DB
}

func (service *KeyService) AddApiKey(ctx context.Context, req *pb.ApiKeyRequest, res *pb.ApiKeyResponse) error {
	return nil
	//device, error := deviceRepo.InsertDevice(service.DB, req)

	//switch {
	//case error == nil:
	//	res.Status = "success"
	//	res.Data = &pb.UserDeviceData{
	//		Device: &pb.Device{
	//			DeviceId:         device.DeviceId,
	//			UserId:           device.UserId,
	//			ExternalDeviceId: device.ExternalDeviceId,
	//			DeviceType:       device.DeviceType,
	//			DeviceToken:      device.DeviceToken,
	//		},
	//	}
	//	return nil

	//default:
	//	res.Status = "error"
	//	res.Message = error.Error()
	//	return error
	//}
}

func (service *KeyService) GetUserApiKey(ctx context.Context, req *pb.GetUserApiKeyRequest, res *pb.ApiKeyResponse) error {
	return nil
	//device, error := deviceRepo.FindDeviceByDeviceId(service.DB, req)

	//if error == nil {
	//	res.Status = "success"
	//	res.Data = &pb.UserDeviceData{
	//		Device: &pb.Device{
	//			DeviceId:         device.DeviceId,
	//			UserId:           device.UserId,
	//			ExternalDeviceId: device.ExternalDeviceId,
	//			DeviceType:       device.DeviceType,
	//			DeviceToken:      device.DeviceToken,
	//		},
	//	}
	//} else {
	//	res.Status = "error"
	//	res.Message = error.Error()
	//}

	//return error
}

func (service *KeyService) GetUserApiKeys(ctx context.Context, req *pb.GetUserApiKeysRequest, res *pb.ApiKeyListResponse) error {
	return nil
	//devices, error := deviceRepo.FindDevicesByUserId(service.DB, req)

	//if error == nil {
	//	res.Status = "success"
	//	res.Data = &pb.UserDevicesData{
	//		Device: devices,
	//	}
	//} else {
	//	res.Status = "error"
	//	res.Message = error.Error()
	//}

	//return error
}

func (service *KeyService) RemoveApiKey(ctx context.Context, req *pb.RemoveApiKeyRequest, res *pb.ApiKeyResponse) error {
	return nil
	//error := deviceRepo.DeleteDevice(service.DB, req.DeviceId)
	//if error == nil {
	//	res.Status = "success"
	//} else {
	//	res.Status = "error"
	//	res.Message = error.Error()
	//}
	//return error
}

func (service *KeyService) UpdateApiKey(ctx context.Context, req *pb.ApiKeyRequest, res *pb.ApiKeyResponse) error {
	return nil
	//error := deviceRepo.DeleteDevice(service.DB, req.DeviceId)
	//if error == nil {
	//	res.Status = "success"
	//} else {
	//	res.Status = "error"
	//	res.Message = error.Error()
	//}
	//return error
}
