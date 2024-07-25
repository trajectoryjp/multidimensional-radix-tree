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

func (nd *Node) searchKey(key *KeyInfo, chop bool, nodeKeys Indexs) Records {

	if nd.value != nil {

		record := &Record{
			indexs: nodeKeys,
			value:  nd.value,
		}

		if chop {
			return Records{record}

		} else {
			r := nd.searchPrefixToChild(key, chop, nodeKeys)
			return append(r, record)
		}

	} else {
		return nd.searchPrefixToChild(key, chop, nodeKeys)
	}

}

func (nd *Node) searchPrefixToChild(key *KeyInfo, chop bool, nodeKeys Indexs) Records {

	if key.ZoomSetLevel > nd.zoomSetLevel {
		branchPath := key.BranchPath(nd.zoomSetLevel)
		next := nd.next[branchPath]
		if next == nil {
			return nil

		} else {
			for dim := len(nodeKeys) - 1; dim >= 0; dim-- {
				zd := key.zoomSetTable.GetZoomDiff(nd.zoomSetLevel, dim)
				mask := 0b01<<(zd+1) - 1
				bp := branchPath & mask
				nodeKeys[dim] = nodeKeys[dim]<<zd | int64(bp)
			}
			return next.searchKey(key, chop, nodeKeys)
		}

	} else {
		// keyとndは同じZoomSetLevel、もしくはndが大きい（子）のZoomSetLevel
		records := make(Records, 0)
		for branchPath, next := range nd.next {
			if next != nil {
				for dim := len(nodeKeys) - 1; dim >= 0; dim-- {
					zd := key.zoomSetTable.GetZoomDiff(nd.zoomSetLevel, dim)
					mask := 0b01<<(zd+1) - 1
					bp := branchPath & mask
					nodeKeys[dim] = nodeKeys[dim]<<zd | int64(bp)
				}
				if r := next.searchKey(key, chop, nodeKeys); len(r) > 0 {
					records = append(records, r...)
				}
			}

		}
		return records
	}
}
