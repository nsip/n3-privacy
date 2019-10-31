package common

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"log"
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

// FailOnErr :
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
