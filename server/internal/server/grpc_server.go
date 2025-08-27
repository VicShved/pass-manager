package server

import (
	"crypto/tls"
	"net"
	"os"

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
	serv        *service.PassManageService
	server      *grpc.Server
	listener    *net.Listener
	certManager *autocert.Manager
}

func getCertManager(serverAddress string) *autocert.Manager {
	manager := autocert.Manager{
		Cache:      autocert.DirCache("certs"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(serverAddress),
	}
	return &manager
}

func getTLSCreds(manager *autocert.Manager, logFile *os.File) grpc.ServerOption {
	// tlsConfig := manager.TLSConfig()
	tlsConfig := manager.TLSConfig()
	if logFile != nil {
		tlsConfig.KeyLogWriter = logFile
	}
	// Create gRPC transport credentials
	creds := credentials.NewTLS(tlsConfig)
	return grpc.Creds(creds)

}

// GetServer return grpc.Server with serverOptions
func GetServer(serv *service.PassManageService, conf *config.ServerConfigStruct) (*GServer, error) {
	gServer := GServer{serv: serv}
	var serverOptions []grpc.ServerOption
	serverOptions = append(serverOptions, grpc.ChainUnaryInterceptor(AuthUnaryInterceptor))
	serverOptions = append(serverOptions, grpc.ChainStreamInterceptor(AuthStreamInterceptor))
	if conf.EnableTLS {
		m := getCertManager(conf.ServerAddress)
		gServer.certManager = m
		var tlsLogFile *os.File
		if conf.LogLevel == "DEBUG" {
			tlsLogFile, _ = os.OpenFile("tls_debug.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		}
		serverOptions = append(serverOptions, getTLSCreds(m, tlsLogFile))
		logger.Log.Debug("GetServer: Add tls creds to server")
	}
	server := grpc.NewServer(
		serverOptions...,
	// creds,
	// grpc.ChainUnaryInterceptor(AuthUnaryInterceptor),
	// grpc.ChainStreamInterceptor(AuthStreamInterceptor),
	// keepAlive,
	// grpc.MaxRecvMsgSize(1024*1024*1024),
	// grpc.MaxSendMsgSize(1024*1024*1024),
	// grpc.ConnectionTimeout(60000),
	)

	gServer.server = server

	pb.RegisterPassManagerServiceServer(server, gServer)
	return &gServer, nil
}

// StartServe run server to serve
func (s *GServer) StartServe(lis *net.Listener, conf *config.ServerConfigStruct) error {
	logger.Log.Debug("StartServe", zap.String("conf.ServerAddress", conf.ServerAddress), zap.String("port", conf.ServerPort))
	s.listener = lis
	if s.certManager != nil {
		newListener := tls.NewListener(*lis, s.certManager.TLSConfig())
		s.listener = &newListener
	}
	return s.server.Serve(*s.listener)
}

// GracefulStop stop server and listener gracefully
func (s *GServer) GracefulStop() {
	(*s.listener).Close()
	s.server.GracefulStop()

}
