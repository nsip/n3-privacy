package db

import (
	"log"

	badger "github.com/dgraph-io/badger"
)

type badgerDB struct {
	pdb *badger.DB
	err error
}

// NewDBByBadger :
func NewDBByBadger(path string) interface{} {
	db := badgerDB{pdb: nil, err: nil}
	return db.init(path)
}

func (db *badgerDB) init(path string) *badgerDB {
	db.pdb, db.err = badger.Open(badger.DefaultOptions(path))
	if db.err != nil {
		log.Fatal(db.err)
	}
	return db
}

func (db *badgerDB) release() {
	if db.pdb != nil {
		db.pdb.Close()
	}
}

func (db *badgerDB) UpdatePolicy(policy, uid, ctx, rw string) (err error) {
	if policy, err = valfmtPolicy(policy); err != nil {
		return err
	}

	id := genPolicyID(policy, uid, ctx, rw)
	txn := db.pdb.NewTransaction(true)
	for k, v := range map[string]string{id: policy} {
		txn.Set([]byte(k), []byte(v))
	}
	txn.Commit()

	// different for hash ...

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
