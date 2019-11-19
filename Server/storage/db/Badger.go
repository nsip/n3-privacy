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

func cmtAllTxn(lsTxn ...*badger.Txn) error {
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
	return (&badgerDB{}).init()
}

// loadIDList :
func (db *badgerDB) loadIDList() int {
	opt := badger.DefaultIteratorOptions
	db.mIDPolicy.View(func(txn *badger.Txn) error {
		itr := txn.NewIterator(opt)
		defer itr.Close()
		for itr.Rewind(); itr.Valid(); itr.Next() {
			item := itr.Item()
			item.Value(func(v []byte) error {
				listID = append(listID, string(item.Key()))
				return nil
			})
		}
		return nil
	})
	return len(listID)
}

// included in New...
func (db *badgerDB) init() *badgerDB {
	path := glb.Cfg.Storage.BadgerDBPath
	if _, db.err = os.Stat(path); os.IsNotExist(db.err) {
		os.MkdirAll(path, os.ModePerm)
	}

	db.mIDPolicy, db.err = badger.Open(badger.DefaultOptions(path + "IDPolicy"))
	cmn.FailOnErr("%v", db.err)
	db.mIDHash, db.err = badger.Open(badger.DefaultOptions(path + "IDHash"))
	cmn.FailOnErr("%v", db.err)
	db.mUIDlkCTX, db.err = badger.Open(badger.DefaultOptions(path + "UIDlkCTX"))
	cmn.FailOnErr("%v", db.err)
	db.mCTXlkUID, db.err = badger.Open(badger.DefaultOptions(path + "CTXlkUID"))
	cmn.FailOnErr("%v", db.err)

	// *** load listID *** //
	CountID := db.loadIDList()
	fPln(CountID, "exist in db")
	return db
}

func (db *badgerDB) close() {
	closeBadgerDB(db.mIDPolicy, db.mIDHash, db.mUIDlkCTX, db.mCTXlkUID)
}

func (db *badgerDB) PolicyCount() int {
	return len(listID)
}

func (db *badgerDB) PolicyID(uid, ctx, rw, object string) []string {
	return getPolicyID(uid, ctx, rw, object)
}

func (db *badgerDB) PolicyIDs(uid, ctx, rw string, objects ...string) []string {
	return getPolicyID(uid, ctx, rw, objects...)
}

func (db *badgerDB) UpdatePolicy(policy, uid, ctx, rw string) (id string, err error) {
	if policy, err = validate(policy); err != nil {
		return "", err
	}

	id = genPolicyID(policy, uid, ctx, rw)

	txIDPolicy := db.mIDPolicy.NewTransaction(true)
	defer txIDPolicy.Discard()
	if err = txIDPolicy.Set([]byte(id), []byte(policy)); err != nil {
		return "", err
	}

	txIDHash := db.mIDHash.NewTransaction(true)
	defer txIDHash.Discard()
	if err = txIDHash.Set([]byte(id), []byte(hash(policy))); err != nil {
		return "", err
	}

	// commit
	if err = cmtAllTxn(txIDPolicy, txIDHash); err == nil {
		if !xin(id, listID) {
			listID = append(listID, id)
		}
	}

	// for extention query
	txUIDlkCTX := db.mUIDlkCTX.NewTransaction(true)
	defer txUIDlkCTX.Discard()
	item, e := txUIDlkCTX.Get([]byte(uid))
	switch e {
	case nil:
		err = item.Value(func(v []byte) error {
			if sIndex(string(v), ctx) < 0 {
				return txUIDlkCTX.Set([]byte(uid), []byte(string(v)+linker+ctx))
			}
			return nil
		})
	case badger.ErrKeyNotFound:
		txUIDlkCTX.Set([]byte(uid), []byte(ctx))
	}
	cmtAllTxn(txUIDlkCTX)

	// for extention query
	txCTXlkUID := db.mCTXlkUID.NewTransaction(true)
	defer txCTXlkUID.Discard()
	item, e = txCTXlkUID.Get([]byte(ctx))
	switch e {
	case nil:
		err = item.Value(func(v []byte) error {
			if sIndex(string(v), uid) < 0 {
				return txCTXlkUID.Set([]byte(ctx), []byte(string(v)+linker+uid))
			}
			return nil
		})
	case badger.ErrKeyNotFound:
		txCTXlkUID.Set([]byte(ctx), []byte(uid))
	}
	cmtAllTxn(txCTXlkUID)

	logMeta(policy, ctx, rw)
	return id, err
}

func (db *badgerDB) PolicyHash(id string) (string, bool) {
	txIDHash := db.mIDHash.NewTransaction(true)
	defer txIDHash.Discard()
	if item, e := txIDHash.Get([]byte(id)); e == nil {
		hashcode := ""
		e = item.Value(func(v []byte) error {
			hashcode = string(v)
			return nil
		})
		return hashcode, true
	}
	return "", false
}

func (db *badgerDB) Policy(id string) (string, bool) {
	txIDPolicy := db.mIDPolicy.NewTransaction(true)
	defer txIDPolicy.Discard()
	if item, e := txIDPolicy.Get([]byte(id)); e == nil {
		policy := ""
		e = item.Value(func(v []byte) error {
			policy = string(v)
			return nil
		})
		return policy, true
	}
	return "", false
}

// ------------------------------------- //

// ListCTXByUID :
func (db *badgerDB) ListCTXByUID(uid string) (lsCTX []string) {
	strCTX := ""
	txUIDlkCTX := db.mUIDlkCTX.NewTransaction(true)
	defer txUIDlkCTX.Discard()
	if item, e := txUIDlkCTX.Get([]byte(uid)); e == nil {
		item.Value(func(v []byte) error {
			strCTX = string(v)
			return nil
		})
	}
	return sSpl(strCTX, linker)
}

// ListUIDByCTX :
func (db *badgerDB) ListUIDByCTX(ctx string) (lsUID []string) {
	strUID := ""
	txCTXlkUID := db.mCTXlkUID.NewTransaction(true)
	defer txCTXlkUID.Discard()
	if item, e := txCTXlkUID.Get([]byte(ctx)); e == nil {
		item.Value(func(v []byte) error {
			strUID = string(v)
			return nil
		})
	}
	return sSpl(strUID, linker)
}

func (db *badgerDB) ListPIDByUID(uid, rw string) (lsPID []string) {
	for _, ctx := range db.ListCTXByUID(uid) {
		for _, pid := range getPolicyID(uid, ctx, rw) {
			lsPID = append(lsPID, pid)
		}
	}
	return
}

func (db *badgerDB) ListPIDByCTX(ctx, rw string) (lsPID []string) {
	for _, uid := range db.ListUIDByCTX(ctx) {
		for _, pid := range getPolicyID(uid, ctx, rw) {
			lsPID = append(lsPID, pid)
		}
	}
	return
}

// func (db *badgerDB) ListPIDByOBJ(obj, rw string) (lsPID []string) {
// 	return nil
// }
