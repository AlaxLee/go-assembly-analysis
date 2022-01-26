package typeinfo

import "fmt"

func (d *displayer) analysisSlice() {
	//type slicetype struct {
	//	typ  _type
	//	elem *_type
	//}
	ti := d.ti

	offset := d.analysisType()

	offset = d.add(offset, 8,
		fmt.Sprintf("elem:\t%s", ti.getRelocateFullName(offset, 8)))

	d.analysisUncommon(offset)
}
