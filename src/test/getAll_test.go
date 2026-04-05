package test

import (
	"fmt"
	"testing"

	tr "github.com/trajectoryjp/multidimensional-radix-tree/src/tree"
)

func TestGetAll(t *testing.T) {
	table := tr.ZoomSetTable{
		tr.ZoomDiffSet{4, 4, 2},
		tr.ZoomDiffSet{4, 4, 2},
	}
	tree := tr.CreateTreeValues(table)

	/*
		indexs12 := tr.Indexs{0, 0, 0}
		indexs22 := tr.Indexs{0, 0, 1}
		indexs23 := tr.Indexs{0, 1, 2}
		indexs11 := tr.Indexs{15, 15, 3}
		indexs21 := tr.Indexs{250, 250, 15}

		tree.Append(indexs12, 1, "1-2")
		tree.Append(indexs11, 1, "1-1")
		tree.Append(indexs21, 2, "2-21a")
		tree.Append(indexs21, 2, "2-21b")
		tree.Append(indexs22, 2, "2-22")
		tree.Append(indexs23, 2, "2-23")
	*/
	recordsInput, recordsOutput := makeRecords()
	for _, v := range recordsInput {
		indexs, zoom, value := v.Get()
		tree.Append(indexs, zoom, value)
	}

	records, page := tree.GetAll(10, nil)
	fmt.Println(records, page)
	if len(records) != 5 {
		t.Errorf("0 length=%v", len(records))
	}
	if !recordsOutput.Equal(records, valueMatchFunc) {
		t.Errorf("0 unmatch=%v", records)
	}

	// page test 1
	page = nil
	pnum := 1
	recordsStack := make(tr.Records, 0)
	for k := 0; k < 5; k++ {
		records, page = tree.GetAll(pnum, page)
		i, z, v := records[0].Get()
		fmt.Printf("1 [%d] index=%v:%d value=%v page=%v\n", k, i, z, v, page)
		recordsStack = append(recordsStack, records...)
	}
	if r, p := tree.GetAll(pnum, page); len(r) != 0 {
		t.Errorf("1 not end=%v", r)

	} else if p != nil {
		t.Errorf("1 not end page=%v", p)
	}

	if len(recordsStack) != 5 {
		t.Errorf("1 length=%v", len(records))
	}
	if !recordsOutput.Equal(recordsStack, valueMatchFunc) {
		t.Errorf("1 unmatch=%v", records)
	}

	// page test 2
	page = nil
	pnum = 2
	recordsStack = make(tr.Records, 0)
	for k := 0; k < 3; k++ {
		records, page = tree.GetAll(pnum, page)
		i, z, v := records[0].Get()
		fmt.Printf("2 [%d] index=%v:%d value=%v page=%v\n", k, i, z, v, page)
		recordsStack = append(recordsStack, records...)
	}
	if page != nil {
		t.Errorf("2 not end page=%v", page)
	}

	if len(recordsStack) != 5 {
		t.Errorf("12 length=%v", len(records))
	}
	if !recordsOutput.Equal(recordsStack, valueMatchFunc) {
		t.Errorf("2 unmatch=%v", records)
	}
}

func valueMatchFunc(a, b any) bool {
	aa := a.([]any)
	bb := b.([]any)
	if len(aa) != len(bb) {
		return false
	}
	if aa[0].(string) != bb[0].(string) {
		return false
	}
	return true
}

func makeRecords() (input, output tr.Records) {
	/*
		indexs12 := tr.Indexs{0, 0, 0}
		indexs22 := tr.Indexs{0, 0, 1}
		indexs23 := tr.Indexs{0, 1, 2}
		indexs11 := tr.Indexs{15, 15, 3}
		indexs21 := tr.Indexs{250, 250, 15}

		tree.Append(indexs12, 1, "1-2")
		tree.Append(indexs11, 1, "1-1")
		tree.Append(indexs21, 2, "2-21a")
		tree.Append(indexs21, 2, "2-21b")
		tree.Append(indexs22, 2, "2-22")
		tree.Append(indexs23, 2, "2-23")
	*/

	r12 := tr.CreateRecord(
		tr.Indexs{0, 0, 0},
		1,
		"1-2",
	)
	r22 := tr.CreateRecord(
		tr.Indexs{0, 0, 1},
		2,
		"2-22",
	)
	r23 := tr.CreateRecord(
		tr.Indexs{0, 1, 2},
		2,
		"2-23",
	)
	r11 := tr.CreateRecord(
		tr.Indexs{15, 15, 3},
		1,
		"1-1",
	)
	r21a := tr.CreateRecord(
		tr.Indexs{250, 250, 15},
		2,
		"2-21a",
	)
	r21b := tr.CreateRecord(
		tr.Indexs{250, 250, 15},
		2,
		"2-21b",
	)

	a12 := tr.CreateRecord(
		tr.Indexs{0, 0, 0},
		1,
		[]any{"1-2"},
	)
	a22 := tr.CreateRecord(
		tr.Indexs{0, 0, 1},
		2,
		[]any{"2-22"},
	)
	a23 := tr.CreateRecord(
		tr.Indexs{0, 1, 2},
		2,
		[]any{"2-23"},
	)
	a11 := tr.CreateRecord(
		tr.Indexs{15, 15, 3},
		1,
		[]any{"1-1"},
	)
	a21 := tr.CreateRecord(
		tr.Indexs{250, 250, 15},
		2,
		[]any{"2-21b", "2-21a"},
	)

	return tr.Records{r12, r22, r23, r11, r21a, r21b}, tr.Records{a12, a22, a23, a11, a21}
}
