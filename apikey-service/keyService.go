package main

import (
	"context"
	"database/sql"
	"log"

	keyRepo "github.com/asciiu/gomo/apikey-service/db/sql"
	keyProto "github.com/asciiu/gomo/apikey-service/proto/apikey"
	"github.com/asciiu/gomo/common/db"
	msg "github.com/asciiu/gomo/common/messages"
	micro "github.com/micro/go-micro"
)

type KeyService struct {
	DB    *sql.DB
	Micro micro.Service
}

func NewKeyService(name, dbUrl string) micro.Service {
	//broker := rabbitmq.NewBroker()

	// Create a new service. Include some options here.
	srv := micro.NewService(
		// This name must match the package name given in your protobuf definition
		micro.Name(name),
		//micro.Broker(broker),
		micro.Version("latest"),
	)

	// Init will parse the command line flags.
	srv.Init()

	gomoDB, err := db.NewDB(dbUrl)

	// TODO read secret from env var
	//dbUrl := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	keyProto.RegisterApiKeyServiceHandler(srv.Server(), &KeyService{gomoDB, srv})

	return srv
}

func (service *KeyService) AddApiKey(ctx context.Context, req *keyProto.ApiKeyRequest, res *keyProto.ApiKeyResponse) error {
	apiKey, error := keyRepo.InsertApiKey(service.DB, req)

	switch {
	case error == nil:
		// publish a new key event
		publisher := micro.NewPublisher(msg.TopicNewKey, service.Micro.Client())
		if err := publisher.Publish(context.Background(), apiKey); err != nil {
			log.Println("could not publish event new.key: ", err)
		}

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
