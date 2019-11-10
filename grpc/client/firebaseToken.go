package client

import (
	"context"
	"fmt"
	"time"

	"github.com/adarosci/brzomessage-common/grpc/communication"
	"google.golang.org/grpc"
)

var conn *grpc.ClientConn

// FirebaseTokenCommunication communication
type FirebaseTokenCommunication struct{}

func init() {
	var err error
	go func() {
		conn, err = grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			fmt.Printf("did not connect: %v", err)
		}
	}()
}

// UpdateToken update token server
func (gc FirebaseTokenCommunication) UpdateToken(keyAccess, firebaseToken string, phones []string) error {
	c := communication.NewFirebaseTokenClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	param := &communication.UpdateFirebaseToken{
		Phones:        phones,
		KeyAccess:     keyAccess,
		FirebaseToken: firebaseToken,
	}
	_, err := c.UpdateApi(ctx, param)
	_, err = c.UpdatePaypal(ctx, param)
	return err
}
