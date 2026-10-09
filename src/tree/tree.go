package tree

import "fmt"

// -----------
// Tree
// -----------

type TreeInterface interface {
	Append(indexs Indexs, zoomSetLevel ZoomSetLevel, value interface{})
	IsOverlap(indexs Indexs, zoomSetLevel ZoomSetLevel) bool
	GetValues(indexs Indexs, zoomSetLevel ZoomSetLevel) Records // indexsに一致するnodeの値を返す。indexsの親や子の値は返さない。
	// GetValuesの拡張版。親や子の値の返し方をフラグで指定する。
	GetValuesEx(indexs Indexs, zoomSetLevel ZoomSetLevel, childMode ChildMode, parentMode ParentMode) Records
	GetAll(size int, page Page) (records Records, nextPage Page)
}

// GetValuesExでindexsの子（子孫）の値をどう返すか
type ChildMode int

const (
	ChildNone  ChildMode = iota // indexsの値のみを返す。子の値は返さない（GetValuesと同じ）
	ChildAll                    // indexsの値と、すべての子孫の値を返す
	ChildCover                  // indexsの全域が子孫の値で隙間なく覆われていれば子孫の値を返す。そうでなければindexsの値のみを返す（indexsに値がなければ子孫の値を返す）
)

// GetValuesExでindexsの親（祖先）の値をどう返すか
type ParentMode int

const (
	ParentNone    ParentMode = iota // 親の値は返さない
	ParentAll                       // ChildModeの結果に加えて、値を持つすべての親を返す
	ParentLargest                   // 値を持つ親があれば、もっとも大きなサイズ（zoomSetLevelが最小）の親のみを返す。なければChildModeの結果を返す
)

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
func (tr *Tree) Append(indexs Indexs, zoomSetLevel ZoomSetLevel, value any) {
	key := CreateKeyInfo(tr.zoomSetTable, indexs, zoomSetLevel, tr.zoomSetOddTable)
	tr.top.append(key, value)
}

func (tr *Tree) IsOverlap(indexs Indexs, zoomSetLevel ZoomSetLevel) bool {
	key := CreateKeyInfo(tr.zoomSetTable, indexs, zoomSetLevel, tr.zoomSetOddTable)
	nodeKeys := make(Indexs, len(indexs))
	indexsArray := tr.top.searchKey(key, true, true, nodeKeys)
	return len(indexsArray) > 0
}

func (tr *Tree) GetValues(indexs Indexs, zoomSetLevel ZoomSetLevel) Records {
	key := CreateKeyInfo(tr.zoomSetTable, indexs, zoomSetLevel, tr.zoomSetOddTable)
	nodeKeys := make(Indexs, len(indexs))
	return tr.top.searchKey(key, false, false, nodeKeys)
}

func (tr *Tree) GetValuesEx(indexs Indexs, zoomSetLevel ZoomSetLevel, childMode ChildMode, parentMode ParentMode) Records {
	key := CreateKeyInfo(tr.zoomSetTable, indexs, zoomSetLevel, tr.zoomSetOddTable)
	return tr.top.getValuesEx(key, childMode, parentMode)
}

func (tr *Tree) GetAll(size int, page Page) (records Records, nextPage Page) {
	dim := len(tr.zoomSetTable[0])
	indexs := make(Indexs, dim)
	return tr.top.GetAll(indexs, tr.zoomSetTable, size, page)
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
