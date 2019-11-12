package storage

import (
	"github.com/nsip/n3-privacy/Server/storage/db"
)

// DB :
type DB interface {
	GenPolicyCode(policy string) string
	GenPolicyID(policy, uid, ctx, rw string) string
	UpdatePolicy(policy, uid, ctx, rw string) error
	GetPolicyID(uid, ctx, object, rw string) []string
	GetPolicyHash(id string) (string, bool)
	GetPolicy(id string) (string, bool)
	RecordMeta(policy, metafile string) bool
}

// NewDB :
func NewDB() DB {
	return db.NewDBByMap().(DB)
}
