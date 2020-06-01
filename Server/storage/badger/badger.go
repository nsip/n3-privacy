package db

import (
	"os"

	badger "github.com/dgraph-io/badger"
	cfg "github.com/nsip/n3-privacy/Server/config"
)

// SetEncPwd :
func (db *badgerDB) SetEncPwd(pwd string) {
	db.encPwd = pwd
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
	Cfg := env2Struct("Cfg", &cfg.Config{}).(*cfg.Config)
	path := Cfg.Storage.BadgerDBPath
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

// ----------------------- Export ----------------------- //

// PolicyCount : tr
func (db *badgerDB) PolicyCount() int {
	return len(listID)
}

// PolicyID : tr
func (db *badgerDB) PolicyID(user, n3ctx, rw, object string) string {
	if lsID := getPolicyID(user, n3ctx, rw, object); len(lsID) > 0 {
		return lsID[0]
	}
	return ""
}

// PolicyIDs : tr
func (db *badgerDB) PolicyIDs(user, n3ctx, rw string, objects ...string) []string {
	return getPolicyID(user, n3ctx, rw, objects...)
}

// UpdatePolicy : tr
func (db *badgerDB) UpdatePolicy(policy, name, user, n3ctx, rw string) (id, obj string, err error) {
	if policy, err = validate(policy); err != nil {
		return "", "", err
	}
	id, obj = genPolicyID(policy, name, user, n3ctx, rw)
	encPolicy := string(encrypt([]byte(policy), db.encPwd))
	err = updateBadgerDB(
		[]*badger.DB{db.mIDPolicy, db.mIDHash, db.mIDUser, db.mIDCtx, db.mIDObject},
		[]string{id, id, hash(user)[:lenOfUID], hash(n3ctx)[:lenOfCTX], hash(obj)[:lenOfOID]},
		[]string{encPolicy, hash(policy), user, n3ctx, obj})
	if err == nil && !xin(id, listID) {
		listID = append(listID, id)
	}
	// logMeta(policy, n3ctx, rw)
	return id, obj, err
}

// DeletePolicy : tr
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

// PolicyHash : tr
func (db *badgerDB) PolicyHash(id string) (string, bool) {
	if values, err := getBadgerDB([]*badger.DB{db.mIDHash}, []string{id}); err == nil {
		return values[0], true
	}
	return "", false
}

// Policy : tr
func (db *badgerDB) Policy(id string) (string, bool) {
	if values, err := getBadgerDB([]*badger.DB{db.mIDPolicy}, []string{id}); err == nil {
		if policy, err := decrypt([]byte(values[0]), db.encPwd); err == nil {
			return string(policy), true
		}
	}
	return "", false
}

// ----------------------- Optional, for management ----------------------- //

func (db *badgerDB) listPolicyID(user, n3ctx string, lsRW ...string) (lsID [][]string) {
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

	if user == "" && n3ctx == "" {
		return allPolicyID(lsRW...)
	}

	uCode := hash(user)[:lenOfUID]
	cCode := hash(n3ctx)[:lenOfCTX]
	lsID = make([][]string, len(lsRW))
	if len(lsRW) == 0 {
		lsID = make([][]string, 1)
	}
	for i, IDs := range allPolicyID(lsRW...) {
		lsID[i] = make([]string, 0)
		for _, id := range IDs {
			switch {
			case user == "" && n3ctx != "":
				if p := sIndex(id, cCode); p == lenOfHash*3/4 {
					lsID[i] = append(lsID[i], id)
				}
			case user != "" && n3ctx == "":
				if p := sIndex(id, uCode); p == lenOfHash/2 {
					lsID[i] = append(lsID[i], id)
				}
			case user != "" && n3ctx != "":
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
	for i, n3ctx := range lsCtx {
		cCode := hash(n3ctx)[:lenOfCTX]
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

func (db *badgerDB) listObject(user, n3ctx string) []string {
	uCode := hash(user)[:lenOfUID]
	cCode := hash(n3ctx)[:lenOfCTX]
	oCodes := []string{}
	for _, id := range listID {
		switch {
		case user == "" && n3ctx == "":
			oCodes = append(oCodes, oCodeByPID(id))
		case user != "" && n3ctx != "":
			if uCodeByPID(id) == uCode && cCodeByPID(id) == cCode {
				oCodes = append(oCodes, oCodeByPID(id))
			}
		case user != "" && n3ctx == "":
			if uCodeByPID(id) == uCode {
				oCodes = append(oCodes, oCodeByPID(id))
			}
		case user == "" && n3ctx != "":
			if cCodeByPID(id) == cCode {
				oCodes = append(oCodes, oCodeByPID(id))
			}
		}
	}
	objList, _ := getOneBadgerDB(db.mIDObject, toSet(oCodes).([]string))
	return objList
}

// ------------------------------------------------- //

// MapRW2lsPID : tr
func (db *badgerDB) MapRW2lsPID(user, n3ctx string, lsRW ...string) map[string][]string {
	rt := make(map[string][]string)
	key := fSf("%s@%s", user, n3ctx)
	for i, IDs := range db.listPolicyID(user, n3ctx, lsRW...) {
		if user == "" && n3ctx == "" {
			rt["all"] = IDs
		} else if user != "" && n3ctx == "" {
			rt[user] = IDs
		} else if user == "" && n3ctx != "" {
			rt[n3ctx] = IDs
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

// MapCtx2lsUser : tr
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

// MapUser2lsCtx : tr
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

// MapUC2lsObject : tr
func (db *badgerDB) MapUC2lsObject(user, n3ctx string) map[string][]string {
	key := user + "@" + n3ctx
	switch {
	case user == "" && n3ctx == "":
		key = "all"
	case user != "" && n3ctx == "":
		key = user
	case user == "" && n3ctx != "":
		key = n3ctx
	}
	return map[string][]string{key: db.listObject(user, n3ctx)}
}
