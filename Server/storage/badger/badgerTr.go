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

// ---------------------- //

// PolicyCount :
func (db *badgerDB) PolicyCount() int {
	if ctx := db.GetContext(); ctx != nil {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			span := db.tracer.StartSpan("PolicyCount", opentracing.ChildOf(span.Context()))
			tags.SpanKindRPCClient.Set(span)
			tags.PeerService.Set(span, "badgerDB")
			span.SetTag("PolicyCount", "")
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	return db.policyCount()
}

// PolicyID :
func (db *badgerDB) PolicyID(user, n3ctx, rw, object string) string {
	if ctx := db.GetContext(); ctx != nil {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			span := db.tracer.StartSpan("PolicyID", opentracing.ChildOf(span.Context()))
			tags.SpanKindRPCClient.Set(span)
			tags.PeerService.Set(span, "badgerDB")
			span.SetTag("PolicyID", fSf("[%s] [%s] [%s] [%s]", user, n3ctx, rw, object))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	return db.policyID(user, n3ctx, rw, object)
}

// PolicyIDs :
func (db *badgerDB) PolicyIDs(user, n3ctx, rw string, objects ...string) []string {
	if ctx := db.GetContext(); ctx != nil {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			span := db.tracer.StartSpan("PolicyIDs", opentracing.ChildOf(span.Context()))
			tags.SpanKindRPCClient.Set(span)
			tags.PeerService.Set(span, "badgerDB")
			span.SetTag("PolicyIDs", fSf("[%s] [%s] [%s] [%v]", user, n3ctx, rw, objects))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	return db.policyIDs(user, n3ctx, rw, objects...)
}

// UpdatePolicy :
func (db *badgerDB) UpdatePolicy(policy, name, user, n3ctx, rw string) (id, obj string, err error) {
	if ctx := db.GetContext(); ctx != nil {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			span := db.tracer.StartSpan("UpdatePolicy", opentracing.ChildOf(span.Context()))
			tags.SpanKindRPCClient.Set(span)
			tags.PeerService.Set(span, "badgerDB")
			span.SetTag("UpdatePolicy", fSf("[%s] [%s] [%s] [%s] [%s]", policy, name, user, n3ctx, rw))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	return db.updatePolicy(policy, name, user, n3ctx, rw)
}

// DeletePolicy :
func (db *badgerDB) DeletePolicy(id string) (err error) {
	if ctx := db.GetContext(); ctx != nil {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			span := db.tracer.StartSpan("DeletePolicy", opentracing.ChildOf(span.Context()))
			tags.SpanKindRPCClient.Set(span)
			tags.PeerService.Set(span, "badgerDB")
			span.SetTag("DeletePolicy", id)
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	return db.deletePolicy(id)
}

// PolicyHash :
func (db *badgerDB) PolicyHash(id string) (string, bool) {
	if ctx := db.GetContext(); ctx != nil {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			span := db.tracer.StartSpan("PolicyHash", opentracing.ChildOf(span.Context()))
			tags.SpanKindRPCClient.Set(span)
			tags.PeerService.Set(span, "badgerDB")
			span.SetTag("PolicyHash", id)
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	return db.policyHash(id)
}

// Policy :
func (db *badgerDB) Policy(id string) (string, bool) {
	if ctx := db.GetContext(); ctx != nil {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			span := db.tracer.StartSpan("Policy", opentracing.ChildOf(span.Context()))
			tags.SpanKindRPCClient.Set(span)
			tags.PeerService.Set(span, "badgerDB")
			span.SetTag("Policy", id)
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	return db.policy(id)
}

// ------------------------------------------- //

// MapRW2lsPID :
func (db *badgerDB) MapRW2lsPID(user, n3ctx string, lsRW ...string) map[string][]string {
	if ctx := db.GetContext(); ctx != nil {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			span := db.tracer.StartSpan("MapRW2lsPID", opentracing.ChildOf(span.Context()))
			tags.SpanKindRPCClient.Set(span)
			tags.PeerService.Set(span, "badgerDB")
			span.SetTag("MapRW2lsPID", fSf("[%s] [%s] [%v]", user, n3ctx, lsRW))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	return db.mapRW2lsPID(user, n3ctx, lsRW...)
}

// MapCtx2lsUser :
func (db *badgerDB) MapCtx2lsUser(lsCtx ...string) map[string][]string {
	if ctx := db.GetContext(); ctx != nil {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			span := db.tracer.StartSpan("MapCtx2lsUser", opentracing.ChildOf(span.Context()))
			tags.SpanKindRPCClient.Set(span)
			tags.PeerService.Set(span, "badgerDB")
			span.SetTag("MapCtx2lsUser", fSf("[%v]", lsCtx))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	return db.mapCtx2lsUser(lsCtx...)
}

// MapUser2lsCtx :
func (db *badgerDB) MapUser2lsCtx(users ...string) map[string][]string {
	if ctx := db.GetContext(); ctx != nil {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			span := db.tracer.StartSpan("MapUser2lsCtx", opentracing.ChildOf(span.Context()))
			tags.SpanKindRPCClient.Set(span)
			tags.PeerService.Set(span, "badgerDB")
			span.SetTag("MapUser2lsCtx", fSf("[%v]", users))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	return db.mapUser2lsCtx(users...)
}

// MapUC2lsObject :
func (db *badgerDB) MapUC2lsObject(user, n3ctx string) map[string][]string {
	if ctx := db.GetContext(); ctx != nil {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			span := db.tracer.StartSpan("MapUC2lsObject", opentracing.ChildOf(span.Context()))
			tags.SpanKindRPCClient.Set(span)
			tags.PeerService.Set(span, "badgerDB")
			span.SetTag("MapUC2lsObject", fSf("[%s] [%s]", user, n3ctx))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	return db.mapUC2lsObject(user, n3ctx)
}
