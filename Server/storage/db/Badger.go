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

func cmtBadgerTxn(lsTxn ...*badger.Txn) error {
	for _, tx := range lsTxn {
		if tx != nil {
			if err := tx.Commit(); err != nil {
				return err
			}
		}
	}
	return nil
}

func updateBadgerDB(dbs []*badger.DB, keys []string, lsValues ...[]string) (err error) {
	lsTxn := []*badger.Txn{}
	for i, db := range dbs {
		txn := db.NewTransaction(true)
		lsTxn = append(lsTxn, txn)
		if len(lsValues) == 0 { //                                    delete
			err = txn.Delete([]byte(keys[i]))
		} else { //                                                   set
			err = txn.Set([]byte(keys[i]), []byte(lsValues[0][i]))
		}
		if err != nil {
			break
		}
	}
	defer func() {
		for _, txn := range lsTxn {
			txn.Discard()
		}
	}()
	if err == nil {
		return cmtBadgerTxn(lsTxn...)
	}
	return err
}

func getBadgerDB(dbs []*badger.DB, keys []string) (values []string, err error) {
	lsTxn := []*badger.Txn{}
	for i, db := range dbs {
		txn := db.NewTransaction(true)
		lsTxn = append(lsTxn, txn)
		switch item, e := txn.Get([]byte(keys[i])); e {
		case nil:
			err = item.Value(func(v []byte) error {
				values = append(values, string(v))
				return nil
			})
			if err != nil {
				return nil, err
			}
		case badger.ErrKeyNotFound:
			return nil, e
		default:
			cmn.FailOnErr("%v", e)
		}
	}
	defer func() {
		for _, txn := range lsTxn {
			txn.Discard()
		}
	}()
	return
}

// NewDBByBadger :
func NewDBByBadger() interface{} {
	return (&badgerDB{}).init()
}

// loadIDList : already invoked by init(), DO NOT call it manually
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

// init : already invoked by New...(), DO NOT call it manually
func (db *badgerDB) init() *badgerDB {
	path := glb.Cfg.Storage.BadgerDBPath
	if _, db.err = os.Stat(path); os.IsNotExist(db.err) {
		os.MkdirAll(path, os.ModePerm)
	}

	db.mIDPolicy, db.err = badger.Open(badger.DefaultOptions(path + "IDPolicy"))
	cmn.FailOnErr("%v", db.err)
	db.mIDHash, db.err = badger.Open(badger.DefaultOptions(path + "IDHash"))
	cmn.FailOnErr("%v", db.err)

	// fPln(db.loadIDList(), "exist in db")
	db.loadIDList()
	return db
}

func (db *badgerDB) close() {
	closeBadgerDB(db.mIDPolicy, db.mIDHash)
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
	err = updateBadgerDB([]*badger.DB{db.mIDPolicy, db.mIDHash}, []string{id, id}, []string{policy, hash(policy)})
	if err == nil && !xin(id, listID) {
		listID = append(listID, id)
	}
	logMeta(policy, ctx, rw)
	return id, err
}

func (db *badgerDB) DeletePolicy(id string) (err error) {
	if err = updateBadgerDB([]*badger.DB{db.mIDHash, db.mIDPolicy}, []string{id, id}); err == nil {
		for i, ID := range listID {
			if ID == id {
				listID = append(listID[:i], listID[i+1:]...)
				break
			}
		}
	}
	return err
}

func (db *badgerDB) PolicyHash(id string) (string, bool) {
	if values, err := getBadgerDB([]*badger.DB{db.mIDHash}, []string{id}); err == nil {
		return values[0], true
	}
	return "", false
}

func (db *badgerDB) Policy(id string) (string, bool) {
	if values, err := getBadgerDB([]*badger.DB{db.mIDPolicy}, []string{id}); err == nil {
		return values[0], true
	}
	return "", false
}

// ---------------------------------------------- //

func (db *badgerDB) AllPolicyID(rw string) (lsID []string) {
	if rw == "" {
		return append(lsID, listID...)
	}
	for _, id := range listID {
		if sHasSuffix(id, rw) {
			lsID = append(lsID, id)
		}
	}
	return
}

func (db *badgerDB) PolicyIDListOfOneUser(uid, rw string) (lsID []string) {
	uCode := hash(uid)[:lenOfUID]
	for _, id := range listID {
		if i := sIndex(id, uCode); i == lenOfHash/2 {
			if rw == "" {
				lsID = append(lsID, id)
				continue
			}
			if sHasSuffix(id, rw) {
				lsID = append(lsID, id)
			}
		}
	}
	return
}

func (db *badgerDB) PolicyIDListOfOneCtx(ctx, rw string) (lsID []string) {
	cCode := hash(ctx)[:lenOfCTX]
	for _, id := range listID {
		if i := sIndex(id, cCode); i == lenOfHash*3/4 {
			if rw == "" {
				lsID = append(lsID, id)
				continue
			}
			if sHasSuffix(id, rw) {
				lsID = append(lsID, id)
			}
		}
	}
	return
}
