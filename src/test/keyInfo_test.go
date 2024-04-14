package test

import (
	"testing"
	tr "trajectory-github/multidementional-radix-tree/src/tree"
)

//　KeyInfo単体試験

func TestKeyInfo(t *testing.T) {
	// 1次元バイナリ
	sf00 := tr.CreateKeyInfo(nil, tr.Indexs{0b101}, 3, nil)
	if sf00.BranchPath(0) != 1 {
		t.Error("sf00-0")
	}
	if sf00.BranchPath(1) != 0 {
		t.Error("sf00-1")
	}
	if sf00.BranchPath(2) != 1 {
		t.Error("sf00-2")
	}
	if sf00.BranchPath(3) != -1 {
		t.Error("sf00-3")
	}

	// 2次元(1)
	table := tr.CreateZoomSetTableForBinary2D()
	sf01 := tr.CreateKeyInfo(table, tr.Indexs{0b101, 0b001}, 3, nil)
	if v := sf01.BranchPath(0); v != 2 {
		t.Errorf("sf01-0 v=%v", v)
	}
	if v := sf01.BranchPath(1); v != 0 {
		t.Errorf("sf01-1 v=%v", v)
	}
	if v := sf01.BranchPath(2); v != 3 {
		t.Errorf("sf01-2 v=%v", v)
	}
	if v := sf01.BranchPath(3); v != -1 {
		t.Errorf("sf01-3 v=%v", v)
	}

	// 2次元(2)
	sf03 := tr.CreateKeyInfo(table, tr.Indexs{0b000, 0b001}, 3, nil)
	if v := sf03.BranchPath(0); v != 0 {
		t.Errorf("sf03-0 v=%v", v)
	}
	if v := sf03.BranchPath(1); v != 0 {
		t.Errorf("sf03-1 v=%v", v)
	}
	if v := sf03.BranchPath(2); v != 1 {
		t.Errorf("sf03-2 v=%v", v)
	}
	if v := sf03.BranchPath(3); v != -1 {
		t.Errorf("sf03-3 v=%v", v)
	}

	// 1次元テーブル
	table2 := tr.ZoomSetTable{tr.ZoomDiffSet{2}, tr.ZoomDiffSet{1}, tr.ZoomDiffSet{2}}
	sf02 := tr.CreateKeyInfo(table2, tr.Indexs{0b01010}, 3, nil)
	if v := sf02.BranchPath(0); v != 1 {
		t.Errorf("sf02-0 v=%v", v)
	}
	if v := sf02.BranchPath(1); v != 0 {
		t.Errorf("sf02-1 v=%v", v)
	}
	if v := sf02.BranchPath(2); v != 2 {
		t.Errorf("sf02-2 v=%v", v)
	}
	if v := sf02.BranchPath(3); v != -1 {
		t.Errorf("sf02-3 v=%v", v)
	}
}
