package typeinfo

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestAnalysisArray(t *testing.T) {
	content := `type.[2]interface {} SRODATA dupok size=72
	0x0000 20 00 00 00 00 00 00 00 20 00 00 00 00 00 00 00   ....... .......
	0x0010 2c 59 a4 f1 02 08 08 11 00 00 00 00 00 00 00 00  ,Y..............
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0040 02 00 00 00 00 00 00 00                          ........
	rel 24+8 t=1 type..eqfunc.[2]interface {}+0
	rel 32+8 t=1 runtime.gcbits.0a+0
	rel 40+4 t=5 type..namedata.*[2]interface {}-+0
	rel 44+4 t=6 type.*[2]interface {}+0
	rel 48+8 t=1 type.interface {}+0
	rel 56+8 t=1 type.[]interface {}+0
`
	result := map[int]analyzedContent{
		9:  {32, 8, "\tgcdata:\t\truntime.gcbits.0a+0"},
		12: {48, 8, "elem:\ttype.interface {}+0"},
		13: {56, 8, "slice:\ttype.[]interface {}+0"},
		14: {64, 8, "len:\t2"},
	}
	if err := check(content, reflect.Array, result); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAnalysisChannel(t *testing.T) {
	content := `type.chan int SRODATA dupok size=64
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 91 55 cb 71 0a 08 08 32 00 00 00 00 00 00 00 00  .U.q...2........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 03 00 00 00 00 00 00 00  ................
	rel 24+8 t=1 runtime.memequal64·f+0
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+4 t=5 type..namedata.*chan int-+0
	rel 44+4 t=6 type.*chan int+0
	rel 48+8 t=1 type.int+0
`
	result := map[int]analyzedContent{
		9:  {32, 8, "\tgcdata:\t\truntime.gcbits.01+0"},
		12: {48, 8, "elem:\ttype.int+0"},
		13: {56, 8, "dir:\t3 (chan)"},
	}
	if err := check(content, reflect.Chan, result); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAnalysisFunction(t *testing.T) {
	content := `type.func("".A, int, string, ...bool) bool SRODATA dupok size=96
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 92 09 cf 39 02 08 08 33 00 00 00 00 00 00 00 00  ...9...3........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 04 00 01 80 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0040 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0050 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+4 t=5 type..namedata.*func(main.A, int, string, ...bool) bool-+0
	rel 44+4 t=6 type.*func("".A, int, string, ...bool) bool+0
	rel 56+8 t=1 type."".A+0
	rel 64+8 t=1 type.int+0
	rel 72+8 t=1 type.string+0
	rel 80+8 t=1 type.[]bool+0
	rel 88+8 t=1 type.bool+0
`
	result := map[int]analyzedContent{
		9:  {32, 8, "\tgcdata:\t\truntime.gcbits.01+0"},
		12: {48, 2, "inCount:\t4"},
		13: {50, 2, "outCount:\n\t(bit:0-14)outCount:\t1\n\t(bit:15)dotdotdot:\t1"},
		14: {56, 8, "in[0]:\ttype.\"\".A+0"},
		18: {88, 8, "out[0]:\ttype.bool+0"},
	}
	if err := check(content, reflect.Func, result); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAnalysisInterface(t *testing.T) {
	content := `type."".I SRODATA size=104
	0x0000 10 00 00 00 00 00 00 00 10 00 00 00 00 00 00 00  ................
	0x0010 d1 cd b2 de 07 08 08 14 00 00 00 00 00 00 00 00  ................
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0040 01 00 00 00 00 00 00 00 01 00 00 00 00 00 00 00  ................
	0x0050 00 00 00 00 00 00 00 00 18 00 00 00 00 00 00 00  ................
	0x0060 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=1 runtime.interequal·f+0
	rel 32+8 t=1 runtime.gcbits.02+0
	rel 40+4 t=5 type..namedata.*main.I.+0
	rel 44+4 t=5 type.*"".I+0
	rel 48+8 t=1 type..importpath."".+0
	rel 56+8 t=1 type."".I+96
	rel 80+4 t=5 type..importpath."".+0
	rel 96+4 t=5 type..namedata.lala-+0
	rel 100+4 t=5 type.func(int, string) bool+0
`
	result := map[int]analyzedContent{
		9:  {32, 8, "\tgcdata:\t\truntime.gcbits.02+0"},
		12: {48, 8, "pkgpath:\ttype..importpath.\"\".+0"},
		14: {56, 8, "\tarray:\ttype.\"\".I+96"},
		18: {80, 4, "\tpkgpath:\ttype..importpath.\"\".+0"},
		25: {100, 4, "\tityp:\ttype.func(int, string) bool+0"},
	}
	if err := check(content, reflect.Interface, result); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAnalysisMap(t *testing.T) {
	content := `type.map[int]string SRODATA dupok size=88
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 26 5c 96 90 02 08 08 35 00 00 00 00 00 00 00 00  &\.....5........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0040 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0050 08 10 d0 00 04 00 00 00                          ........
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+4 t=5 type..namedata.*map[int]string-+0
	rel 44+4 t=6 type.*map[int]string+0
	rel 48+8 t=1 type.int+0
	rel 56+8 t=1 type.string+0
	rel 64+8 t=1 type.noalg.map.bucket[int]string+0
	rel 72+8 t=1 runtime.memhash64·f+0
`
	result := map[int]analyzedContent{
		9:  {32, 8, "\tgcdata:\t\truntime.gcbits.01+0"},
		12: {48, 8, "key:\ttype.int+0"},
		13: {56, 8, "elem:\ttype.string+0"},
		14: {64, 8, "bucket:\ttype.noalg.map.bucket[int]string+0"},
		15: {72, 8, "hasher:\truntime.memhash64·f+0"},
		19: {84, 4, "flags:\t4\n\t(bit:0)indirectkey:\t0\n\t(bit:1)indirectelem:\t0\n\t(bit:2)reflexivekey:\t1\n\t(bit:3)needkeyupdate:\t0\n\t(bit:4)hashMightPanic:\t0"},
	}
	if err := check(content, reflect.Map, result); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAnalysisPointer(t *testing.T) {
	content := `type.*map[int]string SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 56 97 f6 03 08 08 08 36 00 00 00 00 00 00 00 00  V......6........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=1 runtime.memequal64·f+0
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+4 t=5 type..namedata.*map[int]string-+0
	rel 48+8 t=1 type.map[int]string+0
`
	result := map[int]analyzedContent{
		9:  {32, 8, "\tgcdata:\t\truntime.gcbits.01+0"},
		12: {48, 8, "elem:\ttype.map[int]string+0"},
	}
	if err := check(content, reflect.Ptr, result); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAnalysisSlice(t *testing.T) {
	content := `type.[]string SRODATA dupok size=56
	0x0000 18 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 d3 a8 f3 0a 02 08 08 17 00 00 00 00 00 00 00 00  ................
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+4 t=5 type..namedata.*[]string-+0
	rel 44+4 t=6 type.*[]string+0
	rel 48+8 t=1 type.string+0
`
	result := map[int]analyzedContent{
		9:  {32, 8, "\tgcdata:\t\truntime.gcbits.01+0"},
		12: {48, 8, "elem:\ttype.string+0"},
	}
	if err := check(content, reflect.Slice, result); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAnalysisStruct(t *testing.T) {
	content := `type."".A SRODATA size=200
	0x0000 48 00 00 00 00 00 00 00 38 00 00 00 00 00 00 00  H.......8.......
	0x0010 bd 5a 2e 34 07 08 08 19 00 00 00 00 00 00 00 00  .Z.4............
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0040 03 00 00 00 00 00 00 00 03 00 00 00 00 00 00 00  ................
	0x0050 00 00 00 00 02 00 00 00 58 00 00 00 00 00 00 00  ........X.......
	0x0060 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0070 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0080 00 00 00 00 00 00 00 00 10 00 00 00 00 00 00 00  ................
	0x0090 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x00a0 31 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  1...............
	0x00b0 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x00c0 00 00 00 00 00 00 00 00                          ........
	rel 32+8 t=1 runtime.gcbits.4a+0
	rel 40+4 t=5 type..namedata.*main.A.+0
	rel 44+4 t=5 type.*"".A+0
	rel 48+8 t=1 type..importpath."".+0
	rel 56+8 t=1 type."".A+96
	rel 80+4 t=5 type..importpath."".+0
	rel 96+8 t=1 type..namedata.x-+0
	rel 104+8 t=1 type.int+0
	rel 120+8 t=1 type..namedata.Y.+0
	rel 128+8 t=1 type.string+0
	rel 144+8 t=1 type..namedata.Do.+0
	rel 152+8 t=1 type."".Do+0
	rel 168+4 t=5 type..namedata.lala-+0
	rel 172+4 t=27 type.func(int, string) bool+0
	rel 176+4 t=27 "".(*A).lala+0
	rel 180+4 t=27 "".A.lala+0
	rel 184+4 t=5 type..namedata.lele-+0
	rel 188+4 t=27 type.func(int, string, ...bool) bool+0
	rel 192+4 t=27 "".(*A).lele+0
	rel 196+4 t=27 "".A.lele+0
`
	result := map[int]analyzedContent{
		9:  {32, 8, "\tgcdata:\t\truntime.gcbits.4a+0"},
		12: {48, 8, "pkgpath:\ttype..importpath.\"\".+0"},
		14: {56, 8, "\tarray:\ttype.\"\".A+96"},
		18: {80, 4, "\tpkgpath:\ttype..importpath.\"\".+0"},
		24: {96, 8, "\tname:\ttype..namedata.x-+0"},
		34: {160, 8, "\toffsetAnon:\n\t\t(bit:0)embedded:\t1\n\t\t(bit:1-63)offset:\t24"},
		36: {168, 4, "\tname:\ttype..namedata.lala-+0"},
		44: {196, 4, "\ttfn:\t\"\".A.lele+0"},
	}
	if err := check(content, reflect.Struct, result); err != nil {
		t.Errorf(err.Error())
	}
}

func check(content string, typeKind reflect.Kind, result map[int]analyzedContent) error {
	name := getSectionName(content)
	ti := NewTypeInfo(name, strings.Split(content, "\n"))
	d := NewDisplayer(ti)
	d.display()
	if ti.TypeKind != typeKind {
		return errors.New(fmt.Sprintf("content's type `%s` not match `%s`", ti.TypeKind, typeKind))
	}
	for k, v := range result {
		if k >= d.analyzedContents.Len() {
			return errors.New(fmt.Sprintf("result's index %d is out of range 0 ~ %d",
				k, d.analyzedContents.Len()-1))
		}
		if v.offset != d.analyzedContents[k].offset {
			return errors.New(fmt.Sprintf("result[%d].offset %d not equal analyzedContents[%d].offset %d",
				k, v.offset, k, d.analyzedContents[k].offset))
		}
		if v.width != d.analyzedContents[k].width {
			return errors.New(fmt.Sprintf("result[%d].width %d not equal analyzedContents[%d].width %d",
				k, v.width, k, d.analyzedContents[k].width))
		}
		if v.content != d.analyzedContents[k].content {
			return errors.New(fmt.Sprintf("result[%d].content `%s` not equal analyzedContents[%d].content `%s`",
				k, v.content, k, d.analyzedContents[k].content))
		}
	}
	return nil
}

var sectionNameRegex = regexp.MustCompile(`^(\S.+?)\sS[A-Z]+\s.*`)

func getSectionName(s string) string {
	n := sectionNameRegex.FindStringSubmatch(s)
	if len(n) < 2 {
		panic("find name failed in: " + s)
	}
	return n[1]
}
