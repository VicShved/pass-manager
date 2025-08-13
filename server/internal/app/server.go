package app

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/VicShved/pass-manager/server/internal/repository"
	srv "github.com/VicShved/pass-manager/server/internal/server"
	"github.com/VicShved/pass-manager/server/internal/service"
	"github.com/VicShved/pass-manager/server/pkg/config"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func ServerRun() {
	// Get app config
	var conf = config.GetServerConfig()

	// Init cu`stom logger
	logger.InitLogger(conf.LogLevel)

	// repo choice
	var repo repository.RepoInterface
	repo, err := repository.GetGormRepo(conf)
	if err != nil {
		panic(err)
	}
	logger.Log.Info("Connect to db", zap.String("DSN", conf.DBDSN))

	// Bussiness layer
	serv := service.GetService(repo, conf)

	server, err := srv.GetServer(serv, conf)
	if err != nil {
		log.Fatal(err)
	}

	idleChan := make(chan string)
	exitChan := make(chan os.Signal, 3)
	signal.Notify(exitChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func(server *grpc.Server) {
		<-exitChan
		logger.Log.Info("Catch syscall sygnal")

		// Shutdown
		grp := errgroup.Group{}
		grp.Go(func() error {
			server.GracefulStop()
			return nil
		})
		if err := grp.Wait(); err != nil {
			logger.Log.Error("Server shuntdown: %v", zap.Error(err))
		}
		logger.Log.Info("Send message for shutdown gracefully")
		close(idleChan)
	}(server)

	lis, err := net.Listen("tcp", ":443")
	if err != nil {
		log.Fatal(err)
	}
	grp := errgroup.Group{}
	// run grpc server
	grp.Go(func() error {
		return server.Serve(lis)
	})
	err = grp.Wait()

	// Shutdown gracefully
	<-idleChan
	repo.CloseConn()
	logger.Log.Info("Server Shutdown gracefully")

}
