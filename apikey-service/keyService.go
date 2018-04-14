package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	keyRepo "github.com/asciiu/gomo/apikey-service/db/sql"
	keyProto "github.com/asciiu/gomo/apikey-service/proto/apikey"
	"github.com/asciiu/gomo/common/db"
	enums "github.com/asciiu/gomo/common/enums"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-plugins/broker/rabbitmq"
)

type KeyService struct {
	DB     *sql.DB
	PubSub broker.Broker
}

type ExchangeKey struct {
	ApiKeyId string
	ApiKey   string
	Secret   string
	Exchange enums.Exchange
}

const TopicNewKey = "new.key"

func NewKeyService(name, dbUrl string) micro.Service {
	broker := rabbitmq.NewBroker()

	// Create a new service. Include some options here.
	srv := micro.NewService(
		// This name must match the package name given in your protobuf definition
		micro.Name(name),
		micro.Broker(broker),
		micro.Version("latest"),
	)

	// Init will parse the command line flags.
	srv.Init()
	//publisher := micro.NewPublisher("user.created", srv.Client())

	pubsub := srv.Server().Options().Broker

	gomoDB, err := db.NewDB(dbUrl)

	// TODO read secret from env var
	//dbUrl := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	keyProto.RegisterApiKeyServiceHandler(srv.Server(), &KeyService{gomoDB, pubsub})

	return srv
}

func (srv *KeyService) publishEvent(key ExchangeKey) error {

	// Marshal to JSON string
	body, err := json.Marshal(key)
	if err != nil {
		return err
	}

	// Create a broker message
	msg := &broker.Message{
		Header: map[string]string{
			"UserKeyId": key.ApiKeyId,
		},
		Body: body,
	}

	log.Printf("UserKeyId: %s key: %s secret: %s exchange: %s", key.ApiKeyId, key.ApiKey, key.Secret, key.Exchange)
	log.Printf("%s", msg)

	// Publish message to broker
	if err := srv.PubSub.Publish(TopicNewKey, msg); err != nil {
		log.Printf("[pub] failed: %v", err)
	}

	return nil
}

func (service *KeyService) AddApiKey(ctx context.Context, req *keyProto.ApiKeyRequest, res *keyProto.ApiKeyResponse) error {
	apiKey, error := keyRepo.InsertApiKey(service.DB, req)

	switch {
	case error == nil:
		// verify the key in another process
		go service.publishEvent(ExchangeKey{
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
