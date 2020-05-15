package db

import (
	"os"

	badger "github.com/dgraph-io/badger"
	glb "github.com/nsip/n3-privacy/Server/global"
)

type badgerDB struct {
	mIDPolicy *badger.DB
	mIDHash   *badger.DB
	err       error
	mIDUser   *badger.DB
	mIDCtx    *badger.DB
	mIDObject *badger.DB
}

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

// ----------------------- Basic ----------------------- //

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

	lsDB := openBadgerDB(path, "mIDPolicy", "mIDHash", "mIDUser", "mIDCtx", "mIDObject")
	db.mIDPolicy, db.mIDHash, db.mIDUser, db.mIDCtx, db.mIDObject = lsDB[0], lsDB[1], lsDB[2], lsDB[3], lsDB[4]

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

func (db *badgerDB) PolicyID(user, ctx, rw, object string) string {
	if lsID := getPolicyID(user, ctx, rw, object); len(lsID) > 0 {
		return lsID[0]
	}
	return ""
}

func (db *badgerDB) PolicyIDs(user, ctx, rw string, objects ...string) []string {
	return getPolicyID(user, ctx, rw, objects...)
}

func (db *badgerDB) UpdatePolicy(policy, name, user, ctx, rw string) (id, obj string, err error) {
	if policy, err = validate(policy); err != nil {
		return "", "", err
	}
	id, obj = genPolicyID(policy, name, user, ctx, rw)
	encPolicy := string(encrypt([]byte(policy), glb.EncPwd))
	err = updateBadgerDB(
		[]*badger.DB{db.mIDPolicy, db.mIDHash, db.mIDUser, db.mIDCtx, db.mIDObject},
		[]string{id, id, hash(user)[:lenOfUID], hash(ctx)[:lenOfCTX], hash(obj)[:lenOfOID]},
		[]string{encPolicy, hash(policy), user, ctx, obj})
	if err == nil && !xin(id, listID) {
		listID = append(listID, id)
	}
	// logMeta(policy, ctx, rw)
	return id, obj, err
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
		if policy, err := decrypt([]byte(values[0]), glb.EncPwd); err == nil {
			return string(policy), true
		}
	}
	return "", false
}

// ----------------------- Optional, for management ----------------------- //

func (db *badgerDB) listPolicyID(user, ctx string, lsRW ...string) (lsID [][]string) {
	allPolicyID := func(lsRW ...string) (lsID [][]string) {
		if len(lsRW) == 0 {
			return [][]string{append([]string{}, listID...)}
		}
		lsID = make([][]string, len(lsRW))
		for i, rw := range lsRW {
			lsID[i] = make([]string, 0)
			for _, id := range listID {
				if sHasSuffix(id, rw) {
					lsID[i] = append(lsID[i], id)
				}
			}
		}
		return
	}

	if user == "" && ctx == "" {
		return allPolicyID(lsRW...)
	}

	uCode := hash(user)[:lenOfUID]
	cCode := hash(ctx)[:lenOfCTX]
	lsID = make([][]string, len(lsRW))
	if len(lsRW) == 0 {
		lsID = make([][]string, 1)
	}
	for i, IDs := range allPolicyID(lsRW...) {
		lsID[i] = make([]string, 0)
		for _, id := range IDs {
			switch {
			case user == "" && ctx != "":
				if p := sIndex(id, cCode); p == lenOfHash*3/4 {
					lsID[i] = append(lsID[i], id)
				}
			case user != "" && ctx == "":
				if p := sIndex(id, uCode); p == lenOfHash/2 {
					lsID[i] = append(lsID[i], id)
				}
			case user != "" && ctx != "":
				p1, p2 := sIndex(id, uCode), sIndex(id, cCode)
				if p1 == lenOfHash/2 && p2 == lenOfHash*3/4 {
					lsID[i] = append(lsID[i], id)
				}
			}
		}
	}
	return
}

func (db *badgerDB) listUser(lsCtx ...string) (lsUser [][]string) {
	allUsers := func() (users []string) {
		uCodes := []string{}
		for _, id := range listID {
			uCodes = append(uCodes, uCodeByPID(id))
		}
		users, _ = getOneBadgerDB(db.mIDUser, toSet(uCodes).([]string))
		return
	}

	if len(lsCtx) == 0 {
		return [][]string{allUsers()}
	}

	lsUser = make([][]string, len(lsCtx))
	for i, ctx := range lsCtx {
		cCode := hash(ctx)[:lenOfCTX]
		uCodes := []string{}
		for _, id := range listID {
			if cCodeByPID(id) == cCode {
				uCodes = append(uCodes, uCodeByPID(id))
			}
		}
		lsUser[i], _ = getOneBadgerDB(db.mIDUser, toSet(uCodes).([]string))
	}
	return
}

func (db *badgerDB) listCtx(users ...string) (lsCtx [][]string) {
	allCtx := func() (ctxList []string) {
		cCodes := []string{}
		for _, id := range listID {
			cCodes = append(cCodes, cCodeByPID(id))
		}
		ctxList, _ = getOneBadgerDB(db.mIDCtx, toSet(cCodes).([]string))
		return
	}

	if len(users) == 0 {
		return [][]string{allCtx()}
	}

	lsCtx = make([][]string, len(users))
	for i, user := range users {
		uCode := hash(user)[:lenOfUID]
		cCodes := []string{}
		for _, id := range listID {
			if uCodeByPID(id) == uCode {
				cCodes = append(cCodes, cCodeByPID(id))
			}
		}
		lsCtx[i], _ = getOneBadgerDB(db.mIDCtx, toSet(cCodes).([]string))
	}
	return
}

func (db *badgerDB) listObject(user, ctx string) []string {
	uCode := hash(user)[:lenOfUID]
	cCode := hash(ctx)[:lenOfCTX]
	oCodes := []string{}
	for _, id := range listID {
		switch {
		case user == "" && ctx == "":
			oCodes = append(oCodes, oCodeByPID(id))
		case user != "" && ctx != "":
			if uCodeByPID(id) == uCode && cCodeByPID(id) == cCode {
				oCodes = append(oCodes, oCodeByPID(id))
			}
		case user != "" && ctx == "":
			if uCodeByPID(id) == uCode {
				oCodes = append(oCodes, oCodeByPID(id))
			}
		case user == "" && ctx != "":
			if cCodeByPID(id) == cCode {
				oCodes = append(oCodes, oCodeByPID(id))
			}
		}
	}
	objList, _ := getOneBadgerDB(db.mIDObject, toSet(oCodes).([]string))
	return objList
}

// --------- //

func (db *badgerDB) MapRW2lsPID(user, ctx string, lsRW ...string) map[string][]string {
	rt := make(map[string][]string)
	key := fSf("%s@%s", user, ctx)
	for i, IDs := range db.listPolicyID(user, ctx, lsRW...) {
		if user == "" && ctx == "" {
			rt["all"] = IDs
		} else if user != "" && ctx == "" {
			rt[user] = IDs
		} else if user == "" && ctx != "" {
			rt[ctx] = IDs
		} else {
			if len(lsRW) == 0 {
				rt[key] = IDs
			} else {
				rt[key+"@"+lsRW[i]] = IDs
			}
		}
	}
	return rt
}

func (db *badgerDB) MapCtx2lsUser(lsCtx ...string) map[string][]string {
	rt := make(map[string][]string)
	for i, users := range db.listUser(lsCtx...) {
		if len(lsCtx) == 0 {
			rt["all"] = users
		} else {
			rt[lsCtx[i]] = users
		}
	}
	return rt
}

func (db *badgerDB) MapUser2lsCtx(users ...string) map[string][]string {
	rt := make(map[string][]string)
	for i, lsCtx := range db.listCtx(users...) {
		if len(users) == 0 {
			rt["all"] = lsCtx
		} else {
			rt[users[i]] = lsCtx
		}
	}
	return rt
}

func (db *badgerDB) MapUC2lsObject(user, ctx string) map[string][]string {
	key := user + "@" + ctx
	switch {
	case user == "" && ctx == "":
		key = "all"
	case user != "" && ctx == "":
		key = user
	case user == "" && ctx != "":
		key = ctx
	}
	return map[string][]string{key: db.listObject(user, ctx)}
}
