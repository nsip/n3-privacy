package db

import (
	"testing"
)

func TestGenPolicyID(t *testing.T) {
	policy := `{ "policy": "test" }`
	if policy, e := validate(policy); e == nil {
		fPln(genPolicyID(policy, "miao", "ctx123", "w"))
	}
}

func TestGetPolicyID(t *testing.T) {

}

// 9f00fad98bda39a3ee5e93dddd64529740f64e18r
// 9f00fad98bda39a3ee5e983ad048af9740f64e18w
