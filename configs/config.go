package configs

import (
	"context"

	"github.com/bimalabs/framework/v4/messengers"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/olivere/elastic/v7"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var Database *gorm.DB

type (
	Server interface {
		Register(server *grpc.Server)
		Handle(context context.Context, server *runtime.ServeMux, client *grpc.ClientConn) error
		Migrate(db *gorm.DB)
		Consume(messenger *messengers.Messenger)
		Sync(client *elastic.Client)
	}

	Db struct {
		Host     string `json:"host" yaml:"host"`
		User     string `json:"user" yaml:"user"`
		Password string `json:"password" yaml:"password"`
		Name     string `json:"name" yaml:"name"`
		Driver   string `json:"driver" yaml:"driver"`
		Port     int    `json:"port" yaml:"port"`
	}

	Env struct {
		Db            Db     `json:"database" yaml:"database"`
		Name          string `json:"name" yaml:"name"`
		Secret        string `json:"secret" yaml:"secret"`
		Service       string `json:"service" yaml:"service"`
		ApiPrefix     string `json:"api_prefix" yaml:"api_prefix"`
		User          string
		HttpPort      int  `json:"http_port" yaml:"http_port"`
		RpcPort       int  `json:"rpc_port" yaml:"rpc_port"`
		CacheLifetime int  `json:"cache_lifetime" yaml:"cache_lifetime"`
		Debug         bool `json:"debug" yaml:"debug"`
	}
)
