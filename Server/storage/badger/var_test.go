package db

import (
	"testing"
)

func TestGenPolicyID(t *testing.T) {
	policy := `{ "policy": "test" }`
	if policy, e := validate(policy); e == nil {
		fPln(genPolicyID(policy, "", "miao", "ctx123", "w"))
	}
}

func TestGetPolicyID(t *testing.T) {

}
