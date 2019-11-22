package jkv

const (
	ARR JSONTYPE = 1 << iota
	STR
	BOOL
	NUM
	NULL
	OBJ
)

var (
	JT = map[JSONTYPE]string{
		ARR:        "ARR",
		STR:        "STR",
		BOOL:       "BOOL",
		NUM:        "NUM",
		NULL:       "NULL",
		OBJ:        "OBJ",
		ARR | STR:  "ARR_STR",
		ARR | BOOL: "ARR_BOOL",
		ARR | NUM:  "ARR_NUM",
		ARR | NULL: "ARR_NULL",
		ARR | OBJ:  "ARR_OBJ",
	}
)

// Str : JSON Type string
func (jt JSONTYPE) Str() string {
	return JT[jt]
}

// IsArr : is json array type
func (jt JSONTYPE) IsArr() bool {
	return jt&ARR == ARR
}

// IsObj : is json object type
func (jt JSONTYPE) IsObj() bool {
	return jt&OBJ == OBJ && jt&ARR != ARR
}

// IsObjArr : is json object array type
func (jt JSONTYPE) IsObjArr() bool {
	return jt&OBJ == OBJ && jt&ARR == ARR
}

// IsPrimitive : is json primitive type
func (jt JSONTYPE) IsPrimitive() bool {
	return jt&OBJ != OBJ && jt&ARR != ARR
}

// IsLeafValue : is json Primitive OR Primitive array
func (jt JSONTYPE) IsLeafValue() bool {
	return jt&OBJ != OBJ
}
