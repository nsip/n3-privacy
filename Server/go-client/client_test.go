package client

import (
	"io/ioutil"
	"testing"
)

func TestDO(t *testing.T) {

	policy, err := ioutil.ReadFile("./data/policy.json")
	failOnErr("%v", err)

	data, err := ioutil.ReadFile("./data/file.json")
	failOnErr("%v", err)

	str, err := DO(
		"./config.toml",
		"Enforce", // HELP GetID GetHash Get Update Delete Enforce LsID LsContext LsUser LsObject
		Args{
			ID:     "1615307cc4bf38ffcad912dea96fec4024700fd9r", // 1615307cc4bf38ffcad912dea96fec4024700fd9r
			Policy: policy,
			User:   "user",
			Ctx:    "ctx",
			RW:     "r",
			Object: "object",
			Data:   data,
		})
	fPln(str)
	fPln(err)
}
