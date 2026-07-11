package source

import "time"

type Metadata struct {
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	CreatedAt   time.Time `yaml:"created_at"`
	Draft       bool      `yaml:"draft"`
	Template    string    `yaml:"template"`
}

type Document interface {
	Metadata() *Metadata
	Content() (string, error)
	CanonicalPath(base string) (string, error)
}
