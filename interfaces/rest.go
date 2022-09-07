package interfaces

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/bimalabs/framework/v4/configs"
	"github.com/bimalabs/framework/v4/middlewares"
	"github.com/bimalabs/framework/v4/routers"
	"google.golang.org/grpc"
)

type Rest struct {
	HttpPort   int
	Middleware *middlewares.Factory
	Router     *routers.Factory
	GRpcClient *grpc.ClientConn
}

func (r *Rest) Run(ctx context.Context, servers []configs.Server) {
	r.Middleware.Sort()
	r.Router.Sort()

	var httpAddress strings.Builder
	httpAddress.WriteString(":")
	httpAddress.WriteString(strconv.Itoa(r.HttpPort))

	_ = http.ListenAndServe(httpAddress.String(), r.Middleware.Attach(r.Router.Handle(context.Background(), http.NewServeMux(), r.GRpcClient)))
}

func (r *Rest) IsBackground() bool {
	return false
}

func (r *Rest) Priority() int {
	return -253
}
