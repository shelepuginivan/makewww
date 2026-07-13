package resource

import (
	"fmt"
	"time"

	"github.com/goccy/go-yaml"
)

type Metadata struct {
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	CreatedAt   time.Time `yaml:"created_at"`
	Draft       bool      `yaml:"draft"`
	Template    string    `yaml:"template"`
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
