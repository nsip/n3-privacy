package db

import (
	"context"

	"github.com/opentracing/opentracing-go"
	tags "github.com/opentracing/opentracing-go/ext"
)

// SetTracer :
func (db *badgerDB) SetTracer(tracer opentracing.Tracer) {
	db.tracer = tracer
}

// SetContext :
func (db *badgerDB) SetContext(context context.Context) {
	db.context = context
}

func (db *badgerDB) GetContext() context.Context {
	return db.context
}

func (db *badgerDB) UseTracing(operName, spanValue, tagKey, tagValue string) {
	fPln(" ------- db tracing ------- ")

	if ctx := db.GetContext(); ctx != nil {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			span := db.tracer.StartSpan(operName, opentracing.ChildOf(span.Context()))
			tags.SpanKindRPCClient.Set(span)
			tags.PeerService.Set(span, spanValue)
			span.SetTag(tagKey, tagValue)
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
}
