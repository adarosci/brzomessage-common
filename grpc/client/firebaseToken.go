package client

import (
	"context"
	"fmt"
	"time"

	"github.com/adarosci/brzomessage-common/grpc/communication"
	"google.golang.org/grpc"
)

var connAPI *grpc.ClientConn
var connPaypal *grpc.ClientConn

// FirebaseTokenCommunication communication
type FirebaseTokenCommunication struct{}

func init() {
	var err error
	go func() {
		connAPI, err = grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			fmt.Printf("did not connect api: %v", err)
		}
	}()
	go func() {
		connPaypal, err = grpc.Dial("localhost:50052", grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			fmt.Printf("did not connect paypal: %v", err)
		}
	}()
}

// UpdateToken update token server
func (gc FirebaseTokenCommunication) UpdateToken(keyAccess, firebaseToken string, phones []string) error {
	go func() {
		c := communication.NewFirebaseTokenClient(connAPI)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		param := &communication.UpdateFirebaseToken{
			Phones:        phones,
			KeyAccess:     keyAccess,
			FirebaseToken: firebaseToken,
		}
		c.UpdateApi(ctx, param)
	}()
	go func() {
		c := communication.NewFirebaseTokenClient(connPaypal)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		param := &communication.UpdateFirebaseToken{
			Phones:        phones,
			KeyAccess:     keyAccess,
			FirebaseToken: firebaseToken,
		}
		c.UpdatePaypal(ctx, param)
	}()
	return nil
}
