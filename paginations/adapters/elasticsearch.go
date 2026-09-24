package adapters

import (
	"context"
	"strings"

	"github.com/goccy/go-json"

	"github.com/bimalabs/framework/v4/events"
	"github.com/bimalabs/framework/v4/loggers"
	"github.com/bimalabs/framework/v4/paginations"
	"github.com/olivere/elastic/v7"
	"github.com/vcraescu/go-paginator/v2"
)

type (
	ElasticsearchAdapter struct {
		Client     *elastic.Client
		Dispatcher *events.Dispatcher
		Service    string
		Debug      bool
	}

	elasticsearchPaginator struct {
		context    context.Context
		model      any
		client     *elastic.Client
		pageQuery  *elastic.BoolQuery
		totalQuery *elastic.BoolQuery
		index      string
	}
)

func (es *ElasticsearchAdapter) CreateAdapter(ctx context.Context, paginator paginations.Pagination) paginator.Adapter {
	if es.Client == nil {
		loggers.Logger.Error(ctx, "adapter not configured properly")

		return nil
	}

	event := events.ElasticsearchPagination{
		Model:   paginator.Model,
		Query:   elastic.NewBoolQuery(),
		Filters: paginator.Filters,
	}

	if es.Debug {
		var log strings.Builder
		log.WriteString("dispatching ")
		log.WriteString(events.PaginationEvent.String())

		loggers.Logger.Debug(ctx, log.String())
	}

	var index strings.Builder

	index.WriteString(es.Service)
	index.WriteString("_")
	index.WriteString(paginator.Table)

	_ = es.Dispatcher.Dispatch(events.PaginationEvent.String(), &event)

	return newElasticsearchPaginator(ctx, es.Client, index.String(), paginator.Model, event.Query)
}

func newElasticsearchPaginator(context context.Context, client *elastic.Client, index string, model any, query *elastic.BoolQuery) paginator.Adapter {
	totalQuery := query
	paginator := elasticsearchPaginator{
		context:    context,
		client:     client,
		index:      index,
		model:      model,
		pageQuery:  query,
		totalQuery: totalQuery,
	}

	return &paginator
}

func (es *elasticsearchPaginator) Nums() (int64, error) {
	result, err := es.client.Search().Index(es.index).IgnoreUnavailable(true).Query(es.totalQuery).Do(es.context)
	if err != nil {
		return 0, err
	}

	return result.TotalHits(), nil
}

func (es *elasticsearchPaginator) Slice(offset int, length int, data any) error {
	result, err := es.client.Search().Index(es.index).IgnoreUnavailable(true).Query(es.pageQuery).From(offset).Size(length).Do(es.context)
	if err != nil {
		return err
	}

	if result.Hits == nil {
		return nil
	}

	records := make([]map[string]any, 0, result.TotalHits())
	var record map[string]any
	for _, hit := range result.Hits.Hits {
		_ = json.Unmarshal(hit.Source, &record)
		records = append(records, record)
	}

	temp, _ := json.Marshal(records)

	return json.Unmarshal(temp, data)
}
