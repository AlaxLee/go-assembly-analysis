package typeinfo

import "fmt"

func (d *displayer) analysisPointer() {
	//type ptrtype struct {
	//	typ  _type
	//	elem *_type // pointer element (pointed at) type
	//}
	ti := d.ti

	offset := d.analysisType()

	offset = d.add(offset, 8,
		fmt.Sprintf("elem:\t%s", ti.getRelocateFullName(offset, 8)))

	d.analysisUncommon(offset)
}
