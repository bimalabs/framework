package dics

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	bima "github.com/bimalabs/framework/v4"
	"github.com/bimalabs/framework/v4/configs"
	"github.com/bimalabs/framework/v4/drivers"
	"github.com/bimalabs/framework/v4/events"
	"github.com/bimalabs/framework/v4/handlers"
	"github.com/bimalabs/framework/v4/interfaces"
	"github.com/bimalabs/framework/v4/loggers"
	"github.com/bimalabs/framework/v4/middlewares"
	"github.com/bimalabs/framework/v4/models"
	paginations "github.com/bimalabs/framework/v4/paginations"
	"github.com/bimalabs/framework/v4/routers"
	"github.com/bimalabs/framework/v4/routes"
	"github.com/bimalabs/framework/v4/utils"
	"github.com/fatih/color"
	"github.com/kamva/mgm/v3"
	"github.com/sarulabs/dingo/v4"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Application = []dingo.Def{
	{
		Name:  "bima:application",
		Scope: bima.Application,
		Build: func(
			env *configs.Env,
			driver *drivers.Factory,
			extension *loggers.LoggerExtension,
		) (*interfaces.Factory, error) {
			loggers.Configure(env.Debug, env.Service, *extension)
			factory := interfaces.Factory{}
			if env.Db.Driver == "" {
				return &factory, nil
			}

			util := color.New(color.FgGreen)
			util.Print("✓ ")
			fmt.Print("Database configured using ")
			util.Print(env.Db.Driver)
			fmt.Println(" driver")

			switch env.Db.Driver {
			case "mongo":
				var dsn strings.Builder

				dsn.WriteString("mongodb://")
				dsn.WriteString(env.Db.User)
				dsn.WriteString(":")
				dsn.WriteString(env.Db.Password)
				dsn.WriteString("@")
				dsn.WriteString(env.Db.Host)
				dsn.WriteString(":")
				dsn.WriteString(strconv.Itoa(env.Db.Port))

				err := mgm.SetDefaultConfig(nil, env.Db.Name, options.Client().ApplyURI(dsn.String()).SetMonitor(&event.CommandMonitor{
					Started: func(_ context.Context, evt *event.CommandStartedEvent) {
						log.Print(evt.Command)
					},
				}))
				if err != nil {
					err = mgm.SetDefaultConfig(nil, env.Db.Name, options.Client().ApplyURI(env.Db.Host).SetMonitor(&event.CommandMonitor{
						Started: func(_ context.Context, evt *event.CommandStartedEvent) {
							log.Print(evt.Command)
						},
					}))
					if err != nil {
						log.Fatalln(err.Error())
					}
				}

				return &factory, nil
			default:
				configs.Database = driver.Connect(env.Db.Driver, env.Db.Host, env.Db.Port, env.Db.User, env.Db.Password, env.Db.Name)
			}

			return &factory, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config"),
			"1": dingo.Service("bima:driver:factory"),
			"2": dingo.Service("bima:logger:extension"),
		},
	},
	{
		Name:  "bima:config",
		Scope: bima.Application,
		Build: (*configs.Env)(nil),
	},
	{
		Name:  "bima:driver:factory",
		Scope: bima.Application,
		Build: func(env *configs.Env) (*drivers.Factory, error) {
			return drivers.New(env.Debug), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config"),
		},
	},
	{
		Name:  "bima:event:dispatcher",
		Scope: bima.Application,
		Build: (*events.Dispatcher)(nil),
	},
	{
		Name:  "bima:middleware:factory",
		Scope: bima.Application,
		Build: func(env *configs.Env) (*middlewares.Factory, error) {
			middleware := middlewares.Factory{Debug: env.Debug}
			middleware.Add(&middlewares.Header{})

			return &middleware, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config"),
		},
	},
	{
		Name:  "bima:logger:extension",
		Scope: bima.Application,
		Build: (*loggers.LoggerExtension)(nil),
	},
	{
		Name:  "bima:interface:rest",
		Scope: bima.Application,
		Build: func(
			env *configs.Env,
			middleware *middlewares.Factory,
			router *routers.Factory,
		) (*interfaces.Rest, error) {
			return &interfaces.Rest{
				HttpPort:   env.HttpPort,
				Middleware: middleware,
				Router:     router,
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config"),
			"1": dingo.Service("bima:middleware:factory"),
			"2": dingo.Service("bima:router:factory"),
		},
	},
	{
		Name:  "bima:router:factory",
		Scope: bima.Application,
		Build: func(gateway routers.Router, mux routers.Router) (*routers.Factory, error) {
			return &routers.Factory{Routers: []routers.Router{gateway, mux}}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:router:gateway"),
			"1": dingo.Service("bima:router:mux"),
		},
	},
	{
		Name:  "bima:router:mux",
		Scope: bima.Application,
		Build: func(
			env *configs.Env,
			apiDoc routes.Route,
			apiDocRedirection routes.Route,
			health routes.Route,
		) (*routers.MuxRouter, error) {
			routers := routers.MuxRouter{Debug: env.Debug, ApiPrefix: env.ApiPrefix}
			routers.Register([]routes.Route{apiDoc, apiDocRedirection, health})

			return &routers, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config"),
			"1": dingo.Service("bima:route:api-doc"),
			"2": dingo.Service("bima:route:api-doc-redirect"),
			"3": dingo.Service("bima:route:health"),
		},
	},
	{
		Name:  "bima:router:gateway",
		Scope: bima.Application,
		Build: (*routers.GRpcGateway)(nil),
	},
	{
		Name:  "bima:route:api-doc",
		Scope: bima.Application,
		Build: func(env *configs.Env) (*routes.ApiDoc, error) {
			return &routes.ApiDoc{Debug: env.Debug, ApiPrefix: env.ApiPrefix}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config"),
		},
	},
	{
		Name:  "bima:route:api-doc-redirect",
		Scope: bima.Application,
		Build: (*routes.ApiDocRedirect)(nil),
	},
	{
		Name:  "bima:route:health",
		Scope: bima.Application,
		Build: (*routes.Health)(nil),
	},
	{
		Name:  "bima:cache:memory",
		Scope: bima.Application,
		Build: func(env *configs.Env) (*utils.Cache, error) {
			return utils.NewCache(time.Duration(env.CacheLifetime) * time.Second), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config"),
		},
	},
	{
		Name:  "bima:validator",
		Scope: bima.Application,
		Build: func(env *configs.Env, dispatcher *events.Dispatcher) (utils.Validator, error) {
			return utils.NewValidator(env.Debug, dispatcher), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config"),
			"1": dingo.Service("bima:event:dispatcher"),
		},
	},
	{
		Name:  "bima:module",
		Scope: bima.Application,
		Build: func(
			env *configs.Env,
			handler handlers.Handler,
			cache *utils.Cache,
			validator utils.Validator,
		) (bima.Module, error) {
			return bima.NewModule(env.Debug, handler, cache, validator, &paginations.Pagination{}), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config"),
			"1": dingo.Service("bima:handler"),
			"2": dingo.Service("bima:cache:memory"),
			"3": dingo.Service("bima:validator"),
		},
	},
	{
		Name:  "bima:server",
		Scope: bima.Application,
		Build: func(env *configs.Env) (*bima.Server, error) {
			return &bima.Server{Debug: env.Debug}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config"),
		},
	},
	{
		Name:  "bima:model",
		Scope: bima.Application,
		Build: func(env *configs.Env) (*bima.GormModel, error) {
			return &bima.GormModel{
				GormBase: models.GormBase{Env: env},
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config"),
		},
	},
}
