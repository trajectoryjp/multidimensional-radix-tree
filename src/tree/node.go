package tree

type Node struct {
	zoomSetLevel ZoomSetLevel
	next         []*Node
	value        interface{}
}

func createNode(zoomSetLevel ZoomSetLevel) *Node {
	return &Node{
		zoomSetLevel: zoomSetLevel,
	}
}

func (nd *Node) append(key *KeyInfo, value interface{}) {
	if key.ZoomSetLevel == nd.zoomSetLevel {
		nd.value = value // over write
		return

	} else {
		if nd.next == nil {
			var branchNum int
			if key.zoomSetTable == nil {
				//　デフォルト1次元2分木
				branchNum = 2

			} else if key.dimension == 0 {
				//　デフォルト1次元2分木
				branchNum = 2

			} else if len(key.zoomSetTable) == 0 || int(nd.zoomSetLevel) >= len(key.zoomSetTable) {
				//　デフォルト2分木 x 次元
				branchNum = 0b01 << key.dimension

			} else {
				branchNum = 1 // RM#3910（分岐数計算誤り）
				for dim := 0; dim < key.dimension; dim++ {
					zdiff := key.zoomDiff(nd.zoomSetLevel, dim)
					//branchNum += int(math.Pow(2, float64(zdiff)))
					num := 0b01 << int(zdiff)
					branchNum *= num // RM#3910（分岐数計算誤り）
				}
			}
			nd.next = make([]*Node, branchNum)
		}

		// ブランチ番号判定
		// 次のZoomSetLevelのプリフィックスを取り出し
		// 例
		//   node#0 - node#00 - node#000
		//                    - node#001
		//          - node#01 - node#010
		//                    - node#011
		//          - node#10 - node#110
		//                    - node#111
		//   x=01 00 11 10  zoomSetTable=2,2,2
		//   NodeのzoomLevel=1の場合は11を取り出す。11=3がブランチ番号

		branchPath := key.BranchPath(nd.zoomSetLevel)
		if nd.next[branchPath] == nil {
			nd.next[branchPath] = createNode(nd.zoomSetLevel + 1)
		}
		nd.next[branchPath].append(key, value)
	}
}

// chop：値があれば（nilでなければ）探索を打ち切る。発見した値を返す。
// pileup：pileup=falseでindexsの親の値は返さない。indexsの子は探索しない。
func (nd *Node) searchKey(key *KeyInfo, chop, pileup bool, nodeKeys Indexs) Records { // #8636

	if nd.value != nil {

		record := &Record{
			zoom:   nd.zoomSetLevel, // ＃8633
			indexs: nodeKeys,
			value:  nd.value,
		}

		if chop {
			return Records{record}

		} else {
			r := nd.searchPrefixToChild(key, chop, pileup, nodeKeys)

			//#8637
			if pileup {
				return append(r, record)

			} else {
				if key.ZoomSetLevel == nd.zoomSetLevel {
					// 指定されたindexsと一致するnd
					return Records{record}

				} else {
					return r
				}
			}
		}

	} else {
		return nd.searchPrefixToChild(key, chop, pileup, nodeKeys)
	}

}

func (nd *Node) searchPrefixToChild(key *KeyInfo, chop, pileup bool, nodeKeys Indexs) Records {

	if key.ZoomSetLevel > nd.zoomSetLevel {
		branchPath := key.BranchPath(nd.zoomSetLevel)
		if len(nd.next) == 0 { //#8624
			return nil

		} else {
			next := nd.next[branchPath]
			if next == nil {
				return nil

			} else {
				/*
					for dim := len(nodeKeys) - 1; dim >= 0; dim-- {
						zd := key.zoomSetTable.GetZoomDiff(nd.zoomSetLevel, dim)
						mask := 0b01<<(zd+1) - 1
						bp := branchPath & mask
						nodeKeys[dim] = nodeKeys[dim]<<zd | int64(bp)
					}
				*/
				nodeKeys = nd.ConvertBranchToIndexs(branchPath, nodeKeys, key.zoomSetTable)
				return next.searchKey(key, chop, pileup, nodeKeys)
			}
		}

	} else {
		// keyとndは同じZoomSetLevel、もしくはndが大きい（子）のZoomSetLevel
		if pileup {
			records := make(Records, 0)
			for branchPath, next := range nd.next {
				if next != nil {
					/*
						for dim := len(nodeKeys) - 1; dim >= 0; dim-- {
							zd := key.zoomSetTable.GetZoomDiff(nd.zoomSetLevel, dim)
							mask := 0b01<<(zd+1) - 1
							bp := branchPath & mask
							nodeKeys[dim] = nodeKeys[dim]<<zd | int64(bp)
						}
					*/
					nodeKeys = nd.ConvertBranchToIndexs(branchPath, nodeKeys, key.zoomSetTable)
					if r := next.searchKey(key, chop, true, nodeKeys); len(r) > 0 {
						records = append(records, r...)
					}
				}

			}
			return records

		} else {
			// pileup=falseでは指定されたzoomレベルよりも深くまでは探索しない
			// (指定したindexsの値のみを得ることを目的としているため)
			return nil
		}
	}
}

type Page []int

func (nd *Node) GetAll(myIndexs Indexs, zoomSetTable ZoomSetTable, reamainSize int, startPage Page) (records Records, nextPage Page) {

	records = make(Records, 0)
	if nd.value != nil && len(startPage) == 0 {

		record := &Record{
			zoom:   nd.zoomSetLevel, // ＃8633
			indexs: myIndexs,
			value:  nd.value,
		}
		records = append(records, record)
		reamainSize--

		if reamainSize <= 0 { // 負の場合は想定外
			// 打ち切り
			return records, Page{0}
		}
	}

	startBranchNum := 0
	if len(startPage) > 0 {
		startBranchNum = startPage[0]

		if len(startPage) == 1 {
			startPage = nil
		} else {
			startPage = startPage[1:]
		}
	}

	for branchPath := startBranchNum; branchPath < len(nd.next); branchPath++ {
		//for branchPath, next := range nd.next {
		next := nd.next[branchPath]
		if next != nil {
			childIndexs := nd.ConvertBranchToIndexs(branchPath, myIndexs, zoomSetTable)
			if r, nextPage := next.GetAll(childIndexs, zoomSetTable, reamainSize, startPage); len(r) > 0 {
				records = append(records, r...)
				if reamainSize -= len(r); reamainSize <= 0 {
					// 指定量に到達。打ち切り
					return records, append(Page{branchPath}, nextPage...)
				}
			}
			startPage = nil
		}
	}

	return records, nil
}

// myIndexsをbranchPathに対応する子のIndexsに変換する
func (nd *Node) ConvertBranchToIndexs(branchPath int, myIndexs Indexs, zoomSetTable ZoomSetTable) (childIndexs Indexs) {
	//return convertBranchPathToIndexs(myIndexs, zoomSetTable[nd.zoomSetLevel], branchPath)
	if len(zoomSetTable) > 0 {
		return convertBranchPathToIndexs(myIndexs, zoomSetTable.GetZoomDiff(nd.zoomSetLevel), branchPath)

	} else {
		return convertBranchPathToIndexs(myIndexs, zoomDiffSetUnit(len(myIndexs)), branchPath)
	}
	/*
		childIndexs = make(Indexs, len(myIndexs))
		for dim := len(myIndexs) - 1; dim >= 0; dim-- {
			zd := zoomSetTable.GetZoomDiff(nd.zoomSetLevel, dim)
			mask := 0b01<<(zd+1) - 1
			bp := branchPath & mask
			childIndexs[dim] = myIndexs[dim]<<zd | int64(bp)
			branchPath = branchPath >> zd
		}
		return childIndexs
	*/
}
