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
	return ki.zoomSetTable.GetZoomDiff(zsl, dim)
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
		// ki.Indexs= {0b011011,0b010101)（zoomLevel=6
		//               ^        ^ zsl=1
		//     zsbaseSet = 1,1
		//     digit = 1,1
		//     branch = 00

		lengths := ki.zoomSetOdd(ki.ZoomSetLevel)
		for d := 0; d < ki.dimension; d++ {
			if d > 0 {
				branch = branch << ki.zoomSetTable.GetZoomDiff(zsl, d-1)
			}
			digit := ki.zoomSetTable.GetZoomDiff(zsl, d)
			n := pickup(ki.Indexs[d], lengths[d], zsbaseSet[d], digit)
			branch = branch | n
		}
		return branch
	}
}

// zslの累積ズームレベルを返す
//
//	一次元バイナリーツリーの場合
//	zsl=0 zoomSet=0
//	zsl=1 zoomSet=1  <-- MSBから1bit目がZoomLevel=1のビット
//	zsl=2 zoomSet=2  <-- MSBから2bit目がZoomLevel=2のビット
func (ki *KeyInfo) zoomSetOdd(zsl ZoomSetLevel) ZoomSet {
	return ki.zoomSetOddTable.GetZoomSetOdd(zsl)
	/*
		if len(ki.zoomSetOddTable) == 0 {
			//z := 0b01<<(zsl+1) - 1
			zs := make(ZoomSet, ki.dimension)
			for dim := 0; dim < ki.dimension; dim++ {
				//zs[dim] = ZoomLevel(z)
				zs[dim] = ZoomLevel(zsl)
			}
			return zs

		} else if int(zsl) > len(ki.zoomSetOddTable) {
			// テーブルサイズを超えたレベルについては、テーブルの最後のズームレベル設定値を適用する。
			maxZs := len(ki.zoomSetOddTable) - 1
			//diffZoom := ki.zoomSetTable[maxZs]
			diffLevel := int(zsl) - len(ki.zoomSetOddTable)
			zoomSet := make(ZoomSet, ki.dimension)
			for dim := 0; dim < ki.dimension; dim++ {
				diff := ki.zoomSetTable.GetZoomDiff(maxZs, dim)
				zoomSet[dim] = ki.zoomSetOddTable[maxZs][dim] + diffZoom[dim]*ZoomLevel(diffLevel)
			}
			return zoomSet

		} else {
			if zsl <= 0 {
				return make(ZoomSet, ki.dimension)

			} else {
				return ki.zoomSetOddTable[zsl-1]
			}

		}
	*/
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
//	  0
//	  1  base=4
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
//	  1
//	  0 base=6
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
