package sections

import (
	"bufio"
	"io"
	"regexp"
)

type Section interface {
	Display()
}

type Sections struct {
	NameOrder []string
	Map       map[string]Section
}

func (ss *Sections) DisplayAll() {
	for _, name := range ss.NameOrder {
		if s, ok := ss.Map[name]; ok {
			s.Display()
		} else {
			panic("can't find [" + name + "]'s section")
		}
	}
}

func (ss *Sections) Display(name string) {
	if s, ok := ss.Map[name]; ok {
		s.Display()
	} else {
		panic("can't find [" + name + "]'s section")
	}
}

func NewSections(r io.Reader) *Sections {

	reader := bufio.NewReader(r)

	sectionContents := [][]string{}
	sectionContent := []string{}

	var line string
	var err error

	for {
		line, err = reader.ReadString('\n')
		if err == nil {
			if isSectionHeader(line) {
				// record last section and create a new section
				sectionContents = append(sectionContents, sectionContent)
				sectionContent = []string{}
			}
			sectionContent = append(sectionContent, line)
		} else {
			// end of reader
			sectionContents = append(sectionContents, sectionContent)
			break
		}
	}

	if len(sectionContents) <= 1 {
		return nil
	}

	sectionContents = sectionContents[1:]

	sectionNum := len(sectionContents)

	ss := &Sections{
		NameOrder: make([]string, sectionNum),
		Map:       make(map[string]Section, sectionNum),
	}

	for i := 0; i < sectionNum; i++ {
		name, s := createSection(sectionContents[i])
		ss.NameOrder[i] = name
		ss.Map[name] = s
	}

	return ss
}

func createSection(content []string) (name string, s Section) {
	name = getSectionName(content[0])

	bs := &BaseSection{Content: make([]string, 0, 10)}
	bs.Name = name
	bs.Kind = getSkind(name)
	bs.Content = content

	s = sectionCreater[bs.Kind](bs)
	return name, s
}

var sectionHeaderRegex = regexp.MustCompile(`^\S`)

func isSectionHeader(s string) bool {
	return sectionHeaderRegex.MatchString(s)
}

var sectionNameRegex = regexp.MustCompile(`^(\S.+?)\sS[A-Z]+\s.*`)

func getSectionName(s string) string {
	n := sectionNameRegex.FindStringSubmatch(s)
	if len(n) < 2 {
		panic("find name failed in: " + s)
	}
	return n[1]
}
