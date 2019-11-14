package client

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/adarosci/brzomessage-common/grpc/communication"
	"google.golang.org/grpc"
)

// FirebaseTokenCommunication communication
type FirebaseTokenCommunication struct{}

// UpdateToken update token server
func (gc FirebaseTokenCommunication) UpdateToken(keyAccess, firebaseToken string, phones []string) error {
	go func() {
		conn, err := grpc.Dial(os.Getenv("GRPC_HOST_API"), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			fmt.Printf("did not connect api: %v", err)
		}
		defer conn.Close()
		c := communication.NewFirebaseTokenClient(conn)
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
		conn, err := grpc.Dial(os.Getenv("GRPC_HOST_PAYPAL"), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			fmt.Printf("did not connect paypal: %v", err)
		}
		defer conn.Close()
		c := communication.NewFirebaseTokenClient(conn)
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
