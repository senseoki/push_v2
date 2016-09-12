package module

import (
	"container/list"
)

// MakeSliceList : list 슬라이스 준비 ...
func MakeSliceList(numberGoroutines int) []*list.List {
	alistSL := make([]*list.List, numberGoroutines)
	for i := range alistSL {
		alistSL[i] = list.New()
	}
	return alistSL
}
