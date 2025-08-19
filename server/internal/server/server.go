package server

import (
	"net"

	"github.com/VicShved/pass-manager/server/internal/service"
	pb "github.com/VicShved/pass-manager/server/pkg/api/proto"
	"github.com/VicShved/pass-manager/server/pkg/config"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// GServer
type GServer struct {
	pb.UnimplementedPassManagerServiceServer
	serv     *service.PassManageService
	server   *grpc.Server
	listener *net.Listener
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

func GetServer(serv *service.PassManageService, conf *config.ServerConfigStruct) (*GServer, error) {
	// keepAlive := grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionAgeGrace: 84000})
	// creds := getTLSCreds(conf.BaseURL)
	server := grpc.NewServer(
		// creds,
		grpc.ChainUnaryInterceptor(AuthUnaryInterceptor),
		// grpc.ChainStreamInterceptor(AuthStreamInterceptor),
		// keepAlive,
		// grpc.MaxRecvMsgSize(1024*1024*1024),
		// grpc.MaxSendMsgSize(1024*1024*1024),
		// grpc.ConnectionTimeout(60000),
	)
	gServer := GServer{serv: serv, server: server}
	pb.RegisterPassManagerServiceServer(server, gServer)
	return &gServer, nil
}

func (s *GServer) StartServe(lis *net.Listener, conf *config.ServerConfigStruct) error {
	logger.Log.Debug("StartServe", zap.String("conf.ServerAddress", conf.ServerAddress))
	s.listener = lis
	return s.server.Serve(*lis)
}

func (s *GServer) GracefulStop() {
	(*s.listener).Close()
	s.server.GracefulStop()

}
