package test

import (
	"context"

	constKey "github.com/asciiu/gomo/common/constants/key"
	protoKey "github.com/asciiu/gomo/key-service/proto/key"
	"github.com/micro/go-micro/client"
)

// Test clients of the Key service should use this client interface.
type mockKeyService struct {
}

func (m *mockKeyService) GetKeys(ctx context.Context, req *protoKey.GetKeysRequest, opts ...client.CallOption) (*protoKey.KeyListResponse, error) {
	keys := make([]*protoKey.Key, 0)
	keys = append(keys, &protoKey.Key{
		KeyID:       "examplekey",
		UserID:      "testuser",
		Exchange:    "testex",
		Key:         "test_key",
		Secret:      "test_secret",
		Status:      constKey.Verified,
		Description: "Test me!",
	})

	return &protoKey.KeyListResponse{
		Status: "success",
		Data: &protoKey.UserKeysData{
			Keys: keys,
		},
	}, nil
}

func (m *mockKeyService) AddKey(ctx context.Context, in *protoKey.KeyRequest, opts ...client.CallOption) (*protoKey.KeyResponse, error) {
	return &protoKey.KeyResponse{}, nil
}
func (m *mockKeyService) GetUserKey(ctx context.Context, in *protoKey.GetUserKeyRequest, opts ...client.CallOption) (*protoKey.KeyResponse, error) {
	return &protoKey.KeyResponse{}, nil
}
func (m *mockKeyService) GetUserKeys(ctx context.Context, in *protoKey.GetUserKeysRequest, opts ...client.CallOption) (*protoKey.KeyListResponse, error) {
	return &protoKey.KeyListResponse{}, nil
}
func (m *mockKeyService) RemoveKey(ctx context.Context, in *protoKey.RemoveKeyRequest, opts ...client.CallOption) (*protoKey.KeyResponse, error) {
	return &protoKey.KeyResponse{}, nil
}
func (m *mockKeyService) UpdateKeyDescription(ctx context.Context, in *protoKey.KeyRequest, opts ...client.CallOption) (*protoKey.KeyResponse, error) {
	return &protoKey.KeyResponse{}, nil
}

func MockKeyServiceClient() protoKey.KeyServiceClient {
	return new(mockKeyService)
}
