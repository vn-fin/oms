package remote

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/vn-fin/oms/internal/config"
	pb "github.com/vn-fin/xpb/xpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var AuthGrpcClient pb.PermissionServiceClient
var grpcConn *grpc.ClientConn

func InitAuthGrpcClient() error {
	grpcHost := config.PermissionGrpcHost
	var err error
	grpcConn, err = grpc.NewClient(
		grpcHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to auth gRPC: %w", err)
	}

	AuthGrpcClient = pb.NewPermissionServiceClient(grpcConn)
	log.Info().Msgf("Connected to auth gRPC service at %s", grpcHost)

	return nil
}

func CloseAuthGrpcClient() {
	if grpcConn != nil {
		grpcConn.Close()
		log.Info().Msg("Closed auth gRPC connection")
	}
}
