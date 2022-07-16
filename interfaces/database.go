package interfaces

import (
	"github.com/bimalabs/framework/v4/configs"
	"golang.org/x/net/context"
)

type Database struct {
}

func (d *Database) Run(ctx context.Context, servers []configs.Server) {
	if configs.Database == nil {
		return
	}

	for _, server := range servers {
		go server.Migrate(configs.Database)
	}
}

func (d *Database) IsBackground() bool {
	return true
}

func (d *Database) Priority() int {
	return 0
}
