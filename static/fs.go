package static

import (
	"log"
	"path/filepath"

	"github.com/polarpayne/wikix/types"
	"github.com/spf13/afero"
)

type FS struct {
	FS afero.Fs
}

func NewFS(fs afero.Fs) *FS {
	return &FS{fs}
}

func (s *FS) StaticAll() (types.StaticAll, error) {
	files, err := afero.Glob(s.FS, "*")
	if err != nil {
		return nil, err
	}

	var out []string
	for _, file := range files {
		if ok, _ := afero.IsDir(s.FS, file); ok {
			continue
		}

		tmp := filepath.Base(file)
		out = append(out, tmp)
	}
	return out, nil
}

func (s *FS) Static(name string) (types.Static, error) {
	fp, err := s.FS.Open(name)
	if err != nil {
		return nil, err
	}

	return fp, nil
}

func (s *FS) StaticExists(name string) bool {
	out, err := afero.Exists(s.FS, name)
	if err != nil {
		log.Panic(err)
	}

	return out
}
