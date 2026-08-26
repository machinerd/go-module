package gqlgen

import (
	"context"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/machinerd/go-module/log"
)

// LoggingTracer implements dataloader.Tracer[K, V] and logs dataloader
// batch execution (one line per batch, not per key) through go-module/log,
// tagged with Name and the request id GetRequestID extracts from ctx.
// Attach it to a loader with dataloader.WithTracer:
//
//	dataloader.NewBatchedLoader(batchFn, dataloader.WithTracer[string, *model.Company](
//		gqlgen.LoggingTracer[string, *model.Company]{
//			Name:         "CompanyLoader",
//			GetRequestID: myRequestIDFunc,
//		},
//	))
type LoggingTracer[K comparable, V any] struct {
	// Name identifies the loader in the log output, e.g. the loader's field name.
	Name string
	// GetRequestID extracts the request/correlation id from ctx.
	GetRequestID RequestIDFunc
}

func (t LoggingTracer[K, V]) TraceLoad(ctx context.Context, key K) (context.Context, dataloader.TraceLoadFinishFunc[V]) {
	return ctx, func(dataloader.Thunk[V]) {}
}

func (t LoggingTracer[K, V]) TraceLoadMany(ctx context.Context, keys []K) (context.Context, dataloader.TraceLoadManyFinishFunc[V]) {
	return ctx, func(dataloader.ThunkMany[V]) {}
}

func (t LoggingTracer[K, V]) TraceBatch(ctx context.Context, keys []K) (context.Context, dataloader.TraceBatchFinishFunc[V]) {
	requestID := t.GetRequestID(ctx)
	field := ""
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		field = fc.Field.Name
	}
	start := time.Now()
	log.Debug("dataloader batch start", "loader", t.Name, "field", field, "keys", len(keys), "request_id", requestID)

	return ctx, func(results []*dataloader.Result[V]) {
		errCount := 0
		for _, r := range results {
			if r != nil && r.Error != nil {
				errCount++
			}
		}
		duration := time.Since(start)
		if errCount > 0 {
			log.Error("dataloader batch failed", "loader", t.Name, "field", field, "keys", len(keys), "errors", errCount, "duration_ms", duration.Milliseconds(), "request_id", requestID)
		} else {
			log.Info("dataloader batch done", "loader", t.Name, "field", field, "keys", len(keys), "duration_ms", duration.Milliseconds(), "request_id", requestID)
		}
	}
}
