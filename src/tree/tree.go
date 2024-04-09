package tree

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
		return zt[zoomSetLevel][dim]
	}
}

// zoomSetLevelのズームレベル（2の冪数）
func (zt ZoomSetTable) GetZoom(zoomSetLevel ZoomSetLevel, dim int) (zoom ZoomLevel) {
	for k := ZoomSetLevel(0); k <= zoomSetLevel; k++ {
		zoom += zt.GetZoomDiff(k, dim)
	}
	return zoom
}

// -----------
// Tree
// -----------
type Tree struct {
	top          *Node
	zoomSetTable ZoomSetTable
}

func CreateTree(table ZoomSetTable) *Tree {
	return &Tree{
		top:          createNode(0),
		zoomSetTable: table,
	}
}

func (tr *Tree) Append(indexs Indexs, zoomSetLevel ZoomSetLevel, value interface{}) {
	key := CreateKeyInfo(tr.zoomSetTable, indexs, zoomSetLevel)
	tr.top.append(key, value)
}

func (tr *Tree) IsOverlap(indexs Indexs, zoomSetLevel ZoomSetLevel) bool {
	key := CreateKeyInfo(tr.zoomSetTable, indexs, zoomSetLevel)
	nodeKeys := make(Indexs, len(indexs))
	indexsArray := tr.top.searchKey(key, true, nodeKeys)
	return len(indexsArray) > 0
}

/*
func (tr *Tree) makeInitSuffixKey(indexs Indexs, zoomSetLevel int) *SuffixKey {
	return CreateSuffixKey(tr.zoomSetTable, indexs, zoomSetLevel)


		availableDegits := make(ZoomSet, len(indexs))
		if tr.zoomSetTable == nil {
			availableDegits = []ZoomLevel{ZoomLevel(zoomSetLevel)}

		} else {
			for k := 0; k < len(indexs); k++ {
				availableDegits[k] = tr.zoomSetTable.GetZoom(zoomSetLevel, k)
			}
		}

		// indexs補正（上位ビットクリア）
		for k := 0; k < len(indexs); k++ {
			zoom := tr.zoomSetTable.GetZoom(zoomSetLevel, k)
			indexs[k] = makeMaskZoom(indexs[k], zoom)
		}

		return &SuffixKey{
			zoomSetTable:    tr.zoomSetTable,
			dimension:       len(indexs),
			zoomSetLevel:    0,               // suffixの基準
			availableDegits: availableDegits, // suffixの有効桁数
			suffix:          indexs,
		}

}
*/
