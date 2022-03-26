package mr_smart

import "testing"

func TestSlice(t *testing.T) {
	src := []byte{4, 5, 6}
	dest := []byte{1, 2, 3, 4, 5}
	t.Logf("%v ||| %v", src, dest)
	len1 := copy(dest, src)
	t.Logf("%v ||| %v ||| %d", src, dest, len1)
}
