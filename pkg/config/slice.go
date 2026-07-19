package config

import "fmt"

type SliceFlags []string

func (sf *SliceFlags) Set(value string) error {
	*sf = append(*sf, value)
	return nil
}

func (sf *SliceFlags) String() string {
	return fmt.Sprint(*sf)
}
