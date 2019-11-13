package db

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/nsip/n3-privacy/jkv"
	pp "github.com/nsip/n3-privacy/preprocess"
)

// memMap :
type memMap struct {
	mIDMask   map[string]string
	mIDHash   map[string]string
	lsID      []string
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
	// load db.lsID from database
	db.lsID = []string{}
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

	db.lsID = append(db.lsID, id)

	// for extention query
	if sIndex(db.mUIDlkCTX[uid], ctx) < 0 {
		db.mUIDlkCTX[uid] += (linker + ctx)
	}
	if sIndex(db.mCTXlkUID[ctx], uid) < 0 {
		db.mCTXlkUID[ctx] += (linker + uid)
	}

	return nil
}

func (db *memMap) GetPolicyID(uid, ctx, object, rw string) (lsID []string) {
	oid := hash(object)[:lenOfOID]
	sid := hash(uid + ctx + rw)[:lenOfSID]
	for _, id := range db.lsID {
		if sHasPrefix(id, oid) && sHasSuffix(id, sid) {
			lsID = append(lsID, id)
		}
	}
	return lsID
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

func (db *memMap) RecordMeta(policy, metafile string) (updated bool) {
	jkvM := jkv.NewJKV(policy, hash(policy))
	object := jkvM.LsL12Fields[1][0]
	md := &MetaData{Object: object, Fields: jkvM.LsL12Fields[2]}
	if b, e := json.Marshal(md); e == nil {
		newPolicy := pp.FmtJSONStr(string(b))
		// first meta.
		if _, err := os.Stat(metafile); err != nil && os.IsNotExist(err) {
			newPolicy, _ = jkv.Indent(newPolicy, 2, false)
			newPolicy = fSf("[\n%s]", newPolicy)
			ioutil.WriteFile(metafile, []byte(newPolicy), 0666) // make sure meta.json is formated

		} else {
			b, _ = ioutil.ReadFile(metafile)
			policies := jkv.SplitJSONArr(string(b))
			// update
			for i, policy := range policies {
				md := MetaData{}
				json.Unmarshal([]byte(policy), &md)
				if md.Object == object {
					policies[i] = newPolicy
					updated = true
					break
				}
			}
			if !updated {
				policies = append(policies, newPolicy)
			}
			ioutil.WriteFile(metafile, []byte(jkv.MergeJSON(policies...)), 0666)
		}
	}
	return
}
