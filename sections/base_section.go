package sections

type BaseSection struct {
	Name    string
	Kind    Skind
	Content []string
}

func (bs *BaseSection) Display() {}

func NewBaseSection(bs *BaseSection) Section { return bs }
