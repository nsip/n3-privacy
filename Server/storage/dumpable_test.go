package storage

import (
	"testing"

	glb "github.com/nsip/n3-privacy/Server/global"
)

func TestListCTXByUID(t *testing.T) {
	glb.Init()
	db := NewDB("badger").(Dump)
	uid := "qmiao"
	fPln(db.ListCTXByUID(uid))
}

func TestListUIDByCTX(t *testing.T) {
	glb.Init()
	db := NewDB("badger").(Dump)
	ctx := "ctx1234"
	fPln(db.ListUIDByCTX(ctx))
}
