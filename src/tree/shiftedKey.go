package tree

type SuffixKey struct { // 親のビットを0にしたキー
	zoomSetTable    ZoomSetTable
	dimension       int     // availableDegitsの次元数に一致
	zoomSetLevel    int     // suffixの基準。0の場合のsuffixはindexsに一致する
	availableDegits ZoomSet // suffixの有効桁数。zoomSetLevel=0の場合はズームレベルに一致する。
	suffix          Indexs
}

func CreateSuffixKey(table ZoomSetTable, indexs Indexs, zoomSetLevel int) *SuffixKey {
	availableDegits := make(ZoomSet, len(indexs))
	if table == nil {
		availableDegits = []ZoomLevel{ZoomLevel(zoomSetLevel)}

	} else {
		for k := 0; k < len(indexs); k++ {
			availableDegits[k] = table.GetZoom(zoomSetLevel, k)
		}
	}
	return &SuffixKey{
		zoomSetTable:    table,
		dimension:       len(indexs),
		zoomSetLevel:    0,               // suffixの基準
		availableDegits: availableDegits, // suffixの有効桁数
		suffix:          indexs,
	}
}

// bitsのzoomレベル部分のビット列を得る
// 例
//
//	index = 0b10010
//	bits = 3
//	の時、0b010
func mask(index int64, bits int) int64 {
	mask := 0x01<<bits - 1
	return index & int64(mask)
}

// zoomSetLevelを+1し
// SuffixKeyの桁をシフトする（現在のzoomSetLevelのキーを除く）
func (sf *SuffixKey) Shift() {
	for k := 0; k < sf.dimension; k++ {
		sf.zoomSetLevel++
		//zoom := sf.zoom(k)
		sf.availableDegits[k] -= sf.zoomDiff(k)
		sf.suffix[k] = mask(sf.suffix[k], int(sf.availableDegits[k]))
	}
}

// 次の階層への分岐番号
// 例
// sf=0b0100 zでzoomLevelの変化が1固定、zoomlevel=4の時、分岐番号は0
// sf=0b1100 zでzoomLevelの変化が1固定、zoomlevel=4の時、分岐番号は1
func (sf *SuffixKey) branchNumbers() (branchNumbers []int64) {
	branchNumbers = make([]int64, sf.dimension)
	for k := 0; k < sf.dimension; k++ {
		shift := sf.availableDegits[k] - sf.zoomDiff(k)
		branchNumbers[k] = sf.suffix[k] >> shift
	}
	return branchNumbers
}

func (sf *SuffixKey) branchPath() (buranchPath int64, branchNumbers []int64) {
	branchNumbers = sf.branchNumbers()
	buranchPath = branchNumbers[0]
	for dim := 1; dim < sf.dimension; dim++ {
		digit := sf.zoomDiff(dim)
		buranchPath = buranchPath << digit
		buranchPath += branchNumbers[dim]
	}
	return buranchPath, branchNumbers
}

func (sf *SuffixKey) IsEnd() bool {
	if sf.availableDegits == nil {
		return true

	} else {
		for _, d := range sf.availableDegits {
			if d > 0 {
				return false
			}
		}
		return true
	}
}

// 子の階層へのズームレベル差
// 2の冪数
func (sf *SuffixKey) zoomDiff(dim int) ZoomLevel {
	return sf.zoomSetTable.GetZoomDiff(sf.zoomSetLevel, dim)
	/*
		if sf.zoomSetTable == nil {
			return 1

		} else {
			if dim > len(sf.zoomSetTable[sf.zoomSetLevel]) {
				return 1
			} else {
				return sf.zoomSetTable[sf.zoomSetLevel][dim]
			}

		}
	*/
}

func (sf *SuffixKey) zoom(dim int) ZoomLevel {
	return sf.zoomSetTable.GetZoom(sf.zoomSetLevel, dim)
}

func (sf *SuffixKey) DbgEq(indexs Indexs, zoomSetLevel, dimension int, availableDegits ZoomSet) bool {
	if sf.dimension != dimension ||
		sf.zoomSetLevel != zoomSetLevel {
		return false

	} else {
		for k, v := range sf.availableDegits {
			if availableDegits[k] != v {
				return false
			}
		}

		for k, v := range sf.suffix {
			if indexs[k] != v {
				return false
			}
		}
	}
	return true
}
