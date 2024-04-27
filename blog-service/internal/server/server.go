package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jay-bhogayata/blogapi/internal/pb/blogservice"
	"github.com/jay-bhogayata/blogapi/internal/store"
	"google.golang.org/grpc"
)

type Server struct {
	blogservice.UnimplementedBlogServiceServer
}

// func NewServer(cfg *config.Config, router http.Handler) *Server {
// 	return &Server{
// 		httpServer: &http.Server{
// 			Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
// 			Handler: router,
// 		},
// 	}
// }

// func (s *Server) Start() error {

// 	logger.Log.Info("Starting the server...")

// 	err := s.httpServer.ListenAndServe()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func UseGRPCServer() {
	var grpcServer *grpc.Server
	var srv *Server
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		store.BlogStore.Logger.Info("server is starting on", "port", store.BlogStore.Config.Server.Port)
		store.BlogStore.Logger.Info("server is running in ", "environment", store.BlogStore.Config.Env)
		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", store.BlogStore.Config.Server.Port))
		if err != nil {
			store.BlogStore.Logger.Error("failed to listen", "error", err)
			os.Exit(1)
		}
		store.BlogStore.Logger.Info("server is listening on " + lis.Addr().String())

		grpcServer = grpc.NewServer()
		blogservice.RegisterBlogServiceServer(grpcServer, srv.UnimplementedBlogServiceServer)

		if err := grpcServer.Serve(lis); err != nil {
			store.BlogStore.Logger.Error("failed to serve", "error", err)
			os.Exit(1)
		}

	}()

	<-shutdown

	grpcServer.GracefulStop()
	store.BlogStore.Logger.Info("server stopped gracefully")
}
