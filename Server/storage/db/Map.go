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
	db := &memMap{}
	return db.init()
}

func (db *memMap) init() *memMap {
	db.mIDMask = make(map[string]string)
	db.mIDHash = make(map[string]string)
	// load listID from database
	// listID = ......
	//
	db.mUIDlkCTX = make(map[string]string)
	db.mCTXlkUID = make(map[string]string)
	return db
}

// UpdatePolicy :
func (db *memMap) UpdatePolicy(policy, uid, ctx, rw string) (err error) {
	if policy, err = validate(policy); err != nil {
		return err
	}

	id := genPolicyID(policy, uid, ctx, rw)
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

	return nil
}

func (db *memMap) GetPolicyHash(id string) (string, bool) {
	if hashcode, ok := db.mIDHash[id]; ok {
		return hashcode, ok
	}
	return "", false
}

func (db *memMap) GetPolicy(id string) (string, bool) {
	if mask, ok := db.mIDMask[id]; ok {
		return mask, ok
	}
	return "", false
}
