package server

import (
	"github.com/VicShved/pass-manager/server/internal/service"
	pb "github.com/VicShved/pass-manager/server/pkg/api/proto"
	"github.com/VicShved/pass-manager/server/pkg/config"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

// GServer
type GServer struct {
	pb.UnimplementedPassManagerServer
	serv *service.Service
	conf *config.ServerConfigStruct
}

func getTLSCreds(domain string) grpc.ServerOption {
	manager := autocert.Manager{
		Cache:      autocert.DirCache("certs"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
	}
	tlsConfig := manager.TLSConfig()
	// Create gRPC transport credentials
	creds := credentials.NewTLS(tlsConfig)
	return grpc.Creds(creds)

}

func GetServer(serv *service.Service, conf *config.ServerConfigStruct) (*grpc.Server, error) {
	keepAlive := grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionAgeGrace: 84000})
	creds := getTLSCreds(conf.BaseURL)
	server := grpc.NewServer(
		creds,
		grpc.ChainUnaryInterceptor(AuthUnaryInterceptor),
		keepAlive,
		grpc.MaxRecvMsgSize(1024*1024*1024),
		grpc.MaxSendMsgSize(1024*1024*1024),
		grpc.ConnectionTimeout(60000),
	)
	gServer := GServer{serv: serv, conf: conf}
	pb.RegisterPassManagerServer(server, &gServer)
	return server, nil
}
