package morton

type Indexs []int64
type ZoomSetLevel uint8
type ZoomLevel uint8
type ZoomSet []ZoomLevel
type ZoomDiffSet []ZoomLevel    // 次元数分要素。0の要素は、1のズームレベルの差。2のべき数。デフォルト1。
type ZoomSetTable []ZoomDiffSet // keyはzoomSetLevel

// zoomSetLevelとzoomSetLevel+1のズームレベル差（2の冪数）
func (zt ZoomSetTable) GetZoomDiff(zoomSetLevel ZoomLevel, dim int) (diff ZoomLevel) {
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
	for k := 0; k < int(zoomSetLevel); k++ {
		zoom += zt.GetZoomDiff(ZoomLevel(k), dim)
	}
	return zoom
}
