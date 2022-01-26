package typeinfo

import (
	"fmt"
)

func (d *displayer) analysisStruct() {
	//type structtype struct {
	//	typ     _type
	//	pkgPath name
	//	fields  []structfield
	//}
	//type structfield struct {
	//	name       name
	//	typ        *_type
	//	offsetAnon uintptr
	//}
	//type name struct {
	//	bytes *byte
	//}
	ti := d.ti

	offset := d.analysisType()

	offset = d.add(offset, 8,
		fmt.Sprintf("pkgpath:\t%s", ti.getRelocateFullName(offset, 8)))

	d.add(offset, 24, "fields:")

	var arrayPointerOffset int
	if r := ti.getRelocate(offset, 8); r != nil && r.Name == ti.TypeName {
		arrayPointerOffset = r.NameOffset
	}

	offset = d.add(offset, 8,
		fmt.Sprintf("\tarray:\t%s", ti.getRelocateFullName(offset, 8)))

	var structfieldCount = ti.getUintFromData(offset, 8)
	offset = d.add(offset, 8,
		fmt.Sprintf("\tlen:\t%d", ti.getUintFromData(offset, 8)))

	offset = d.add(offset, 8,
		fmt.Sprintf("\tcap:\t%d", ti.getUintFromData(offset, 8)))

	d.analysisUncommon(offset)

	offset = arrayPointerOffset
	for structfieldNum := uint(0); structfieldNum < structfieldCount; structfieldNum++ {
		d.add(offset, 24,
			fmt.Sprintf("structfield[%d]:", structfieldNum))

		offset = d.add(offset, 8,
			fmt.Sprintf("\tname:\t%s", ti.getRelocateFullName(offset, 8)))

		offset = d.add(offset, 8,
			fmt.Sprintf("\ttyp:\t%s", ti.getRelocateFullName(offset, 8)))

		var offsetAnon = uint64(ti.getUintFromData(offset, 8))
		offset = d.add(offset, 8,
			fmt.Sprintf("\toffsetAnon:\n")+
				fmt.Sprintf("\t\t(bit:0)embedded:\t%d\n", offsetAnon&1)+
				fmt.Sprintf("\t\t(bit:1-63)offset:\t%d", offsetAnon>>1))
	}
}
