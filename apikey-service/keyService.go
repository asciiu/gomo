package main

import (
	"context"
	"database/sql"

	keyRepo "github.com/asciiu/gomo/apikey-service/db/sql"
	keyProto "github.com/asciiu/gomo/apikey-service/proto/apikey"
)

type KeyService struct {
	DB *sql.DB
}

func (service *KeyService) AddApiKey(ctx context.Context, req *keyProto.ApiKeyRequest, res *keyProto.ApiKeyResponse) error {
	apiKey, error := keyRepo.InsertApiKey(service.DB, req)

	switch {
	case error == nil:
		res.Status = "success"
		res.Data = &keyProto.UserApiKeyData{
			ApiKey: &keyProto.ApiKey{
				ApiKeyId:    apiKey.ApiKeyId,
				UserId:      apiKey.UserId,
				Exchange:    apiKey.Exchange,
				Key:         apiKey.Key,
				Description: apiKey.Description,
				Status:      apiKey.Status,
			},
		}
		return nil

	default:
		res.Status = "error"
		res.Message = error.Error()
		return error
	}
}

func (service *KeyService) GetUserApiKey(ctx context.Context, req *keyProto.GetUserApiKeyRequest, res *keyProto.ApiKeyResponse) error {
	apiKey, error := keyRepo.FindApiKeyById(service.DB, req)

	switch {
	case error == nil:
		res.Status = "success"
		res.Status = "success"
		res.Data = &keyProto.UserApiKeyData{
			ApiKey: &keyProto.ApiKey{
				ApiKeyId:    apiKey.ApiKeyId,
				UserId:      apiKey.UserId,
				Exchange:    apiKey.Exchange,
				Key:         apiKey.Key,
				Description: apiKey.Description,
				Status:      apiKey.Status,
			},
		}
		return nil
	default:
		res.Status = "error"
		res.Message = error.Error()
		return error
	}
}

func (service *KeyService) GetUserApiKeys(ctx context.Context, req *keyProto.GetUserApiKeysRequest, res *keyProto.ApiKeyListResponse) error {
	keys, error := keyRepo.FindApiKeysByUserId(service.DB, req)

	switch {
	case error == nil:
		res.Status = "success"
		res.Data = &keyProto.UserApiKeysData{
			ApiKey: keys,
		}
		return nil
	default:
		res.Status = "error"
		res.Message = error.Error()
		return error
	}
}

func (service *KeyService) RemoveApiKey(ctx context.Context, req *keyProto.RemoveApiKeyRequest, res *keyProto.ApiKeyResponse) error {
	error := keyRepo.DeleteApiKey(service.DB, req.ApiKeyId)
	switch {
	case error == nil:
		res.Status = "success"
		return nil
	default:
		res.Status = "error"
		res.Message = error.Error()
		return error
	}
}

func (service *KeyService) UpdateApiKeyDescription(ctx context.Context, req *keyProto.ApiKeyRequest, res *keyProto.ApiKeyResponse) error {
	apiKey, error := keyRepo.UpdateApiKeyDescription(service.DB, req)
	switch {
	case error == nil:
		res.Status = "success"
		res.Data = &keyProto.UserApiKeyData{
			ApiKey: &keyProto.ApiKey{
				ApiKeyId:    apiKey.ApiKeyId,
				UserId:      apiKey.UserId,
				Exchange:    apiKey.Exchange,
				Key:         apiKey.Key,
				Description: apiKey.Description,
				Status:      apiKey.Status,
			},
		}
		return nil
	default:
		res.Status = "error"
		res.Message = error.Error()
		return error
	}
}
