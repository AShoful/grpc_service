package grpcserver

import (
	"grpc/proto"
	"grpc/server/pkg/handler"
	"grpc/server/pkg/service"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func RunServer(h *handler.Handler, s *service.Service) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(handler.UnaryAuthInterceptor(s)),
	)

	proto.RegisterUserServiceServer(grpcServer, h.AuthHandler)
	proto.RegisterBookServiceServer(grpcServer, h.BookHandler)

	// Запуск сервера в горутине
	go func() {
		log.Println("gRPC server running on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Ловим сигналы OS
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gRPC server...")

	grpcServer.GracefulStop()

	log.Println("gRPC server stopped gracefully")
}
