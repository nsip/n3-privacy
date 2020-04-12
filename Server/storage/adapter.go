package storage

import (
	"log"

	"github.com/nsip/n3-privacy/Server/storage/db"
)

// DB :
type DB interface {
	UpdatePolicy(policy, user, ctx, rw string) (string, string, error)
	PolicyCount() int
	PolicyID(user, ctx, rw, object string) string
	PolicyIDs(user, ctx, rw string, objects ...string) []string
	PolicyHash(id string) (string, bool)
	Policy(id string) (string, bool)
	DeletePolicy(id string) error

	// Optional, for management
	MapRW2lsPID(user, ctx string, lsRW ...string) map[string][]string
	MapCtx2lsUser(lsCtx ...string) map[string][]string
	MapUser2lsCtx(users ...string) map[string][]string
	MapUC2lsObject(user, ctx string) map[string][]string
}

// NewDB :
func NewDB(dbType string) DB {
	switch dbType {
	case "badger":
		return db.NewDBByBadger().(DB)
	// case "map":
	// 	return db.NewDBByMap().(DB)
	default:
		log.Fatalf("%s is not supported", dbType)
		return nil
	}
}
