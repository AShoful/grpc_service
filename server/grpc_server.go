package grpcserver

import (
	"grpc/proto"
	"grpc/server/pkg/handler"
	"grpc/server/pkg/service"
	"log"
	"net"

	"google.golang.org/grpc"
)

func RunServer(h *handler.Handler, s *service.Service) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(handler.UnaryAuthInterceptor(s)))

	proto.RegisterUserServiceServer(grpcServer, h.AuthHandler)
	proto.RegisterBookServiceServer(grpcServer, h.BookHandler)

	log.Println("gRPC server running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
