package pubsubservice

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"strconv"
	"vk-test-assignment/internal/config"
	pb "vk-test-assignment/internal/proto/gen"
	"vk-test-assignment/pkg/subpub"
)

type PubSubServer struct {
	pb.UnimplementedPubSubServer
	bus        subpub.SubPub
	config     *config.Config
	grpcServer *grpc.Server
	log        *zap.Logger
}

func NewPubSubServer(bus subpub.SubPub, config *config.Config, log *zap.Logger) *PubSubServer {
	return &PubSubServer{
		bus:    bus,
		config: config,
		log:    log,
	}
}

func (s *PubSubServer) Start() error {
	listener, err := net.Listen("tcp", strconv.Itoa(s.config.GRPCServer.Port))
	if err != nil {
		s.log.Fatal("Error making listener", zap.Error(err))
		return err
	}
	s.grpcServer = grpc.NewServer()
	pb.RegisterPubSubServer(s.grpcServer, s)
	s.log.Info("starting gRPC server", zap.Int("port", s.config.GRPCServer.Port))
	go func() {
		if err := s.grpcServer.Serve(listener); err != nil {
			s.log.Error("gRPC server error", zap.Error(err))
		}
	}()
	return nil
}

func (s *PubSubServer) Subscribe(req *pb.SubscribeRequest, stream pb.PubSub_SubscribeServer) error {
	return nil
}

func (s *PubSubServer) Publish(ctx context.Context, req *pb.PublishRequest) (*emptypb.Empty, error) {
	return nil, nil
}
