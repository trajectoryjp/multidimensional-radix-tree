package tree

import "fmt"

// -----------
// Tree
// -----------

type TreeInterface interface {
	Append(indexs Indexs, zoomSetLevel ZoomSetLevel, value interface{})
	IsOverlap(indexs Indexs, zoomSetLevel ZoomSetLevel) bool
}

type Tree struct {
	top          *Node
	zoomSetTable ZoomSetTable

	zoomSetOddTable ZoomSetOddTable
}

func CreateTree(table ZoomSetTable) TreeInterface {
	return &Tree{
		top:             createNode(0),
		zoomSetTable:    table,
		zoomSetOddTable: createZommSetOddTable(table),
	}
}

// indexsの次元はCreateTreeで与えたテーブルの次元数と一致しなければならない
// 処理能力向上のため、次元チェックは行わない。不一致の場合はpanicが発生する。
// valueはnil以外を設定すること（nilはセルなしと扱われる）
func (tr *Tree) Append(indexs Indexs, zoomSetLevel ZoomSetLevel, value interface{}) {
	key := CreateKeyInfo(tr.zoomSetTable, indexs, zoomSetLevel, tr.zoomSetOddTable)
	tr.top.append(key, value)
}

func (tr *Tree) IsOverlap(indexs Indexs, zoomSetLevel ZoomSetLevel) bool {
	key := CreateKeyInfo(tr.zoomSetTable, indexs, zoomSetLevel, tr.zoomSetOddTable)
	nodeKeys := make(Indexs, len(indexs))
	indexsArray := tr.top.searchKey(key, true, nodeKeys)
	return len(indexsArray) > 0
}

// ----------------
//  Tree For Debug
// ----------------
// デバッグ用Tree
// パラメータチェックを実施する
// デバッグ後はCreateDebugTreeをCreateTreeに変えることを推奨する

type DebugTree struct {
	Tree
	exception func(message string)
}

func CreateDebugTree(table ZoomSetTable, exception func(message string)) TreeInterface {
	return &DebugTree{
		Tree: Tree{
			top:             createNode(0),
			zoomSetTable:    table,
			zoomSetOddTable: createZommSetOddTable(table),
		},
		exception: exception,
	}
}

func (tr *DebugTree) Append(indexs Indexs, zoomSetLevel ZoomSetLevel, value interface{}) {
	dim := 1
	if len(tr.zoomSetTable) > 0 {
		dim = len(tr.zoomSetTable[0])
	}
	if len(indexs) != dim {
		emsg := fmt.Sprintf("Append indexs dimension[%v] is unmatch dimension for table[%v]", len(indexs), dim)
		tr.exception(emsg)
	}

	if value == nil {
		emsg := "value shoud not be nil"
		tr.exception(emsg)
	}

	if err := indexs.validate(zoomSetLevel, tr.zoomSetOddTable); err != nil {
		tr.exception(err.Error())
	}
	tr.Tree.Append(indexs, zoomSetLevel, value)
}

func (tr *DebugTree) IsOverlap(indexs Indexs, zoomSetLevel ZoomSetLevel) bool {
	dim := 1
	if len(tr.zoomSetTable) > 0 {
		dim = len(tr.zoomSetTable[0])
	}
	if len(indexs) != dim {
		emsg := fmt.Sprintf("Append indexs dimension[%v] is unmatch dimension for table[%v]", len(indexs), dim)
		tr.exception(emsg)
	}
	if err := indexs.validate(zoomSetLevel, tr.zoomSetOddTable); err != nil {
		tr.exception(err.Error())
	}

	return tr.Tree.IsOverlap(indexs, zoomSetLevel)
}
