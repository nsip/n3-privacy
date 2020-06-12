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
	PolicyCount() int
	PolicyID(user, n3ctx, rw, object string) string
	PolicyIDs(user, n3ctx, rw string, objects ...string) []string
	PolicyHash(id string) (string, bool)
	Policy(id string) (string, bool)
	DeletePolicy(id string) error
	MapRW2lsPID(user, n3ctx string, lsRW ...string) map[string][]string
	MapCtx2lsUser(lsCtx ...string) map[string][]string
	MapUser2lsCtx(users ...string) map[string][]string
	MapUC2lsObject(user, n3ctx string) map[string][]string

	// EnCode
	SetEncPwd(pwd string)
}

// DBTr :
type DBTr interface {
	DB
	// Tracer
	SetTracer(tracer opentracing.Tracer)
	SetContext(ctx context.Context)
	GetContext() context.Context
}

// NewDB :
func NewDB(dbType string, tracing bool) interface{} {

	var tracer opentracing.Tracer
	if tracing {
		initTracer := func(serviceName string) opentracing.Tracer {
			cfg, err := config.FromEnv()
			failOnErr("%v: ", err)
			cfg.ServiceName = serviceName
			cfg.Sampler.Type = "const"
			cfg.Sampler.Param = 1
			tracer, _, err := cfg.NewTracer()
			failOnErr("%v: ", err)
			return tracer
		}
		tracer = initTracer(dbType)
	}

	switch dbType {
	case "badger", "BADGER":
		if tracing {
			n3db := db.NewDBByBadger().(DBTr)
			n3db.SetEncPwd(dbType)
			n3db.SetTracer(tracer)
			return n3db
		}
		n3db := db.NewDBByBadger().(DB)
		n3db.SetEncPwd(dbType)
		return n3db

	default:
		failOnErr("%v: [%s]", eg.PARAM_NOT_SUPPORTED, dbType)
		return nil
	}
}
