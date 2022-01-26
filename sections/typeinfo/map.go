package typeinfo

import "fmt"

func (d *displayer) analysisMap() {
	//type maptype struct {
	//	typ    _type
	//	key    *_type
	//	elem   *_type
	//	bucket *_type // internal type representing a hash bucket
	//	// function for hashing keys (ptr to key, seed) -> hash
	//	hasher     func(unsafe.Pointer, uintptr) uintptr
	//	keysize    uint8  // size of key slot
	//	elemsize   uint8  // size of elem slot
	//	bucketsize uint16 // size of bucket
	//	flags      uint32
	//}
	//func (mt *maptype) indirectkey() bool { // store ptr to key instead of key itself
	//	return mt.flags&1 != 0
	//}
	//func (mt *maptype) indirectelem() bool { // store ptr to elem instead of elem itself
	//	return mt.flags&2 != 0
	//}
	//func (mt *maptype) reflexivekey() bool { // true if k==k for all keys
	//	return mt.flags&4 != 0
	//}
	//func (mt *maptype) needkeyupdate() bool { // true if we need to update key on an overwrite
	//	return mt.flags&8 != 0
	//}
	//func (mt *maptype) hashMightPanic() bool { // true if hash function might panic
	//	return mt.flags&16 != 0
	//}
	ti := d.ti

	offset := d.analysisType()

	offset = d.add(offset, 8,
		fmt.Sprintf("key:\t%s", ti.getRelocateFullName(offset, 8)))

	offset = d.add(offset, 8,
		fmt.Sprintf("elem:\t%s", ti.getRelocateFullName(offset, 8)))

	offset = d.add(offset, 8,
		fmt.Sprintf("bucket:\t%s", ti.getRelocateFullName(offset, 8)))

	offset = d.add(offset, 8,
		fmt.Sprintf("hasher:\t%s", ti.getRelocateFullName(offset, 8)))

	offset = d.add(offset, 1,
		fmt.Sprintf("keysize:\t%d", ti.getUintFromData(offset, 1)))

	offset = d.add(offset, 1,
		fmt.Sprintf("elemsize:\t%d", ti.getUintFromData(offset, 1)))

	offset = d.add(offset, 2,
		fmt.Sprintf("bucketsize:\t%d", ti.getUintFromData(offset, 2)))

	var flags = ti.getUintFromData(offset, 4)
	offset = d.add(offset, 4,
		fmt.Sprintf("flags:\t%d\n", flags)+
			fmt.Sprintf("\t(bit:0)indirectkey:\t%d\n", flags&1)+
			fmt.Sprintf("\t(bit:1)indirectelem:\t%d\n", (flags>>1)&1)+
			fmt.Sprintf("\t(bit:2)reflexivekey:\t%d\n", (flags>>2)&1)+
			fmt.Sprintf("\t(bit:3)needkeyupdate:\t%d\n", (flags>>4)&1)+
			fmt.Sprintf("\t(bit:4)hashMightPanic:\t%d", (flags>>8)&1))

	d.analysisUncommon(offset)
}
