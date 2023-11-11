package interfaces

import (
	"sort"
	"time"

	"github.com/bimalabs/framework/v4/configs"
	"golang.org/x/net/context"
)

type (
	Application interface {
		Run(ctx context.Context, servers []configs.Server)
		IsBackground() bool
		Priority() int
	}

	Factory struct {
		applications []Application
	}
)

func (f *Factory) Register(applications ...Application) {
	for _, application := range applications {
		f.Add(application)
	}
}

func (f *Factory) Add(application Application) {
	f.applications = append(f.applications, application)
}

func (f *Factory) Run(servers []configs.Server) {
	sort.Slice(f.applications, func(i int, j int) bool {
		return f.applications[i].Priority() > f.applications[j].Priority()
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, application := range f.applications {
		if !application.IsBackground() {
			time.Sleep(100 * time.Millisecond)
			application.Run(ctx, servers)

			continue
		}

		go application.Run(ctx, servers)
	}
}
