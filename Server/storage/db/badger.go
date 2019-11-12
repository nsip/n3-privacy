package db

type badgerDB struct {
}

// NewDBByBadger :
func NewDBByBadger() interface{} {
	db := badgerDB{}
	return db.init()
}

func (db *badgerDB) init() *badgerDB {
	return db
}

func (db *badgerDB) GenPolicyCode(policy string) string {
	return ""
}

func (db *badgerDB) GenPolicyID(policy, uid, ctx, rw string) string {
	return ""
}

func (db *badgerDB) UpdatePolicy(policy, uid, ctx, rw string) error {
	return nil
}

func (db *badgerDB) GetPolicyID(uid, ctx, object, rw string) []string {
	return nil
}

func (db *badgerDB) GetPolicyHash(id string) (string, bool) {
	return "", false
}

func (db *badgerDB) GetPolicy(id string) (string, bool) {
	return "", false
}

func (db *badgerDB) RecordMeta(policy, metafile string) bool {
	return false
}
