package go_assembly_analysis

import (
	"github.com/AlaxLee/go-assembly-analysis/sections"
	"io"
	"os"
	"strings"
)

func AnalysisCompiledFile(file string) *sections.Sections {
	r, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	return Analysis(r)
}

func AnalysisCompiledFileContent(content string) *sections.Sections {
	return Analysis(strings.NewReader(content))
}

func Analysis(r io.Reader) *sections.Sections {
	return sections.NewSections(r)
}
