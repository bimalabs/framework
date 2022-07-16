package interfaces

import (
	"github.com/bimalabs/framework/v4/configs"
	"github.com/bimalabs/framework/v4/messengers"
	"golang.org/x/net/context"
)

type Consumer struct {
	Messenger *messengers.Messenger
}

func (q *Consumer) Run(ctx context.Context, servers []configs.Server) {
	if q.Messenger == nil {
		return
	}

	for _, server := range servers {
		go server.Consume(q.Messenger)
	}
}

func (q *Consumer) IsBackground() bool {
	return true
}

func (q *Consumer) Priority() int {
	return 0
}
