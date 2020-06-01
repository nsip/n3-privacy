package storage

import (
	"context"

	eg "github.com/cdutwhu/n3-util/n3errs"
	db "github.com/nsip/n3-privacy/Server/storage/badger"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

// DB :
type DB interface {
	UpdatePolicy(policy, name, user, n3ctx, rw string) (string, string, error)
	UpdatePolicyTr(ctx context.Context, policy, name, user, n3ctx, rw string) (string, string, error)

	PolicyCount() int
	PolicyCountTr(ctx context.Context) int

	PolicyID(user, n3ctx, rw, object string) string
	PolicyIDTr(ctx context.Context, user, n3ctx, rw, object string) string

	PolicyIDs(user, n3ctx, rw string, objects ...string) []string
	PolicyIDsTr(ctx context.Context, user, n3ctx, rw string, objects ...string) []string

	PolicyHash(id string) (string, bool)
	PolicyHashTr(ctx context.Context, id string) (string, bool)

	Policy(id string) (string, bool)
	PolicyTr(ctx context.Context, id string) (string, bool)

	DeletePolicy(id string) error
	DeletePolicyTr(ctx context.Context, id string) error

	// for management
	MapRW2lsPID(user, n3ctx string, lsRW ...string) map[string][]string
	MapRW2lsPIDTr(ctx context.Context, user, n3ctx string, lsRW ...string) map[string][]string

	MapCtx2lsUser(lsCtx ...string) map[string][]string
	MapCtx2lsUserTr(ctx context.Context, lsCtx ...string) map[string][]string

	MapUser2lsCtx(users ...string) map[string][]string
	MapUser2lsCtxTr(ctx context.Context, users ...string) map[string][]string

	MapUC2lsObject(user, n3ctx string) map[string][]string
	MapUC2lsObjectTr(ctx context.Context, user, n3ctx string) map[string][]string

	// Tracer
	SetTracer(tracer opentracing.Tracer)

	// EnCode
	SetEncPwd(pwd string)
}

// NewDB :
func NewDB(dbType string) DB {
	switch dbType {
	case "badger":
		n3db := db.NewDBByBadger().(DB)
		n3db.SetTracer(initTracer(dbType))
		n3db.SetEncPwd(dbType)
		return n3db

	case "map":
		failOnErr("%v: [%s]", eg.NOT_IMPLEMENTED, dbType)
		return nil

	default:
		failOnErr("%v: [%s]", eg.PARAM_NOT_SUPPORTED, dbType)
		return nil
	}
}

func initTracer(serviceName string) opentracing.Tracer {
	cfg, err := config.FromEnv()
	failOnErr("%v: ", err)
	cfg.ServiceName = serviceName
	cfg.Sampler.Type = "const"
	cfg.Sampler.Param = 1

	tracer, _, err := cfg.NewTracer()
	failOnErr("%v: ", err)
	return tracer
}
