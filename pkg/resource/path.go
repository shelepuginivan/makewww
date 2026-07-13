package resource

import (
	"path/filepath"
	"strings"
)

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

func (p *Path) Normalized() string {
	abs := p.Absolute()
	ext := filepath.Ext(abs)
	return strings.TrimSuffix(abs, ext)
}

func (p *Path) RelativeNormalized() string {
	ext := filepath.Ext(p.relative)
	return strings.TrimSuffix(p.relative, ext)
}

func (p *Path) Stem() string {
	base := p.Base()
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}
