package db

import (
	"os"

	badger "github.com/dgraph-io/badger"
	glb "github.com/nsip/n3-privacy/Server/global"
	cmn "github.com/nsip/n3-privacy/common"
)

type badgerDB struct {
	mIDPolicy *badger.DB
	mIDHash   *badger.DB
	lsID      []string
	mUIDlkCTX *badger.DB
	mCTXlkUID *badger.DB
	err       error
}

func closeBadgerDB(lsDB ...*badger.DB) error {
	for _, db := range lsDB {
		if db != nil {
			if err := db.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

func commitAllTxn(lsTxn ...*badger.Txn) error {
	for _, tx := range lsTxn {
		if tx != nil {
			if err := tx.Commit(); err != nil {
				return err
			}
		}
	}
	return nil
}

// NewDBByBadger :
func NewDBByBadger() interface{} {
	db := badgerDB{}
	return db.init()
}

func (db *badgerDB) init() *badgerDB {
	path := glb.Cfg.Storage.BadgerDBPath
	if _, db.err = os.Stat(path); os.IsNotExist(db.err) {
		os.MkdirAll(path, os.ModePerm)
	}

	db.mIDPolicy, db.err = badger.Open(badger.DefaultOptions(path + "IDPolicy"))
	cmn.FailOnErr("%v", db.err)
	db.mIDHash, db.err = badger.Open(badger.DefaultOptions(path + "IDHash"))
	cmn.FailOnErr("%v", db.err)

	// *** load lsID *** //
	countID := func() int {
		opt := badger.DefaultIteratorOptions
		db.mIDPolicy.View(func(txn *badger.Txn) error {
			itr := txn.NewIterator(opt)
			defer itr.Close()
			for itr.Rewind(); itr.Valid(); itr.Next() {
				item := itr.Item()
				item.Value(func(v []byte) error {
					db.lsID = append(db.lsID, string(item.Key()))
					return nil
				})
			}
			return nil
		})
		return len(db.lsID)
	}
	fPln(countID())

	//
	db.mUIDlkCTX, db.err = badger.Open(badger.DefaultOptions(path + "UIDlkCTX"))
	cmn.FailOnErr("%v", db.err)
	db.mCTXlkUID, db.err = badger.Open(badger.DefaultOptions(path + "CTXlkUID"))
	cmn.FailOnErr("%v", db.err)

	return db
}

func (db *badgerDB) close() {
	closeBadgerDB(db.mIDPolicy, db.mIDHash, db.mUIDlkCTX, db.mCTXlkUID)
}

func (db *badgerDB) UpdatePolicy(policy, uid, ctx, rw string) (err error) {
	if policy, err = validate(policy); err != nil {
		return err
	}

	id := genPolicyID(policy, uid, ctx, rw)

	txIDPolicy := db.mIDPolicy.NewTransaction(true)
	defer txIDPolicy.Discard()
	if err = txIDPolicy.Set([]byte(id), []byte(policy)); err != nil {
		return err
	}

	txIDHash := db.mIDHash.NewTransaction(true)
	defer txIDHash.Discard()
	if err = txIDHash.Set([]byte(id), []byte(hash(policy))); err != nil {
		return err
	}

	// commit
	if err = commitAllTxn(txIDPolicy, txIDHash); err == nil {
		db.lsID = append(db.lsID, id)
	}

	// for extention query
	txUIDlkCTX := db.mUIDlkCTX.NewTransaction(true)
	defer txUIDlkCTX.Discard()
	if item, e := txUIDlkCTX.Get([]byte(uid)); e == nil {
		lkCTX := ""
		err = item.Value(func(v []byte) error {
			lkCTX = string(v) + linker + ctx
			return nil
		})
		txUIDlkCTX.Set([]byte(uid), []byte(lkCTX))
		err = commitAllTxn(txUIDlkCTX)
	}

	// for extention query
	txCTXlkUID := db.mCTXlkUID.NewTransaction(true)
	defer txCTXlkUID.Discard()
	if item, e := txCTXlkUID.Get([]byte(ctx)); e == nil {
		lkUID := ""
		err = item.Value(func(v []byte) error {
			lkUID = string(v) + linker + uid
			return nil
		})
		txCTXlkUID.Set([]byte(ctx), []byte(lkUID))
		err = commitAllTxn(txCTXlkUID)
	}

	return err
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
