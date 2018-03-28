package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/antonlindstrom/pgstore"
	pb "github.com/asciiu/gomo/session-service/proto/session"
	micro "github.com/micro/go-micro"
)

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
type service struct {
	store *pgstore.PGStore
}

// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (s *service) CreateSession(ctx context.Context, req *pb.SessionRequest, res *pb.Response) error {

	// Get a session.
	sess, err := s.store.Get(nil, "session-key")
	if err != nil {
		log.Fatalf(err.Error())
	}

	session := pb.SessionResponse{
		SessionId: sess.ID,
	}

	res.Success = true
	res.Response = &session

	// Return matching the `Response` message we created in our
	// protobuf definition.
	return nil
}

func (s *service) GetSession(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	// Save our consignment
	session := pb.SessionResponse{
		SessionId: "id",
	}

	res.Success = true
	res.Response = &session
	// Return matching the `Response` message we created in our
	// protobuf definition.
	return nil
}

func (s *service) DeleteSession(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	// Save our consignment
	session := pb.SessionResponse{
		SessionId: "id",
	}

	res.Success = true
	res.Response = &session
	// Return matching the `Response` message we created in our
	// protobuf definition.
	return nil
}

// ExampleHandler is an example that displays the usage of PGStore.
// func ExampleHandler(w http.ResponseWriter, r *http.Request) {
// 	// Fetch new store.
// 	store, err := pgstore.NewPGStore("postgres://postgres@127.0.0.1:5432/gomo_dev?sslmode=disable", []byte("secret-key"))
// 	if err != nil {
// 		log.Fatalf(err.Error())
// 	}
// 	defer store.Close()

// 	// Run a background goroutine to clean up expired sessions from the database.
// 	defer store.StopCleanup(store.Cleanup(time.Minute * 5))

// 	// Get a session.
// 	session, err := store.Get(r, "session-key")
// 	if err != nil {
// 		log.Fatalf(err.Error())
// 	}

// 	// Add a value.
// 	session.Values["foo"] = "bar"

// 	// Save.
// 	if err = session.Save(r, w); err != nil {
// 		log.Fatalf("Error saving session: %v", err)
// 	}

// 	// Delete session.
// 	// session.Options.MaxAge = -1
// 	// if err = session.Save(r, w); err != nil {
// 	// 	log.Fatalf("Error saving session: %v", err)
// 	// }
// }

func main() {
	// Create a new service. Include some options here.
	srv := micro.NewService(
		// This name must match the package name given in your protobuf definition
		micro.Name("go.micro.srv.session"),
		micro.Version("latest"),
	)
	// Init will parse the command line flags.
	srv.Init()

	// TODO read DB from env var
	// TODO read secret from env var
	dbUrl := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	store, err := pgstore.NewPGStore(dbUrl, []byte("secret-key"))
	if err != nil {
		log.Fatalf(err.Error())
	}
	store.StopCleanup(store.Cleanup(time.Minute * 5))

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	pb.RegisterSessionServiceHandler(srv.Server(), &service{store})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
