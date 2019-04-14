package types

import (
	"html/template"
	"io"
	"time"
)

type Error struct {
	Error string
	Code  int
}

type PageAll []string
type TagAll []string
type StaticAll []string

type Page struct {
	Name         string
	Content      template.HTML
	Tags         []string
	LastModified time.Time
	Created      time.Time
}

type Tag struct {
	Name  string
	Pages []string
}

type Static io.ReadCloser

type StorePage interface {
	PageAll() (PageAll, error)
	Page(string) (Page, error)
	PageExists(string) bool

	TagAll() (TagAll, error)
	Tag(string) (Tag, error)
}

type StoreStatic interface {
	StaticAll() (StaticAll, error)
	Static(string) (Static, error)
	StaticExists(string) bool
}
