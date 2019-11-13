package storage

import (
	"log"

	"github.com/nsip/n3-privacy/Server/storage/db"
)

// DB :
type DB interface {
	UpdatePolicy(policy, uid, ctx, rw string) error
	GetPolicyID(uid, ctx, object, rw string) []string
	GetPolicyHash(id string) (string, bool)
	GetPolicy(id string) (string, bool)
	RecordMeta(policy, metafile string) bool
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
