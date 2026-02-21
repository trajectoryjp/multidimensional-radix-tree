package test

import (
	"testing"

	tr "github.com/trajectoryjp/multidimensional-radix-tree/src/tree"
)

//　3次元16x164分岐
//  GetValue

func TestTrheeDimGetValue(t *testing.T) {

	table := tr.ZoomSetTable{
		tr.ZoomDiffSet{4, 4, 2},
		tr.ZoomDiffSet{4, 4, 2},
	}
	tree := tr.CreateTreeValues(table)

	indexs11 := tr.Indexs{15, 15, 3}
	indexs12 := tr.Indexs{0, 0, 0}
	indexs21 := tr.Indexs{250, 250, 15}
	indexs22 := tr.Indexs{0, 0, 1}
	indexs23 := tr.Indexs{0, 1, 2}

	tree.Append(indexs12, 1, "1-2")
	tree.Append(indexs11, 1, "1-1")
	tree.Append(indexs21, 2, "2-21a")
	tree.Append(indexs21, 2, "2-21b")
	tree.Append(indexs22, 2, "2-22")
	tree.Append(indexs23, 2, "2-23")

	// 比較
	checkValue(t, tree, indexs21, 2, []string{"2-21a", "2-21b"})
	checkValue(t, tree, indexs22, 2, []string{"2-22"})
	checkValue(t, tree, indexs23, 2, []string{"2-23"})
	checkValue(t, tree, indexs11, 1, []string{"1-1"})
	checkValue(t, tree, indexs12, 1, []string{"1-2"})

}

func checkValue(t *testing.T, tree tr.TreeInterface, indexs tr.Indexs, zoom tr.ZoomSetLevel, values []string) {
	records := tree.GetValues(indexs, zoom)
	if len(records) != len(values) {
		t.Errorf("unmatch len. indexs=%v zoom=%v records=%v values=%v", indexs, zoom, records, values)
	}

	for _, v := range values {
		match := false
	next:
		for _, vv := range records {
			rIndexs, rValue := vv.Get()
			for k, id := range indexs {
				if rIndexs[k] != id {
					t.Errorf("unmatch index. [%d] indexs=%v zoom=%v values=%v", k, indexs, zoom, values)
				}
			}
			if rValue == v {
				match = true
				break next
			}
		}
		if !match {
			t.Errorf("unmatch value. indexs=%v zoom=%v values=%v", indexs, zoom, values)
		}
	}
}
