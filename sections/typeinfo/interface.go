package typeinfo

import "fmt"

func (d *displayer) analysisInterface() {
	//type interfacetype struct {
	//	typ     _type
	//	pkgpath name
	//	mhdr    []imethod
	//}
	//
	//type imethod struct {
	//	name nameOff
	//	ityp typeOff
	//}
	//
	//type name struct {
	//	bytes *byte
	//}
	//type nameOff int32
	//type typeOff int32

	ti := d.ti

	offset := d.analysisType()

	offset = d.add(offset, 8,
		fmt.Sprintf("pkgpath:\t%s", ti.getRelocateFullName(offset, 8)))

	d.add(offset, 24, "mhdr:")

	var arrayPointerOffset int
	if r := ti.getRelocate(offset, 8); r != nil && r.Name == ti.TypeName {
		arrayPointerOffset = r.NameOffset
	} else {
		//fmt.Println(r, offset)
		panic("unexpect")
	}
	offset = d.add(offset, 8,
		fmt.Sprintf("\tarray:\t%s", ti.getRelocateFullName(offset, 8)))

	var imethodCount = ti.getUintFromData(offset, 8)
	offset = d.add(offset, 8,
		fmt.Sprintf("\tlen:\t%d", ti.getUintFromData(offset, 8)))

	offset = d.add(offset, 8,
		fmt.Sprintf("\tcap:\t%d", ti.getUintFromData(offset, 8)))

	d.analysisUncommon(offset)

	offset = arrayPointerOffset
	for imethodNum := uint(0); imethodNum < imethodCount; imethodNum++ {
		d.add(offset, 8,
			fmt.Sprintf("imethod[%d]:", imethodNum))

		offset = d.add(offset, 4,
			fmt.Sprintf("\tname:\t%s", ti.getRelocateFullName(offset, 4)))

		offset = d.add(offset, 4,
			fmt.Sprintf("\tityp:\t%s", ti.getRelocateFullName(offset, 4)))
	}
}
