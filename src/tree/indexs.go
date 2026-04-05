package tree

import (
	"errors"
	"fmt"
)

type Indexs []int64

func (id Indexs) Equal(t Indexs) bool {
	if len(id) != len(t) {
		return false
	}

	for k, v := range id {
		if v != t[k] {
			return false
		}
	}
	return true
}

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
	zoom   ZoomSetLevel // ＃8633
	indexs Indexs
	value  any
}

func CreateRecord(indexs Indexs, zoom ZoomSetLevel, value any) *Record {
	return &Record{
		zoom:   zoom,
		indexs: indexs,
		value:  value,
	}
}

type Records []*Record

func (r *Record) Get() (Indexs, ZoomSetLevel, any) {
	return r.indexs, r.zoom, r.value
}

func (rs Records) Equal(t Records, equalValue func(a, b any) bool) bool {
	if len(rs) != len(t) {
		return false
	}

	for _, v := range rs {
		match := false
	insideLoop:
		for _, vv := range t {
			if v.zoom == vv.zoom && v.indexs.Equal(vv.indexs) && equalValue(v.value, vv.value) {
				match = true
				break insideLoop
			}
		}
		if !match {
			return false
		}
	}
	return true
}
