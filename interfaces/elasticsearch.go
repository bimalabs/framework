package interfaces

import (
	"github.com/bimalabs/framework/v4/configs"
	"github.com/olivere/elastic/v7"
	"golang.org/x/net/context"
)

type Elasticsearch struct {
	Client *elastic.Client
}

func (e *Elasticsearch) Run(ctx context.Context, servers []configs.Server) {
	if e.Client == nil {
		return
	}

	for _, server := range servers {
		go server.Sync(e.Client)
	}
}

func (e *Elasticsearch) IsBackground() bool {
	return true
}

func (e *Elasticsearch) Priority() int {
	return 0
}
