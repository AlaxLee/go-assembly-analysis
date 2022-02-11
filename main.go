package main

import (
	"flag"
	"github.com/AlaxLee/go-assembly-analysis/analysis"
	"github.com/AlaxLee/go-assembly-analysis/sections"
	"log"
	"os"
)

//文件输入有三种方式，url、file、pipe
var url string
var file string
var pipe bool

func init() {
	flag.StringVar(&url, "code_by_url", "", "code by url")
	flag.StringVar(&url, "u", "", "code by url")

	flag.StringVar(&file, "code_by_file", "", "code by file")
	flag.StringVar(&file, "f", "", "code by file")

	flag.BoolVar(&pipe, "code_by_pipe", false, "code by pipe")
	flag.BoolVar(&pipe, "p", false, "code by pipe")
}

//文件格式两种，默认是go源码，可以是编译后的文件
var assembled bool

func init() {
	flag.BoolVar(&assembled, "assembled", false, "code has assembled")
	flag.BoolVar(&assembled, "a", false, "code has assembled")
}

//结果展示，目前有三种，展示所有section name，展示所有分析后的结果，按section name展示分析结果
var displaySectionNames bool
var displayAllAnalyzedInfo bool
var bySectionName string

func init() {
	flag.BoolVar(&displaySectionNames, "display_section_names", false, "display all section names")
	flag.BoolVar(&displaySectionNames, "ns", false, "display all section names")

	flag.BoolVar(&displayAllAnalyzedInfo, "display_all_analyzed_info", false, "display all analyzed info")
	flag.BoolVar(&displayAllAnalyzedInfo, "aa", false, "display all analyzed info")

	flag.StringVar(&bySectionName, "display_analyzed_info_by_section_name", "", "display analyzed info by section name")
	flag.StringVar(&bySectionName, "an", "", "display analyzed info by section name")
}

func main() {
	flag.Parse()
	var ss *sections.Sections
	var err error
	if url != "" {
		ss, err = analysis.AnalysisUrl(url, assembled)
	} else if file != "" {
		ss, err = analysis.AnalysisFile(file, assembled)
	} else if pipe {
		ss, err = analysis.AnalysisPipe(os.Stdin, assembled)
	} else {
		log.Fatal("param of input error")
	}

	if err != nil {
		log.Fatalln(err)
	}

	if displaySectionNames {
		ss.DisplayNames()
	} else if displayAllAnalyzedInfo {
		ss.DisplayAll()
	} else if bySectionName != "" {
		ss.Display(bySectionName)
	} else {
		log.Fatal("param of display error")
	}
}
