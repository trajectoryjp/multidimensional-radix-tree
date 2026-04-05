package tree

type KeyInfo struct { // 親のビットを0にしたキー
	zoomSetTable ZoomSetTable
	dimension    int // availableDegitsの次元数に一致
	//availableDegits ZoomSet // suffixの有効桁数。zoomSetLevel=0の場合はズームレベルに一致する。
	zoomSetOddTable ZoomSetOddTable

	ZoomSetLevel ZoomSetLevel
	Indexs       Indexs
}

func CreateKeyInfo(table ZoomSetTable, indexs Indexs, zoomSetLevel ZoomSetLevel, zoomOddTable ZoomSetOddTable) *KeyInfo {
	if zoomOddTable == nil {
		zoomOddTable = createZommSetOddTable(table)
	}

	return &KeyInfo{
		zoomSetTable:    table,
		dimension:       len(indexs),
		zoomSetOddTable: zoomOddTable,

		ZoomSetLevel: zoomSetLevel,
		Indexs:       indexs,
	}
}

// 子の階層へのズームレベル差
// 2の冪数
func (ki *KeyInfo) zoomDiff(zsl ZoomSetLevel, dim int) ZoomLevel {
	return ki.zoomSetTable.GetZoomDiffDim(zsl, dim)
}

// zslからzsl+1へのブランチの分岐番号
//
// quadkeyの例
// 3120  ki.indexs = 0b11011000 ki.zoomSetLevel=4  ki.dimension=1
// zsl=0のbranchは3（zsbase=2,digit=2)
// zsl=1のbranchは1（zsbase=4,digit=2)
// zsl=2のbranchは2（zsbase=6,degit=2)
func (ki *KeyInfo) BranchPath(zsl ZoomSetLevel) (branch int) {
	if ki.ZoomSetLevel <= zsl {
		return -1

	} else {
		zsbaseSet := ki.zoomSetOdd(zsl + 1) // zslの桁の上位からのビット数（一番上位のビットは1）
		// 1次元バイナリーの場合
		//                zoomSetOdd
		//          zsl=0 0
		//          zsl=1 1  (2)
		//          zsl=2 2  (4)
		//          zsl=3 3  (8)
		// ki.Indexs = 0b011011（zoomSetLevel=length=6 ）
		//              ^ zsl　＝ 0
		//     zsbaseSet = 1
		//     digit = 1
		//     branch = 0
		//
		// ki.Indexs = 0b011011（zoomSetLevel=length=6 ）
		//               ^ zsl　＝ 1
		//     zsbaseSet = 2
		//     digit = 1
		//     branch = 1
		//
		// 2次元バイナリーの場合
		//                zoomSetOdd
		//          zsl=0 0,0
		//          zsl=1 1,1  (2x2)
		//          zsl=2 2,2  (4x4)
		//          zsl=3 3,3  (8x8)
		// ki.Indexs= {0b011011,0b010101}（zoomLevel=6)
		//               ^        ^ zsl=1
		//     zsbaseSet = 1,1
		//     digit = 1,1
		//     branch = 00

		lengths := ki.zoomSetOdd(ki.ZoomSetLevel) // 有効桁数
		/*
				for d := 0; d < ki.dimension; d++ {
					digit := ki.zoomSetTable.GetZoomDiff(zsl, d)
					n := pickup(ki.Indexs[d], lengths[d], zsbaseSet[d], digit)
					branch = branch<<digit | n
				}

			return branch
		*/
		//return convertIndexsToBranchPath(ki.Indexs, lengths, zsbaseSet, ki.zoomSetTable[zsl])
		return convertIndexsToBranchPath(ki.Indexs, lengths, zsbaseSet, ki.GetZoomDiff(zsl))
	}
}

// zslの累積ズームレベルを返す
//
//	一次元バイナリーツリーの場合
//	zsl=0 zoomSet=0
//	zsl=1 zoomSet=1
//	zsl=2 zoomSet=2
func (ki *KeyInfo) zoomSetOdd(zsl ZoomSetLevel) ZoomSet {
	return ki.zoomSetOddTable.GetZoomSetOdd(zsl)
}

func (ki *KeyInfo) GetZoomDiff(zoomSetLevel ZoomSetLevel) ZoomDiffSet {
	if ki.zoomSetTable == nil {
		return zoomDiffSetUnit(len(ki.Indexs))

	} else if int(zoomSetLevel) >= len(ki.zoomSetTable) {
		return zoomDiffSetUnit(len(ki.Indexs))

	} else if ki.zoomSetTable[zoomSetLevel] == nil {
		return zoomDiffSetUnit(len(ki.Indexs))

	} else {
		return ki.zoomSetTable[zoomSetLevel]
	}
}

func zoomDiffSetUnit(dim int) ZoomDiffSet {
	zd := make(ZoomDiffSet, dim)
	for d := 0; d < dim; d++ {
		zd[d] = 1
	}
	return zd
}
