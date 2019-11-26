package storage

import (
	"log"

	"github.com/nsip/n3-privacy/Server/storage/db"
)

// DB :
type DB interface {
	UpdatePolicy(policy, uid, ctx, rw string) (string, error)
	PolicyCount() int
	PolicyID(uid, ctx, rw, object string) []string
	PolicyIDs(uid, ctx, rw string, objects ...string) []string
	PolicyHash(id string) (string, bool)
	Policy(id string) (string, bool)
	DeletePolicy(id string) error
	// Optional
	AllPolicyIDs(rw string) []string               // if rw is "", get all
	AllPolicyIDsOfOneUser(uid, rw string) []string // if rw is "", get all
	AllPolicyIDsOfOneCtx(ctx, rw string) []string  // if rw is "", get all
}

// NewDB :
func NewDB(dbType string) DB {
	switch dbType {
	case "badger":
		return db.NewDBByBadger().(DB)
	case "map":
		return db.NewDBByMap().(DB)
	default:
		log.Fatalf("%s is not supported", dbType)
		return nil
	}
}
