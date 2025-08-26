package client

import (
	"context"
	"fmt"

	"github.com/VicShved/pass-manager/client/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// // AuthorizationTokenName is name of header
// const AuthorizationTokenName string = "Authorization"
// const serverAddress string = "localhost:7777"

type GClient struct {
	conf config.ClientConfigStruct
}

var gClient GClient

func GetgClient(conf config.ClientConfigStruct) *GClient {
	gClient.conf = conf
	return &gClient
}

func (c GClient) getConnection() (*grpc.ClientConn, error) {
	// todo tls creds up
	fmt.Printf("server address %s", c.conf.ServerAddress)
	return grpc.NewClient(c.conf.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func (c GClient) addToken2Context(ctx context.Context, tokenStr string) context.Context {
	md := metadata.Pairs(config.AuthorizationTokenName, tokenStr)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx
}
