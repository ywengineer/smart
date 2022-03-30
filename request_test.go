package mr_smart

import (
	"net/url"
	"testing"
)

func TestSlice(t *testing.T) {
	src := []byte{4, 5, 6}
	dest := []byte{1, 2, 3, 4, 5}
	t.Logf("%v ||| %v", src, dest)
	len1 := copy(dest, src)
	t.Logf("%v ||| %v ||| %d", src, dest, len1)
}

func TestUrl(t *testing.T) {
	p := "/Users/yangwei/.m2/conf.json"
	if u, err := url.Parse(p); err != nil {
		t.Logf("parse url failed. %v", err)
	} else {
		t.Log(u.Opaque)
		t.Log(u.Path)
		t.Log(u.Scheme)
	}
}
