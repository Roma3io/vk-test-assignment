package pubsubservice

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"strconv"
	"strings"
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

func (s *PubSubServer) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(s.config.GRPCServer.Port))
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
	go func() {
		<-ctx.Done()
		s.log.Info("Shutting down gRPC server")
		s.grpcServer.GracefulStop()
		if err := s.bus.Close(context.Background()); err != nil {
			s.log.Error("Failed to close subpub", zap.Error(err))
		}
	}()
	return nil
}

func (s *PubSubServer) Subscribe(req *pb.SubscribeRequest, stream pb.PubSub_SubscribeServer) error {
	if req.Key == "" || strings.TrimSpace(req.Key) == "" {
		return status.Error(codes.InvalidArgument, "Key cannot be empty")
	}

	s.log.Info("New subscription", zap.String("key", req.Key))
	ch := make(chan interface{}, 10)
	sub, err := s.bus.Subscribe(req.Key, func(msg interface{}) {
		if data, ok := msg.(string); ok {
			ch <- data
		}
	})
	if err != nil {
		s.log.Error("Failed to subscribe", zap.Error(err))
		return status.Error(codes.Internal, "Subscription error")
	}
	defer sub.Unsubscribe()
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return status.Error(codes.Aborted, "Subscription closed")
			}
			if err := stream.Send(&pb.Event{Data: msg.(string)}); err != nil {
				s.log.Error("Couldn't send event", zap.Error(err))
				return err
			}
		case <-stream.Context().Done():
			s.log.Info("Client disconnected", zap.String("key", req.Key))
			return nil
		}
	}
}

func (s *PubSubServer) Publish(ctx context.Context, req *pb.PublishRequest) (*emptypb.Empty, error) {
	return nil, nil
}
