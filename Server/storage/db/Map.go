package db

// memMap :
type memMap struct {
	mIDMask   map[string]string
	mIDHash   map[string]string
	mUIDlkCTX map[string]string
	mCTXlkUID map[string]string
}

// NewDBByMap :
func NewDBByMap() interface{} {
	return (&memMap{}).init()
}

func (db *memMap) loadIDList() int {
	listID = []string{}
	return len(listID)
}

func (db *memMap) init() *memMap {
	db.mIDMask = make(map[string]string)
	db.mIDHash = make(map[string]string)
	db.mUIDlkCTX = make(map[string]string)
	db.mCTXlkUID = make(map[string]string)
	// load listID from database
	db.loadIDList()
	//
	return db
}

// PolicyCount :
func (db *memMap) PolicyCount() int {
	return len(listID)
}

// PolicyID :
func (db *memMap) PolicyID(uid, ctx, rw, object string) []string {
	return getPolicyID(uid, ctx, rw, object)
}

func (db *memMap) PolicyIDs(uid, ctx, rw string, objects ...string) []string {
	return getPolicyID(uid, ctx, rw, objects...)
}

// UpdatePolicy :
func (db *memMap) UpdatePolicy(policy, uid, ctx, rw string) (id string, err error) {
	if policy, err = validate(policy); err != nil {
		return "", err
	}

	id = genPolicyID(policy, uid, ctx, rw)
	db.mIDMask[id] = policy
	db.mIDHash[id] = hash(policy)

	if !xin(id, listID) {
		listID = append(listID, id)
	}

	// for extention query
	if sIndex(db.mUIDlkCTX[uid], ctx) < 0 {
		db.mUIDlkCTX[uid] += (linker + ctx)
	}
	if sIndex(db.mCTXlkUID[ctx], uid) < 0 {
		db.mCTXlkUID[ctx] += (linker + uid)
	}

	logMeta(policy, ctx, rw)
	return id, nil
}

func (db *memMap) PolicyHash(id string) (string, bool) {
	if hashcode, ok := db.mIDHash[id]; ok {
		return hashcode, ok
	}
	return "", false
}

func (db *memMap) Policy(id string) (string, bool) {
	if mask, ok := db.mIDMask[id]; ok {
		return mask, ok
	}
	return "", false
}
