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
)

func RunServer() {
	// Get app config
	var conf = config.GetServerConfig()

	// Init cu`stom logger
	logger.InitLogger(conf.LogLevel)

	// filestorage repo
	var fileRepo repository.FileStoragerRepoInterface
	fileRepo, err := repository.GetFileStorageRepo("")
	// repo choice
	var repo repository.RepoInterface
	repo, err = repository.GetGormRepo(conf, fileRepo)
	if err != nil {
		panic(err)
	}
	logger.Log.Info("Connect to db", zap.String("DSN", conf.DBDSN))

	// Bussiness layer
	serv := service.GetService(repo, conf)

	gserver, err := srv.GetServer(serv, conf)
	if err != nil {
		log.Fatal(err)
	}

	idleChan := make(chan string)
	exitChan := make(chan os.Signal, 3)
	signal.Notify(exitChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func(gserver *srv.GServer) {
		<-exitChan
		logger.Log.Info("Catch syscall sygnal")

		// Shutdown
		grp := errgroup.Group{}
		grp.Go(func() error {
			gserver.GracefulStop()
			return nil
		})
		if err := grp.Wait(); err != nil {
			logger.Log.Error("Server shuntdown: %v", zap.Error(err))
		}
		logger.Log.Info("Send message for shutdown gracefully")
		close(idleChan)
	}(gserver)

	listenAddress := ":" + conf.ServerPort
	logger.Log.Debug("RunServer", zap.String("listenAddress", listenAddress))
	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatal(err)
	}
	grp := errgroup.Group{}
	// run grpc server
	grp.Go(func() error {
		return gserver.StartServe(&lis, conf)
	})
	err = grp.Wait()

	// Shutdown gracefully
	<-idleChan
	gserver.GracefulStop()
	repo.CloseConn()
	logger.Log.Info("Server Shutdown gracefully")

}
