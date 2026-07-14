package resource

import (
	"fmt"
	"time"

	"github.com/goccy/go-yaml"
)

type Metadata struct {
	Title     string    `yaml:"title"`
	Summary   string    `yaml:"summary"`
	CreatedAt time.Time `yaml:"created_at"`
	Order     int       `yaml:"order"`
	Draft     bool      `yaml:"draft"`
	Layout    string    `yaml:"layout"`
}

type WithMetadata interface {
	Metadata() *Metadata
}

func metadataFromYAML(data []byte) (*Metadata, error) {
	var metadata Metadata

	err := yaml.Unmarshal(data, &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to parse yaml: %w", err)
	}

	return &metadata, nil
}
