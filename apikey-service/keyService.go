package main

import (
	"context"
	"database/sql"
	"log"

	keyRepo "github.com/asciiu/gomo/apikey-service/db/sql"
	keyProto "github.com/asciiu/gomo/apikey-service/proto/apikey"
	enums "github.com/asciiu/gomo/common/enums"
)

type KeyService struct {
	DB *sql.DB
}

type ExchangeKey struct {
	ApiKeyId string
	ApiKey   string
	Secret   string
	Exchange enums.Exchange
}

// This private function is responsible for verifying the integrity of an exchange key
// It will update the status of a key to 'verified'.
func (service *KeyService) verifyKey(key ExchangeKey) {
	log.Printf("keyId: %s key: %s secret: %s exchange: %s", key.ApiKeyId, key.ApiKey, key.Secret, key.Exchange)
}

func (service *KeyService) AddApiKey(ctx context.Context, req *keyProto.ApiKeyRequest, res *keyProto.ApiKeyResponse) error {
	apiKey, error := keyRepo.InsertApiKey(service.DB, req)

	switch {
	case error == nil:
		// verify the key in another process
		go service.verifyKey(ExchangeKey{
			ApiKeyId: apiKey.ApiKeyId,
			ApiKey:   apiKey.Key,
			Secret:   apiKey.Secret,
			Exchange: enums.NewExchange(apiKey.Exchange),
		})

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
