package analysis

import (
	"github.com/AlaxLee/go-assembly-analysis/sections"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func AnalysisUrl(url string, isAssemble bool) (*sections.Sections, error) {
	// deal url
	if isGithubUrl.MatchString(url) {
		url = getRawCodeUrlFromGithub(url)
	}

	// download code to a tmp file
	// support proxy from env, example: http_proxy=xxxx;https_proxy=xxxx
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return Analysis(resp.Body, isAssemble)

}

func AnalysisFile(file string, isAssemble bool) (*sections.Sections, error) {
	r, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	return Analysis(r, isAssemble)
}

func AnalysisPipe(pipe io.Reader, isAssemble bool) (*sections.Sections, error) {
	return Analysis(pipe, isAssemble)
}

func AnalysisContent(content string, isAssemble bool) (*sections.Sections, error) {
	r := strings.NewReader(content)

	return Analysis(r, isAssemble)
}

func Analysis(input io.Reader, isAssemble bool) (*sections.Sections, error) {
	r := input
	if !isAssemble {
		dir, err := os.MkdirTemp("", "example")
		if err != nil {
			return nil, err
		}
		defer os.RemoveAll(dir) // clean up

		f, err := os.CreateTemp(dir, "example")
		if err != nil {
			return nil, err
		}
		defer os.Remove(f.Name()) // clean up

		_, err = io.Copy(f, r)
		if err != nil {
			return nil, err
		}

		codeFile := f.Name()
		assembleFile := codeFile + ".o"
		defer os.Remove(assembleFile)

		//get local go
		goBin := runtime.GOROOT() + "/bin/go"

		//go tool compile -S -N -l main.go
		cmd := exec.Command(goBin, "tool", "compile", "-S", "-N", "-l", "-o", assembleFile, codeFile)
		r, err = cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}
		if err := cmd.Start(); err != nil {
			return nil, err
		}
	}
	return sections.NewSections(r)
}
