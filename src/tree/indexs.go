package tree

type ZoomSetOddTable []ZoomSet

type KeyInfo struct { // 親のビットを0にしたキー
	zoomSetTable ZoomSetTable
	dimension    int // availableDegitsの次元数に一致
	//availableDegits ZoomSet // suffixの有効桁数。zoomSetLevel=0の場合はズームレベルに一致する。
	zoomSetOddTable ZoomSetOddTable

	ZoomSetLevel ZoomSetLevel
	Indexs       Indexs
}

func CreateKeyInfo(table ZoomSetTable, indexs Indexs, zoomSetLevel ZoomSetLevel) *KeyInfo {
	/*
		availableDegits := make(ZoomSet, len(indexs))
		if table == nil {
			availableDegits = []ZoomLevel{ZoomLevel(zoomSetLevel)}

		} else {
			for k := 0; k < len(indexs); k++ {
				availableDegits[k] = table.GetZoom(zoomSetLevel, k)
			}
		}
	*/

	zoomSetOddTable := make(ZoomSetOddTable, 0)
	for zsl := range table {
		zso := make(ZoomSet, len(indexs))
		for dim := range indexs {
			zso[dim] = table.GetZoom(ZoomSetLevel(zsl), dim)
		}
		zoomSetOddTable = append(zoomSetOddTable, zso)
	}

	return &KeyInfo{
		zoomSetTable: table,
		dimension:    len(indexs),
		//availableDegits: availableDegits, // suffixの有効桁数
		zoomSetOddTable: zoomSetOddTable,

		ZoomSetLevel: zoomSetLevel, // suffixの基準
		Indexs:       indexs,
	}
}

// 子の階層へのズームレベル差
// 2の冪数
func (ki *KeyInfo) zoomDiff(zsl ZoomSetLevel, dim int) ZoomLevel {
	return ki.zoomSetTable.GetZoomDiff(zsl, dim)
}

// zslからzsl+1へのブランチの分岐番号
func (ki *KeyInfo) BranchPath(zsl ZoomSetLevel) (branch int) {
	if ki.ZoomSetLevel <= zsl {
		return -1

	} else {
		zsmsb := ki.zoomSetOdd(zsl)     // 上位ビット数
		zslsb := ki.zoomSetOdd(zsl + 1) // 次のレベルまでのビット数
		zoomSetOdd := ki.zoomSetOdd(ki.ZoomSetLevel)
		for d := 0; d < ki.dimension; d++ {
			if d > 0 {
				branch = branch << ki.zoomSetTable.GetZoomDiff(zsl, d-1)
			}
			n := pickup(ki.Indexs[d], zoomSetOdd[d], zsmsb[d], zslsb[d])
			branch = branch | n

		}
		return branch
	}
}

// zslの累積ズームレベルを返す
//
//	2を基数とするZoomSetLevel
//	zsl=0 zoom=0
//	zsl=1 zoom=1
//	zsl=2 zoom=2
func (ki *KeyInfo) zoomSetOdd(zsl ZoomSetLevel) ZoomSet {
	if len(ki.zoomSetOddTable) == 0 {
		//z := 0b01<<(zsl+1) - 1
		zs := make(ZoomSet, ki.dimension)
		for dim := 0; dim < ki.dimension; dim++ {
			//zs[dim] = ZoomLevel(z)
			zs[dim] = ZoomLevel(zsl)
		}
		return zs

	} else if int(zsl) >= len(ki.zoomSetOddTable) {
		maxZs := len(ki.zoomSetOddTable) - 1
		zs := ki.zoomSetOddTable[maxZs]
		for l := maxZs + 1; l < int(zsl); l++ {
			for k := 0; k < ki.dimension; k++ {
				zs = append(zs, 1)
			}
		}
		return zs

	} else {
		if zsl <= 0 {
			return make(ZoomSet, ki.dimension)

		} else {
			return ki.zoomSetOddTable[zsl-1]
		}

	}

}

func pickup(index int64, zoomLevel, zsMsb, zsLsb ZoomLevel) int {
	// zsMsbまでビットクリア
	cmask := 0b1<<(zoomLevel-zsMsb) - 1
	index = index & int64(cmask)
	// zsLsbまで捨てる
	index = index >> (zoomLevel - zsLsb)
	return int(index)
}
