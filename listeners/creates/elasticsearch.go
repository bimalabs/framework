package creates

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	configs "github.com/KejawenLab/bima/v2/configs"
	events "github.com/KejawenLab/bima/v2/events"
	handlers "github.com/KejawenLab/bima/v2/handlers"
	elastic "github.com/olivere/elastic/v7"
)

type Elasticsearch struct {
	Env           *configs.Env
	Context       context.Context
	Elasticsearch *elastic.Client
	Logger        *handlers.Logger
}

func (c *Elasticsearch) Handle(event interface{}) interface{} {
	e := event.(*events.Model)
	m := e.Data.(configs.Model)
	data, _ := json.Marshal(e.Data)

	c.Elasticsearch.Index().Index(fmt.Sprintf("%s_%s", c.Env.Service.ConnonicalName, m.TableName())).BodyJson(string(data)).Do(c.Context)

	m.SetSyncedAt(time.Now())
	e.Repository.Update(m)

	return e
}

func (u *Elasticsearch) Listen() string {
	return events.AFTER_CREATE_EVENT
}

func (c *Elasticsearch) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
