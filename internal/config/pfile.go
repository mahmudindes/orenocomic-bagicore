package config

import "os"

type pfile struct {
	path string
}

func (c pfile) ReadBytes() ([]byte, error) {
	return os.ReadFile(c.path)
}

func (c pfile) Read() (map[string]any, error) {
	return nil, nil
}
