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

// ---------------------- //

// PolicyCountTr :
func (db *badgerDB) PolicyCountTr(ctx context.Context) int {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := db.tracer.StartSpan("PolicyCount", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, "badgerDB")
		span.SetTag("PolicyCount", "")
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	return db.PolicyCount()
}

// PolicyIDTr :
func (db *badgerDB) PolicyIDTr(ctx context.Context, user, n3ctx, rw, object string) string {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := db.tracer.StartSpan("PolicyID", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, "badgerDB")
		span.SetTag("PolicyID", fSf("[%s] [%s] [%s] [%s]", user, n3ctx, rw, object))
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	return db.PolicyID(user, n3ctx, rw, object)
}

// PolicyIDsTr :
func (db *badgerDB) PolicyIDsTr(ctx context.Context, user, n3ctx, rw string, objects ...string) []string {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := db.tracer.StartSpan("PolicyIDs", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, "badgerDB")
		span.SetTag("PolicyIDs", fSf("[%s] [%s] [%s] [%v]", user, n3ctx, rw, objects))
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	return db.PolicyIDs(user, n3ctx, rw, objects...)
}

// UpdatePolicyTr :
func (db *badgerDB) UpdatePolicyTr(ctx context.Context, policy, name, user, n3ctx, rw string) (id, obj string, err error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := db.tracer.StartSpan("UpdatePolicy", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, "badgerDB")
		span.SetTag("UpdatePolicy", fSf("[%s] [%s] [%s] [%s] [%s]", policy, name, user, n3ctx, rw))
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	return db.UpdatePolicy(policy, name, user, n3ctx, rw)
}

// DeletePolicyTr :
func (db *badgerDB) DeletePolicyTr(ctx context.Context, id string) (err error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := db.tracer.StartSpan("DeletePolicy", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, "badgerDB")
		span.SetTag("DeletePolicy", id)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	return db.DeletePolicy(id)
}

// PolicyHashTr :
func (db *badgerDB) PolicyHashTr(ctx context.Context, id string) (string, bool) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := db.tracer.StartSpan("PolicyHash", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, "badgerDB")
		span.SetTag("PolicyHash", id)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	return db.PolicyHash(id)
}

// PolicyTr :
func (db *badgerDB) PolicyTr(ctx context.Context, id string) (string, bool) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := db.tracer.StartSpan("Policy", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, "badgerDB")
		span.SetTag("Policy", id)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	return db.Policy(id)
}

// ------------------------------------------- //

// MapRW2lsPIDTr :
func (db *badgerDB) MapRW2lsPIDTr(ctx context.Context, user, n3ctx string, lsRW ...string) map[string][]string {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := db.tracer.StartSpan("MapRW2lsPID", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, "badgerDB")
		span.SetTag("MapRW2lsPID", fSf("[%s] [%s] [%v]", user, n3ctx, lsRW))
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	return db.MapRW2lsPID(user, n3ctx, lsRW...)
}

// MapCtx2lsUserTr :
func (db *badgerDB) MapCtx2lsUserTr(ctx context.Context, lsCtx ...string) map[string][]string {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := db.tracer.StartSpan("MapCtx2lsUser", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, "badgerDB")
		span.SetTag("MapCtx2lsUser", fSf("[%v]", lsCtx))
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	return db.MapCtx2lsUser(lsCtx...)
}

// MapUser2lsCtxTr :
func (db *badgerDB) MapUser2lsCtxTr(ctx context.Context, users ...string) map[string][]string {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := db.tracer.StartSpan("MapUser2lsCtx", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, "badgerDB")
		span.SetTag("MapUser2lsCtx", fSf("[%v]", users))
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	return db.MapUser2lsCtx(users...)
}

// MapUC2lsObjectTr :
func (db *badgerDB) MapUC2lsObjectTr(ctx context.Context, user, n3ctx string) map[string][]string {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := db.tracer.StartSpan("MapUC2lsObject", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, "badgerDB")
		span.SetTag("MapUC2lsObject", fSf("[%s] [%s]", user, n3ctx))
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	return db.MapUC2lsObject(user, n3ctx)
}
