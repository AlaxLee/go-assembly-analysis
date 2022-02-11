package analysis

import (
	"regexp"
	"strings"
)

/*
from "https://github.com/AlaxLee/go-assembly-analysis/blob/master/doc/example/typeinfo/struct/struct.go"
to "https://raw.githubusercontent.com/AlaxLee/go-assembly-analysis/master/doc/example/typeinfo/struct/struct.go"
*/
var isGithubUrl *regexp.Regexp = regexp.MustCompile(`https://github\.com/`)

func getRawCodeUrlFromGithub(url string) string {
	newUrl := isGithubUrl.ReplaceAllString(url, `https://raw.githubusercontent.com/`)
	newUrl = strings.ReplaceAll(newUrl, `/blob/`, `/`)
	return newUrl
}
