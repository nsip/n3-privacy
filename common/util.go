package common

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"time"
)

// SHA1Str :
func SHA1Str(s string) string {
	return fSf("%x", sha1.Sum([]byte(s)))
}

// SHA256Str :
func SHA256Str(s string) string {
	return fSf("%x", sha256.Sum256([]byte(s)))
}

// MD5Str :
func MD5Str(s string) string {
	return fSf("%x", md5.Sum([]byte(s)))
}

// TmTrack :
func TmTrack(start time.Time) {
	elapsed := time.Since(start)
	fPf("took %s\n", elapsed)
}

// FailOnErr : error holder use "%v"
func FailOnErr(format string, v ...interface{}) {
	for _, p := range v {
		switch p.(type) {
		case error:
			{
				if p != nil {
					log.Fatalf(format, v...)
				}
			}
		}
	}
}

// WrapOnErr :
func WrapOnErr(format string, v ...interface{}) error {
	for _, p := range v {
		switch p.(type) {
		case error:
			{
				if p != nil {
					return fmt.Errorf(format, v...)
				}
			}
		}
	}
	return nil
}

// LocalIP returns the non loopback local IP of the host
func LocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ""
}

// CanSetCover : check if setA contains setB ? return the first B-Index of which item is not in setA
func CanSetCover(setA, setB interface{}) (bool, int) {
	tA, tB := reflect.TypeOf(setA), reflect.TypeOf(setB)
	if tA != tB || (tA.Kind() != reflect.Slice && tA.Kind() != reflect.Array) {
		FailOnErr("%v", errors.New("parameters only can be [slice] or [array]"))
	}
	vA, vB := reflect.ValueOf(setA), reflect.ValueOf(setB)
	if vA.Len() < vB.Len() {
		return false, -1
	}
NEXT:
	for j := 0; j < vB.Len(); j++ {
		b := vB.Index(j).Interface()
		for i := 0; i < vA.Len(); i++ {
			if reflect.DeepEqual(b, vA.Index(i).Interface()) {
				continue NEXT
			}
			if i == vA.Len()-1 { // if b falls down to the last vA item position, which means vA doesn't have b item, return false
				return false, j
			}
		}
	}
	return true, -1
}

// SetIntersect :
func SetIntersect(setA, setB interface{}) interface{} {
	tA, tB := reflect.TypeOf(setA), reflect.TypeOf(setB)
	if tA != tB || (tA.Kind() != reflect.Slice && tA.Kind() != reflect.Array) {
		FailOnErr("%v", errors.New("parameters only can be [slice] or [array]"))
	}
	vA, vB := reflect.ValueOf(setA), reflect.ValueOf(setB)
	set := reflect.MakeSlice(tA, 0, vA.Len())
NEXT:
	for j := 0; j < vB.Len(); j++ {
		b := vB.Index(j)
		for i := 0; i < vA.Len(); i++ {
			if reflect.DeepEqual(b.Interface(), vA.Index(i).Interface()) {
				set = reflect.Append(set, b)
				continue NEXT
			}
		}
	}
	return set.Interface()
}

// SetUnion :
func SetUnion(setA, setB interface{}) interface{} {
	tA, tB := reflect.TypeOf(setA), reflect.TypeOf(setB)
	if tA != tB || (tA.Kind() != reflect.Slice && tA.Kind() != reflect.Array) {
		FailOnErr("%v", errors.New("parameters only can be [slice] or [array]"))
	}
	vA, vB := reflect.ValueOf(setA), reflect.ValueOf(setB)
	set := reflect.MakeSlice(tA, 0, vA.Len()+vB.Len())
	set = reflect.AppendSlice(set, vA)
	set = reflect.AppendSlice(set, vB)
	return ToSet(set.Interface())
}

// ToSet : convert slice / array to set. i.e. remove duplicated items
func ToSet(slc interface{}) interface{} {
	t := reflect.TypeOf(slc)
	if t.Kind() != reflect.Slice && t.Kind() != reflect.Array {
		FailOnErr("%v", errors.New("parameter only can be [slice] or [array]"))
	}
	v := reflect.ValueOf(slc)
	if v.Len() == 0 {
		return slc
	}

	set := reflect.MakeSlice(t, 0, v.Len())
	set = reflect.Append(set, v.Index(0))
NEXT:
	for i := 1; i < v.Len(); i++ {
		vItem := v.Index(i)
		for j := 0; j < set.Len(); j++ {
			if reflect.DeepEqual(vItem.Interface(), set.Index(j).Interface()) {
				continue NEXT
			}
			if j == set.Len()-1 { // if vItem falls down to the last set position, which means set doesn't have this item, then add it.
				set = reflect.Append(set, vItem)
			}
		}
	}
	return set.Interface()
}
