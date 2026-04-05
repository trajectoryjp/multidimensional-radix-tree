package tree

// indexsに対応するbranchPathを得る
// zoomSet：indexsのズームレベル
// brachPathはindexsのtargetZoomSet-1 から　targetZoomSetへの分岐
// digit：indexsのbranchに該当する桁数（2のべき数）。targetZoomSet-1 から　targetZoomSetの桁数。
func convertIndexsToBranchPath(indexs Indexs, zoomSet, targetZoomSet ZoomSet, digit ZoomDiffSet) (branchPath int) {
	dim := len(indexs)
	//lengths := ki.zoomSetOdd(ki.ZoomSetLevel) // 有効桁数
	for d := 0; d < dim; d++ {
		//digit := ki.zoomSetTable.GetZoomDiff(zsl, d)
		n := pickup(indexs[d], zoomSet[d], targetZoomSet[d], digit[d])
		branchPath = branchPath<<digit[d] | n
	}
	return branchPath
}

// parentIndexsにbranchPathのビットを付加する
// 付加する桁数はdigit
// 付加処理のみなのでズームレベルは不要　（indexsのズームレベルはわからない）
func convertBranchPathToIndexs(parentIndexs Indexs, digit ZoomDiffSet, branchPath int) (indexs Indexs) {
	indexs = make(Indexs, len(parentIndexs))
	for dim := len(parentIndexs) - 1; dim >= 0; dim-- {
		zd := digit[dim]
		mask := 0b01<<zd - 1
		bp := branchPath & mask
		indexs[dim] = parentIndexs[dim]<<zd | int64(bp)
		branchPath = branchPath >> zd
	}
	return indexs
}

// indexのbaseを基準にしたdigit数分のビットを取り出す
// baseは上位ビットからの桁数
//
// 例
//
//	index = 0b00011011 length=8
//	base  = 4
//	digit = 2
//	の時
//	int = 01
//
//	 index
//	  0
//	  0
//	  0          <-
//	  1  base=4  <-
//	  1
//	  0
//	  1
//	  1
//
//	index = 0b00011011 length=8
//	base  = 6
//	digit = 2
//	の時
//	int = 10
//
//	 index
//	  0
//	  0
//	  0
//	  1
//	  1          <-
//	  0 base=6   <-
//	  1
//	  1
func pickup(index int64, length, base, digit ZoomLevel) int {
	// msbまでビットクリア
	cmask := 0b1<<(length-base+digit) - 1
	index = index & int64(cmask)
	// lsbまで捨てる
	index = index >> (length - base)
	return int(index)
}
