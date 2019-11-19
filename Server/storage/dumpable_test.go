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
	ctx := "ctx123456"
	fPln(db.ListUIDByCTX(ctx))
}

func TestListPIDByUID(t *testing.T) {
	glb.Init()
	db := NewDB("badger").(Dump)
	uid := "qmiao"
	fPln(db.ListPIDByUID(uid, "r"))
}

func TestListPIDByCTX(t *testing.T) {
	glb.Init()
	db := NewDB("badger").(Dump)
	ctx := "ctx123"
	fPln(db.ListPIDByCTX(ctx, "r"))
}
