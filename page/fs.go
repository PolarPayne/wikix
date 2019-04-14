package page

import (
	"bufio"
	"bytes"
	"log"
	"path/filepath"
	"sort"

	"github.com/pelletier/go-toml"
	"github.com/polarpayne/wikix/render"
	"github.com/polarpayne/wikix/types"
	"github.com/spf13/afero"
)

type FS struct {
	FS           afero.Fs
	Unsafe       bool
	pageRenderer *render.Page
}

func NewFS(fs, macroFS afero.Fs, unsafe bool) *FS {
	return &FS{fs, unsafe, render.NewPage(macroFS)}
}

func (s *FS) pageHeader(name string) (types.Page, error) {
	fp, err := s.FS.Open(name + ".md")
	if err != nil {
		return types.Page{}, err
	}

	r := bufio.NewScanner(fp)

	var header bytes.Buffer
	out := types.Page{}

	for r.Scan() {
		line := r.Text()
		if line == "+++" {
			err = toml.Unmarshal(header.Bytes(), &out)
			if err != nil {
				return types.Page{}, err
			}

			out.Name = name
			return out, nil
		}

		header.WriteString(line + "\n")
	}
	if err = r.Err(); err != nil {
		return types.Page{}, err
	}

	out.Name = name
	return out, nil
}

func (s *FS) PageAll() (types.PageAll, error) {
	files, err := afero.Glob(s.FS, "*.md")
	if err != nil {
		return nil, err
	}

	var out []string
	for _, file := range files {
		if ok, _ := afero.IsDir(s.FS, file); ok {
			continue
		}

		tmp := filepath.Base(file)
		// the last three letters are .md and are not part of the name
		tmp = tmp[:len(tmp)-3]
		out = append(out, tmp)
	}

	sort.Strings(out)
	return out, nil
}

func (s *FS) Page(name string) (types.Page, error) {
	fp, err := s.FS.Open(name + ".md")
	if err != nil {
		return types.Page{}, err
	}

	r := bufio.NewScanner(fp)

	var (
		inContent bool
		header    bytes.Buffer
		content   bytes.Buffer
	)

	for r.Scan() {
		line := r.Text()
		if !inContent && line == "+++" {
			inContent = true
			header.WriteString(content.String())
			content.Reset()
			continue
		}

		content.WriteString(line + "\n")
	}
	if err = r.Err(); err != nil {
		return types.Page{}, err
	}

	out := types.Page{}
	err = toml.Unmarshal(header.Bytes(), &out)
	if err != nil {
		return types.Page{}, err
	}

	out.Name = name
	out.Content, err = s.pageRenderer.Render(content.Bytes(), s.Unsafe)
	if err != nil {
		return types.Page{}, err
	}
	sort.Strings(out.Tags)

	return out, nil
}

func (s *FS) PageExists(name string) bool {
	out, err := afero.Exists(s.FS, name+".md")
	if err != nil {
		log.Panic(err)
	}

	return out
}

func (s *FS) TagAll() (types.TagAll, error) {
	pages, err := s.PageAll()
	if err != nil {
		return types.TagAll{}, err
	}

	tags := make(map[string]bool)

	for _, page := range pages {
		pageData, err := s.pageHeader(page)
		if err != nil {
			return types.TagAll{}, err
		}

		for _, tag := range pageData.Tags {
			tags[tag] = true
		}
	}

	out := types.TagAll{}
	for tag := range tags {
		out = append(out, tag)
	}
	sort.Strings(out)

	return out, nil
}

func (s *FS) Tag(name string) (types.Tag, error) {
	pages, err := s.PageAll()
	if err != nil {
		return types.Tag{}, err
	}

	out := types.Tag{}
	out.Name = name

	for _, page := range pages {
		pageData, err := s.pageHeader(page)
		if err != nil {
			return types.Tag{}, err
		}
		for _, tag := range pageData.Tags {
			if tag == name {
				out.Pages = append(out.Pages, pageData.Name)
				break
			}
		}
	}
	sort.Strings(out.Pages)

	return out, nil
}
