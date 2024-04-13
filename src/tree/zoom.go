package tree

import "log"

//-------------
// helper
//-------------

func CreateZoomSetTableForBinary1D() ZoomSetTable {
	return ZoomSetTable{ZoomDiffSet{1}}
}

func CreateZoomSetTableForBinary2D() ZoomSetTable {
	return ZoomSetTable{ZoomDiffSet{1, 1}}
}

func CreateZoomSetTableForBinary3D() ZoomSetTable {
	return ZoomSetTable{ZoomDiffSet{1, 1, 1}}
}

func CreateZoomSetTableForQuad() ZoomSetTable {
	return ZoomSetTable{ZoomDiffSet{2}}
}

//-------------
// Zoom
//-------------

// zoom level=0は 値なし
// zoom level=1はindexの1ビット目（0(On),1(Off)の二値）
// zoomSetLevel=0は、zoomLevel=0

type ZoomLevel uint8
type ZoomSetLevel uint8
type ZoomSet []ZoomLevel
type ZoomDiffSet []ZoomLevel    // 次元数分要素。0の要素は、1のズームレベルの差。2のべき数。デフォルト1。
type ZoomSetTable []ZoomDiffSet // keyはzoomSetLevel

// zoomSetLevelとzoomSetLevel+1のズームレベル差（2の冪数）
func (zt ZoomSetTable) GetZoomDiff(zoomSetLevel ZoomSetLevel, dim int) (diff ZoomLevel) {
	if zt == nil {
		return 1

	} else if int(zoomSetLevel) >= len(zt) {
		return 1

	} else {
		if zt[zoomSetLevel] == nil {
			return 1

		} else {
			return zt[zoomSetLevel][dim]
		}
	}
}

// zoomSetLevelのズームレベル（2の冪数）
// zoomSetLevel=0はズームレベル0
// バイナリツリーの場合、zoomSetLevel=1はズームレベル1
func (zt ZoomSetTable) GetZoom(zoomSetLevel ZoomSetLevel, dim int) (zoom ZoomLevel) {
	for k := ZoomSetLevel(0); k < zoomSetLevel; k++ {
		zoom += zt.GetZoomDiff(k, dim)
	}
	return zoom
}

// -----------------
// ZoomSetOddTable
// -----------------
type ZoomSetOddTable []ZoomSet

func createZommSetOddTable(table ZoomSetTable) ZoomSetOddTable {
	dim := 1
	if len(table) > 0 {
		dim = len(table[0])
	}
	zoomSetOddTable := make(ZoomSetOddTable, 0)
	for zsl := 0; zsl <= len(table)+1; zsl++ { // tableサイズのひとつ大きいレベルまでOddTableを作ることで最後のレコードにデフォルトのレベルが入るようにする）
		zso := make(ZoomSet, dim)
		for d := 0; d < dim; d++ {
			zso[d] = table.GetZoom(ZoomSetLevel(zsl), d)
		}
		zoomSetOddTable = append(zoomSetOddTable, zso)
	}
	return zoomSetOddTable
}

func (zo ZoomSetOddTable) GetZoomSetOdd(zoomSetLevel ZoomSetLevel) ZoomSet {
	if len(zo) > int(zoomSetLevel) {
		return zo[zoomSetLevel]

	} else {
		if len(zo) == 0 {
			// 1次元バイナリーツリー
			zoomLevel := ZoomLevel(zoomSetLevel)
			return ZoomSet{zoomLevel}

		} else {
			var dim int
			if dim = len(zo[0]); dim == 0 {
				dim = 1
			}
			zsDiff := make(ZoomSet, dim)

			if lastLevel := len(zo) - 1; lastLevel > 0 {

				for d := 0; d < dim; d++ {
					zsDiff[d] = zo[lastLevel][d] - zo[lastLevel-1][d]
				}

				zs := make(ZoomSet, dim)
				for d := 0; d < dim; d++ {
					zs[d] = zo[lastLevel][d] + zsDiff[d]*ZoomLevel(zoomSetLevel-ZoomSetLevel(lastLevel))
				}
				return zs

			} else {
				// createZommSetOddTableで
				log.Fatal("illegal ZoomSetOddTable")
				return nil
			}
		}

	}
}
