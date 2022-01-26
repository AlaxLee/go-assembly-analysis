package typeinfo

import "fmt"

func (d *displayer) analysisChannel() {
	//type chantype struct {
	//	typ  _type
	//	elem *_type
	//	dir  uintptr
	//}
	// ChanDir represents a channel type's direction.
	//type ChanDir int
	//
	//const (
	//	RecvDir ChanDir             = 1 << iota // <-chan
	//	SendDir                                 // chan<-
	//	BothDir = RecvDir | SendDir             // chan
	//)
	ti := d.ti

	offset := d.analysisType()

	offset = d.add(offset, 8,
		fmt.Sprintf("elem:\t%s", ti.getRelocateFullName(offset, 8)))

	var dir = ti.getUintFromData(offset, 8)
	var dirString string
	switch dir {
	case 1:
		dirString = "<-chan"
	case 2:
		dirString = "chan<-"
	case 3:
		dirString = "chan"
	default:
		dirString = "unkown"
	}
	offset = d.add(offset, 8,
		fmt.Sprintf("dir:\t%d (%s)", dir, dirString))

	d.analysisUncommon(offset)
}
