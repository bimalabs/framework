package dics

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	bima "github.com/KejawenLab/bima/v2"
	configs "github.com/KejawenLab/bima/v2/configs"
	drivers "github.com/KejawenLab/bima/v2/configs/drivers"
	events "github.com/KejawenLab/bima/v2/events"
	generators "github.com/KejawenLab/bima/v2/generators"
	handlers "github.com/KejawenLab/bima/v2/handlers"
	interfaces "github.com/KejawenLab/bima/v2/interfaces"
	creates "github.com/KejawenLab/bima/v2/listeners/creates"
	deletes "github.com/KejawenLab/bima/v2/listeners/deletes"
	filters "github.com/KejawenLab/bima/v2/listeners/paginations"
	updates "github.com/KejawenLab/bima/v2/listeners/updates"
	middlewares "github.com/KejawenLab/bima/v2/middlewares"
	paginations "github.com/KejawenLab/bima/v2/paginations"
	parsers "github.com/KejawenLab/bima/v2/parsers"
	routes "github.com/KejawenLab/bima/v2/routes"
	services "github.com/KejawenLab/bima/v2/services"
	utils "github.com/KejawenLab/bima/v2/utils"
	"github.com/ThreeDotsLabs/watermill"
	amqp "github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/fatih/color"
	"github.com/gadelkareem/cachita"
	"github.com/gertd/go-pluralize"
	elastic "github.com/olivere/elastic/v7"
	"github.com/sarulabs/dingo/v4"
	"github.com/sirupsen/logrus"
	mongodb "github.com/weekface/mgorus"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var Container = []dingo.Def{
	{
		Name:  "bima:config:parser:listener",
		Build: (*parsers.Listeners)(nil),
	},
	{
		Name:  "bima:config:parser:logger",
		Build: (*parsers.Logger)(nil),
	},
	{
		Name:  "bima:config:parser:middleware",
		Build: (*parsers.Middleware)(nil),
	},
	{
		Name:  "bima:config:parser:module",
		Build: (*parsers.Module)(nil),
	},
	{
		Name:  "bima:config:parser:route",
		Build: (*parsers.Route)(nil),
	},
	{
		Name:  "bima:config:user",
		Build: (*configs.User)(nil),
	},
	{
		Name: "bima:config:type",
		Build: func() (*configs.Type, error) {
			return &configs.Type{
				Map: map[string]string{
					"double":   "float64",
					"float":    "float32",
					"int32":    "int32",
					"int64":    "int64",
					"uint32":   "uint32",
					"uint64":   "uint64",
					"sint32":   "int32",
					"sint64":   "int64",
					"fixed32":  "uint32",
					"fixed64":  "uint64",
					"sfixed32": "int32",
					"sfixed64": "int64",
					"bool":     "bool",
					"string":   "string",
					"bytes":    "[]byte",
				},
			}, nil
		},
	},
	{
		Name:  "bima:config:template",
		Build: (*configs.Template)(nil),
	},
	{
		Name:  "bima:template:module",
		Build: (*configs.ModuleTemplate)(nil),
	},
	{
		Name:  "bima:template:field",
		Build: (*configs.FieldTemplate)(nil),
	},
	{
		Name: "bima:config:env",
		Build: func(
			user *configs.User,
			word *utils.Word,
		) (*configs.Env, error) {
			env := configs.Env{}

			env.Version = os.Getenv("APP_VERSION")
			env.ApiVersion = os.Getenv("API_VERSION")
			env.Debug, _ = strconv.ParseBool(os.Getenv("APP_DEBUG"))
			env.HtppPort, _ = strconv.Atoi(os.Getenv("APP_PORT"))
			env.RpcPort, _ = strconv.Atoi(os.Getenv("GRPC_PORT"))

			sName := os.Getenv("APP_NAME")
			env.Service = configs.Service{
				Name:           sName,
				ConnonicalName: word.Underscore(sName),
				Host:           os.Getenv("APP_HOST"),
			}

			dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
			env.Db = configs.Db{
				Host:     os.Getenv("DB_HOST"),
				Port:     dbPort,
				User:     os.Getenv("DB_USER"),
				Password: os.Getenv("DB_PASSWORD"),
				Name:     os.Getenv("DB_NAME"),
				Driver:   os.Getenv("DB_DRIVER"),
			}

			esPort, _ := strconv.Atoi(os.Getenv("ELASTICSEARCH_PORT"))
			env.Elasticsearch = configs.Elasticsearch{
				Host:  os.Getenv("ELASTICSEARCH_HOST"),
				Port:  esPort,
				Index: env.Db.Name,
			}

			mgdbPort, _ := strconv.Atoi(os.Getenv("MONGODB_PORT"))
			env.MongoDb = configs.MongoDb{
				Host:     os.Getenv("MONGODB_HOST"),
				Port:     mgdbPort,
				Database: "data_logs",
			}

			amqpPort, _ := strconv.Atoi(os.Getenv("AMQP_PORT"))
			env.Amqp = configs.Amqp{
				Host:     os.Getenv("AMQP_HOST"),
				Port:     amqpPort,
				User:     os.Getenv("AMQP_USER"),
				Password: os.Getenv("AMQP_PASSWORD"),
			}

			maxRole, _ := strconv.Atoi(os.Getenv("AUTH_HEADER_MAX_ROLE"))
			env.AuthHeader = configs.AuthHeader{
				Id:        os.Getenv("AUTH_HEADER_ID"),
				Email:     os.Getenv("AUTH_HEADER_EMAIL"),
				Role:      os.Getenv("AUTH_HEADER_ROLE"),
				Whitelist: os.Getenv("AUTH_HEADER_WHITELIST"),
				MaxRole:   maxRole,
			}

			env.CacheLifetime, _ = strconv.Atoi(os.Getenv("CACHE_LIFETIME"))
			env.User = user
			env.TemplateLocation = generators.TEMPLATE_PATH

			return &env, nil
		},
	},
	{
		Name: "bima:module:generator",
		Build: func(
			dic configs.Generator,
			model configs.Generator,
			module configs.Generator,
			proto configs.Generator,
			provider configs.Generator,
			server configs.Generator,
			validation configs.Generator,
			swagger configs.Generator,
			env *configs.Env,
			pluralizer *pluralize.Client,
			template *configs.Template,
			word *utils.Word,
		) (*generators.Factory, error) {
			return &generators.Factory{
				Env:        env,
				Pluralizer: pluralizer,
				Template:   template,
				Word:       word,
				Generators: []configs.Generator{
					dic,
					model,
					module,
					proto,
					provider,
					server,
					validation,
					swagger,
				},
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:generator:dic"),
			"1": dingo.Service("bima:generator:model"),
			"2": dingo.Service("bima:generator:module"),
			"3": dingo.Service("bima:generator:proto"),
			"4": dingo.Service("bima:generator:provider"),
			"5": dingo.Service("bima:generator:server"),
			"6": dingo.Service("bima:generator:validation"),
			"7": dingo.Service("bima:generator:swagger"),
		},
	},
	{
		Name:  "bima:generator:dic",
		Build: (*generators.Dic)(nil),
	},
	{
		Name:  "bima:generator:model",
		Build: (*generators.Model)(nil),
	},
	{
		Name:  "bima:generator:module",
		Build: (*generators.Module)(nil),
		Params: dingo.Params{
			"Config": dingo.Service("bima:config:parser:module"),
		},
	},
	{
		Name:  "bima:generator:proto",
		Build: (*generators.Proto)(nil),
	},
	{
		Name:  "bima:generator:provider",
		Build: (*generators.Provider)(nil),
	},
	{
		Name:  "bima:generator:server",
		Build: (*generators.Server)(nil),
	},
	{
		Name:  "bima:generator:validation",
		Build: (*generators.Validation)(nil),
	},
	{
		Name:  "bima:generator:swagger",
		Build: (*generators.Swagger)(nil),
	},
	{
		Name:  "bima:database:driver:mysql",
		Build: (*drivers.Mysql)(nil),
	},
	{
		Name:  "bima:database:driver:postgresql",
		Build: (*drivers.PostgreSql)(nil),
	},
	{
		Name: "bima:application",
		Build: func(
			database configs.Application,
			elasticsearch configs.Application,
			grpc configs.Application,
			queue configs.Application,
			rest configs.Application,
		) (*interfaces.Application, error) {
			return &interfaces.Application{
				Applications: []configs.Application{database, elasticsearch, grpc, queue, rest},
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:interface:database"),
			"1": dingo.Service("bima:interface:elasticsearch"),
			"2": dingo.Service("bima:interface:grpc"),
			"3": dingo.Service("bima:interface:queue"),
			"4": dingo.Service("bima:interface:rest"),
		},
	},
	{
		Name:  "bima:event:dispatcher",
		Build: (*events.Dispatcher)(nil),
	},
	{
		Name:  "bima:handler:middleware",
		Build: (*handlers.Middleware)(nil),
		Params: dingo.Params{
			"Dispatcher": dingo.Service("bima:event:dispatcher"),
		},
	},
	{
		Name:  "bima:logger:extension",
		Build: (*configs.LoggerExtension)(nil),
	},
	{
		Name: "bima:connection:database",
		Build: func(
			env *configs.Env,
			mysql configs.Driver,
			postgresql configs.Driver,
		) (*gorm.DB, error) {
			var db configs.Driver

			switch env.Db.Driver {
			case "mysql":
				db = mysql
			case "postgresql":
				db = postgresql
			default:
				return nil, errors.New("Unknown database driver")
			}

			util := color.New(color.FgCyan, color.Bold)

			util.Printf("✓ ")
			fmt.Printf("Database configured using '%s' driver...\n", env.Db.Driver)
			time.Sleep(100 * time.Millisecond)

			return db.Connect(
				env.Db.Host,
				env.Db.Port,
				env.Db.User,
				env.Db.Password,
				env.Db.Name,
				env.Debug,
			), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config:env"),
			"1": dingo.Service("bima:database:driver:mysql"),
			"2": dingo.Service("bima:database:driver:postgresql"),
		},
	},
	{
		Name: "bima:connection:elasticsearch",
		Build: func(env *configs.Env) (*elastic.Client, error) {
			client, err := elastic.NewClient(
				elastic.SetURL(fmt.Sprintf("%s:%d", env.Elasticsearch.Host, env.Elasticsearch.Port)),
				elastic.SetSniff(false),
				elastic.SetHealthcheck(false),
				elastic.SetGzip(true),
			)

			if err != nil {
				return nil, err
			}

			color.New(color.FgCyan, color.Bold).Printf("✓ ")
			fmt.Println("Elasticsearch configured...")
			time.Sleep(100 * time.Millisecond)

			return client, nil
		},
	},
	{
		Name:  "bima:listener:create:elasticsearch",
		Build: (*creates.Elasticsearch)(nil),
		Params: dingo.Params{
			"Env":           dingo.Service("bima:config:env"),
			"Context":       dingo.Service("bima:context:background"),
			"Elasticsearch": dingo.Service("bima:connection:elasticsearch"),
			"Logger":        dingo.Service("bima:handler:logger"),
		},
	},
	{
		Name:  "bima:listener:update:elasticsearch",
		Build: (*updates.Elasticsearch)(nil),
		Params: dingo.Params{
			"Env":           dingo.Service("bima:config:env"),
			"Context":       dingo.Service("bima:context:background"),
			"Elasticsearch": dingo.Service("bima:connection:elasticsearch"),
			"Logger":        dingo.Service("bima:handler:logger"),
		},
	},
	{
		Name:  "bima:listener:delete:elasticsearch",
		Build: (*deletes.Elasticsearch)(nil),
		Params: dingo.Params{
			"Env":           dingo.Service("bima:config:env"),
			"Context":       dingo.Service("bima:context:background"),
			"Elasticsearch": dingo.Service("bima:connection:elasticsearch"),
			"Logger":        dingo.Service("bima:handler:logger"),
		},
	},
	{
		Name:  "bima:listener:create:created_by",
		Build: (*creates.CreatedBy)(nil),
		Params: dingo.Params{
			"Env": dingo.Service("bima:config:env"),
		},
	},
	{
		Name:  "bima:listener:update:updated_by",
		Build: (*updates.UpdatedBy)(nil),
		Params: dingo.Params{
			"Env": dingo.Service("bima:config:env"),
		},
	},
	{
		Name:  "bima:listener:delete:deleted_by",
		Build: (*deletes.DeletedBy)(nil),
		Params: dingo.Params{
			"Env": dingo.Service("bima:config:env"),
		},
	},
	{
		Name:  "bima:listener:filter:elasticsearch",
		Build: (*filters.Filter)(nil),
	},
	{
		Name:  "bima:interface:database",
		Build: (*interfaces.Database)(nil),
	},
	{
		Name:  "bima:interface:elasticsearch",
		Build: (*interfaces.Elasticsearch)(nil),
	},
	{
		Name:  "bima:interface:grpc",
		Build: (*interfaces.GRpc)(nil),
		Params: dingo.Params{
			"Env":  dingo.Service("bima:config:env"),
			"GRpc": dingo.Service("bima:grpc:server"),
		},
	},
	{
		Name:  "bima:interface:queue",
		Build: (*interfaces.Queue)(nil),
	},
	{
		Name:  "bima:interface:rest",
		Build: (*interfaces.Rest)(nil),
		Params: dingo.Params{
			"Middleware": dingo.Service("bima:handler:middleware"),
			"Router":     dingo.Service("bima:handler:router"),
			"Server":     dingo.Service("bima:http:mux"),
			"Context":    dingo.Service("bima:context:background"),
		},
	},
	{
		Name: "bima:handler:logger",
		Build: func(
			env *configs.Env,
			logger *logrus.Logger,
			extension *configs.LoggerExtension,
		) (*handlers.Logger, error) {
			logger.SetFormatter(&logrus.JSONFormatter{})
			for _, e := range extension.Extensions {
				logger.AddHook(e)
			}

			return &handlers.Logger{
				Env:    env,
				Logger: logger,
			}, nil
		},
	},
	{
		Name:  "bima:handler:messager",
		Build: (*handlers.Messenger)(nil),
		Params: dingo.Params{
			"Logger":    dingo.Service("bima:handler:logger"),
			"Publisher": dingo.Service("bima:message:publisher"),
			"Consumer":  dingo.Service("bima:message:consumer"),
		},
	},
	{
		Name:  "bima:handler:handler",
		Build: (*handlers.Handler)(nil),
		Params: dingo.Params{
			"Env":           dingo.Service("bima:config:env"),
			"Context":       dingo.Service("bima:context:background"),
			"Elasticsearch": dingo.Service("bima:connection:elasticsearch"),
			"Dispatcher":    dingo.Service("bima:event:dispatcher"),
			"Repository":    dingo.Service("bima:service:repository"),
		},
	},
	{
		Name: "bima:handler:router",
		Build: func(
			gateway configs.Router,
			mux configs.Router,
		) (*handlers.Router, error) {
			return &handlers.Router{
				Routes: []configs.Router{
					gateway,
					mux,
				},
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:router:gateway"),
			"1": dingo.Service("bima:router:mux"),
		},
	},
	{
		Name:  "bima:middleware:auth",
		Build: (*middlewares.Auth)(nil),
		Params: dingo.Params{
			"Env": dingo.Service("bima:config:env"),
		},
	},
	{
		Name:  "bima:middleware:recovery",
		Build: (*middlewares.Recovery)(nil),
		Params: dingo.Params{
			"Logger": dingo.Service("bima:handler:logger"),
		},
	},
	{
		Name:  "bima:middleware:requestid",
		Build: (*middlewares.RequestID)(nil),
		Params: dingo.Params{
			"Logger": dingo.Service("bima:handler:logger"),
		},
	},
	{
		Name:  "bima:router:mux",
		Build: (*routes.MuxRouter)(nil),
	},
	{
		Name:  "bima:router:gateway",
		Build: (*routes.GRpcGateway)(nil),
	},
	{
		Name:  "bima:routes:api-doc",
		Build: (*routes.ApiDoc)(nil),
	},
	{
		Name:  "bima:routes:health",
		Build: (*routes.Health)(nil),
		Params: dingo.Params{
			"Logger": dingo.Service("bima:handler:logger"),
		},
	},
	{
		Name: "bima:http:mux",
		Build: func() (*http.ServeMux, error) {
			return http.NewServeMux(), nil
		},
	},
	{
		Name: "bima:grpc:server",
		Build: func() (*grpc.Server, error) {
			return grpc.NewServer(), nil
		},
	},
	{
		Name: "bima:log:logger",
		Build: func() (*logrus.Logger, error) {
			color.New(color.FgCyan, color.Bold).Printf("✓ ")
			fmt.Println("Logger configured...")
			time.Sleep(100 * time.Millisecond)

			return logrus.New(), nil
		},
	},
	{
		Name: "bima:logger:extension:mongodb",
		Build: func(env *configs.Env) (logrus.Hook, error) {
			color.New(color.FgCyan, color.Bold).Printf("✓ ")
			fmt.Println("MongoDB Logger Extension configured...")
			time.Sleep(100 * time.Millisecond)

			mongodb, err := mongodb.NewHooker(fmt.Sprintf("%s:%d", env.MongoDb.Host, env.MongoDb.Port), env.MongoDb.Database, "logs")
			if err != nil {
				return nil, err
			}

			return mongodb, nil
		},
	},
	{
		Name: "bima:context:background",
		Build: func() (context.Context, error) {
			return context.Background(), nil
		},
	},
	{
		Name: "bima:message:config",
		Build: func(env *configs.Env) (amqp.Config, error) {
			color.New(color.FgCyan, color.Bold).Printf("✓ ")
			fmt.Println("Pub/Sub configured...")
			time.Sleep(100 * time.Millisecond)

			return amqp.NewDurableQueueConfig(fmt.Sprintf("amqp://%s:%s@%s:%d/", env.Amqp.User, env.Amqp.Password, env.Amqp.Host, env.Amqp.Port)), nil
		},
	},
	{
		Name: "bima:message:publisher",
		Build: func(env *configs.Env, config amqp.Config) (*amqp.Publisher, error) {
			publisher, err := amqp.NewPublisher(config, watermill.NewStdLogger(env.Debug, env.Debug))
			if err != nil {
				return nil, err
			}

			return publisher, nil
		},
	},
	{
		Name: "bima:message:consumer",
		Build: func(env *configs.Env, config amqp.Config) (*amqp.Subscriber, error) {
			consumer, err := amqp.NewSubscriber(config, watermill.NewStdLogger(env.Debug, env.Debug))
			if err != nil {
				return nil, err
			}

			return consumer, nil
		},
	},
	{
		Name:  "bima:pagination:paginator",
		Build: (*paginations.Pagination)(nil),
	},
	{
		Name:  "bima:pagination:request",
		Build: (*paginations.Request)(nil),
	},
	{
		Name:  "bima:service:repository",
		Build: (*services.Repository)(nil),
		Params: dingo.Params{
			"Env":      dingo.Service("bima:config:env"),
			"Database": dingo.Service("bima:connection:database"),
		},
	},
	{
		Name:  "bima:cache:memory",
		Build: (*utils.Cache)(nil),
		Params: dingo.Params{
			"Pool": dingo.Service("bima:cachita:cache"),
		},
	},
	{
		Name:  "bima:util:number",
		Build: (*utils.Number)(nil),
	},
	{
		Name:  "bima:util:word",
		Build: (*utils.Word)(nil),
	},
	{
		Name: "bima:util:cli",
		Build: func() (*color.Color, error) {
			return color.New(color.FgCyan, color.Bold), nil
		},
	},
	{
		Name: "bima:util:pluralizer",
		Build: func() (*pluralize.Client, error) {
			return pluralize.NewClient(), nil
		},
	},
	{
		Name: "bima:cachita:cache",
		Build: func() (cachita.Cache, error) {
			return cachita.Memory(), nil
		},
	},
	{
		Name:  "bima:module",
		Build: (*bima.Module)(nil),
		Params: dingo.Params{
			"Context":       dingo.Service("bima:context:background"),
			"Elasticsearch": dingo.Service("bima:connection:elasticsearch"),
			"Handler":       dingo.Service("bima:handler:handler"),
			"Logger":        dingo.Service("bima:handler:logger"),
			"Messenger":     dingo.Service("bima:handler:messager"),
			"Cache":         dingo.Service("bima:cache:memory"),
			"Paginator":     dingo.Service("bima:pagination:paginator"),
		},
	},
	{
		Name:  "bima:server",
		Build: (*bima.Server)(nil),
		Params: dingo.Params{
			"Env":      dingo.Service("bima:config:env"),
			"Database": dingo.Service("bima:connection:database"),
		},
	},
	{
		Name: "bima:model",
		Build: func() (*bima.Model, error) {
			return &bima.Model{Base: configs.Base{}}, nil
		},
	},
}
