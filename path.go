package tree

import "strings"

type Path struct {
	Raw string
}

func ParsePath(path string) *Path {
	p := &Path{}
	p.Raw = strings.TrimSpace(path)
	p.Raw = strings.TrimSuffix(p.Raw, "status")
	p.Raw = strings.Trim(p.Raw, "/")
	//	log.Printf("ParsePath in=%s out=%s\n", path, p.Raw)
	return p
}

func (p *Path) Levels() []string {
	return strings.Split(p.Raw, "/")
}
