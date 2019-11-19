package storage

// Dump :
type Dump interface {
	ListCTXByUID(string) []string
	ListUIDByCTX(string) []string
	ListPIDByUID(string, string) []string
	ListPIDByCTX(string, string) []string
	// ListPIDByOBJ(string, string) []string
}
