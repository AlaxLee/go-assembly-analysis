# go-assembly-analysis

用于通过观察go编译后的汇编代码来学习go的类型系统

## 以一个 slice 为例
slice 的类型声明如下（见源码 runtime/type.go）
```go
type slicetype struct {
	typ  _type
	elem *_type
}
type _type struct {
	size       uintptr
	ptrdata    uintptr // size of memory prefix holding all pointers
	hash       uint32
	tflag      tflag
	align      uint8
	fieldAlign uint8
	kind       uint8
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal func(unsafe.Pointer, unsafe.Pointer) bool
	// gcdata stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, gcdata is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	gcdata    *byte
	str       nameOff
	ptrToThis typeOff
}
```
如下是一个用了 slice 的小段代码
```go
package main

import "fmt"

func main() {
	a := []int{}
	fmt.Println(a)
}
```
观察它的go汇编（这里只保留了slice类型声明相关的部分）
```text
test alax$ go tool compile -S -N -l main.go
...
type.[]int SRODATA dupok size=56
        0x0000 18 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
        0x0010 8e 66 f9 1b 02 08 08 17 00 00 00 00 00 00 00 00  .f..............
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00                          ........
        rel 32+8 t=1 runtime.gcbits.01+0
        rel 40+4 t=5 type..namedata.*[]int-+0
        rel 44+4 t=6 type.*[]int+0
        rel 48+8 t=1 type.int+0
...
```
此时，可以用 go-assembly-analysis 来分析，代码如下
```go
package main

import (
	"github.com/AlaxLee/go-assembly-analysis/analysis"
	"log"
)

func main() {

	a, err := analysis.AnalysisContent(`
type.[]int SRODATA dupok size=56
        0x0000 18 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
        0x0010 8e 66 f9 1b 02 08 08 17 00 00 00 00 00 00 00 00  .f..............
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00                          ........
        rel 32+8 t=1 runtime.gcbits.01+0
        rel 40+4 t=5 type..namedata.*[]int-+0
        rel 44+4 t=6 type.*[]int+0
        rel 48+8 t=1 type.int+0
`, true)
	if err != nil {
		log.Fatalln(err)
	}
	a.DisplayAll()
}
```
执行结果如下：
```text
type.[]int is a slice, and type size is 56
0+48:   _type:
0+8:            size:           24
8+8:            ptrdata:        8
16+4:           hash(0x):       8e 66 f9 1b
20+1:           tflag(0b):              00000010
                (bit:0)tflagUncommon:   0
                (bit:1)tflagExtraStar:  1
                (bit:2)tflagNamed:      0
                (bit:3)tflagRegularMemory:      0
21+1:           align:          8
22+1:           fieldAlign:     8
23+1:           kind(0b):               00010111
                (bit:0-4)typekind:      23  (slice)
                (bit:5)kindDirectIface: 0
                (bit:6)kindGCProg:      0
24+8:           equal:          unknown
32+8:           gcdata:         runtime.gcbits.01+0
40+4:           str:            type..namedata.*[]int-+0
44+4:           ptrToThis:      type.*[]int+0
48+8:   elem:   type.int+0
```