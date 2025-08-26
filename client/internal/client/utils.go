package client

import (
	"context"
	"crypto/x509"
	"log"

	"github.com/VicShved/pass-manager/client/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type GClient struct {
	conf           config.ClientConfigStruct
	transportCreds credentials.TransportCredentials
}

var gClient GClient

func getTLSCreds(serverAddress string) credentials.TransportCredentials {
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("failed to load system root CAs: %v", err)
	}
	creds := credentials.NewClientTLSFromCert(systemRoots, serverAddress)
	return creds
}

func GetgClient(conf config.ClientConfigStruct) *GClient {
	gClient.conf = conf
	gClient.transportCreds = insecure.NewCredentials()
	if conf.EnableTLS {
		gClient.transportCreds = getTLSCreds(conf.ServerAddress)
	}
	return &gClient
}

func (c GClient) getConnection() (*grpc.ClientConn, error) {
	// todo tls creds up
	serverAddress := c.conf.ServerAddress + ":" + c.conf.ServerPort
	return grpc.NewClient(serverAddress, grpc.WithTransportCredentials(c.transportCreds))
}

func (c GClient) addToken2Context(ctx context.Context, tokenStr string) context.Context {
	md := metadata.Pairs(config.AuthorizationTokenName, tokenStr)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx
}
