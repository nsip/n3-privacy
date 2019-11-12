package db

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sort"

	u "github.com/cdutwhu/go-util"
	"github.com/nsip/n3-privacy/jkv"
	pp "github.com/nsip/n3-privacy/preprocess"
)

// memMap :
type memMap struct {
	mMIDMask  map[string]string   //
	mMIDHash  map[string]string   //
	lsMID     []string            //
	mUIDlsCTX map[string][]string //
	mCTXlsUID map[string][]string //
}

// NewDBByMap :
func NewDBByMap() interface{} {
	db := &memMap{}
	return db.init()
}

func (db *memMap) init() *memMap {
	db.mMIDMask = make(map[string]string)
	db.mMIDHash = make(map[string]string)
	db.mUIDlsCTX = make(map[string][]string)
	db.mCTXlsUID = make(map[string][]string)
	return db
}

// GenPolicyCode :
func (db *memMap) GenPolicyCode(policy string) string {
	jkvM := jkv.NewJKV(policy, hash(policy))
	object := jkvM.LsL12Fields[1][0]
	fields := jkvM.LsL12Fields[2]
	sort.Strings(fields)
	oCode := hash(object)[:lenOfOID]
	fCode := hash(sJoin(fields, ""))[:lenOfFID]
	return oCode + fCode
}

// GenPolicyID :
func (db *memMap) GenPolicyID(policy, uid, ctx, rw string) string {
	code := db.GenPolicyCode(policy)
	suffix := hash(uid + ctx + rw)[:lenOfSID]
	return code + suffix
}

// UpdatePolicy :
func (db *memMap) UpdatePolicy(policy, uid, ctx, rw string) error {

	// check & format policy
	if !jkv.IsJSON(policy) {
		return errors.New("Not a valid JSON")
	}

	policy = pp.FmtJSONStr(policy)
	//

	id := db.GenPolicyID(policy, uid, ctx, rw)
	db.mMIDMask[id] = policy
	db.mMIDHash[id] = hash(policy)
	db.lsMID = u.MapKeys(db.mMIDMask).([]string)

	// for further query
	lsCTX := db.mUIDlsCTX[uid]
	if !xin(ctx, lsCTX) {
		db.mUIDlsCTX[uid] = append(lsCTX, ctx)
	}
	lsUID := db.mCTXlsUID[ctx]
	if !xin(uid, lsUID) {
		db.mCTXlsUID[ctx] = append(lsUID, uid)
	}

	return nil
}

func (db *memMap) GetPolicyID(uid, ctx, object, rw string) (lsID []string) {
	oid := hash(object)[:lenOfOID]
	sid := hash(uid + ctx + rw)[:lenOfSID]
	for _, id := range db.lsMID {
		if sHasPrefix(id, oid) && sHasSuffix(id, sid) {
			lsID = append(lsID, id)
		}
	}
	return lsID
}

func (db *memMap) GetPolicyHash(id string) (string, bool) {
	if hashcode, ok := db.mMIDHash[id]; ok {
		return hashcode, ok
	}
	return "", false
}

func (db *memMap) GetPolicy(id string) (string, bool) {
	if mask, ok := db.mMIDMask[id]; ok {
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
