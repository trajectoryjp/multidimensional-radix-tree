package test

import (
	"testing"
	tr "trajectory-github/multidementional-radix-tree/src/tree"
)

//　一次元バイナリーツリー

func TestOneDimBinary(t *testing.T) {

	// 前方一致検索テーブル作成
	tree := tr.CreateTree(nil)

	id0 := tr.Indexs{0b0}
	tree.Append(id0, 1, "0b0")

	id001 := tr.Indexs{0b001}
	tree.Append(id001, 3, "0b001")

	id00 := tr.Indexs{0b00}
	tree.Append(id00, 2, "0b00")

	id101 := tr.Indexs{0b101}
	tree.Append(id101, 3, "0b101")

	// 検索
	target0 := tr.Indexs{0b01}
	if !tree.IsOverlap(target0, 2) {
		t.Error("target0")
	}

	target1 := tr.Indexs{0b001}
	if !tree.IsOverlap(target1, 3) {
		t.Error("target1")
	}

	target2 := tr.Indexs{0b00}
	if !tree.IsOverlap(target2, 2) {
		t.Error("target2")
	}

	target3 := tr.Indexs{0b11}
	if tree.IsOverlap(target3, 2) {
		t.Error("target3")
	}

	target4 := tr.Indexs{0b1}
	if !tree.IsOverlap(target4, 1) {
		t.Error("target4")
	}

	target5 := tr.Indexs{0b11} // 不正データ（2ビット目は無視される）
	if !tree.IsOverlap(target5, 1) {
		t.Error("target5")
	}

}
