package server

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/adarosci/brzomessage-common/grpc/communication"
	"google.golang.org/grpc"
)

var registerEmail *EmailSendServer

type serverEmail struct {
	communication.UnimplementedEmailServiceServer
}

func (*serverEmail) ConfirmEmail(ctx context.Context, req *communication.UserInfoRequest) (*communication.ResultMessages, error) {
	err := registerEmail.SendConfirmEmail(req.GetPersonID())
	return &communication.ResultMessages{}, err
}

func (*serverEmail) ResetPassword(ctx context.Context, req *communication.UserInfoRequest) (*communication.ResultMessages, error) {
	err := registerEmail.SendResetPassword(req.GetPersonID())
	return &communication.ResultMessages{}, err
}

func (*serverEmail) PurchaseCredits(ctx context.Context, req *communication.PurchaseCreditsRequest) (*communication.ResultMessages, error) {
	err := registerEmail.SendPurchaseCredits(req.GetID(), req.GetPhoneNumber(), req.GetDate(), req.GetPlanSelect(), req.GetDescription(), req.GetValue())
	return &communication.ResultMessages{}, err
}

// EmailSendServer server
type EmailSendServer struct {
	SendConfirmEmail    func(personID int32) error
	SendResetPassword   func(personID int32) error
	SendPurchaseCredits func(id, phoneNumber, date, planSelect, description string, value int32) error
}

// Start serve
func (f *EmailSendServer) Start() *EmailSendServer {
	if registerEmail == nil {
		registerEmail = f
		go func() {
			lis, err := net.Listen("tcp", os.Getenv("GRPC_HOST_API"))
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			s := grpc.NewServer()
			communication.RegisterEmailServiceServer(s, &serverEmail{})
			if err := s.Serve(lis); err != nil {
				log.Fatalf("failed to serve: %v", err)
			}
		}()
	}
	return f
}
