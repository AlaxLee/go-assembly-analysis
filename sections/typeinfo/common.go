package typeinfo

import (
	"fmt"
	"reflect"
)

func (d *displayer) analysisType() (offset int) {
	ti := d.ti
	data := ti.Data
	//type _type struct {
	//	size       uintptr
	//	ptrdata    uintptr // size of memory prefix holding all pointers
	//	hash       uint32  // hash of type; avoids computation in hash tables
	//	tflag      tflag   // extra type information flags
	//	align      uint8   // alignment of variable with this type
	//	fieldAlign uint8   // alignment of struct field with this type
	//	kind       uint8   // enumeration for C
	//	// function for comparing objects of this type
	//	// (ptr to object A, ptr to object B) -> ==?
	//	equal func(unsafe.Pointer, unsafe.Pointer) bool
	//	// gcdata stores the GC type data for the garbage collector.
	//	// If the KindGCProg bit is set in kind, gcdata is a GC program.
	//	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	//	gcdata    *byte
	//	str       nameOff   // string form
	//	ptrToThis typeOff   // type for pointer to this type, may be zero
	//}
	typeWidth := 48

	d.add(offset, typeWidth, "_type:")

	offset = d.add(offset, 8,
		fmt.Sprintf("\tsize:\t\t%d", ti.getUintFromData(offset, 8)))

	offset = d.add(offset, 8,
		fmt.Sprintf("\tptrdata:\t%d", ti.getUintFromData(offset, 8)))

	offset = d.add(offset, 4,
		fmt.Sprintf("\thash(0x):\t%02x %02x %02x %02x", data[16], data[17], data[18], data[19]))

	d.flag.tflagUncommon = data[20] & 1
	offset = d.add(offset, 1,
		fmt.Sprintf("\ttflag(0b):\t\t%08b\n", data[20])+
			fmt.Sprintf("\t\t(bit:0)tflagUncommon:\t%d\n", d.flag.tflagUncommon)+
			fmt.Sprintf("\t\t(bit:1)tflagExtraStar:\t%d\n", (data[20]>>1)&1)+
			fmt.Sprintf("\t\t(bit:2)tflagNamed:\t%d\n", (data[20]>>2)&1)+
			fmt.Sprintf("\t\t(bit:3)tflagRegularMemory:\t%d", (data[20]>>3)&1))

	offset = d.add(offset, 1,
		fmt.Sprintf("\talign:\t\t%d", data[21]))

	offset = d.add(offset, 1,
		fmt.Sprintf("\tfieldAlign:\t%d", data[22]))

	typekind := reflect.Kind(data[23] & 0b11111)
	kindGCProg := (data[23] >> 6) & 1
	offset = d.add(offset, 1,
		fmt.Sprintf("\tkind(0b):\t\t%08b\n", data[23])+
			fmt.Sprintf("\t\t(bit:0-4)typekind:\t%d  (%s)\n", typekind, typekind.String())+
			fmt.Sprintf("\t\t(bit:5)kindDirectIface:\t%d\n", (data[23]>>5)&1)+
			fmt.Sprintf("\t\t(bit:6)kindGCProg:\t%d", kindGCProg))

	offset = d.add(offset, 8,
		fmt.Sprintf("\tequal:\t\t%s", ti.getRelocateFullName(offset, 8)))

	if kindGCProg == 1 {
		offset = d.add(offset, 8,
			fmt.Sprintf("\tgcdata:\t\t指向一个 gc 程序")) // 这类案例还没有成功构建过，本行需要在成功构造案例后，再补充完整
	} else {
		offset = d.add(offset, 8,
			fmt.Sprintf("\tgcdata:\t\t%s", ti.getRelocateFullName(offset, 8)))
	}

	offset = d.add(offset, 4,
		fmt.Sprintf("\tstr:\t\t%s", ti.getRelocateFullName(offset, 4)))

	offset = d.add(offset, 4,
		fmt.Sprintf("\tptrToThis:\t%s", ti.getRelocateFullName(offset, 4)))

	return offset
}

var uncommonSize = 16 //byte
func (d *displayer) analysisUncommon(offset int) {
	//type uncommontype struct {
	//	pkgpath nameOff // import path; empty for built-in types like int, string
	//	mcount  uint16 // number of methods
	//	xcount  uint16 // number of exported methods
	//	moff    uint32 // offset from this uncommontype to [mcount]method
	//	_       uint32 // unused
	//}
	//
	//// Method on non-interface type
	//type method struct {
	//	name nameOff // name of method
	//	mtyp typeOff // method type (without receiver)
	//	ifn  textOff // fn used in interface call (one-word receiver)
	//	tfn  textOff // fn used for normal method call
	//}
	//
	//type nameOff int32
	//type typeOff int32
	//type textOff int32
	ti := d.ti

	if d.flag.tflagUncommon != 1 {
		return
	}

	uncommontypeBaseOffset := offset

	d.add(offset, 16, "uncommontype:")

	offset = d.add(offset, 4,
		fmt.Sprintf("\tpkgpath:\t%s", ti.getRelocateFullName(offset, 4)))

	var methodCount = ti.getUintFromData(offset, 2)
	offset = d.add(offset, 2,
		fmt.Sprintf("\tmcount:\t%d", methodCount))

	offset = d.add(offset, 2,
		fmt.Sprintf("\txcount:\t%d", ti.getUintFromData(offset, 2)))

	var methodOffset = ti.getUintFromData(offset, 4)
	offset = d.add(offset, 4,
		fmt.Sprintf("\tmoff:\t%d", ti.getUintFromData(offset, 4)))

	offset = d.add(offset, 4,
		fmt.Sprintf("\t_:\t%d", ti.getUintFromData(offset, 4)))

	offset = uncommontypeBaseOffset + int(methodOffset)
	for methodNum := uint(0); methodNum < methodCount; methodNum++ {
		d.add(offset, 16,
			fmt.Sprintf("method[%d]:", methodNum))

		offset = d.add(offset, 4,
			fmt.Sprintf("\tname:\t%s", ti.getRelocateFullName(offset, 4)))

		offset = d.add(offset, 4,
			fmt.Sprintf("\tmtyp:\t%s", ti.getRelocateFullName(offset, 4)))

		offset = d.add(offset, 4,
			fmt.Sprintf("\tifn:\t%s", ti.getRelocateFullName(offset, 4)))

		offset = d.add(offset, 4,
			fmt.Sprintf("\ttfn:\t%s", ti.getRelocateFullName(offset, 4)))
	}
}
