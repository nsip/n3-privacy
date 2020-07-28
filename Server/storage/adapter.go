package storage

import (
	"github.com/cdutwhu/n3-util/n3err"
	"github.com/cdutwhu/n3-util/n3tracing"
	db "github.com/nsip/n3-privacy/Server/storage/badger"
)

// DB :
type DB interface {
	//
	n3tracing.ITrace

	// DB functions
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

// NewDB :
func NewDB(dbType string) DB {
	failP1OnErrWhen(
		!exist(sToUpper(dbType), "BADGER"),
		"%v: [%s]", n3err.PARAM_NOT_SUPPORTED, dbType,
	)
	db := db.NewDBByBadger().(DB)
	db.SetEncPwd(dbType)
	db.SetTracer(n3tracing.InitTracer(dbType))
	return db
}
