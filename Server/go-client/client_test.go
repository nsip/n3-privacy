package client

import "testing"

func TestDO(t *testing.T) {
	str, err := DO(
		"./config.toml",
		"Enforce", // HELP GetID GetHash Get Update Delete Enforce LsID LsContext LsUser LsObject
		Args{
			// ID:     "8e6f6c7cb618b369ed0a12dea96fec4024700fd9r", // 1615307cc4bf38ffcad912dea96fec4024700fd9r
			Policy: "./data/policy.json",
			User:   "user",
			Ctx:    "ctx",
			RW:     "r",
			Object: "object",
			File:   "./data/file.json",
		})
	fPln(str)
	fPln(err)
}
