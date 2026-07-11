package source

import "path/filepath"

type Path struct {
	relative string
}

func (p *Path) Absolute() string {
	return "/" + p.relative
}

func (p *Path) String() string {
	return p.Absolute()
}

func (p *Path) Relative() string {
	return p.relative
}

func (p *Path) Base() string {
	return filepath.Base(p.Absolute())
}

func (p *Path) Dir() string {
	return filepath.Dir(p.Absolute())
}
