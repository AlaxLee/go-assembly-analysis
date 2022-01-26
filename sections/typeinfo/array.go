package typeinfo

import "fmt"

func (d *displayer) analysisArray() {
	//type arraytype struct {
	//	typ   _type
	//	elem  *_type
	//	slice *_type
	//	len   uintptr
	//}
	ti := d.ti

	offset := d.analysisType()

	offset = d.add(offset, 8,
		fmt.Sprintf("elem:\t%s", ti.getRelocateFullName(offset, 8)))

	offset = d.add(offset, 8,
		fmt.Sprintf("slice:\t%s", ti.getRelocateFullName(offset, 8)))

	offset = d.add(offset, 8,
		fmt.Sprintf("len:\t%d", ti.getUintFromData(offset, 8)))

	d.analysisUncommon(offset)
}
