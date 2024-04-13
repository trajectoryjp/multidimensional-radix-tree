package tree

import (
	"errors"
	"fmt"
)

type Indexs []int64

//----------------
// for debug mode
//----------------

// zoomleveよりも上位にビットがセットされていないかチェックする
func (id Indexs) validate(zoomSetLevel ZoomSetLevel, zoomSetOddTable ZoomSetOddTable) error {
	// idの最上位ビットの桁を調べる
	zls := zoomSetOddTable.GetZoomSetOdd(zoomSetLevel)
	for d, v := range id {
		if !checkDigit(v, zls[d]) {
			emsg := fmt.Sprintf("id[%v]:%d over digit %v", id, d, zls[d])
			return errors.New(emsg)
		}
	}
	return nil
}

func checkDigit(v int64, z ZoomLevel) bool {
	var c ZoomLevel
	for ; v != 0; v = v >> 1 {
		if c > z {
			return false
		}
		c++
	}
	return true
}

type Record struct {
	indexs Indexs
	value  interface{}
}

type Records []*Record
