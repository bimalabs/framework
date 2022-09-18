package routers

import (
	"context"

	"github.com/bimalabs/framework/v4/configs"
	"github.com/bimalabs/framework/v4/loggers"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type GRpcGateway struct {
	servers []configs.Server
}

func (g *GRpcGateway) Register(servers []configs.Server) {
	g.servers = servers
}

func (g *GRpcGateway) Handle(ctx context.Context, server *runtime.ServeMux, client *grpc.ClientConn) {
	cLog := context.WithValue(ctx, loggers.ScopeKey, "gateway")
	for _, handler := range g.servers {
		err := handler.Handle(ctx, server, client)
		if err != nil {
			loggers.Logger.Error(cLog, err.Error())

			break
		}
	}
}

func (a *GRpcGateway) Priority() int {
	return 255
}
