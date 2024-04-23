package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jay-bhogayata/notifyHub/docs"
	pb "github.com/jay-bhogayata/notifyHub/proto"
	"google.golang.org/grpc"
)

func (app *application) ServerInit() {
	var grpcServer *grpc.Server

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("server is starting on", "port", app.config.port)
		logger.Info("server is running in ", "environment", app.config.env)
		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", app.config.port))
		if err != nil {
			logger.Error("failed to listen", "error", err)
			os.Exit(1)
		}
		logger.Info("server is listening on " + lis.Addr().String())

		grpcServer = grpc.NewServer()
		pb.RegisterNotificationServiceServer(grpcServer, app)

		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("failed to serve", "error", err)
			os.Exit(1)
		}

	}()

	<-shutdown

	grpcServer.GracefulStop()
	logger.Info("server stopped gracefully")
}
