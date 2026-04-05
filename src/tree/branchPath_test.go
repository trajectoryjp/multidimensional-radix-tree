package tree

import (
	"testing"
	//tr "github.com/trajectoryjp/multidimensional-radix-tree/src/tree"
)

func TestBranchPath(t *testing.T) {
	/*
		table := ZoomSetTable{
			ZoomDiffSet{4, 4, 2},
			ZoomDiffSet{4, 4, 2},
		}

			tree := tr.CreateTreeValues(table)

			indexs12 := Indexs{0, 0, 0}
			indexs22 := Indexs{0, 0, 1}
			indexs23 := Indexs{0, 1, 2}
			indexs11 := Indexs{15, 15, 3}
			indexs21 := Indexs{250, 250, 15}
	*/

	indexs0 := Indexs{0, 0, 0}
	indexs1_1 := Indexs{0, 0, 1}
	zoomSet1 := ZoomSet{4, 4, 2}
	zoomDiffSet1 := ZoomDiffSet{4, 4, 2}
	branchPath1_1 := convertIndexsToBranchPath(indexs1_1, zoomSet1, zoomSet1, zoomDiffSet1)
	if branchPath1_1 != 1 {
		t.Errorf("1_1 %v", branchPath1_1)
	}
	indexs1_1_r := convertBranchPathToIndexs(indexs0, zoomDiffSet1, branchPath1_1)
	if indexs1_1_r[0] != indexs1_1[0] || indexs1_1_r[1] != indexs1_1[1] || indexs1_1_r[2] != indexs1_1[2] {
		t.Errorf("1_1r %v", branchPath1_1)
	}

	indexs2_2 := Indexs{0x11, 0x22, 0x7} // 0b00010001,0b00100010,0b0111 (17,34,7)
	index1_2 := Indexs{1, 2, 1}          // indexs2_2 の親 0b0001,0b0010,0b01
	zoomSet2 := ZoomSet{8, 8, 4}
	zoomDiffSet2 := ZoomDiffSet{4, 4, 2}
	branchPath2_2 := convertIndexsToBranchPath(indexs2_2, zoomSet2, zoomSet2, zoomDiffSet2)
	if branchPath2_2 != 75 {
		t.Errorf("2_2 %v", branchPath2_2)
	}
	indexs2_2_r := convertBranchPathToIndexs(index1_2, zoomDiffSet2, branchPath2_2)
	if indexs2_2_r[0] != indexs2_2[0] || indexs2_2_r[1] != indexs2_2[1] || indexs2_2_r[2] != indexs2_2[2] {
		t.Errorf("2_2r %v", branchPath2_2)
	}

	// indexs2_2のzoomSetLevelは2、zoomSetLevel1でのブランチ番号
	branchPath2_1 := convertIndexsToBranchPath(indexs2_2, zoomSet2, zoomSet1, zoomDiffSet1)
	if branchPath2_1 != 73 {
		t.Errorf("2_1r %v", branchPath2_1)
	}
	indexs2_1_r := convertBranchPathToIndexs(indexs0, zoomDiffSet1, branchPath2_1)
	if indexs2_1_r[0] != indexs2_2[0]>>4 || indexs2_1_r[1] != indexs2_2[1]>>4 || indexs2_1_r[2] != indexs2_2[2]>>2 {
		t.Errorf("2_1r %v", branchPath2_1)
	}

}
