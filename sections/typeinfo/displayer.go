package typeinfo

import (
	"fmt"
	"reflect"
	"sort"
)

type displayer struct {
	flag
	analyzedContents
	ti *TypeInfo
}

func NewDisplayer(ti *TypeInfo) *displayer {
	d := &displayer{ti: ti}
	d.analysis()
	return d
}

func (a *displayer) analysis() {

	switch a.ti.TypeKind {
	case reflect.Interface:
		a.analysisInterface()
	case reflect.Ptr:
		a.analysisPointer()
	case reflect.Array:
		a.analysisArray()
	case reflect.Struct:
		a.analysisStruct()
	case reflect.Slice:
		a.analysisSlice()
	case reflect.Func:
		a.analysisFunction()
	case reflect.Map:
		a.analysisMap()
	case reflect.Chan:
		a.analysisChannel()
	default:
	}
	a.sort()
}

type analyzedContent struct {
	offset  int
	width   int
	content string
}

type analyzedContents []analyzedContent

func (acs *analyzedContents) add(offset int, width int, content string) (newOffset int) {
	*acs = append(*acs, analyzedContent{offset, width, content})
	newOffset = offset + width
	return newOffset
}

func (acs *analyzedContents) display() {
	for _, ac := range *acs {
		fmt.Printf("%d+%d:\t%s\n", ac.offset, ac.width, ac.content)
	}
}

func (acs *analyzedContents) sort() {
	//offset 小的优先，offset 相同时 width 大的优先
	sort.Sort(acs)
}

// Len is part of sort.Interface.
func (acs *analyzedContents) Len() int {
	return len(*acs)
}

// Swap is part of sort.Interface.
func (acs *analyzedContents) Swap(i, j int) {
	(*acs)[i], (*acs)[j] = (*acs)[j], (*acs)[i]
}

// Less is part of sort.Interface.
func (acs *analyzedContents) Less(i, j int) bool {
	//offset 小的优先，offset 相同时 width 大的优先
	if (*acs)[i].offset < (*acs)[j].offset {
		return true
	} else if (*acs)[i].offset == (*acs)[j].offset {
		if (*acs)[i].width > (*acs)[j].width {
			return true
		}
	}
	return false
}

type flag struct {
	tflagUncommon byte
}
