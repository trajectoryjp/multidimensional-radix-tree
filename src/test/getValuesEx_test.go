package test

import (
	"fmt"
	"testing"

	tr "github.com/trajectoryjp/multidimensional-radix-tree/src/tree"
)

// GetValuesEx
// 一次元バイナリーツリー
//
//	0*   - 00*  - 000*
//	            - 001*      <-- 00は子孫で覆われている
//	     - 01   - 010*      <-- 01は子孫で覆われていない（011なし）
//	1    - 10*
//	     - 11*              <-- 1は値なし、子孫で覆われている
//
// *は値あり
func TestGetValuesEx(t *testing.T) {
	tree := tr.CreateTree(tr.Create1DTable())

	r0 := tr.CreateRecord(tr.Indexs{0b0}, 1, "0")
	r00 := tr.CreateRecord(tr.Indexs{0b00}, 2, "00")
	r000 := tr.CreateRecord(tr.Indexs{0b000}, 3, "000")
	r001 := tr.CreateRecord(tr.Indexs{0b001}, 3, "001")
	r010 := tr.CreateRecord(tr.Indexs{0b010}, 3, "010")
	r10 := tr.CreateRecord(tr.Indexs{0b10}, 2, "10")
	r11 := tr.CreateRecord(tr.Indexs{0b11}, 2, "11")
	for _, r := range (tr.Records{r0, r00, r000, r001, r010, r10, r11}) {
		indexs, zoom, value := r.Get()
		tree.Append(indexs, zoom, value)
	}

	tests := []struct {
		name       string
		indexs     tr.Indexs
		zoom       tr.ZoomSetLevel
		childMode  tr.ChildMode
		parentMode tr.ParentMode
		expected   tr.Records
	}{
		// 1) 子の値は返さない
		{"1 child none", tr.Indexs{0b0}, 1, tr.ChildNone, tr.ParentNone, tr.Records{r0}},
		{"1 no value", tr.Indexs{0b01}, 2, tr.ChildNone, tr.ParentNone, tr.Records{}},
		{"1 no node", tr.Indexs{0b011}, 3, tr.ChildNone, tr.ParentNone, tr.Records{}},

		// 2) 親を含むすべての子の値を返す
		{"2 child all", tr.Indexs{0b0}, 1, tr.ChildAll, tr.ParentNone, tr.Records{r0, r00, r000, r001, r010}},
		{"2 child all no value", tr.Indexs{0b1}, 1, tr.ChildAll, tr.ParentNone, tr.Records{r10, r11}},

		// 3) 子孫で覆われていれば子を返し、そうでなければ親のみを返す
		{"3 covered", tr.Indexs{0b00}, 2, tr.ChildCover, tr.ParentNone, tr.Records{r000, r001}},
		{"3 not covered", tr.Indexs{0b0}, 1, tr.ChildCover, tr.ParentNone, tr.Records{r0}},
		{"3 covered no value", tr.Indexs{0b1}, 1, tr.ChildCover, tr.ParentNone, tr.Records{r10, r11}},
		{"3 not covered no value", tr.Indexs{0b01}, 2, tr.ChildCover, tr.ParentNone, tr.Records{r010}},
		{"3 leaf", tr.Indexs{0b000}, 3, tr.ChildCover, tr.ParentNone, tr.Records{r000}},

		// 4) すべての親を返す
		{"4 parent all", tr.Indexs{0b001}, 3, tr.ChildNone, tr.ParentAll, tr.Records{r001, r00, r0}},
		{"4 parent all no node", tr.Indexs{0b011}, 3, tr.ChildNone, tr.ParentAll, tr.Records{r0}},
		{"4 parent all no parent", tr.Indexs{0b10}, 2, tr.ChildNone, tr.ParentAll, tr.Records{r10}},

		// 5) もっとも大きなサイズの親のみを返す
		{"5 parent largest", tr.Indexs{0b001}, 3, tr.ChildNone, tr.ParentLargest, tr.Records{r0}},
		{"5 parent largest no parent", tr.Indexs{0b10}, 2, tr.ChildNone, tr.ParentLargest, tr.Records{r10}},

		// 組み合わせ
		{"2+4", tr.Indexs{0b00}, 2, tr.ChildAll, tr.ParentAll, tr.Records{r0, r00, r000, r001}},
		{"3+4", tr.Indexs{0b00}, 2, tr.ChildCover, tr.ParentAll, tr.Records{r0, r000, r001}},
		{"2+5", tr.Indexs{0b00}, 2, tr.ChildAll, tr.ParentLargest, tr.Records{r0}},
	}

	for _, tt := range tests {
		records := tree.GetValuesEx(tt.indexs, tt.zoom, tt.childMode, tt.parentMode)
		if !tt.expected.Equal(records, stringMatchFunc) {
			t.Errorf("%s: unmatch. expected=%v records=%v", tt.name, recordsString(tt.expected), recordsString(records))
		}
	}
}

// TreeValuesは1つのindexsの複数の値を展開して返す
func TestGetValuesExTreeValues(t *testing.T) {
	tree := tr.CreateTreeValues(tr.Create1DTable())

	tree.Append(tr.Indexs{0b0}, 1, "0a")
	tree.Append(tr.Indexs{0b0}, 1, "0b")
	tree.Append(tr.Indexs{0b001}, 3, "001a")
	tree.Append(tr.Indexs{0b001}, 3, "001b")

	expected := tr.Records{
		tr.CreateRecord(tr.Indexs{0b0}, 1, "0a"),
		tr.CreateRecord(tr.Indexs{0b0}, 1, "0b"),
		tr.CreateRecord(tr.Indexs{0b001}, 3, "001a"),
		tr.CreateRecord(tr.Indexs{0b001}, 3, "001b"),
	}
	records := tree.GetValuesEx(tr.Indexs{0b001}, 3, tr.ChildNone, tr.ParentAll)
	if !expected.Equal(records, stringMatchFunc) {
		t.Errorf("unmatch. expected=%v records=%v", recordsString(expected), recordsString(records))
	}
}

func recordsString(records tr.Records) []string {
	s := make([]string, 0, len(records))
	for _, r := range records {
		indexs, zoom, value := r.Get()
		s = append(s, fmt.Sprintf("%v:%d=%v", indexs, zoom, value))
	}
	return s
}

func stringMatchFunc(a, b any) bool {
	return a.(string) == b.(string)
}
