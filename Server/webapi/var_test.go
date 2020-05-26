package webapi

import (
	"net/url"
	"testing"
)

func TestUrlParam3(t *testing.T) {

	values := map[string][]string{
		"1": {"11", "111"},
		"2": {"22", "222"},
		"3": {"33", "333"},
	}

	ok, p1, p2, p3 := url3Values(url.Values(values), 1, "1", "2", "3")
	fPln(ok, p1, p2, p3)
}
