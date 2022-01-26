package sections

import (
	"github.com/AlaxLee/go-assembly-analysis/sections/typeinfo"
)

func NewTypeInfoSection(bs *BaseSection) Section {
	ts := &TypeInfoSection{BaseSection: bs}
	//fmt.Println(bs.Name)
	ts.TypeInfo = typeinfo.NewTypeInfo(bs.Name, bs.Content)
	return ts
}

type TypeInfoSection struct {
	*BaseSection
	*typeinfo.TypeInfo
}

func (ts *TypeInfoSection) Display() {
	//if ts.Name == `type."".A ` {
	//	fmt.Println(ts.Content)
	//}
	ts.TypeInfo.Display()
	return
}
