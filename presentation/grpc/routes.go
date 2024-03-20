package grpc

import (
	pb "github.com/kingstonduy/mcs-money-transfer/presentation/grpc/pb"
	"github.com/kingstonduy/mcs-money-transfer/presentation/grpc/server"
)

func (g *GrpcServer) WithAccountServer() GrpcServerStartOption {
	return func(s *GrpcServer) error {
		tSrv := server.NewAccountServer()
		pb.RegisterAccountServiceServer(s.gSrv, tSrv)
		return nil
	}
}
