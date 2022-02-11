package typeinfo

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

type Relocate struct {
	Addr       int
	Width      int
	Reloctype  int
	Name       string
	NameOffset int
}

type TypeInfo struct {
	TypeName    string
	TypeKind    reflect.Kind
	Data        []byte
	RelocateMap map[string]*Relocate
}

func NewTypeInfo(name string, Content []string) *TypeInfo {
	// get size
	sizeString := typeSizeRegex.FindStringSubmatch(Content[0])[1]
	size, err := strconv.Atoi(sizeString)
	if err != nil {
		panic(err)
	}
	//fmt.Println(size)

	// get binary data
	data := make([]byte, 0, size)
	relocates := make([]*Relocate, 0, 10)
	relocateMap := make(map[string]*Relocate, 10) // key: Relocate.addr+Relocate.width, example 0+8
	//relocates := make([]Relocate, 0, 10)
	for i := 1; i < len(Content); i++ {
		switch {
		case typeBinaryDataHeaderRegex.MatchString(Content[i]):
			bs := typeBinaryDataRegex.FindStringSubmatch(Content[i])
			//fmt.Println(bs)
			for j := 0; j < len(bs[1]); j += 3 {
				d, err := strconv.ParseUint(bs[1][j:j+2], 16, 8)
				if err != nil {
					panic(err)
				}
				data = append(data, byte(d))
			}
		case typeRelocateHeaderRegex.MatchString(Content[i]):
			rs := typeRelocateRegex.FindStringSubmatch(Content[i])
			//fmt.Println(Content[i], rs)
			r := &Relocate{
				Addr:       stringToInt(rs[1]),
				Width:      stringToInt(rs[2]),
				Reloctype:  stringToInt(rs[3]),
				Name:       rs[4],
				NameOffset: stringToInt(rs[5]),
			}
			relocates = append(relocates, r)
			position := fmt.Sprintf("%d+%d", r.Addr, r.Width)
			relocateMap[position] = r
		}
	}

	//fmt.Println(data)
	//fmt.Println(relocates)

	typeKind := getTypeKindFromTypeData(data)

	ti := &TypeInfo{TypeName: name, TypeKind: typeKind, Data: data, RelocateMap: relocateMap}
	return ti
}

func (ti *TypeInfo) getRelocate(addr int, width int) *Relocate {
	return ti.RelocateMap[fmt.Sprintf("%d+%d", addr, width)]
}

func (ti *TypeInfo) getRelocateFullName(addr int, width int) string {
	if r, ok := ti.RelocateMap[fmt.Sprintf("%d+%d", addr, width)]; ok {
		return fmt.Sprintf("%s+%d", r.Name, r.NameOffset)
	} else {
		return "unknown"
	}
}

func (ti *TypeInfo) getUintFromData(offset, length int) uint {
	return byteToInt(ti.Data[offset : offset+length])
}

func (ti *TypeInfo) Display() {
	fmt.Printf("\n%s is a %s, and type size is %d\n", ti.TypeName, ti.TypeKind, len(ti.Data))
	d := NewDisplayer(ti)
	d.display()
}

func getTypeKindFromTypeData(data []byte) reflect.Kind {
	return reflect.Kind(data[23] & 0b11111)
}

func byteToInt(s []byte) (u uint) {
	if len(s) == 0 {
		return
	}
	lastIndex := len(s) - 1
	u = uint(s[lastIndex])
	for i := lastIndex - 1; i >= 0; i-- {
		u = u<<8 + uint(s[i])
	}
	return u
}

func stringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

var typeSizeRegex = regexp.MustCompile(`^\S.+?size=(\d+)`)
var typeBinaryDataHeaderRegex = regexp.MustCompile(`^\s+0x`)
var typeBinaryDataRegex = regexp.MustCompile(`0x[abcdef\d]+\s((?:[abcdef\d][abcdef\d]\s)+)\s`)
var typeRelocateHeaderRegex = regexp.MustCompile(`^\s+rel`)
var typeRelocateRegex = regexp.MustCompile(`^\s+rel\s(\d+)\+(\d+)\st=(-?\d+)\s(.+)\+(\d+)`)
