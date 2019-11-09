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

var register *FirebaseTokenServer

type server struct {
	communication.UnimplementedFirebaseTokenServer
}

// FirebaseTokenServer struct
type FirebaseTokenServer struct {
	UpdateAPI    func(key, token string, phones []string)
	UpdateAdmin  func(key, token string, phones []string)
	UpdatePaypal func(key, token string, phones []string)
}

// UpdateApi update
func (s *server) UpdateApi(ctx context.Context, in *communication.UpdateFirebaseToken) (*communication.ResultMessages, error) {
	if register.UpdateAPI != nil {
		register.UpdateAPI(in.GetKeyAccess(), in.GetFirebaseToken(), in.GetPhones())
	}
	return &communication.ResultMessages{Success: true}, nil
}

// UpdateApi update
func (s *server) UpdateAdmin(ctx context.Context, in *communication.UpdateFirebaseToken) (*communication.ResultMessages, error) {
	if register.UpdateAdmin != nil {
		register.UpdateAdmin(in.GetKeyAccess(), in.GetFirebaseToken(), in.GetPhones())
	}
	return &communication.ResultMessages{Success: true}, nil
}

// UpdateApi update
func (s *server) UpdatePaypal(ctx context.Context, in *communication.UpdateFirebaseToken) (*communication.ResultMessages, error) {
	if register.UpdatePaypal != nil {
		register.UpdatePaypal(in.GetKeyAccess(), in.GetFirebaseToken(), in.GetPhones())
	}
	return &communication.ResultMessages{Success: true}, nil
}

// Start serve
func (f *FirebaseTokenServer) Start() {
	if register == nil {
		register = f
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
