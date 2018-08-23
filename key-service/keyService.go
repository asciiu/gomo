package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	protoActivity "github.com/asciiu/gomo/activity-bulletin/proto"
	constExch "github.com/asciiu/gomo/common/constants/exchange"
	constRes "github.com/asciiu/gomo/common/constants/response"
	repoKey "github.com/asciiu/gomo/key-service/db/sql"
	protoKey "github.com/asciiu/gomo/key-service/proto/key"
	micro "github.com/micro/go-micro"
)

type KeyService struct {
	DB        *sql.DB
	KeyPub    micro.Publisher
	NotifyPub micro.Publisher
}

// IMPORTANT! When adding support for new exchanges we must add there names here!
// TODO: read these from a config or from the DB
var supported = [...]string{constExch.Binance}

func stringInSupported(a string) bool {
	for _, b := range supported {
		if b == a {
			return true
		}
	}
	return false
}

func (service *KeyService) HandleVerifiedKey(ctx context.Context, key *protoKey.Key) error {
	log.Println("received verified key ", key.KeyID)

	_, error := repoKey.UpdateKeyStatus(service.DB, key)

	description := fmt.Sprintf("%s key verified", key.Exchange)
	notification := protoActivity.Activity{
		UserID:      key.UserID,
		Type:        "key",
		ObjectID:    key.KeyID,
		Title:       "Exchange Setup",
		Description: description,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	// publish verify key event
	if err := service.NotifyPub.Publish(context.Background(), &notification); err != nil {
		log.Println("could not publish verified key event: ", err)
	}

	return error
}

// AddKey returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *KeyService) AddKey(ctx context.Context, req *protoKey.KeyRequest, res *protoKey.KeyResponse) error {

	// supported exchange keys check
	if !stringInSupported(req.Exchange) {
		res.Status = constRes.Fail
		res.Message = fmt.Sprintf("%s is not a supported exchange", req.Exchange)
		return nil
	}

	apiKey, error := repoKey.InsertKey(service.DB, req)

	switch {
	case error == nil:
		// ignore in test
		if service.KeyPub != nil {
			// publish a new key event
			if err := service.KeyPub.Publish(context.Background(), apiKey); err != nil {
				log.Println("could not publish event new.key: ", err)
			}
		}

		res.Status = constRes.Success
		res.Data = &protoKey.UserKeyData{
			Key: &protoKey.Key{
				KeyID:       apiKey.KeyID,
				UserID:      apiKey.UserID,
				Exchange:    apiKey.Exchange,
				Key:         apiKey.Key,
				Description: apiKey.Description,
				Status:      apiKey.Status,
			},
		}

	case strings.Contains(error.Error(), "user_keys_api_key_secret_key"):
		res.Status = constRes.Fail
		res.Message = "key/secret pair already exists"

	default:
		res.Status = constRes.Error
		res.Message = error.Error()
	}
	return nil
}

func (service *KeyService) GetKeys(ctx context.Context, req *protoKey.GetKeysRequest, res *protoKey.KeyListResponse) error {
	keys, error := repoKey.FindKeys(service.DB, req)

	switch {
	case error == nil:
		res.Status = constRes.Success
		res.Data = &protoKey.UserKeysData{
			Keys: keys,
		}
	default:
		res.Status = constRes.Error
		res.Message = error.Error()
	}
	return nil
}

// GetUserKey returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *KeyService) GetUserKey(ctx context.Context, req *protoKey.GetUserKeyRequest, res *protoKey.KeyResponse) error {
	apiKey, error := repoKey.FindKeyByID(service.DB, req)

	switch {
	case error == nil:
		res.Status = constRes.Success
		res.Data = &protoKey.UserKeyData{
			Key: &protoKey.Key{
				KeyID:       apiKey.KeyID,
				UserID:      apiKey.UserID,
				Exchange:    apiKey.Exchange,
				Key:         apiKey.Key,
				Secret:      apiKey.Secret,
				Description: apiKey.Description,
				Status:      apiKey.Status,
			},
		}
	case strings.Contains(error.Error(), "no rows in result set"):
		res.Status = constRes.Nonentity
		res.Message = fmt.Sprintf("key not found")
	default:
		res.Status = constRes.Error
		res.Message = error.Error()
	}

	return nil
}

// GetUserKeys returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *KeyService) GetUserKeys(ctx context.Context, req *protoKey.GetUserKeysRequest, res *protoKey.KeyListResponse) error {
	keys, error := repoKey.FindKeysByUserID(service.DB, req)

	switch {
	case error == nil:
		res.Status = constRes.Success
		res.Data = &protoKey.UserKeysData{
			Keys: keys,
		}
	default:
		res.Status = constRes.Error
		res.Message = error.Error()
	}
	return nil
}

// RemoveKey returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *KeyService) RemoveKey(ctx context.Context, req *protoKey.RemoveKeyRequest, res *protoKey.KeyResponse) error {
	error := repoKey.DeleteKey(service.DB, req.KeyID)
	switch {
	case error == nil:
		res.Status = constRes.Success
	default:
		res.Status = constRes.Error
		res.Message = error.Error()
	}
	return nil
}

// UpdateKeyDescription returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *KeyService) UpdateKeyDescription(ctx context.Context, req *protoKey.KeyRequest, res *protoKey.KeyResponse) error {
	apiKey, error := repoKey.UpdateKeyDescription(service.DB, req)
	switch {
	case error == nil:
		res.Status = constRes.Success
		res.Data = &protoKey.UserKeyData{
			Key: &protoKey.Key{
				KeyID:       apiKey.KeyID,
				UserID:      apiKey.UserID,
				Exchange:    apiKey.Exchange,
				Key:         apiKey.Key,
				Description: apiKey.Description,
				Status:      apiKey.Status,
			},
		}
	default:
		res.Status = constRes.Error
		res.Message = error.Error()
	}
	return nil
}
