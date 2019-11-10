package server

import (
	"context"
	"log"
	"net"

	"github.com/adarosci/brzomessage-common/grpc/communication"
	"google.golang.org/grpc"
)

const (
	portAPI    = ":50051"
	portPaypal = ":50052"
)

var started bool

var registerAPI *FirebaseTokenAPIServer
var registerPaypal *FirebaseTokenPaypalServer

type serverAPI struct {
	communication.UnimplementedFirebaseTokenServer
}

type serverPaypal struct {
	communication.UnimplementedFirebaseTokenServer
}

// FirebaseTokenAPIServer struct
type FirebaseTokenAPIServer struct {
	Update func(key, token string, phones []string)
}

// FirebaseTokenPaypalServer struct
type FirebaseTokenPaypalServer struct {
	Update func(key, token string, phones []string)
}

// UpdateApi update
func (s *serverAPI) UpdateApi(ctx context.Context, in *communication.UpdateFirebaseToken) (*communication.ResultMessages, error) {
	if registerAPI != nil {
		registerAPI.Update(in.GetKeyAccess(), in.GetFirebaseToken(), in.GetPhones())
	}
	return &communication.ResultMessages{Success: true}, nil
}

// UpdateApi update
func (s *serverPaypal) UpdatePaypal(ctx context.Context, in *communication.UpdateFirebaseToken) (*communication.ResultMessages, error) {
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

func start() {
	if !started {
		started = true
		go func() {
			lis, err := net.Listen("tcp", portAPI)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			s := grpc.NewServer()
			communication.RegisterFirebaseTokenServer(s, &serverAPI{})
			if err := s.Serve(lis); err != nil {
				log.Fatalf("failed to serve: %v", err)
			}
		}()
		go func() {
			lis, err := net.Listen("tcp", portPaypal)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			s := grpc.NewServer()
			communication.RegisterFirebaseTokenServer(s, &serverPaypal{})
			if err := s.Serve(lis); err != nil {
				log.Fatalf("failed to serve: %v", err)
			}
		}()
	}
}
