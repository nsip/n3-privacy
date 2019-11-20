package db

// memMap :
type memMap struct {
	mIDMask map[string]string
	mIDHash map[string]string
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

func (db *memMap) DeletePolicy(id string) error {
	delete(db.mIDMask, id)
	delete(db.mIDHash, id)

	// listID handle...

	return nil
}
