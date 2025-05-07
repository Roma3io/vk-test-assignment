package pubsubservice

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"vk-test-assignment/internal/config"
	proto "vk-test-assignment/internal/proto/gen"
	"vk-test-assignment/pkg/subpub"
)

type PubSubServer struct {
	proto.UnimplementedPubSubServer
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

func (s *PubSubServer) Subscribe() {}

func (s *PubSubServer) Publish() {}
