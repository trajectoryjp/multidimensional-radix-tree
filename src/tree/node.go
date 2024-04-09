package tree

import (
	"math"
)

type Node struct {
	zoomSetLevel int
	next         []*Node
	value        interface{}
}

func createNode(zoomSetLevel int) *Node {
	return &Node{
		zoomSetLevel: zoomSetLevel,
	}
}

func (nd *Node) append(suffixKey *SuffixKey, value interface{}) {
	if suffixKey.IsEnd() {
		nd.value = value // over write
		return

	} else {
		if nd.next == nil {
			var branchNum int
			if suffixKey.zoomSetTable == nil {
				//　デフォルト1次元2分木
				branchNum = 2

			} else if suffixKey.dimension == 0 {
				//　デフォルト1次元2分木
				branchNum = 2

			} else if len(suffixKey.zoomSetTable) == 0 || suffixKey.zoomSetLevel >= len(suffixKey.zoomSetTable) {
				//　デフォルト2分木
				if pow := len(suffixKey.zoomSetTable[0]); pow == 0 {
					branchNum = 2
				} else {
					branchNum = int(math.Pow(2, float64(pow)))
				}

			} else {
				//for zs := range suffixKey.zoomSetTable[suffixKey.zoomSetLevel] {
				//	branchNum += int(math.Pow(2, float64(zs)))
				//}
				for dim := 0; dim < suffixKey.dimension; dim++ {
					zdiff := suffixKey.zoomDiff(dim)
					branchNum += int(math.Pow(2, float64(zdiff)))
				}
			}
			nd.next = make([]*Node, branchNum)
		}

		// ブランチ番号判定
		// 次のZoomSetLevelのプリフィックスを取り出し
		// 例
		//   x=01 00 11 10  zoomSetTable=2,2,2
		//   NodeのzoomLevel=1の場合は11を取り出す。11=3がブランチ番号
		/*
			branchNumbers := suffixKey.shift()
			var branchIndex int64
			for dim := 0; dim < suffixKey.dimension; dim++ {
				branchIndex += branchNumbers[dim]
				digit := suffixKey.zoomSetTable[suffixKey.zoomSetLevel][dim]
				branchIndex = branchIndex << digit
			}
		*/

		buranchPath, _ := suffixKey.branchPath()
		if nd.next[buranchPath] == nil {
			nd.next[buranchPath] = createNode(suffixKey.zoomSetLevel + 1)
		}

		suffixKey.Shift()
		nd.next[buranchPath].append(suffixKey, value)

	}
}

func (nd *Node) searchPrefix(suffixKey *SuffixKey, prefixes []int64, chop bool) Records {

	if nd.value != nil {
		record := &Record{
			indexs: prefixes,
			value:  nd.value,
		}

		if chop {
			return Records{record}

		} else {
			r := nd.searchPrefixToChild(suffixKey, prefixes, chop)
			return append(r, record)
		}

	} else {
		return nd.searchPrefixToChild(suffixKey, prefixes, chop)
	}

}

func (nd *Node) searchPrefixToChild(suffixKey *SuffixKey, prefixes []int64, chop bool) Records {
	buranchPath, branchNumbers := suffixKey.branchPath()
	for k, bn := range branchNumbers {
		prefixes[k] = prefixes[k]<<suffixKey.zoomDiff(k) + bn
	}

	if suffixKey.IsEnd() {
		// 子（すべてのpathが対象）
		rr := make(Records, 0)
		for n := range nd.next {
			if next := nd.next[n]; next != nil {
				r := next.searchPrefix(suffixKey, prefixes, chop)
				if len(r) > 0 {
					if chop {
						return r
					} else {
						rr = append(rr, r...)
					}
				}
			}
		}
		return rr

	} else {
		//buranchPath, branchNumbers := suffixKey.branchPath()
		//for k, bn := range branchNumbers {
		//	prefixes[k] = prefixes[k]<<suffixKey.zoom(k) + bn
		//}

		if nd.next != nil && nd.next[buranchPath] != nil {
			suffixKey.Shift()
			return nd.next[buranchPath].searchPrefix(suffixKey, prefixes, chop)

		} else {
			return nil
		}
	}

}
