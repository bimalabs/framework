package interfaces

import (
	"context"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/bimalabs/framework/v4/configs"
	"github.com/bimalabs/framework/v4/loggers"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type GRpc struct {
	GRpcPort int
	Debug    bool
}

func (g *GRpc) Run(ctx context.Context, servers []configs.Server) {
	var gRpcAddress strings.Builder
	gRpcAddress.WriteString(":")
	gRpcAddress.WriteString(strconv.Itoa(g.GRpcPort))

	listen, err := net.Listen("tcp", gRpcAddress.String())
	if err != nil {
		log.Fatalf("Port %d is not available. %v", g.GRpcPort, err)
	}

	streams := make([]grpc.StreamServerInterceptor, 0, 2)
	unaries := make([]grpc.UnaryServerInterceptor, 0, 2)

	streams = append(streams, grpc_recovery.StreamServerInterceptor())
	unaries = append(unaries, grpc_recovery.UnaryServerInterceptor())
	if g.Debug {
		options := []grpc_logrus.Option{
			grpc_logrus.WithLevels(grpc_logrus.DefaultCodeToLevel),
		}
		streams = append(streams, grpc_logrus.StreamServerInterceptor(logrus.NewEntry(loggers.Logger.Engine), options...))
		unaries = append(unaries, grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(loggers.Logger.Engine), options...))
	}

	gRpc := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streams...)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaries...)),
	)

	for _, server := range servers {
		server.Register(gRpc)
	}

	_ = gRpc.Serve(listen)
}

func (g *GRpc) IsBackground() bool {
	return true
}

func (g *GRpc) Priority() int {
	return 257
}
