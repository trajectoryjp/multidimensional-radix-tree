package test

import (
	"testing"

	tr "github.com/trajectoryjp/multidimensional-radix-tree/src/tree"
)

//　一次元クアッドキー・バイナリーキー混在ツリー

func TestOneDimMixedQuad(t *testing.T) {
	//ZoomSetTable
	table := tr.ZoomSetTable{
		tr.ZoomDiffSet{2},
		tr.ZoomDiffSet{1},
		tr.ZoomDiffSet{2},
	}

	// 前方一致検索テーブル作成
	tree := tr.CreateTree(table)

	id3 := tr.Indexs{0b11}
	tree.Append(id3, 1, "3")

	id303 := tr.Indexs{0b11011} // 0x1b (27)
	tree.Append(id303, 3, "303")

	id31 := tr.Indexs{0b111}
	tree.Append(id31, 2, "31")

	id002 := tr.Indexs{0b00010}
	tree.Append(id002, 3, "002")

	// 検索
	target0 := tr.Indexs{0b110}      // 30
	if !tree.IsOverlap(target0, 2) { // 3,303が存在するのでOverlab=trueが正解
		t.Error("target0")
	}

	target1 := tr.Indexs{0b111}      // 31
	if !tree.IsOverlap(target1, 2) { // 3,31が存在するのでOverlag=trueが正解
		t.Error("target1")
	}

	target2 := tr.Indexs{0b100}     // 20
	if tree.IsOverlap(target2, 2) { // 20は存在しないのでOverlap=falseが正解
		t.Error("target2")
	}

	target3 := tr.Indexs{0b11010}    // 302
	if !tree.IsOverlap(target3, 3) { // 3が存在するのでOverlag=trueが正解
		t.Error("target3")
	}

	target4 := tr.Indexs{0b11}       // 3
	if !tree.IsOverlap(target4, 1) { // 3が存在するのでOverlag=trueが正解
		t.Error("target4")
	}

	target5 := tr.Indexs{0b101}     // 21
	if tree.IsOverlap(target5, 2) { // 2,21が存在しないのでOverlag=flaseが正解
		t.Error("target5")
	}

	target6 := tr.Indexs{0b111}      // 不正データ（3ビット目は無視） -> 0b11 = 3
	if !tree.IsOverlap(target6, 1) { // 3が存在するのでOverlag=trueが正解
		t.Error("target6")
	}

}
