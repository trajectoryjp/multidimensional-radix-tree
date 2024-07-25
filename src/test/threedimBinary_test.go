package test

import (
	"testing"

	tr "github.com/trajectoryjp/multidimensional-radix-tree/src/tree"
)

//　三次元バイナリーツリー

func TestThreeeDimBinary(t *testing.T) {

	//Table
	//二次元デフォルトテーブル（バイナリツリー）
	table := tr.Create3DTable()

	// 前方一致検索テーブル作成
	tree := tr.CreateDebugTree(table, exception)

	id0 := tr.Indexs{0b0, 0b0, 0b0}
	tree.Append(id0, 1, "0b0,0b0, 0b0")

	id001 := tr.Indexs{0b000, 0b001, 0b00}
	tree.Append(id001, 3, "0b000,0b001, 0b00")

	id00 := tr.Indexs{0b00, 0b00, 0b00}
	tree.Append(id00, 2, "0b00,0b00,0b00")

	id123 := tr.Indexs{0b011, 0b101, 0b000}
	tree.Append(id123, 3, "0b011,0b101, 0b000")

	// 検索
	target002 := tr.Indexs{0b001, 0b000, 0b000}
	if !tree.IsOverlap(target002, 2) { // id0xと重複
		t.Error("target002")
	}

	target200 := tr.Indexs{0b100, 0b000, 0b000}
	if tree.IsOverlap(target200, 3) { // id2xは存在しない
		t.Error("target100")
	}

	target120 := tr.Indexs{0b110, 0b000, 0b000}
	if tree.IsOverlap(target120, 3) { // 120は存在しない
		t.Error("target120")
	}
}
