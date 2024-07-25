package test

import (
	"testing"

	tr "github.com/trajectoryjp/multidimensional-radix-tree/src/tree"
)

//　3次元バイナリーツリー

func TestThreeDimBinary2(t *testing.T) {

	//Table
	//3次元デフォルトテーブル（バイナリツリー）
	table := tr.Create3DTable()

	// 前方一致検索テーブル作成
	tree := tr.CreateDebugTree(table, exception)

	id0 := tr.Indexs{4096, 7274, 3225}
	tree.Append(id0, 13, "13/4096/7224/3225") // panic here

	// 検索
	target002 := tr.Indexs{32768, 58198, 25804} // 16/0/58198/25804
	if !tree.IsOverlap(target002, 16) {         // id0xと重複
		t.Error("target002")
	}
}
