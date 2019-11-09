package server

import (
	"context"
	"log"
	"net"

	"github.com/adarosci/brzomessage-common/grpc/communication"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

var started bool

var registerAPI *FirebaseTokenAPIServer
var registerPaypal *FirebaseTokenPaypalServer
var registerAdmin *FirebaseTokenAdminServer

type server struct {
	communication.UnimplementedFirebaseTokenServer
}

// FirebaseTokenAPIServer struct
type FirebaseTokenAPIServer struct {
	Update func(key, token string, phones []string)
}

// FirebaseTokenAdminServer struct
type FirebaseTokenAdminServer struct {
	Update func(key, token string, phones []string)
}

// FirebaseTokenPaypalServer struct
type FirebaseTokenPaypalServer struct {
	Update func(key, token string, phones []string)
}

// UpdateApi update
func (s *server) UpdateApi(ctx context.Context, in *communication.UpdateFirebaseToken) (*communication.ResultMessages, error) {
	if registerAPI != nil {
		registerAPI.Update(in.GetKeyAccess(), in.GetFirebaseToken(), in.GetPhones())
	}
	return &communication.ResultMessages{Success: true}, nil
}

// UpdateApi update
func (s *server) UpdateAdmin(ctx context.Context, in *communication.UpdateFirebaseToken) (*communication.ResultMessages, error) {
	if registerAdmin != nil {
		registerAdmin.Update(in.GetKeyAccess(), in.GetFirebaseToken(), in.GetPhones())
	}
	return &communication.ResultMessages{Success: true}, nil
}

// UpdateApi update
func (s *server) UpdatePaypal(ctx context.Context, in *communication.UpdateFirebaseToken) (*communication.ResultMessages, error) {
	if registerPaypal != nil {
		registerPaypal.Update(in.GetKeyAccess(), in.GetFirebaseToken(), in.GetPhones())
	}
	return &communication.ResultMessages{Success: true}, nil
}

// Start serve
func (f *FirebaseTokenAPIServer) Start() *FirebaseTokenAPIServer {
	if registerAPI == nil {
		registerAPI = f
		start()
	}
	return f
}

// Start serve
func (f *FirebaseTokenPaypalServer) Start() *FirebaseTokenPaypalServer {
	if registerPaypal == nil {
		registerPaypal = f
		start()
	}
	return f
}

// Start serve
func (f *FirebaseTokenAdminServer) Start() *FirebaseTokenAdminServer {
	if registerAdmin == nil {
		registerAdmin = f
		start()
	}
	return f
}

func start() {
	if !started {
		started = true
		go func() {
			exit := make(chan bool)
			lis, err := net.Listen("tcp", port)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			s := grpc.NewServer()
			communication.RegisterFirebaseTokenServer(s, &server{})
			if err := s.Serve(lis); err != nil {
				log.Fatalf("failed to serve: %v", err)
			}
			<-exit
			<-exit
		}()
	}
}
