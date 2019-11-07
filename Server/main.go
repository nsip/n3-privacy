package main

import (
	"fmt"

	db "github.com/nsip/n3-privacy/Server/storage"
	pp "github.com/nsip/n3-privacy/preprocess"
)

func main() {
	policy := pp.FmtJSONFile("../../Server/config/mask.json", "../preprocess/utils")
	db.UpdatePolicy("qm", "ctx1", "r", policy)
	policy, ok := db.GetPolicy("qm", "ctx1", "inquiry_skills", "r")
	fmt.Println(policy, ok)
}
