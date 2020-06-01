package db

import badger "github.com/dgraph-io/badger"

func openBadgerDB(dbPath string, lsDBName ...string) (lsDB []*badger.DB) {
	lsDB = make([]*badger.DB, len(lsDBName))
	for i, name := range lsDBName {
		db, err := badger.Open(badger.DefaultOptions(dbPath + name))
		failOnErr("%v", err)
		lsDB[i] = db
	}
	return
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
	values = make([]string, 0)
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
			failOnErr("%v", e)
		}
	}
	defer func() {
		for _, txn := range lsTxn {
			txn.Discard()
		}
	}()
	return
}

func getOneBadgerDB(db *badger.DB, keys []string) (values []string, err error) {
	values = make([]string, 0)
	txn := db.NewTransaction(true)
	defer txn.Discard()
	for _, key := range keys {
		switch item, e := txn.Get([]byte(key)); e {
		case nil:
			err = item.Value(func(v []byte) error {
				values = append(values, string(v))
				return nil
			})
			if err != nil {
				return nil, err
			}
		case badger.ErrKeyNotFound:
			values = append(values, "")
		default:
			failOnErr("%v", e)
		}
	}
	return
}

// NewDBByBadger :
func NewDBByBadger() interface{} {
	return (&badgerDB{}).init()
}
