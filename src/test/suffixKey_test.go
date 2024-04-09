package test

import (
	"testing"
	tr "trajectory-github/multidementional-radix-tree/src/tree"
)

//　SuffixKey単体試験

func TestSuffixKey(t *testing.T) {
	sf00 := tr.CreateSuffixKey(nil, tr.Indexs{0b111}, 3)
	sf00.Shift()
	if !sf00.DbgEq(tr.Indexs{0b11}, 1, 1, tr.ZoomSet{2}) {
		t.Error("sf00a-1")
	}
	sf00.Shift()
	if !sf00.DbgEq(tr.Indexs{0b1}, 2, 1, tr.ZoomSet{1}) {
		t.Error("sf00a-2")
	}
	sf00.Shift()
	if !sf00.DbgEq(tr.Indexs{0}, 3, 1, tr.ZoomSet{0}) {
		t.Error("sf00a-3")
	}
	sf00.Shift()
	if !sf00.DbgEq(tr.Indexs{0}, 4, 1, tr.ZoomSet{0}) {
		t.Error("sf00a-4")
	}
}
