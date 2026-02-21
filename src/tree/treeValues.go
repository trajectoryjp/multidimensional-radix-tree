package tree

// ----------------
// Tree For Values
// ----------------

// 1つのインデックスに複数の値を保持するツリー

type TreeValues struct { // nodeのvalueは[]any
	Tree
}

func CreateTreeValues(table ZoomSetTable) TreeInterface {
	return &TreeValues{
		Tree: Tree{
			top:             createNode(0),
			zoomSetTable:    table,
			zoomSetOddTable: createZommSetOddTable(table),
		},
	}
}

// すでに存在するindexsであればvalueを追加する
func (tr *TreeValues) Append(indexs Indexs, zoomSetLevel ZoomSetLevel, value any) {
	values := []any{value}
	key := CreateKeyInfo(tr.zoomSetTable, indexs, zoomSetLevel, tr.zoomSetOddTable)
	records := tr.GetValues(indexs, zoomSetLevel)
	for _, r := range records {
		values = append(values, r.value)
	}
	tr.top.append(key, values)
}

func (tr *TreeValues) GetValues(indexs Indexs, zoomSetLevel ZoomSetLevel) (records Records) {
	key := CreateKeyInfo(tr.zoomSetTable, indexs, zoomSetLevel, tr.zoomSetOddTable)
	nodeKeys := make(Indexs, len(indexs))
	rs := tr.top.searchKey(key, false, nodeKeys)

	records = make(Records, 0)
	for _, r := range rs {
		values := r.value.([]any)
		for _, v := range values {
			rec := &Record{
				zoom:   r.zoom, // #8633
				indexs: indexs,
				value:  v,
			}
			records = append(records, rec)
		}
	}
	return records
}
