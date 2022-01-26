package sections

import "regexp"

type Skind int

const (
	Base Skind = iota
	TypeInfo
)

var skindRegexMap = map[Skind]*regexp.Regexp{
	TypeInfo: regexp.MustCompile(`^type\.[^.]`),
}

var sectionCreater = map[Skind]func(*BaseSection) Section{
	Base:     NewBaseSection,
	TypeInfo: NewTypeInfoSection,
}

//var skindName = []string {
//	Base: "base",
//	TypeInfo: "typeinfo",
//}

func getSkind(name string) (sk Skind) {
	if name != "" {
		for k, v := range skindRegexMap {
			if v.MatchString(name) {
				return k
			}
		}
	}
	return sk
}
