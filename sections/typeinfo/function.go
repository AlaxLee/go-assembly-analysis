package typeinfo

import "fmt"

func (d *displayer) analysisFunction() {
	//type functype struct {
	//	typ      _type
	//	inCount  uint16
	//	outCount uint16 // top bit is set if last input parameter is ...
	//}
	//
	//// 从如下的 in 和 out 方法可知，functype 后面会跟着记录输入参数和输出参数类型
	//
	//func (t *functype) in() []*_type {
	//	// See funcType in reflect/type.go for details on data layout.
	//	uadd := uintptr(unsafe.Sizeof(functype{}))
	//	if t.typ.tflag&tflagUncommon != 0 {
	//	  uadd += unsafe.Sizeof(uncommontype{})
	//  }
	//	return (*[1 << 20]*_type)(add(unsafe.Pointer(t), uadd))[:t.inCount]
	//}
	//
	//func (t *functype) out() []*_type {
	//	// See funcType in reflect/type.go for details on data layout.
	//	uadd := uintptr(unsafe.Sizeof(functype{}))
	//	if t.typ.tflag&tflagUncommon != 0 {
	//	  uadd += unsafe.Sizeof(uncommontype{})
	//  }
	//	outCount := t.outCount & (1<<15 - 1)
	//	return (*[1 << 20]*_type)(add(unsafe.Pointer(t), uadd))[t.inCount : t.inCount+outCount]
	//}
	//
	//func (t *functype) dotdotdot() bool { // outCount's top bit is set if last input parameter is ...
	//	return t.outCount&(1<<15) != 0
	//}
	ti := d.ti

	offset := d.analysisType()

	var inCount = ti.getUintFromData(offset, 2)
	offset = d.add(offset, 2,
		fmt.Sprintf("inCount:\t%d", ti.getUintFromData(offset, 2)))

	var outCountAndDot = uint16(ti.getUintFromData(offset, 2))
	var outCount = outCountAndDot & 0x7F
	var dotdotdot = outCountAndDot >> 15
	offset = d.add(offset, 2,
		fmt.Sprintf("outCount:\n")+
			fmt.Sprintf("\t(bit:0-14)outCount:\t%d\n", outCount)+
			fmt.Sprintf("\t(bit:15)dotdotdot:\t%d", dotdotdot))

	offset = 56 // 48 + 2 + 2 按 8 对齐
	d.analysisUncommon(offset)
	if d.flag.tflagUncommon == 1 {
		offset += uncommonSize
	}

	for i := uint(0); i < inCount; i++ {
		offset = d.add(offset, 8,
			fmt.Sprintf("in[%d]:\t%s", i, ti.getRelocateFullName(offset, 8)))
	}

	for i := uint16(0); i < outCount; i++ {
		offset = d.add(offset, 8,
			fmt.Sprintf("out[%d]:\t%s", i, ti.getRelocateFullName(offset, 8)))
	}
}
