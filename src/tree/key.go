package tree

/*
type Key struct {
	zoomSetLevel int
	indexs       Indexs // 次元数はtableと一致していること。不一致の場合の動作保証はしない
}
*/

/*
func (key *Key) Parent(zoomSetLevel int, zoomSetTable ZoomSetTable) (parent *Key, suffix Indexs) {

	var prefix Indexs
	prefix, suffix = key.shiftIndexToNext(zoomSetTable)
	parent = &Key{
		zoomSetLevel: zoomSetLevel - 1,
		indexs:       prefix,
	}
	return parent, suffix
}
*/

/*
// keyをzoomDiffTableだけシフトする
func (key *Key) shiftIndexToNext(zoomSetTable ZoomSetTable) (prefix, suffix Indexs) {
	prefix = make(Indexs, len(key.indexs))
	suffix = make(Indexs, len(key.indexs))

	zoomDiffTable := zoomSetTable[key.zoomSetLevel]
	for k := 0; k < len(key.indexs); k++ {
		prefix[k] = key.indexs[k] >> zoomDiffTable[k]
		mask := int64(0b01 << zoomDiffTable[k])
		suffix[k] = key.indexs[k] & mask
	}
	return prefix, suffix
}
*/

type Indexs []int64

//type IndexsArray []Indexs

/*
func (id Indexs) ShiftToNext(zoomdiff ZoomDiffTable) Indexs {

}
*/

type Record struct {
	indexs Indexs
	value  interface{}
}

type Records []*Record
